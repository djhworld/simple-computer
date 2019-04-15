package components

import (
	"fmt"

	"github.com/djhworld/simple-computer/circuit"
)

type Stepper struct {
	bits           [12]Bit
	reset          circuit.Wire
	resetNotGate   circuit.NOTGate
	clockIn        circuit.Wire
	clockInNotGate circuit.NOTGate
	inputOrGates   [2]circuit.ORGate
	outputs        [7]circuit.Wire
	outputAndGates [5]circuit.ANDGate
	outputOrGate   circuit.ORGate
	outputNotGates [6]circuit.NOTGate
}

func NewStepper() *Stepper {
	s := new(Stepper)

	for i, _ := range s.bits {
		s.bits[i] = *NewBit()
	}

	s.resetNotGate = *circuit.NewNOTGate()
	s.clockInNotGate = *circuit.NewNOTGate()
	for i, _ := range s.inputOrGates {
		s.inputOrGates[i] = *circuit.NewORGate()
	}

	for i, _ := range s.outputAndGates {
		s.outputAndGates[i] = *circuit.NewANDGate()
	}
	s.outputOrGate = *circuit.NewORGate()

	for i, _ := range s.outputNotGates {
		s.outputNotGates[i] = *circuit.NewNOTGate()
	}

	return s
}

func (s *Stepper) GetOutputWire(index int) bool {
	return s.outputs[index].Get()
}

func (s *Stepper) String() string {
	result := ""
	for i := 0; i < len(s.outputs)-1; i++ {
		if s.outputs[i].Get() {
			result += fmt.Sprint("* ")
		} else {
			result += fmt.Sprint("- ")
		}
	}
	return result
}

func (s *Stepper) Update(clockIn bool) {
	s.clockIn.Update(clockIn)
	s.reset.Update(s.outputs[6].Get())

	s.step()

	// reset is instant so should do it immediately
	if s.outputs[6].Get() {
		s.reset.Update(s.outputs[6].Get())
		s.step()
	}
}

func (s *Stepper) step() {
	s.clockInNotGate.Update(s.clockIn.Get())
	s.resetNotGate.Update(s.reset.Get())

	s.inputOrGates[0].Update(s.reset.Get(), s.clockInNotGate.Output())
	s.inputOrGates[1].Update(s.reset.Get(), s.clockIn.Get())

	s.bits[0].Update(s.resetNotGate.Output(), s.inputOrGates[0].Output())
	s.bits[1].Update(s.bits[0].Get(), s.inputOrGates[1].Output())
	s.outputNotGates[0].Update(s.bits[1].Get())
	s.outputOrGate.Update(s.outputNotGates[0].Output(), s.reset.Get())

	s.bits[2].Update(s.bits[1].Get(), s.inputOrGates[0].Output())
	s.bits[3].Update(s.bits[2].Get(), s.inputOrGates[1].Output())
	s.outputNotGates[1].Update(s.bits[3].Get())
	s.outputAndGates[0].Update(s.outputNotGates[1].Output(), s.bits[1].Get())

	s.bits[4].Update(s.bits[3].Get(), s.inputOrGates[0].Output())
	s.bits[5].Update(s.bits[4].Get(), s.inputOrGates[1].Output())
	s.outputNotGates[2].Update(s.bits[5].Get())
	s.outputAndGates[1].Update(s.outputNotGates[2].Output(), s.bits[3].Get())

	s.bits[6].Update(s.bits[5].Get(), s.inputOrGates[0].Output())
	s.bits[7].Update(s.bits[6].Get(), s.inputOrGates[1].Output())
	s.outputNotGates[3].Update(s.bits[7].Get())
	s.outputAndGates[2].Update(s.outputNotGates[3].Output(), s.bits[5].Get())

	s.bits[8].Update(s.bits[7].Get(), s.inputOrGates[0].Output())
	s.bits[9].Update(s.bits[8].Get(), s.inputOrGates[1].Output())
	s.outputNotGates[4].Update(s.bits[9].Get())
	s.outputAndGates[3].Update(s.outputNotGates[4].Output(), s.bits[7].Get())

	s.bits[10].Update(s.bits[9].Get(), s.inputOrGates[0].Output())
	s.bits[11].Update(s.bits[10].Get(), s.inputOrGates[1].Output())
	s.outputNotGates[5].Update(s.bits[11].Get())
	s.outputAndGates[4].Update(s.outputNotGates[5].Output(), s.bits[9].Get())

	s.outputs[0].Update(s.outputOrGate.Output())
	s.outputs[1].Update(s.outputAndGates[0].Output())
	s.outputs[2].Update(s.outputAndGates[1].Output())
	s.outputs[3].Update(s.outputAndGates[2].Output())
	s.outputs[4].Update(s.outputAndGates[3].Output())
	s.outputs[5].Update(s.outputAndGates[4].Output())
	s.outputs[6].Update(s.bits[11].Get())
}
