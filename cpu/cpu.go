package cpu

import (
	"fmt"

	"github.com/djhworld/simple-computer/alu"
	"github.com/djhworld/simple-computer/circuit"
	"github.com/djhworld/simple-computer/components"
	"github.com/djhworld/simple-computer/io"
	"github.com/djhworld/simple-computer/memory"
)

// LOADS
// ----------------------
// arg A = memory address to load from
// arg B = register to store value in
// 0x0000 = LD R0, R0
// 0x0001 = LD R0, R1
// 0x0002 = LD R0, R2
// 0x0003 = LD R0, R3

// 0x0004 = LD R1, R0
// 0x0005 = LD R1, R1
// 0x0006 = LD R1, R2
// 0x0007 = LD R1, R3

// 0x0008 = LD R2, R0
// 0x0009 = LD R2, R1
// 0x000A = LD R2, R2
// 0x000B = LD R2, R3

// 0x000C = LD R3, R0
// 0x000D = LD R3, R1
// 0x000E = LD R3, R2
// 0x000F = LD R3, R3

// STORES
// ----------------------
// arg A = memory address for value
// arg B = value to store in memory
// 0x0010 = ST R0, R0
// 0x0011 = ST R0, R1
// 0x0012 = ST R0, R2
// 0x0013 = ST R0, R3

// 0x0014 = ST R1, R0
// 0x0015 = ST R1, R1
// 0x0016 = ST R1, R2
// 0x0017 = ST R1, R3

// 0x0018 = ST R2, R0
// 0x0019 = ST R2, R1
// 0x001A = ST R2, R2
// 0x001B = ST R2, R3

// 0x001C = ST R3, R0
// 0x001D = ST R3, R1
// 0x001E = ST R3, R2
// 0x001F = ST R3, R3

// DATA
// put value in memory into register (2 byte instruction)
// ----------------------
// 0x0020 = DATA R0
// 0x0021 = DATA R1
// 0x0022 = DATA R2
// 0x0023 = DATA R3

// JR
// set instruction address register to value in register
// ----------------------
// 0x0030 = JR R0
// 0x0031 = JR R1
// 0x0032 = JR R2
// 0x0033 = JR R3

// JMP
// set instruction address register to next byte (2 byte instruction)
// ----------------------
// 0x0040 = JMP <value>

// JMP(CAEZ)
// set instruction address register to next byte (2 byte instruction)
// jump if <flag(s)> are true
// ----------------------
// 0x0051 = JMPZ <value>
// 0x0052 = JMPE <value>
// 0x0053 = JMPEZ <value>
// 0x0054 = JMPA <value>
// 0x0055 = JMPAZ <value>
// 0x0056 = JMPAE <value>
// 0x0057 = JMPAEZ <value>
// 0x0058 = JMPC <value>
// 0x0059 = JMPCZ <value>
// 0x005A = JMPCE <value>
// 0x005B = JMPCEZ <value>
// 0x005C = JMPCA <value>
// 0x005D = JMPCAZ <value>
// 0x005E = JMPCAE <value>
// 0x005F = JMPCAEZ <value>

// CLF (CLEAR FLAGS)
// ----------------------
// 0x0060 CLF

// IN
// ----------------------
// 0x0070 = IN Data, R0
// 0x0071 = IN Data, R1
// 0x0072 = IN Data, R2
// 0x0073 = IN Data, R3
// 0x0074 = IN Addr, R0
// 0x0075 = IN Addr, R1
// 0x0076 = IN Addr, R2
// 0x0077 = IN Addr, R3

// OUT
// ----------------------
// 0x0078 = OUT Data, R0
// 0x0079 = OUT Data, R1
// 0x007A = OUT Data, R2
// 0x007B = OUT Data, R3
// 0x007C = OUT Addr, R0
// 0x007D = OUT Addr, R1
// 0x007E = OUT Addr, R2
// 0x007F = OUT Addr, R3

// ADDS
// ----------------------
// 0x0080 = ADD R0, R0
// 0x0081 = ADD R0, R1
// 0x0082 = ADD R0, R2
// 0x0083 = ADD R0, R3

// 0x0084 = ADD R1, R0
// 0x0085 = ADD R1, R1
// 0x0086 = ADD R1, R2
// 0x0087 = ADD R1, R3

// 0x0088 = ADD R2, R0
// 0x0089 = ADD R2, R1
// 0x008A = ADD R2, R2
// 0x008B = ADD R2, R3

// 0x008C = ADD R3, R0
// 0x008D = ADD R3, R1
// 0x008E = ADD R3, R2
// 0x008F = ADD R3, R3

// SHR
// ----------------------
// 0x0090 = SHR R0
// 0x0095 = SHR R1
// 0x009A = SHR R2
// 0x009F = SHR R3

// SHL
// ----------------------
// 0x00A0 = SHL R0
// 0x00A5 = SHL R1
// 0x00AA = SHL R2
// 0x00AF = SHL R3

// NOT
// ----------------------
// 0x00B0 = NOT R0
// 0x00B5 = NOT R1
// 0x00BA = NOT R2
// 0x00BF = NOT R3

// ANDS
// ----------------------
// 0x00C0 = AND R0, R0
// 0x00C1 = AND R0, R1
// 0x00C2 = AND R0, R2
// 0x00C3 = AND R0, R3

// 0x00C4 = AND R1, R0
// 0x00C5 = AND R1, R1
// 0x00C6 = AND R1, R2
// 0x00C7 = AND R1, R3

// 0x00C8 = AND R2, R0
// 0x00C9 = AND R2, R1
// 0x00CA = AND R2, R2
// 0x00CB = AND R2, R3

// 0x00CC = AND R3, R0
// 0x00CD = AND R3, R1
// 0x00CE = AND R3, R2
// 0x00CF = AND R3, R3

// ORS
// ----------------------
// 0x00D0 = OR R0, R0
// 0x00D1 = OR R0, R1
// 0x00D2 = OR R0, R2
// 0x00D3 = OR R0, R3

// 0x00D4 = OR R1, R0
// 0x00D5 = OR R1, R1
// 0x00D6 = OR R1, R2
// 0x00D7 = OR R1, R3

// 0x00D8 = OR R2, R0
// 0x00D9 = OR R2, R1
// 0x00DA = OR R2, R2
// 0x00DB = OR R2, R3

// 0x00DC = OR R3, R0
// 0x00DD = OR R3, R1
// 0x00DE = OR R3, R2
// 0x00DF = OR R3, R3

// XORS
// ----------------------
// 0x00E0 = XOR R0, R0
// 0x00E1 = XOR R0, R1
// 0x00E2 = XOR R0, R2
// 0x00E3 = XOR R0, R3

// 0x00E4 = XOR R1, R0
// 0x00E5 = XOR R1, R1
// 0x00E6 = XOR R1, R2
// 0x00E7 = XOR R1, R3

// 0x00E8 = XOR R2, R0
// 0x00E9 = XOR R2, R1
// 0x00EA = XOR R2, R2
// 0x00EB = XOR R2, R3

// 0x00EC = XOR R3, R0
// 0x00ED = XOR R3, R1
// 0x00EE = XOR R3, R2
// 0x00EF = XOR R3, R3

// CMP
// ----------------------
// 0x00F0 = CMP R0, R0
// 0x00F1 = CMP R0, R1
// 0x00F2 = CMP R0, R2
// 0x00F3 = CMP R0, R3

// 0x00F4 = CMP R1, R0
// 0x00F5 = CMP R1, R1
// 0x00F6 = CMP R1, R2
// 0x00F7 = CMP R1, R3

// 0x00F8 = CMP R2, R0
// 0x00F9 = CMP R2, R1
// 0x00FA = CMP R2, R2
// 0x00FB = CMP R2, R3

// 0x00FC = CMP R3, R0
// 0x00FD = CMP R3, R1
// 0x00FE = CMP R3, R2
// 0x00FF = CMP R3, R3

const BUS_WIDTH = 16

const (
	FLAGS_BUS_CARRY    = 0
	FLAGS_BUS_A_LARGER = 1
	FLAGS_BUS_EQUAL    = 2
	FLAGS_BUS_ZERO     = 3
)

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

	clockState bool
	memory     *memory.Memory64K
	alu        *alu.ALU
	stepper    *components.Stepper
	busOne     components.BusOne

	mainBus       *components.Bus
	tmpBus        *components.Bus
	busOneOutput  *components.Bus
	controlBus    *components.Bus
	accBus        *components.Bus
	aluToFlagsBus *components.Bus
	flagsBus      *components.Bus
	ioBus         *components.IOBus

	// CONTROL UNIT
	// inc. gates, wiring, instruction decoding etc
	step4Gates     [8]circuit.ANDGate
	step4Gate3And  components.ANDGate3
	step5Gates     [6]circuit.ANDGate
	step5Gate3And  components.ANDGate3
	step6Gates     [2]components.ANDGate3
	step6Gates2And circuit.ANDGate

	instrDecoder3x8              InstructionDecoder3x8
	instructionDecoderEnables2x4 [2]components.Decoder2x4
	instructionDecoderSet2x4     components.Decoder2x4

	irInstructionANDGate components.ANDGate3
	irInstructionNOTGate circuit.NOTGate

	ioBusEnableGate       circuit.ANDGate
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

	ioBusSetGate    circuit.ANDGate
	irBit4NOTGate   circuit.NOTGate
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

	aluOpAndGates [3]components.ANDGate3

	carryTemp    components.Bit
	carryANDGate circuit.ANDGate

	peripherals []io.Peripheral
}

func NewCPU(mainBus *components.Bus, memory *memory.Memory64K) *CPU {
	c := new(CPU)

	c.clockState = false
	c.stepper = components.NewStepper()
	c.memory = memory

	// REGISTERS
	c.controlBus = components.NewBus(BUS_WIDTH)
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
	c.aluToFlagsBus = components.NewBus(BUS_WIDTH)
	c.flagsBus = components.NewBus(BUS_WIDTH)
	c.flags = *components.NewRegister("FLAGS", c.aluToFlagsBus, c.flagsBus)
	// flags register is always enabled, and we initialise it with value 0
	updateEnableStatus(&c.flags, true)
	updateSetStatus(&c.flags, true)
	runUpdateOn(&c.flags)
	updateSetStatus(&c.flags, false)

	// TMP
	c.tmpBus = components.NewBus(BUS_WIDTH)
	c.tmp = *components.NewRegister("TMP", c.mainBus, c.tmpBus)
	// tmp register is always enabled, and we initialise it with value 0
	updateEnableStatus(&c.tmp, true)
	updateSetStatus(&c.tmp, true)
	runUpdateOn(&c.tmp)
	updateSetStatus(&c.tmp, false)

	// BUS 1
	c.busOneOutput = components.NewBus(BUS_WIDTH)
	c.busOne = *components.NewBusOne(c.tmpBus, c.busOneOutput)

	// ACC
	c.accBus = components.NewBus(BUS_WIDTH)
	c.acc = *components.NewRegister("ACC", c.accBus, c.mainBus)

	// ALU
	c.alu = alu.NewALU(c.mainBus, c.busOneOutput, c.accBus, c.aluToFlagsBus)

	// IR Register Output
	c.irInstructionANDGate = *components.NewANDGate3()
	c.irInstructionNOTGate = *circuit.NewNOTGate()
	c.irBit4NOTGate = *circuit.NewNOTGate()

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

	c.carryTemp = *components.NewBit()
	c.carryANDGate = *circuit.NewANDGate()

	for i := range c.step4Gates {
		c.step4Gates[i] = *circuit.NewANDGate()
	}
	c.step4Gate3And = *components.NewANDGate3()

	for i := range c.step5Gates {
		c.step5Gates[i] = *circuit.NewANDGate()
	}
	c.step5Gate3And = *components.NewANDGate3()

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

	c.ioBus = components.NewIOBus()
	c.ioBusEnableGate = *circuit.NewANDGate()
	c.ioBusSetGate = *circuit.NewANDGate()

	c.peripherals = make([]io.Peripheral, 0)

	return c
}

func (c *CPU) ConnectPeripheral(p io.Peripheral) {
	p.Connect(c.ioBus, c.mainBus)
	c.peripherals = append(c.peripherals, p)
}

// Jump IAR
func (c *CPU) SetIAR(address uint16) {
	c.mainBus.SetValue(address)

	updateSetStatus(&c.iar, true)
	runUpdateOn(&c.iar)
	updateSetStatus(&c.iar, false)
	runUpdateOn(&c.iar)

	c.clearMainBus()
}

func (c *CPU) Step() {
	for i := 0; i < 2; i++ {
		if c.clockState {
			c.clockState = false
		} else {
			c.clockState = true
		}

		c.step(c.clockState)
	}
}

func (c *CPU) String() string {
	return fmt.Sprintf("STEPPER: %s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\nBUS1: %s\n%s\n%s\n",
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
		c.flags.String(),
		c.alu.String(),
	)
}

func (c *CPU) step(clockState bool) {
	c.stepper.Update(clockState)
	c.runStep4Gates()
	c.runStep5Gates()
	c.runStep6Gates()

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

	// main bus should not have residual data!
	c.clearMainBus()
}

func (c *CPU) clearMainBus() {
	for i := 0; i < BUS_WIDTH; i++ {
		c.mainBus.SetInputWire(i, false)
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

	// ALU
	c.updateALU()

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
	c.updateIOBus()
	c.updatePeripherals()
}

func (c *CPU) updatePeripherals() {
	for _, p := range c.peripherals {
		p.Update()
	}
}

func (c *CPU) updateIOBus() {
	c.ioBus.Update(c.ir.Bit(12), c.ir.Bit(13))
}

func (c *CPU) updateALU() {
	//update ALU operation based on instruction register
	c.aluOpAndGates[2].Update(c.ir.Bit(9), c.ir.Bit(8), c.stepper.GetOutputWire(4))
	c.aluOpAndGates[1].Update(c.ir.Bit(10), c.ir.Bit(8), c.stepper.GetOutputWire(4))
	c.aluOpAndGates[0].Update(c.ir.Bit(11), c.ir.Bit(8), c.stepper.GetOutputWire(4))

	c.alu.Op[2].Update(c.aluOpAndGates[2].Output())
	c.alu.Op[1].Update(c.aluOpAndGates[1].Output())
	c.alu.Op[0].Update(c.aluOpAndGates[0].Output())

	c.alu.CarryIn.Update(c.carryANDGate.Output())
	c.alu.Update()

}

func (c *CPU) updateInstructionDecoder3x8() {
	c.instrDecoder3x8.bit0NOTGate.Update(c.ir.Bit(8))

	c.instrDecoder3x8.decoder.Update(c.ir.Bit(9), c.ir.Bit(10), c.ir.Bit(11))

	for i := 0; i < 8; i++ {
		c.instrDecoder3x8.selectorGates[i].Update(c.instrDecoder3x8.decoder.GetOutputWire(i), c.instrDecoder3x8.bit0NOTGate.Output())
	}
}

func (c *CPU) runStep4Gates() {
	c.step4Gates[0].Update(c.stepper.GetOutputWire(3), c.ir.Bit(8))

	gate := 1
	for selector := 0; selector < 7; selector++ {
		c.step4Gates[gate].Update(c.stepper.GetOutputWire(3), c.instrDecoder3x8.selectorGates[selector].Output())
		gate++
	}

	c.step4Gate3And.Update(c.stepper.GetOutputWire(3), c.instrDecoder3x8.selectorGates[7].Output(), c.ir.Bit(12))
	c.irBit4NOTGate.Update(c.ir.Bit(12))
}

func (c *CPU) runStep5Gates() {
	c.step5Gates[0].Update(c.stepper.GetOutputWire(4), c.ir.Bit(8))
	c.step5Gates[1].Update(c.stepper.GetOutputWire(4), c.instrDecoder3x8.selectorGates[0].Output())
	c.step5Gates[2].Update(c.stepper.GetOutputWire(4), c.instrDecoder3x8.selectorGates[1].Output())
	c.step5Gates[3].Update(c.stepper.GetOutputWire(4), c.instrDecoder3x8.selectorGates[2].Output())

	c.step5Gates[4].Update(c.stepper.GetOutputWire(4), c.instrDecoder3x8.selectorGates[4].Output())
	c.step5Gates[5].Update(c.stepper.GetOutputWire(4), c.instrDecoder3x8.selectorGates[5].Output())

	c.step5Gate3And.Update(c.stepper.GetOutputWire(4), c.instrDecoder3x8.selectorGates[7].Output(), c.irBit4NOTGate.Output())
}

func (c *CPU) runStep6Gates() {
	c.step6Gates[0].Update(c.stepper.GetOutputWire(5), c.ir.Bit(8), c.irInstructionNOTGate.Output())
	c.step6Gates2And.Update(c.stepper.GetOutputWire(5), c.instrDecoder3x8.selectorGates[2].Output())
	c.step6Gates[1].Update(c.stepper.GetOutputWire(5), c.instrDecoder3x8.selectorGates[5].Output(), c.flagStateORGate.Output())
}

func (c *CPU) runEnable(state bool) {
	c.runEnableOnIO(state)
	c.runEnableOnIAR(state)
	c.runEnableOnBusOne(state)
	c.runEnableOnACC(state)
	c.runEnableOnRAM(state)
	c.runEnableOnRegisterB()
	c.runEnableOnRegisterA()
	c.runEnableGeneralPurposeRegisters(state)
}

func (c *CPU) runEnableOnIO(state bool) {
	c.ioBusEnableGate.Update(state, c.step5Gate3And.Output())
	updateEnableStatus(c.ioBus, c.ioBusEnableGate.Output())
}

func (c *CPU) runEnableOnRegisterB() {
	c.registerBEnableORGate.Update(c.step4Gates[0].Output(), c.step5Gates[2].Output(), c.step4Gates[4].Output(), c.step4Gate3And.Output())
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

	c.instructionDecoderEnables2x4[0].Update(c.ir.Bit(14), c.ir.Bit(15))
	c.instructionDecoderEnables2x4[1].Update(c.ir.Bit(12), c.ir.Bit(13))

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
	c.irInstructionANDGate.Update(c.ir.Bit(11), c.ir.Bit(10), c.ir.Bit(9))
	c.irInstructionNOTGate.Update(c.irInstructionANDGate.Output())

	c.refreshFlagStateGates()

	c.runSetOnIO(state)
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
	/*
		// C
		c.flagStateGates[0].Update(c.ir.Bit(12), c.flagsBus.GetOutputWire(FLAGS_BUS_CARRY))
		// A
		c.flagStateGates[1].Update(c.ir.Bit(13), c.flagsBus.GetOutputWire(FLAGS_BUS_A_LARGER))
		// E
		c.flagStateGates[2].Update(c.ir.Bit(14), c.flagsBus.GetOutputWire(FLAGS_BUS_EQUAL))
		// Z
		c.flagStateGates[3].Update(c.ir.Bit(15), c.flagsBus.GetOutputWire(FLAGS_BUS_ZERO))
	*/
	// C
	c.flagStateGates[0].Update(c.ir.Bit(12), c.flagsBus.GetOutputWire(FLAGS_BUS_CARRY))
	// A
	c.flagStateGates[1].Update(c.ir.Bit(13), c.flagsBus.GetOutputWire(FLAGS_BUS_A_LARGER))
	// E
	c.flagStateGates[2].Update(c.ir.Bit(14), c.flagsBus.GetOutputWire(FLAGS_BUS_EQUAL))
	// Z
	c.flagStateGates[3].Update(c.ir.Bit(15), c.flagsBus.GetOutputWire(FLAGS_BUS_ZERO))

	c.flagStateORGate.Update(
		c.flagStateGates[0].Output(),
		c.flagStateGates[1].Output(),
		c.flagStateGates[2].Output(),
		c.flagStateGates[3].Output(),
	)
}

func (c *CPU) runSetOnIO(state bool) {
	c.ioBusSetGate.Update(state, c.step4Gate3And.Output())
	updateSetStatus(c.ioBus, c.ioBusSetGate.Output())
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

	// We will add a new memory bit called
	// "Carry Temp" that goes between the Carry Flag and the enabler
	// we just added above. It will be set in step 4, the same time that the TMP register gets
	// set. Thus, the ALU instruction will have a carry input that cannot change during step 5.
	c.carryTemp.Update(c.flagsBus.GetOutputWire(FLAGS_BUS_CARRY), c.tmpSetANDGate.Output())
	c.carryANDGate.Update(c.carryTemp.Get(), c.step5Gates[0].Output())
}

func (c *CPU) runSetOnRegisterB() {
	c.registerBSetORGate.Update(
		c.step5Gates[1].Output(),
		c.step6Gates[0].Output(),
		c.step5Gates[3].Output(),
		c.step5Gate3And.Output(),
	)

	c.registerBSet.Update(c.registerBSetORGate.Output())
}

func (c *CPU) runSetGeneralPurposeRegisters(state bool) {
	c.instructionDecoderSet2x4.Update(c.ir.Bit(14), c.ir.Bit(15))

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
