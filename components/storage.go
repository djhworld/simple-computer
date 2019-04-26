package components

import (
	"github.com/djhworld/simple-computer/circuit"
)

type Bit struct {
	gates [4]circuit.NANDGate
	wireO circuit.Wire
}

func NewBit() *Bit {
	wireO := *circuit.NewWire("O", false)

	gates := [4]circuit.NANDGate{
		*circuit.NewNANDGate(),
		*circuit.NewNANDGate(),
		*circuit.NewNANDGate(),
		*circuit.NewNANDGate(),
	}

	return &Bit{
		wireO: wireO,
		gates: gates,
	}
}

func (m *Bit) Get() bool {
	return m.wireO.Get()
}

func (m *Bit) Update(wireI bool, wireS bool) {
	for i := 0; i < 2; i++ {
		m.gates[0].Update(wireI, wireS)
		m.gates[1].Update(m.gates[0].Output(), wireS)
		m.gates[2].Update(m.gates[0].Output(), m.gates[3].Output())
		m.gates[3].Update(m.gates[2].Output(), m.gates[1].Output())
		m.wireO.Update(m.gates[2].Output())
	}
}

type Word struct {
	inputs  [16]circuit.Wire
	bits    [16]Bit
	outputs [16]circuit.Wire
	next    Component
}

func NewWord() *Word {
	w := new(Word)
	for i, _ := range w.bits {
		w.bits[i] = *NewBit()
	}
	return w
}

func (e *Word) ConnectOutput(b Component) {
	e.next = b
}

func (e *Word) GetOutputWire(index int) bool {
	return e.outputs[index].Get()
}

func (e *Word) SetInputWire(index int, value bool) {
	e.inputs[index].Update(value)
}

func (e *Word) Update(set bool) {
	for i := 0; i < len(e.inputs); i++ {
		e.bits[i].Update(e.inputs[i].Get(), set)
		e.outputs[i].Update(e.bits[i].Get())
	}

	if e.next != nil {
		for i := 0; i < len(e.outputs); i++ {
			e.next.SetInputWire(i, e.outputs[i].Get())
		}
	}
}
