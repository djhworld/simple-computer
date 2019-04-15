package cpu

import (
	"fmt"
	"time"

	"github.com/djhworld/simple-computer/alu"
	"github.com/djhworld/simple-computer/circuit"
	"github.com/djhworld/simple-computer/components"
	"github.com/djhworld/simple-computer/memory"
	"github.com/paulbellamy/ratecounter"
)

// ADDS
// ----------------------
// 0x80 = ADD R0, R0
// 0x81 = ADD R0, R1
// 0x82 = ADD R0, R2
// 0x83 = ADD R0, R3

// 0x84 = ADD R1, R0
// 0x85 = ADD R1, R1
// 0x86 = ADD R1, R2
// 0x87 = ADD R1, R3

// 0x88 = ADD R2, R0
// 0x89 = ADD R2, R1
// 0x8A = ADD R2, R2
// 0x8B = ADD R2, R3

// 0x8C = ADD R3, R0
// 0x8D = ADD R3, R1
// 0x8E = ADD R3, R2
// 0x8F = ADD R3, R3

// SHR
// ----------------------
// 0x90 R0
// 0x95 R1
// 0x9A R2
// 0x9F R3

// SHL
// ----------------------
// 0xA0 R0
// 0xA5 R1
// 0xAA R2
// 0xAF R3

// NOT
// ----------------------
// 0xB0 R0
// 0xB5 R1
// 0xBA R2
// 0xBF R3

// ANDS
// ----------------------
// 0xC0 = AND R0, R0
// 0xC1 = AND R0, R1
// 0xC2 = AND R0, R2
// 0xC3 = AND R0, R3

// 0xC4 = AND R1, R0
// 0xC5 = AND R1, R1
// 0xC6 = AND R1, R2
// 0xC7 = AND R1, R3

// 0xC8 = AND R2, R0
// 0xC9 = AND R2, R1
// 0xCA = AND R2, R2
// 0xCB = AND R2, R3

// 0xCC = AND R3, R0
// 0xCD = AND R3, R1
// 0xCE = AND R3, R2
// 0xCF = AND R3, R3

// ORS
// ----------------------
// 0xD0 = OR R0, R0
// 0xD1 = OR R0, R1
// 0xD2 = OR R0, R2
// 0xD3 = OR R0, R3

// 0xD4 = OR R1, R0
// 0xD5 = OR R1, R1
// 0xD6 = OR R1, R2
// 0xD7 = OR R1, R3

// 0xD8 = OR R2, R0
// 0xD9 = OR R2, R1
// 0xDA = OR R2, R2
// 0xDB = OR R2, R3

// 0xDC = OR R3, R0
// 0xDD = OR R3, R1
// 0xDE = OR R3, R2
// 0xDF = OR R3, R3

// XORS
// ----------------------
// 0xE0 = XOR R0, R0
// 0xE1 = XOR R0, R1
// 0xE2 = XOR R0, R2
// 0xE3 = XOR R0, R3

// 0xE4 = XOR R1, R0
// 0xE5 = XOR R1, R1
// 0xE6 = XOR R1, R2
// 0xE7 = XOR R1, R3

// 0xE8 = XOR R2, R0
// 0xE9 = XOR R2, R1
// 0xEA = XOR R2, R2
// 0xEB = XOR R2, R3

// 0xEC = XOR R3, R0
// 0xED = XOR R3, R1
// 0xEE = XOR R3, R2
// 0xEF = XOR R3, R3

// CMP
// ----------------------
// 0xF0 = CMP R0, R0
// 0xF1 = CMP R0, R1
// 0xF2 = CMP R0, R2
// 0xF3 = CMP R0, R3

// 0xF4 = CMP R1, R0
// 0xF5 = CMP R1, R1
// 0xF6 = CMP R1, R2
// 0xF7 = CMP R1, R3

// 0xF8 = CMP R2, R0
// 0xF9 = CMP R2, R1
// 0xFA = CMP R2, R2
// 0xFB = CMP R2, R3

// 0xFC = CMP R3, R0
// 0xFD = CMP R3, R1
// 0xFE = CMP R3, R2
// 0xFF = CMP R3, R3

type Enableable interface {
	Enable()
	Disable()
}

type Settable interface {
	Set()
	Unset()
}

type Updatable interface {
	Update()
}

type CPU struct {
	gpReg0 components.Register
	gpReg1 components.Register
	gpReg2 components.Register
	gpReg3 components.Register
	tmp    components.Register
	acc    components.Register
	ir     components.Register
	iar    components.Register

	instructionDecoderEnables [2]components.Decoder2x4
	instructionDecoderSet     components.Decoder2x4

	regEnableGates [3]circuit.ANDGate
	regSetGates    [6]circuit.ANDGate

	gpRegEnableANDGates [8]components.ANDGate3
	gpRegEnableORGates  [4]circuit.ORGate
	gpRegSetANDGates    [4]components.ANDGate3

	registerBSetANDGate1 components.ANDGate3
	registerBSetANDGate2 components.ANDGate3
	registerBSetNOTGate  circuit.NOTGate
	registerBSet         circuit.Wire

	registerAEnable        circuit.Wire
	registerAEnableANDGate circuit.ANDGate
	registerBEnable        circuit.Wire
	registerBEnableANDGate circuit.ANDGate

	aluOpAndGates [3]components.ANDGate3

	busOne       components.BusOne
	mainBus      *components.Bus
	tmpBus       *components.Bus
	busOneOutput *components.Bus
	controlBus   *components.Bus
	accBus       *components.Bus
	memory       *memory.Memory256
	alu          *alu.ALU
	stepper      *components.Stepper
	clock        *components.Clock
}

func NewCPU(mainBus *components.Bus, memory *memory.Memory256) *CPU {
	c := new(CPU)

	c.controlBus = components.NewBus()

	c.mainBus = mainBus
	c.gpReg0 = *components.NewRegister("R0", c.mainBus, c.mainBus)
	c.gpReg1 = *components.NewRegister("R1", c.mainBus, c.mainBus)
	c.gpReg2 = *components.NewRegister("R2", c.mainBus, c.mainBus)
	c.gpReg3 = *components.NewRegister("R3", c.mainBus, c.mainBus)
	c.ir = *components.NewRegister("IR", c.mainBus, c.controlBus)
	c.ir.Disable()

	c.iar = *components.NewRegister("IAR", c.mainBus, c.mainBus)
	c.instructionDecoderEnables[0] = *components.NewDecoder2x4()
	c.instructionDecoderEnables[1] = *components.NewDecoder2x4()
	c.instructionDecoderSet = *components.NewDecoder2x4()

	c.tmpBus = components.NewBus()

	c.tmp = *components.NewRegister("TMP", c.mainBus, c.tmpBus)
	c.tmp.Enable()

	//TODO this needs a receiving bus
	c.busOneOutput = components.NewBus()
	c.busOne = *components.NewBusOne(c.tmpBus, c.busOneOutput)

	c.accBus = components.NewBus()
	c.acc = *components.NewRegister("ACC", c.accBus, c.mainBus)

	c.stepper = components.NewStepper()
	c.clock = &components.Clock{}
	c.alu = alu.NewALU(c.mainBus, c.busOneOutput, c.accBus)
	c.memory = memory

	c.registerBSetANDGate1 = *components.NewANDGate3()
	c.registerBSetANDGate2 = *components.NewANDGate3()
	c.registerBSetNOTGate = *circuit.NewNOTGate()

	c.registerAEnableANDGate = *circuit.NewANDGate()
	c.registerBEnableANDGate = *circuit.NewANDGate()

	for i := range c.regEnableGates {
		c.regEnableGates[i] = *circuit.NewANDGate()
	}

	for i := range c.regSetGates {
		c.regSetGates[i] = *circuit.NewANDGate()
	}

	for i := range c.gpRegEnableORGates {
		c.gpRegEnableORGates[i] = *circuit.NewORGate()
	}

	for i := range c.gpRegEnableANDGates {
		c.gpRegEnableANDGates[i] = *components.NewANDGate3()
	}

	for i := range c.gpRegSetANDGates {
		c.gpRegSetANDGates[i] = *components.NewANDGate3()
	}

	for i := range c.aluOpAndGates {
		c.aluOpAndGates[i] = *components.NewANDGate3()
	}

	return c
}

func (c *CPU) Run() {
	setBus(c.mainBus, 0x00)

	c.iar.Set()
	c.iar.Update()
	c.iar.Unset()
	c.iar.Update()

	var i byte
	var q byte = 0xFF
	for i = 0x00; i < 0xFF; i++ {
		c.memory.AddressRegister.Set()
		setBus(c.mainBus, i)
		c.memory.Update()

		c.memory.AddressRegister.Unset()
		c.memory.Update()

		setBus(c.mainBus, 0x8B)
		c.memory.Set()
		c.memory.Update()

		c.memory.Unset()
		c.memory.Update()

		q--
	}

	c.gpReg0.Set()
	c.gpReg0.Update()
	setBus(c.mainBus, 0x01)
	c.gpReg0.Update()
	c.gpReg0.Unset()
	c.gpReg0.Update()

	c.gpReg1.Set()
	c.gpReg1.Update()
	setBus(c.mainBus, 0x02)
	c.gpReg1.Update()
	c.gpReg1.Unset()
	c.gpReg1.Update()

	c.gpReg2.Set()
	c.gpReg2.Update()
	setBus(c.mainBus, 0x03)
	c.gpReg2.Update()
	c.gpReg2.Unset()
	c.gpReg2.Update()

	c.gpReg3.Set()
	c.gpReg3.Update()
	setBus(c.mainBus, 0x01)
	c.gpReg3.Update()
	c.gpReg3.Unset()
	c.gpReg3.Update()


	r := 0
	counter := ratecounter.NewRateCounter(1 * time.Second)
	c.clock.Start()
	clockState := false
	for {
		select {
		case <-c.clock.BaseClock.C:
			counter.Incr(1)
			r++
			if clockState {
				clockState = false
			} else {
				clockState = true
			}
			c.step(clockState)

			if r%3000 == 0 {
				fmt.Println("CPU Speed:", counter.Rate(), "hz")
			}

		}
	}

}

func (c *CPU) String() string {
	return fmt.Sprintf("stepper: %s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\nbus1: %s\n%s\n",
		c.stepper.String(),
		c.iar.String(),
		c.memory.AddressRegister.String(),
		c.ir.String(),
		c.acc.String(),
		c.tmp.String(),
		c.gpReg0.String(),
		c.gpReg1.String(),
		c.gpReg2.String(),
		c.gpReg3.String(),
		c.busOne.String(),
		c.alu.String(),
	)
}

func (c *CPU) step(clockState bool) {

	c.stepper.Update(clockState)

	//TODO this is problematic
	c.runEnable(clockState)
	c.updateStates()
	if clockState {
		c.runEnable(false)
		c.updateStates()
	}

	c.runSet(clockState)
	c.updateStates()
	if clockState {
		c.runSet(false)
		c.updateStates()
	}

//	fmt.Println(c)
}

func (c *CPU) updateStates() {
	// IAR
	runUpdateOn(&c.iar)

	// MAR
	runUpdateOn(&c.memory.AddressRegister)

	// IR
	runUpdateOn(&c.ir)

	// ACC
	runUpdateOn(&c.acc)

	// RAM
	runUpdateOn(c.memory)

	// TMP
	runUpdateOn(&c.tmp)

	// BUS1
	runUpdateOn(&c.busOne)

	// R0
	runUpdateOn(&c.gpReg0)

	// R1
	runUpdateOn(&c.gpReg1)

	// R2
	runUpdateOn(&c.gpReg2)

	// R3
	runUpdateOn(&c.gpReg3)

	c.updateALU()
}

func (c *CPU) updateALU() {
	//update ALU operation base on instruction register
	c.aluOpAndGates[2].Update(c.ir.Bit(1), c.ir.Bit(0), c.stepper.GetOutputWire(4))
	c.aluOpAndGates[1].Update(c.ir.Bit(2), c.ir.Bit(0), c.stepper.GetOutputWire(4))
	c.aluOpAndGates[0].Update(c.ir.Bit(3), c.ir.Bit(0), c.stepper.GetOutputWire(4))

	c.alu.Op[2].Update(c.aluOpAndGates[2].Output())
	c.alu.Op[1].Update(c.aluOpAndGates[1].Output())
	c.alu.Op[0].Update(c.aluOpAndGates[0].Output())
	c.alu.Update()

}

func (c *CPU) runEnable(state bool) {
	// IAR
	c.regEnableGates[0].Update(state, c.stepper.GetOutputWire(0))
	updateEnableStatus(&c.iar, c.regEnableGates[0].Output())

	// RAM
	c.regEnableGates[1].Update(state, c.stepper.GetOutputWire(1))
	updateEnableStatus(c.memory, c.regEnableGates[1].Output())

	// ACC
	c.regEnableGates[2].Update(state, c.stepper.GetOutputWire(2))
	updateEnableStatus(&c.acc, c.regEnableGates[2].Output())

	// BUS1
	updateEnableStatus(&c.busOne, c.stepper.GetOutputWire(0))

	c.runEnableOnRegisterB()
	c.runEnableOnRegisterA()
	c.runEnableGeneralPurposeRegisters(state)
}

func (c *CPU) runEnableOnRegisterB() {
	c.registerBEnableANDGate.Update(c.ir.Bit(0), c.stepper.GetOutputWire(3))
	c.registerBEnable.Update(c.registerBEnableANDGate.Output())

	// TMP - hmmmm
	updateSetStatus(&c.tmp, c.registerBEnableANDGate.Output())
}

func (c *CPU) runEnableOnRegisterA() {
	c.registerAEnableANDGate.Update(c.ir.Bit(0), c.stepper.GetOutputWire(4))
	c.registerAEnable.Update(c.registerAEnableANDGate.Output())

	// ACC - hmmmm
	updateSetStatus(&c.acc, c.registerAEnableANDGate.Output())
}

func (c *CPU) runEnableGeneralPurposeRegisters(state bool) {

	c.instructionDecoderEnables[0].Update(c.ir.Bit(6), c.ir.Bit(7))
	c.instructionDecoderEnables[1].Update(c.ir.Bit(4), c.ir.Bit(5))

	// R0
	c.gpRegEnableANDGates[0].Update(state, c.registerBEnable.Get(), c.instructionDecoderEnables[0].GetOutputWire(0))
	c.gpRegEnableANDGates[4].Update(state, c.registerAEnable.Get(), c.instructionDecoderEnables[1].GetOutputWire(0))
	c.gpRegEnableORGates[0].Update(c.gpRegEnableANDGates[4].Output(), c.gpRegEnableANDGates[0].Output())
	updateEnableStatus(&c.gpReg0, c.gpRegEnableORGates[0].Output())

	// R1
	c.gpRegEnableANDGates[1].Update(state, c.registerBEnable.Get(), c.instructionDecoderEnables[0].GetOutputWire(1))
	c.gpRegEnableANDGates[5].Update(state, c.registerAEnable.Get(), c.instructionDecoderEnables[1].GetOutputWire(1))
	c.gpRegEnableORGates[1].Update(c.gpRegEnableANDGates[5].Output(), c.gpRegEnableANDGates[1].Output())
	updateEnableStatus(&c.gpReg1, c.gpRegEnableORGates[1].Output())

	// R2
	// this register should be enabled at some point but isn't....
	c.gpRegEnableANDGates[2].Update(state, c.registerBEnable.Get(), c.instructionDecoderEnables[0].GetOutputWire(2))
	c.gpRegEnableANDGates[6].Update(state, c.registerAEnable.Get(), c.instructionDecoderEnables[1].GetOutputWire(2))
	c.gpRegEnableORGates[2].Update(c.gpRegEnableANDGates[6].Output(), c.gpRegEnableANDGates[2].Output())
	updateEnableStatus(&c.gpReg2, c.gpRegEnableORGates[2].Output())

	// R3
	c.gpRegEnableANDGates[3].Update(state, c.registerBEnable.Get(), c.instructionDecoderEnables[0].GetOutputWire(3))
	c.gpRegEnableANDGates[7].Update(state, c.registerAEnable.Get(), c.instructionDecoderEnables[1].GetOutputWire(3))
	c.gpRegEnableORGates[3].Update(c.gpRegEnableANDGates[7].Output(), c.gpRegEnableANDGates[3].Output())
	updateEnableStatus(&c.gpReg3, c.gpRegEnableORGates[3].Output())
}

func (c *CPU) runSet(state bool) {

	// IR
	c.regSetGates[0].Update(state, c.stepper.GetOutputWire(1))
	updateSetStatus(&c.ir, c.regSetGates[0].Output())

	// MAR
	c.regSetGates[1].Update(state, c.stepper.GetOutputWire(0))
	updateSetStatus(&c.memory.AddressRegister, c.regSetGates[1].Output())

	// IAR
	c.regSetGates[2].Update(state, c.stepper.GetOutputWire(2))
	updateSetStatus(&c.iar, c.regSetGates[2].Output())

	// ACC
	c.regSetGates[3].Update(state, c.stepper.GetOutputWire(0))
	updateSetStatus(&c.acc, c.regSetGates[3].Output())

	// TMP
	// - not in original document...
	c.regSetGates[5].Update(state, c.stepper.GetOutputWire(0))
	updateSetStatus(&c.tmp, c.regSetGates[5].Output())

	// RAM
	c.regSetGates[4].Update(state, false)
	updateSetStatus(c.memory, c.regSetGates[4].Output())

	c.runSetOnRegisterB()
	c.runSetGeneralPurposeRegisters(state)
}

func (c *CPU) runSetOnRegisterB() {
	c.registerBSetANDGate1.Update(c.ir.Bit(3), c.ir.Bit(2), c.ir.Bit(1))
	c.registerBSetNOTGate.Update(c.registerBSetANDGate1.Output())
	c.registerBSetANDGate2.Update(c.stepper.GetOutputWire(5), c.ir.Bit(0), c.registerBSetNOTGate.Output())
	c.registerBSet.Update(c.registerBSetANDGate2.Output())

	// ACC - hmmmmm
	updateEnableStatus(&c.acc, c.registerBSetANDGate2.Output())
}

func (c *CPU) runSetGeneralPurposeRegisters(state bool) {
	c.instructionDecoderSet.Update(c.ir.Bit(6), c.ir.Bit(7))

	// R0
	c.gpRegSetANDGates[0].Update(state, c.registerBSet.Get(), c.instructionDecoderSet.GetOutputWire(0))
	updateSetStatus(&c.gpReg0, c.gpRegSetANDGates[0].Output())

	// R1
	c.gpRegSetANDGates[1].Update(state, c.registerBSet.Get(), c.instructionDecoderSet.GetOutputWire(1))
	updateSetStatus(&c.gpReg1, c.gpRegSetANDGates[1].Output())

	// R2
	c.gpRegSetANDGates[2].Update(state, c.registerBSet.Get(), c.instructionDecoderSet.GetOutputWire(2))
	updateSetStatus(&c.gpReg2, c.gpRegSetANDGates[2].Output())

	// R3
	c.gpRegSetANDGates[3].Update(state, c.registerBSet.Get(), c.instructionDecoderSet.GetOutputWire(3))
	updateSetStatus(&c.gpReg3, c.gpRegSetANDGates[3].Output())
}

func runUpdateOn(component Updatable) {
	component.Update()
}

func updateEnableStatus(component Enableable, state bool) {
	if state {
		component.Enable()
	} else {
		component.Disable()
	}
}

func updateSetStatus(component Settable, state bool) {
	if state {
		component.Set()
	} else {
		component.Unset()
	}
}

func setBus(b *components.Bus, value byte) {
	var x = 0
	for i := 7; i >= 0; i-- {
		r := (value & (1 << byte(x)))
		if r != 0 {
			b.SetInputWire(i, true)
		} else {
			b.SetInputWire(i, false)
		}

		x++
	}
}
