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

// LOADS
// ----------------------
// arg A = memory address to load from
// arg B = register to store value in
// 0x00 = LD R0, R0
// 0x01 = LD R0, R1
// 0x02 = LD R0, R2
// 0x03 = LD R0, R3

// 0x04 = LD R1, R0
// 0x05 = LD R1, R1
// 0x06 = LD R1, R2
// 0x07 = LD R1, R3

// 0x08 = LD R2, R0
// 0x09 = LD R2, R1
// 0x0A = LD R2, R2
// 0x0B = LD R2, R3

// 0x0C = LD R3, R0
// 0x0D = LD R3, R1
// 0x0E = LD R3, R2
// 0x0F = LD R3, R3

// STORES
// ----------------------
// arg A = memory address for value
// arg B = value to store in memory
// 0x10 = ST R0, R0
// 0x11 = ST R0, R1
// 0x12 = ST R0, R2
// 0x13 = ST R0, R3

// 0x14 = ST R1, R0
// 0x15 = ST R1, R1
// 0x16 = ST R1, R2
// 0x17 = ST R1, R3

// 0x18 = ST R2, R0
// 0x19 = ST R2, R1
// 0x1A = ST R2, R2
// 0x1B = ST R2, R3

// 0x1C = ST R3, R0
// 0x1D = ST R3, R1
// 0x1E = ST R3, R2
// 0x1F = ST R3, R3

// DATA
// put value in memory into register (2 byte instruction)
// ----------------------
// 0x20 = DATA R0
// 0x21 = DATA R1
// 0x22 = DATA R2
// 0x23 = DATA R3

// JMPR
// set instruction address register to value in register
// ----------------------
// 0x30 = JMPR R0
// 0x31 = JMPR R1
// 0x32 = JMPR R2
// 0x33 = JMPR R3

// JMP
// set instruction address register to next byte (2 byte instruction)
// ----------------------
// 0x40 = JMP <value>

// JMP(CAEZ)
// set instruction address register to next byte (2 byte instruction)
// jump if <flag(s)> are true
// ----------------------
// 0x51 = JMPZ <value>
// 0x52 = JMPE <value>
// 0x53 = JMPEZ <value>
// 0x54 = JMPA <value>
// 0x55 = JMPAZ <value>
// 0x56 = JMPAE <value>
// 0x57 = JMPAEZ <value>
// 0x58 = JMPC <value>
// 0x59 = JMPCZ <value>
// 0x5A = JMPCE <value>
// 0x5B = JMPCEZ <value>
// 0x5C = JMPCA <value>
// 0x5D = JMPCAZ <value>
// 0x5E = JMPCAE <value>
// 0x5F = JMPCAEZ <value>

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
// 0x90 = SHR R0
// 0x95 = SHR R1
// 0x9A = SHR R2
// 0x9F = SHR R3

// SHL
// ----------------------
// 0xA0 = SHL R0
// 0xA5 = SHL R1
// 0xAA = SHL R2
// 0xAF = SHL R3

// NOT
// ----------------------
// 0xB0 = NOT R0
// 0xB5 = NOT R1
// 0xBA = NOT R2
// 0xBF = NOT R3

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

type InstructionDecoder3x8 struct {
	decoder       components.Decoder3x8
	selectorGates [8]circuit.ANDGate
	bit0NOTGate   circuit.NOTGate
}

func NewInstructionDecoder3x8() *InstructionDecoder3x8 {
	i := new(InstructionDecoder3x8)
	i.decoder = *components.NewDecoder3x8()
	for n := range i.selectorGates {
		i.selectorGates[n] = *circuit.NewANDGate()
	}
	i.bit0NOTGate = *circuit.NewNOTGate()
	return i

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
	flags  components.Register

	instrDecoder3x8              InstructionDecoder3x8
	instructionDecoderEnables2x4 [2]components.Decoder2x4
	instructionDecoderSet2x4     components.Decoder2x4

	irInstructionANDGate components.ANDGate3
	irInstructionNOTGate circuit.NOTGate

	registerAEnableORGate components.ORGate3
	registerBEnableORGate components.ORGate4
	registerBSetORGate    components.ORGate4
	registerAEnable       circuit.Wire
	registerBEnable       circuit.Wire
	accEnableORGate       components.ORGate4
	accEnableANDGate      circuit.ANDGate
	busOneEnableORGate    components.ORGate4
	iarEnableORGate       components.ORGate4
	iarEnableANDGate      circuit.ANDGate
	ramEnableORGate       components.ORGate5
	ramEnableANDGate      circuit.ANDGate
	gpRegEnableANDGates   [8]components.ANDGate3
	gpRegEnableORGates    [4]circuit.ORGate
	gpRegSetANDGates      [4]components.ANDGate3

	irSetANDGate    circuit.ANDGate
	marSetORGate    components.ORGate6
	marSetANDGate   circuit.ANDGate
	iarSetORGate    components.ORGate6
	iarSetANDGate   circuit.ANDGate
	accSetORGate    components.ORGate4
	accSetANDGate   circuit.ANDGate
	ramSetANDGate   circuit.ANDGate
	tmpSetANDGate   circuit.ANDGate
	flagsSetORGate  circuit.ORGate
	flagsSetANDGate circuit.ANDGate
	registerBSet    circuit.Wire

	flagStateGates  [4]circuit.ANDGate
	flagStateORGate components.ORGate4

	step4Gates     [9]circuit.ANDGate
	step5Gates     [7]circuit.ANDGate
	step6Gates     [2]components.ANDGate3
	step6Gates2And circuit.ANDGate

	aluOpAndGates [3]components.ANDGate3

	busOne        components.BusOne
	mainBus       *components.Bus
	tmpBus        *components.Bus
	busOneOutput  *components.Bus
	controlBus    *components.Bus
	accBus        *components.Bus
	aluToFlagsBus *components.Bus
	flagsBus      *components.Bus
	memory        *memory.Memory256
	alu           *alu.ALU
	stepper       *components.Stepper
	clock         *components.Clock
}

func NewCPU(mainBus *components.Bus, memory *memory.Memory256) *CPU {
	c := new(CPU)

	// REGISTERS
	c.controlBus = components.NewBus()
	c.mainBus = mainBus
	c.gpReg0 = *components.NewRegister("R0", c.mainBus, c.mainBus)
	c.gpReg1 = *components.NewRegister("R1", c.mainBus, c.mainBus)
	c.gpReg2 = *components.NewRegister("R2", c.mainBus, c.mainBus)
	c.gpReg3 = *components.NewRegister("R3", c.mainBus, c.mainBus)
	c.ir = *components.NewRegister("IR", c.mainBus, c.controlBus)
	c.ir.Disable()
	c.iar = *components.NewRegister("IAR", c.mainBus, c.mainBus)

	// Decoders
	c.instructionDecoderEnables2x4[0] = *components.NewDecoder2x4()
	c.instructionDecoderEnables2x4[1] = *components.NewDecoder2x4()
	c.instructionDecoderSet2x4 = *components.NewDecoder2x4()

	c.instrDecoder3x8 = *NewInstructionDecoder3x8()

	// FLAGS
	c.aluToFlagsBus = components.NewBus()
	c.flagsBus = components.NewBus()
	c.flags = *components.NewRegister("FLAGS", c.aluToFlagsBus, c.flagsBus)
	c.flags.Enable()

	// TMP
	c.tmpBus = components.NewBus()
	c.tmp = *components.NewRegister("TMP", c.mainBus, c.tmpBus)
	c.tmp.Enable()

	// BUS 1
	c.busOneOutput = components.NewBus()
	c.busOne = *components.NewBusOne(c.tmpBus, c.busOneOutput)

	// ACC
	c.accBus = components.NewBus()
	c.acc = *components.NewRegister("ACC", c.accBus, c.mainBus)

	// IR Register Output
	c.irInstructionANDGate = *components.NewANDGate3()
	c.irInstructionNOTGate = *circuit.NewNOTGate()

	// Enables
	c.registerAEnableORGate = *components.NewORGate3()
	c.registerBEnableORGate = *components.NewORGate4()
	c.registerBSetORGate = *components.NewORGate4()
	c.accEnableORGate = *components.NewORGate4()
	c.accEnableANDGate = *circuit.NewANDGate()
	c.busOneEnableORGate = *components.NewORGate4()
	c.iarEnableORGate = *components.NewORGate4()
	c.iarEnableANDGate = *circuit.NewANDGate()
	c.ramEnableORGate = *components.NewORGate5()
	c.ramEnableANDGate = *circuit.NewANDGate()

	// Sets
	c.irSetANDGate = *circuit.NewANDGate()
	c.marSetORGate = *components.NewORGate6()
	c.marSetANDGate = *circuit.NewANDGate()
	c.iarSetORGate = *components.NewORGate6()
	c.iarSetANDGate = *circuit.NewANDGate()
	c.accSetORGate = *components.NewORGate4()
	c.accSetANDGate = *circuit.NewANDGate()
	c.ramSetANDGate = *circuit.NewANDGate()
	c.tmpSetANDGate = *circuit.NewANDGate()
	c.flagsSetORGate = *circuit.NewORGate()
	c.flagsSetANDGate = *circuit.NewANDGate()

	for i := range c.step4Gates {
		c.step4Gates[i] = *circuit.NewANDGate()
	}

	for i := range c.step5Gates {
		c.step5Gates[i] = *circuit.NewANDGate()
	}

	for i := range c.step6Gates {
		c.step6Gates[i] = *components.NewANDGate3()
	}
	c.step6Gates2And = *circuit.NewANDGate()

	for i := range c.flagStateGates {
		c.flagStateGates[i] = *circuit.NewANDGate()
	}
	c.flagStateORGate = *components.NewORGate4()

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

	c.stepper = components.NewStepper()
	c.clock = &components.Clock{}
	c.alu = alu.NewALU(c.mainBus, c.busOneOutput, c.accBus, c.aluToFlagsBus)
	c.memory = memory

	return c
}

func (c *CPU) Run() {
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
	return fmt.Sprintf("stepper: %s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\nbus1: %s\n%s\n%s\n",
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
		c.flags.String(),
	)
}

func (c *CPU) step(clockState bool) {

	c.stepper.Update(clockState)
	c.runStep4Gates()
	c.runStep5Gates()
	c.runStep6Gates()

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
}

func (c *CPU) updateStates() {
	// IAR
	runUpdateOn(&c.iar)

	// MAR
	runUpdateOn(&c.memory.AddressRegister)

	// IR
	runUpdateOn(&c.ir)

	// RAM
	runUpdateOn(c.memory)

	// TMP
	runUpdateOn(&c.tmp)

	// FLAGS
	runUpdateOn(&c.flags)

	// BUS1
	runUpdateOn(&c.busOne)

	// ACC
	runUpdateOn(&c.acc)

	// R0
	runUpdateOn(&c.gpReg0)

	// R1
	runUpdateOn(&c.gpReg1)

	// R2
	runUpdateOn(&c.gpReg2)

	// R3
	runUpdateOn(&c.gpReg3)

	c.updateInstructionDecoder3x8()
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

func (c *CPU) updateInstructionDecoder3x8() {
	c.instrDecoder3x8.bit0NOTGate.Update(c.ir.Bit(0))

	c.instrDecoder3x8.decoder.Update(c.ir.Bit(1), c.ir.Bit(2), c.ir.Bit(3))

	for i := 0; i < 8; i++ {
		c.instrDecoder3x8.selectorGates[i].Update(c.instrDecoder3x8.decoder.GetOutputWire(i), c.instrDecoder3x8.bit0NOTGate.Output())
	}
}

func (c *CPU) runStep4Gates() {
	c.step4Gates[0].Update(c.stepper.GetOutputWire(3), c.ir.Bit(0))

	gate := 1
	for selector := 0; selector < 8; selector++ {
		c.step4Gates[gate].Update(c.stepper.GetOutputWire(3), c.instrDecoder3x8.selectorGates[selector].Output())
		gate++
	}
}

func (c *CPU) runStep5Gates() {
	c.step5Gates[0].Update(c.stepper.GetOutputWire(4), c.ir.Bit(0))
	c.step5Gates[1].Update(c.stepper.GetOutputWire(4), c.instrDecoder3x8.selectorGates[0].Output())
	c.step5Gates[2].Update(c.stepper.GetOutputWire(4), c.instrDecoder3x8.selectorGates[1].Output())
	c.step5Gates[3].Update(c.stepper.GetOutputWire(4), c.instrDecoder3x8.selectorGates[2].Output())

	c.step5Gates[4].Update(c.stepper.GetOutputWire(4), c.instrDecoder3x8.selectorGates[4].Output())
	c.step5Gates[5].Update(c.stepper.GetOutputWire(4), c.instrDecoder3x8.selectorGates[5].Output())

	//TODO - this has a not gate attached to it
	//	c.step5Gates[6].Update(c.stepper.GetOutputWire(4), c.instrDecoder3x8.selectorGates[5].Output())
}

func (c *CPU) runStep6Gates() {
	c.step6Gates[0].Update(c.stepper.GetOutputWire(5), c.ir.Bit(0), c.irInstructionNOTGate.Output())
	c.step6Gates2And.Update(c.stepper.GetOutputWire(5), c.instrDecoder3x8.selectorGates[2].Output())
	c.step6Gates[1].Update(c.stepper.GetOutputWire(5), c.instrDecoder3x8.selectorGates[5].Output(), c.flagStateORGate.Output())
}

func (c *CPU) runEnable(state bool) {
	c.runEnableOnIAR(state)
	c.runEnableOnBusOne(state)
	c.runEnableOnACC(state)
	c.runEnableOnRAM(state)
	c.runEnableOnRegisterB()
	c.runEnableOnRegisterA()
	c.runEnableGeneralPurposeRegisters(state)
}

func (c *CPU) runEnableOnRegisterB() {
	c.registerBEnableORGate.Update(c.step4Gates[0].Output(), c.step5Gates[2].Output(), c.step4Gates[4].Output(), c.step4Gates[8].Output())
	c.registerBEnable.Update(c.registerBEnableORGate.Output())
}

func (c *CPU) runEnableOnRegisterA() {
	c.registerAEnableORGate.Update(c.step4Gates[1].Output(), c.step4Gates[2].Output(), c.step5Gates[0].Output())
	c.registerAEnable.Update(c.registerAEnableORGate.Output())
}

func (c *CPU) runEnableOnBusOne(state bool) {
	c.busOneEnableORGate.Update(c.stepper.GetOutputWire(0), c.step4Gates[7].Output(), c.step4Gates[6].Output(), c.step4Gates[3].Output())
	updateEnableStatus(&c.busOne, c.busOneEnableORGate.Output())
}

func (c *CPU) runEnableOnACC(state bool) {
	c.accEnableORGate.Update(c.stepper.GetOutputWire(2), c.step5Gates[5].Output(), c.step6Gates2And.Output(), c.step6Gates[0].Output())
	c.accEnableANDGate.Update(state, c.accEnableORGate.Output())

	updateEnableStatus(&c.acc, c.accEnableANDGate.Output())
}

func (c *CPU) runEnableOnIAR(state bool) {
	c.iarEnableORGate.Update(c.stepper.GetOutputWire(0), c.step4Gates[3].Output(), c.step4Gates[5].Output(), c.step4Gates[6].Output())
	c.iarEnableANDGate.Update(state, c.iarEnableORGate.Output())
	updateEnableStatus(&c.iar, c.iarEnableANDGate.Output())
}

func (c *CPU) runEnableOnRAM(state bool) {
	c.ramEnableORGate.Update(
		c.stepper.GetOutputWire(1),
		c.step6Gates[1].Output(),
		c.step5Gates[4].Output(),
		c.step5Gates[3].Output(),
		c.step5Gates[1].Output(),
	)
	c.ramEnableANDGate.Update(state, c.ramEnableORGate.Output())
	updateEnableStatus(c.memory, c.ramEnableANDGate.Output())
}

func (c *CPU) runEnableGeneralPurposeRegisters(state bool) {

	c.instructionDecoderEnables2x4[0].Update(c.ir.Bit(6), c.ir.Bit(7))
	c.instructionDecoderEnables2x4[1].Update(c.ir.Bit(4), c.ir.Bit(5))

	// R0
	c.gpRegEnableANDGates[0].Update(state, c.registerBEnable.Get(), c.instructionDecoderEnables2x4[0].GetOutputWire(0))
	c.gpRegEnableANDGates[4].Update(state, c.registerAEnable.Get(), c.instructionDecoderEnables2x4[1].GetOutputWire(0))
	c.gpRegEnableORGates[0].Update(c.gpRegEnableANDGates[4].Output(), c.gpRegEnableANDGates[0].Output())
	updateEnableStatus(&c.gpReg0, c.gpRegEnableORGates[0].Output())

	// R1
	c.gpRegEnableANDGates[1].Update(state, c.registerBEnable.Get(), c.instructionDecoderEnables2x4[0].GetOutputWire(1))
	c.gpRegEnableANDGates[5].Update(state, c.registerAEnable.Get(), c.instructionDecoderEnables2x4[1].GetOutputWire(1))
	c.gpRegEnableORGates[1].Update(c.gpRegEnableANDGates[5].Output(), c.gpRegEnableANDGates[1].Output())
	updateEnableStatus(&c.gpReg1, c.gpRegEnableORGates[1].Output())

	// R2
	// this register should be enabled at some point but isn't....
	c.gpRegEnableANDGates[2].Update(state, c.registerBEnable.Get(), c.instructionDecoderEnables2x4[0].GetOutputWire(2))
	c.gpRegEnableANDGates[6].Update(state, c.registerAEnable.Get(), c.instructionDecoderEnables2x4[1].GetOutputWire(2))
	c.gpRegEnableORGates[2].Update(c.gpRegEnableANDGates[6].Output(), c.gpRegEnableANDGates[2].Output())
	updateEnableStatus(&c.gpReg2, c.gpRegEnableORGates[2].Output())

	// R3
	c.gpRegEnableANDGates[3].Update(state, c.registerBEnable.Get(), c.instructionDecoderEnables2x4[0].GetOutputWire(3))
	c.gpRegEnableANDGates[7].Update(state, c.registerAEnable.Get(), c.instructionDecoderEnables2x4[1].GetOutputWire(3))
	c.gpRegEnableORGates[3].Update(c.gpRegEnableANDGates[7].Output(), c.gpRegEnableANDGates[3].Output())
	updateEnableStatus(&c.gpReg3, c.gpRegEnableORGates[3].Output())
}

func (c *CPU) runSet(state bool) {
	c.irInstructionANDGate.Update(c.ir.Bit(3), c.ir.Bit(2), c.ir.Bit(1))
	c.irInstructionNOTGate.Update(c.irInstructionANDGate.Output())

	c.refreshFlagStateGates()

	c.runSetOnMAR(state)
	c.runSetOnIAR(state)
	c.runSetOnIR(state)
	c.runSetOnACC(state)
	c.runSetOnRAM(state)
	c.runSetOnTMP(state)
	c.runSetOnFLAGS(state)
	c.runSetOnRegisterB()
	c.runSetGeneralPurposeRegisters(state)
}

func (c *CPU) refreshFlagStateGates() {
	// C
	c.flagStateGates[0].Update(c.ir.Bit(4), c.flagsBus.GetOutputWire(0))
	// A
	c.flagStateGates[1].Update(c.ir.Bit(5), c.flagsBus.GetOutputWire(1))
	// E
	c.flagStateGates[2].Update(c.ir.Bit(6), c.flagsBus.GetOutputWire(2))
	// Z
	c.flagStateGates[3].Update(c.ir.Bit(7), c.flagsBus.GetOutputWire(3))

	c.flagStateORGate.Update(
		c.flagStateGates[0].Output(),
		c.flagStateGates[1].Output(),
		c.flagStateGates[2].Output(),
		c.flagStateGates[3].Output(),
	)
}

func (c *CPU) runSetOnMAR(state bool) {
	c.marSetORGate.Update(
		c.stepper.GetOutputWire(0),
		c.step4Gates[3].Output(),
		c.step4Gates[6].Output(),
		c.step4Gates[1].Output(),
		c.step4Gates[2].Output(),
		c.step4Gates[5].Output(),
	)
	c.marSetANDGate.Update(state, c.marSetORGate.Output())
	updateSetStatus(&c.memory.AddressRegister, c.marSetANDGate.Output())
}

func (c *CPU) runSetOnIAR(state bool) {
	c.iarSetORGate.Update(
		c.stepper.GetOutputWire(2),
		c.step4Gates[4].Output(),
		c.step5Gates[4].Output(),
		c.step5Gates[5].Output(),
		c.step6Gates2And.Output(),
		c.step6Gates[1].Output(),
	)
	c.iarSetANDGate.Update(state, c.iarSetORGate.Output())
	updateSetStatus(&c.iar, c.iarSetANDGate.Output())
}

func (c *CPU) runSetOnIR(state bool) {
	c.irSetANDGate.Update(state, c.stepper.GetOutputWire(1))
	updateSetStatus(&c.ir, c.irSetANDGate.Output())
}

func (c *CPU) runSetOnACC(state bool) {
	c.accSetORGate.Update(
		c.stepper.GetOutputWire(0),
		c.step4Gates[3].Output(),
		c.step4Gates[6].Output(),
		c.step5Gates[0].Output(),
	)
	c.accSetANDGate.Update(state, c.accSetORGate.Output())
	updateSetStatus(&c.acc, c.accSetANDGate.Output())
}

func (c *CPU) runSetOnFLAGS(state bool) {
	c.flagsSetORGate.Update(
		c.step5Gates[0].Output(),
		c.step4Gates[7].Output(),
	)
	c.flagsSetANDGate.Update(state, c.flagsSetORGate.Output())
	updateSetStatus(&c.flags, c.flagsSetANDGate.Output())
}

func (c *CPU) runSetOnRAM(state bool) {
	c.ramSetANDGate.Update(state, c.step5Gates[2].Output())
	updateSetStatus(c.memory, c.ramSetANDGate.Output())
}

func (c *CPU) runSetOnTMP(state bool) {
	c.tmpSetANDGate.Update(state, c.step4Gates[0].Output())
	updateSetStatus(&c.tmp, c.tmpSetANDGate.Output())
}

func (c *CPU) runSetOnRegisterB() {
	c.registerBSetORGate.Update(
		c.step5Gates[1].Output(),
		c.step6Gates[0].Output(),
		c.step5Gates[3].Output(),
		c.step5Gates[6].Output(),
	)

	c.registerBSet.Update(c.registerBSetORGate.Output())
}

func (c *CPU) runSetGeneralPurposeRegisters(state bool) {
	c.instructionDecoderSet2x4.Update(c.ir.Bit(6), c.ir.Bit(7))

	// R0
	c.gpRegSetANDGates[0].Update(state, c.registerBSet.Get(), c.instructionDecoderSet2x4.GetOutputWire(0))
	updateSetStatus(&c.gpReg0, c.gpRegSetANDGates[0].Output())

	// R1
	c.gpRegSetANDGates[1].Update(state, c.registerBSet.Get(), c.instructionDecoderSet2x4.GetOutputWire(1))
	updateSetStatus(&c.gpReg1, c.gpRegSetANDGates[1].Output())

	// R2
	c.gpRegSetANDGates[2].Update(state, c.registerBSet.Get(), c.instructionDecoderSet2x4.GetOutputWire(2))
	updateSetStatus(&c.gpReg2, c.gpRegSetANDGates[2].Output())

	// R3
	c.gpRegSetANDGates[3].Update(state, c.registerBSet.Get(), c.instructionDecoderSet2x4.GetOutputWire(3))
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
