package alu

import (
	"fmt"

	"github.com/djhworld/simple-computer/circuit"
	"github.com/djhworld/simple-computer/components"
	"github.com/djhworld/simple-computer/utils"
)

const (
	ADD = iota
	SHR
	SHL
	NOT
	AND
	OR
	XOR
	CMP
)

const BUS_WIDTH = 16

type ALU struct {
	inputABus      *components.Bus
	inputBBus      *components.Bus
	outputBus      *components.Bus
	flagsOutputBus *components.Bus

	Op      [3]circuit.Wire
	CarryIn circuit.Wire

	carryOut  circuit.Wire
	aIsLarger circuit.Wire
	isEqual   circuit.Wire

	opDecoder components.Decoder3x8

	comparator  components.Comparator
	xorer       components.XORer
	orer        components.ORer
	ander       components.ANDer
	notter      components.NOTer
	leftShifer  components.LeftShifter
	rightShifer components.RightShifter
	adder       components.Adder
	isZero      components.IsZero
	enablers    [7]components.Enabler
	andGates    [3]circuit.ANDGate
}

func NewALU(inputABus, inputBBus, outputBus, flagsOutputBus *components.Bus) *ALU {
	a := new(ALU)
	a.inputABus = inputABus
	a.inputBBus = inputBBus
	a.outputBus = outputBus
	a.flagsOutputBus = flagsOutputBus

	a.opDecoder = *components.NewDecoder3x8()

	a.comparator = *components.NewComparator()
	a.xorer = *components.NewXORer()
	a.orer = *components.NewORer()
	a.ander = *components.NewANDer()
	a.notter = *components.NewNOTer()
	a.leftShifer = *components.NewLeftShifter()
	a.rightShifer = *components.NewRightShifter()
	a.adder = *components.NewAdder()
	a.isZero = *components.NewIsZero()
	a.andGates[0] = *circuit.NewANDGate()
	a.andGates[1] = *circuit.NewANDGate()
	a.andGates[2] = *circuit.NewANDGate()

	for i := range a.enablers {
		a.enablers[i] = *components.NewEnabler()
	}

	return a
}

func (a *ALU) updateOpDecoder() {
	a.opDecoder.Update(a.Op[2].Get(), a.Op[1].Get(), a.Op[0].Get())
}

func (a *ALU) updateComparator() {
	// comparator is not wired to an enabler and runs all the time
	a.setWireOnComponent(&a.comparator)
	a.comparator.Update()
	a.aIsLarger.Update(a.comparator.Larger())
	a.isEqual.Update(a.comparator.Equal())
}

func (a *ALU) updateXorer() {
	a.setWireOnComponent(&a.xorer)
	a.xorer.Update()
	a.wireToEnabler(&a.xorer, 6)
}

func (a *ALU) updateOrer() {
	a.setWireOnComponent(&a.orer)
	a.orer.Update()
	a.wireToEnabler(&a.orer, 5)
}

func (a *ALU) updateAnder() {
	a.setWireOnComponent(&a.ander)
	a.ander.Update()
	a.wireToEnabler(&a.ander, 4)
}

func (a *ALU) updateNotter() {
	for i := (BUS_WIDTH - 1); i >= 0; i-- {
		a.notter.SetInputWire(i, a.inputABus.GetOutputWire(i))
	}
	a.notter.Update()
	a.wireToEnabler(&a.notter, 3)
}

func (a *ALU) updateLeftShifter() {
	for i := (BUS_WIDTH - 1); i >= 0; i-- {
		a.leftShifer.SetInputWire(i, a.inputABus.GetOutputWire(i))
	}
	a.leftShifer.Update(a.CarryIn.Get())
	a.wireToEnabler(&a.leftShifer, 2)
}

func (a *ALU) updateRightShifter() {
	for i := (BUS_WIDTH - 1); i >= 0; i-- {
		a.rightShifer.SetInputWire(i, a.inputABus.GetOutputWire(i))
	}
	a.rightShifer.Update(a.CarryIn.Get())
	a.wireToEnabler(&a.rightShifer, 1)
}

func (a *ALU) updateAdder() {
	a.setWireOnComponent(&a.adder)
	a.adder.Update(a.CarryIn.Get())
	a.wireToEnabler(&a.adder, 0)
}

func (a *ALU) wireToEnabler(b components.Component, enablerIndex int) {
	for i := 0; i < BUS_WIDTH; i++ {
		a.enablers[enablerIndex].SetInputWire(i, b.GetOutputWire(i))
	}
}

func (a *ALU) setWireOnComponent(b components.Component) {
	for i := BUS_WIDTH - 1; i >= 0; i-- {
		b.SetInputWire(i, a.inputABus.GetOutputWire(i))
	}

	for i := (BUS_WIDTH * 2) - 1; i >= BUS_WIDTH; i-- {
		b.SetInputWire(i, a.inputBBus.GetOutputWire(i-16))
	}
}

func (a *ALU) String() string {
	s := ""
	for i := 2; i >= 0; i-- {
		if a.Op[i].Get() {
			s += fmt.Sprint("1")
		} else {
			s += fmt.Sprint("0")
		}
	}

	var inputA uint16
	var inputB uint16
	var output uint16
	var x uint16 = 0
	for i := BUS_WIDTH - 1; i >= 0; i-- {
		if a.inputABus.GetOutputWire(i) {
			inputA = inputA | (1 << x)
		} else {
			inputA = inputA & ^(1 << x)
		}

		if a.inputBBus.GetOutputWire(i) {
			inputB = inputB | (1 << x)
		} else {
			inputB = inputB & ^(1 << x)
		}

		if a.outputBus.GetOutputWire(i) {
			output = output | (1 << x)
		} else {
			output = output & ^(1 << x)
		}
		x++
	}
	return fmt.Sprintf(
		"ALU OP: %s, A: %s, B: %s, OUT: %s, carryin: %v, carryout: %v, larger: %v, eq: %v, zero: %v",
		s,
		utils.ValueToString(inputA),
		utils.ValueToString(inputB),
		utils.ValueToString(output),
		a.CarryIn.Get(),
		a.flagsOutputBus.GetOutputWire(0),
		a.flagsOutputBus.GetOutputWire(1),
		a.flagsOutputBus.GetOutputWire(2),
		a.flagsOutputBus.GetOutputWire(3),
	)
}

func (a *ALU) Update() {
	a.updateOpDecoder()
	enabler := a.opDecoder.Index()

	a.updateComparator()

	switch enabler {
	case ADD:
		a.updateAdder()
	case XOR:
		a.updateXorer()
	case OR:
		a.updateOrer()
	case AND:
		a.updateAnder()
	case NOT:
		a.updateNotter()
	case SHL:
		a.updateLeftShifter()
	case SHR:
		a.updateRightShifter()
	}

	if enabler != CMP {
		a.enablers[enabler].Update(true)

		switch enabler {
		case ADD:
			a.andGates[0].Update(a.adder.Carry(), a.opDecoder.GetOutputWire(ADD))
			a.carryOut.Update(a.andGates[0].Output())
		case SHR:
			a.andGates[1].Update(a.rightShifer.ShiftOut(), a.opDecoder.GetOutputWire(SHR))
			a.carryOut.Update(a.andGates[1].Output())
		case SHL:
			a.andGates[2].Update(a.leftShifer.ShiftOut(), a.opDecoder.GetOutputWire(SHL))
			a.carryOut.Update(a.andGates[2].Output())
		}

		for i := 0; i < BUS_WIDTH; i++ {
			a.isZero.SetInputWire(i, a.enablers[enabler].GetOutputWire(i))
			a.outputBus.SetInputWire(i, a.enablers[enabler].GetOutputWire(i))
		}
	} else {
		for i := 0; i < BUS_WIDTH; i++ {
			a.isZero.SetInputWire(i, true)
			a.outputBus.SetInputWire(i, false)
		}
	}
	a.isZero.Update()

	a.flagsOutputBus.SetInputWire(0, a.carryOut.Get())
	a.flagsOutputBus.SetInputWire(1, a.aIsLarger.Get())
	a.flagsOutputBus.SetInputWire(2, a.isEqual.Get())
	a.flagsOutputBus.SetInputWire(3, a.isZero.GetOutputWire(0))
}
