package alu

import (
	"fmt"

	"github.com/djhworld/simple-computer/circuit"
	"github.com/djhworld/simple-computer/components"
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

type ALU struct {
	inputABus      *components.Bus
	inputBBus      *components.Bus
	outputBus      *components.Bus
	flagsOutputBus *components.Bus

	inputA [8]circuit.Wire
	inputB [8]circuit.Wire

	Op        [3]circuit.Wire
	CarryIn   circuit.Wire
	CarryOut  circuit.Wire
	AisLarger circuit.Wire
	IsEqual   circuit.Wire

	opDecoder components.Decoder3x8

	output [8]circuit.Wire

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
	a.AisLarger.Update(a.comparator.Larger())
	a.IsEqual.Update(a.comparator.Equal())
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
	for i := (8 - 1); i >= 0; i-- {
		a.notter.SetInputWire(i, a.inputA[i].Get())
	}
	a.notter.Update()
	a.wireToEnabler(&a.notter, 3)
}

func (a *ALU) updateLeftShifter() {
	for i := (8 - 1); i >= 0; i-- {
		a.leftShifer.SetInputWire(i, a.inputA[i].Get())
	}
	a.leftShifer.Update(a.CarryIn.Get())
	a.wireToEnabler(&a.leftShifer, 2)
}

func (a *ALU) updateRightShifter() {
	for i := (8 - 1); i >= 0; i-- {
		a.rightShifer.SetInputWire(i, a.inputA[i].Get())
	}
	a.rightShifer.Update(a.CarryIn.Get())
	a.wireToEnabler(&a.rightShifer, 1)
}

func (a *ALU) updateAdder() {
	a.setWireOnComponent(&a.adder)
	a.adder.Update(a.CarryIn.Get())
	a.wireToEnabler(&a.adder, 0)
}

func (a *ALU) wireToEnabler(b components.ByteComponent, enablerIndex int) {
	for i := 0; i < 8; i++ {
		a.enablers[enablerIndex].SetInputWire(i, b.GetOutputWire(i))
	}
}

func (a *ALU) setWireOnComponent(b components.ByteComponent) {
	for i := 8 - 1; i >= 0; i-- {
		b.SetInputWire(i, a.inputA[i].Get())
	}

	for i := 16 - 1; i >= 8; i-- {
		b.SetInputWire(i, a.inputB[i-8].Get())
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

	var inputA byte
	var inputB byte
	var x int = 0
	for i := 7; i >= 0; i-- {
		if a.inputA[i].Get() {
			inputA = inputA | (1 << byte(x))
		} else {
			inputA = inputA & ^(1 << byte(x))
		}

		if a.inputB[i].Get() {
			inputB = inputB | (1 << byte(x))
		} else {
			inputB = inputB & ^(1 << byte(x))
		}
		x++
	}
	return fmt.Sprintf("ALU OP: %s, A: 0x%X, B: 0x%X, eq: %v, larger: %v, zero: %v", s, inputA, inputB, a.IsEqual.Get(), a.AisLarger.Get(), a.isZero.GetOutputWire(0))
}

func (a *ALU) resetOutputs() {
	a.CarryOut.Update(false)
	a.isZero.Reset()
	a.AisLarger.Update(false)
	a.IsEqual.Update(false)
	for i := 0; i < 8; i++ {
		a.output[i].Update(false)
		if i < 7 {
			a.enablers[i].Update(false)
		}
	}

	for i := range a.andGates {
		a.andGates[i].Update(false, false)
	}
}

func (a *ALU) Update() {
	for i := 8 - 1; i >= 0; i-- {
		a.inputA[i].Update(a.inputABus.GetOutputWire(i))
		a.inputB[i].Update(a.inputBBus.GetOutputWire(i))
	}

	a.resetOutputs()
	a.updateOpDecoder()

	a.updateComparator()
	a.updateXorer()
	a.updateOrer()
	a.updateAnder()
	a.updateNotter()
	a.updateLeftShifter()
	a.updateRightShifter()
	a.updateAdder()

	enabler := a.opDecoder.Index()

	if enabler != CMP {
		a.enablers[enabler].Update(true)

		switch enabler {
		case ADD:
			a.andGates[0].Update(a.adder.Carry(), true)
			a.CarryOut.Update(a.andGates[0].Output())
		case SHR:
			a.andGates[1].Update(a.rightShifer.ShiftOut(), true)
			a.CarryOut.Update(a.andGates[1].Output())
		case SHL:
			a.andGates[2].Update(a.leftShifer.ShiftOut(), true)
			a.CarryOut.Update(a.andGates[2].Output())
		}

		for i := 0; i < 8; i++ {
			a.isZero.SetInputWire(i, a.enablers[enabler].GetOutputWire(i))
			a.output[i].Update(a.enablers[enabler].GetOutputWire(i))
		}

		a.isZero.Update()
	}

	for i := 8 - 1; i >= 0; i-- {
		a.outputBus.SetInputWire(i, a.output[i].Get())
	}

	a.flagsOutputBus.SetInputWire(0, a.CarryOut.Get())
	a.flagsOutputBus.SetInputWire(1, a.AisLarger.Get())
	a.flagsOutputBus.SetInputWire(2, a.IsEqual.Get())
	a.flagsOutputBus.SetInputWire(3, a.isZero.GetOutputWire(0))
	a.flagsOutputBus.SetInputWire(4, false)
	a.flagsOutputBus.SetInputWire(5, false)
	a.flagsOutputBus.SetInputWire(6, false)
	a.flagsOutputBus.SetInputWire(7, false)
}
