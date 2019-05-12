package io

import (
	"log"
	"time"

	"github.com/djhworld/simple-computer/circuit"
	"github.com/djhworld/simple-computer/components"
)

// [cpu] -------> display adapter --------> display RAM <--------- screen control ---------> [screenChannel]
//       write                     write                   read                     write
type DisplayAdapter struct {
	ioBus      *components.IOBus
	mainBus    *components.Bus
	screenBus  *components.Bus
	displayRAM *displayRAM

	displayAdapterActiveBit *components.Bit
	inputMAROutBus          *components.Bus

	addressSelectAndGate  components.ANDGate8
	addressSelectNOTGates [5]circuit.NOTGate

	isAddressOutputModeGate components.ANDGate3

	inputMARSetGate      components.ANDGate5
	inputMarSetNOTGates  [2]circuit.NOTGate
	writeToRAM           *components.Bit
	writeToRAMToggleGate circuit.NOTGate

	displayRAMSetGate components.ANDGate5
}

func NewDisplaydAdapter() *DisplayAdapter {
	d := new(DisplayAdapter)
	return d
}

func (k *DisplayAdapter) Connect(ioBus *components.IOBus, mainBus *components.Bus) {
	k.ioBus = ioBus
	k.mainBus = mainBus
	k.screenBus = components.NewBus(BUS_WIDTH)
	k.displayRAM = newDisplayRAM(k.mainBus, k.screenBus)

	k.displayAdapterActiveBit = components.NewBit()
	k.displayAdapterActiveBit.Update(false, true)
	k.displayAdapterActiveBit.Update(false, false)
	k.addressSelectAndGate = *components.NewANDGate8()
	k.isAddressOutputModeGate = *components.NewANDGate3()
	k.inputMARSetGate = *components.NewANDGate5()
	k.displayRAMSetGate = *components.NewANDGate5()
	k.inputMAROutBus = components.NewBus(BUS_WIDTH)

	k.writeToRAM = components.NewBit()
	k.writeToRAM.Update(false, true)
	k.writeToRAM.Update(false, false)
	k.writeToRAMToggleGate = *circuit.NewNOTGate()

	for i := range k.addressSelectNOTGates {
		k.addressSelectNOTGates[i] = *circuit.NewNOTGate()
	}

	for i := range k.inputMarSetNOTGates {
		k.inputMarSetNOTGates[i] = *circuit.NewNOTGate()
	}
}

func (k *DisplayAdapter) Update() {
	// check if bus = 0x0007
	k.addressSelectNOTGates[0].Update(k.mainBus.GetOutputWire(8))
	k.addressSelectNOTGates[1].Update(k.mainBus.GetOutputWire(9))
	k.addressSelectNOTGates[2].Update(k.mainBus.GetOutputWire(10))
	k.addressSelectNOTGates[3].Update(k.mainBus.GetOutputWire(11))
	k.addressSelectNOTGates[4].Update(k.mainBus.GetOutputWire(12))
	k.addressSelectAndGate.Update(
		k.addressSelectNOTGates[0].Output(),
		k.addressSelectNOTGates[1].Output(),
		k.addressSelectNOTGates[2].Output(),
		k.addressSelectNOTGates[3].Output(),
		k.addressSelectNOTGates[4].Output(),
		k.mainBus.GetOutputWire(13),
		k.mainBus.GetOutputWire(14),
		k.mainBus.GetOutputWire(15),
	)

	k.isAddressOutputModeGate.Update(
		k.ioBus.IsSet(),
		k.ioBus.IsAddressMode(),
		k.ioBus.IsOutputMode(),
	)

	k.displayAdapterActiveBit.Update(k.addressSelectAndGate.Output(), k.isAddressOutputModeGate.Output())

	// switch between writing to display RAM and writing to address register
	if k.writeToRAM.Get() {
		k.writeToDisplayRAM()
	} else {
		k.writeToInputMAR()
	}

}

func (k *DisplayAdapter) toggleWriteToRAM() {
	k.writeToRAMToggleGate.Update(k.writeToRAM.Get())
	k.writeToRAM.Update(k.writeToRAMToggleGate.Output(), true)
	k.writeToRAM.Update(false, false)
}

func (k *DisplayAdapter) writeToInputMAR() {
	// if writeToRAM == false then put bus contents in Input-MAR
	k.inputMarSetNOTGates[0].Update(k.writeToRAM.Get())

	k.inputMARSetGate.Update(
		k.ioBus.IsDataMode(),
		k.ioBus.IsSet(),
		k.ioBus.IsOutputMode(),
		k.displayAdapterActiveBit.Get(),
		k.inputMarSetNOTGates[0].Output(),
	)

	if k.inputMARSetGate.Output() {
		k.displayRAM.InputAddressRegister.Set()
		k.displayRAM.InputAddressRegister.Update()
		k.displayRAM.InputAddressRegister.Unset()
		k.displayRAM.InputAddressRegister.Update()
		k.toggleWriteToRAM()
	}
}

func (k *DisplayAdapter) writeToDisplayRAM() {
	// if writeToRAM == true then put bus contents in RAM
	k.displayRAMSetGate.Update(
		k.ioBus.IsDataMode(),
		k.ioBus.IsSet(),
		k.ioBus.IsOutputMode(),
		k.displayAdapterActiveBit.Get(),
		k.writeToRAM.Get(),
	)

	if k.displayRAMSetGate.Output() {
		k.displayRAM.Set()
		k.displayRAM.UpdateIncoming()
		k.displayRAM.Unset()
		k.displayRAM.UpdateIncoming()
		k.toggleWriteToRAM()
	}
}

func (k *DisplayAdapter) String() string {
	return ""
}

type ScreenControl struct {
	adapter    *DisplayAdapter
	inputBus   *components.Bus
	outputChan chan *[160][240]byte

	clock <-chan time.Time
	quit  chan bool

	//y, x
	output [160][240]byte
}

func NewScreenControl(adapter *DisplayAdapter, outputChan chan *[160][240]byte, quit chan bool) *ScreenControl {
	s := new(ScreenControl)
	s.adapter = adapter
	s.clock = time.Tick(33 * time.Millisecond)
	s.quit = quit
	s.outputChan = outputChan
	return s
}

func (s *ScreenControl) Run() {
	for {
		select {
		case <-s.quit:
			log.Println("Stopping screen control")
			return
		default:
			<-s.clock
			s.Update()
			s.outputChan <- &s.output
		}
	}
}

func (s *ScreenControl) Update() {
	widthInBytes := uint16(30) // 30 * 8 = 240

	y := uint16(0)
	for vertical := uint16(0x0000); vertical < 0x12C0; vertical += widthInBytes {
		x := uint16(0)
		for horizontal := uint16(0x0000); horizontal < widthInBytes; horizontal++ {
			s.setOutputRAMAddress(vertical + horizontal)
			s.renderPixelsFromRAM(y, x)
			x += 8
		}
		y++
	}
}

func (s *ScreenControl) setOutputRAMAddress(address uint16) {
	s.adapter.screenBus.SetValue(address)
	s.adapter.displayRAM.OutputAddressRegister.Set()
	s.adapter.displayRAM.OutputAddressRegister.Update()
	s.adapter.displayRAM.OutputAddressRegister.Unset()
	s.adapter.displayRAM.OutputAddressRegister.Update()
}

func (s *ScreenControl) renderPixelsFromRAM(y, x uint16) {
	s.adapter.displayRAM.Enable()
	s.adapter.displayRAM.UpdateOutgoing()

	for b := 8; b < 16; b++ {
		if s.adapter.screenBus.GetOutputWire(b) {
			s.output[y][x] = 0x01
		} else {
			s.output[y][x] = 0x00
		}
		x++
	}

	s.adapter.displayRAM.Disable()
	s.adapter.displayRAM.UpdateOutgoing()
}
