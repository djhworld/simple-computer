package peripherals

import (
	"github.com/djhworld/simple-computer/circuit"
	"github.com/djhworld/simple-computer/components"
	"github.com/eiannone/keyboard"
)

const BUS_WIDTH = 16

type Keyboard struct {
	adapter *KeyboardAdapter
}

func NewKeyboard(ka *KeyboardAdapter) *Keyboard {
	k := new(Keyboard)
	k.adapter = ka
	return k
}

func (k *Keyboard) Update() {
	err := keyboard.Open()
	if err != nil {
		panic(err)
	}
	defer keyboard.Close()

	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		} else if key == keyboard.KeyEsc {
			break
		} else if key > 0 {
			setBus(k.adapter.keyboardInBus, uint16(key))
		} else {
			setBus(k.adapter.keyboardInBus, uint16(char))
		}
	}
}

type KeyboardAdapter struct {
	ioBus   *components.IOBus
	mainBus *components.Bus

	memoryBit       *components.Bit
	keyboardInBus   *components.Bus
	keycodeRegister components.Register

	andGate1            components.ANDGate8
	notGatesForAndGate1 [4]circuit.NOTGate

	andGate2            components.ANDGate3
	andGate3            components.ANDGate3
	notGatesForAndGate3 [2]circuit.NOTGate

	andGate4 circuit.ANDGate
}

func NewKeyboardAdapter() *KeyboardAdapter {
	return new(KeyboardAdapter)
}

func (k *KeyboardAdapter) Connect(ioBus *components.IOBus, mainBus *components.Bus) {
	k.ioBus = ioBus
	k.mainBus = mainBus
	k.memoryBit = components.NewBit()
	k.memoryBit.Update(false, true)
	k.memoryBit.Update(false, false)
	k.keyboardInBus = components.NewBus(BUS_WIDTH)
	k.andGate1 = *components.NewANDGate8()
	k.andGate2 = *components.NewANDGate3()
	k.andGate3 = *components.NewANDGate3()
	k.andGate4 = *circuit.NewANDGate()
	k.keycodeRegister = *components.NewRegister("KCR", k.keyboardInBus, k.mainBus)

	for i := range k.notGatesForAndGate1 {
		k.notGatesForAndGate1[i] = *circuit.NewNOTGate()
	}

	for i := range k.notGatesForAndGate3 {
		k.notGatesForAndGate3[i] = *circuit.NewNOTGate()
	}
}

func (k *KeyboardAdapter) Update() {
	k.updateKeycodeReg()
	k.update()
}

func (k *KeyboardAdapter) update() {
	k.notGatesForAndGate1[0].Update(k.mainBus.GetOutputWire(8))
	k.notGatesForAndGate1[1].Update(k.mainBus.GetOutputWire(9))
	k.notGatesForAndGate1[2].Update(k.mainBus.GetOutputWire(10))
	k.notGatesForAndGate1[3].Update(k.mainBus.GetOutputWire(11))

	k.andGate1.Update(
		k.notGatesForAndGate1[0].Output(),
		k.notGatesForAndGate1[1].Output(),
		k.notGatesForAndGate1[2].Output(),
		k.notGatesForAndGate1[3].Output(),
		k.mainBus.GetOutputWire(12),
		k.mainBus.GetOutputWire(13),
		k.mainBus.GetOutputWire(14),
		k.mainBus.GetOutputWire(15),
	)

	k.andGate2.Update(
		k.ioBus.GetOutputWire(components.CLOCK_SET),
		k.ioBus.GetOutputWire(components.DATA_OR_ADDRESS),
		k.ioBus.GetOutputWire(components.MODE),
	)

	k.notGatesForAndGate3[0].Update(k.ioBus.GetOutputWire(components.DATA_OR_ADDRESS))
	k.notGatesForAndGate3[1].Update(k.ioBus.GetOutputWire(components.MODE))

	k.andGate3.Update(
		k.ioBus.GetOutputWire(components.CLOCK_ENABLE),
		k.notGatesForAndGate3[0].Output(),
		k.notGatesForAndGate3[1].Output(),
	)

	k.memoryBit.Update(k.andGate1.Output(), k.andGate2.Output())
	k.andGate4.Update(k.memoryBit.Get(), k.andGate3.Output())
}

func (k *KeyboardAdapter) updateKeycodeReg() {
	if k.andGate4.Output() {
		k.keycodeRegister.Set()

		k.keycodeRegister.Enable()
		k.keycodeRegister.Update()
		k.keycodeRegister.Disable()

		// clear the register once everything is out
		setBus(k.keyboardInBus, 0x00)
		k.keycodeRegister.Update()
		k.keycodeRegister.Unset()
		k.keycodeRegister.Update()
	}
}

func setBus(b *components.Bus, value uint16) {
	var x = 0
	for i := BUS_WIDTH - 1; i >= 0; i-- {
		r := (value & (1 << uint16(x)))
		if r != 0 {
			b.SetInputWire(i, true)
		} else {
			b.SetInputWire(i, false)
		}

		x++
	}
}
