package components

import (
	"github.com/djhworld/simple-computer/circuit"
)

type Bit struct {
	gates []circuit.NANDGate
	wireI circuit.Wire
	wireS circuit.Wire
	wireO circuit.Wire
}

func NewBit() *Bit {
	wireI := *circuit.NewWire("I", false)
	wireS := *circuit.NewWire("S", false)
	wireO := *circuit.NewWire("O", false)

	gates := []circuit.NANDGate{
		*circuit.NewNANDGate(),
		*circuit.NewNANDGate(),
		*circuit.NewNANDGate(),
		*circuit.NewNANDGate(),
	}

	return &Bit{
		wireI: wireI,
		wireS: wireS,
		wireO: wireO,
		gates: gates,
	}
}

func (m *Bit) Get() bool {
	return m.wireO.Get()
}

func (m *Bit) Refresh() {
	currentValue := m.Get()
	m.Update(currentValue, true)
}

func (m *Bit) Update(i bool, s bool) {
	m.wireI.Update(i)
	m.wireS.Update(s)
	for i := 0; i < 2; i++ {
		m.gates[0].Update(m.wireI.Get(), m.wireS.Get())
		m.gates[1].Update(m.gates[0].Output(), m.wireS.Get())
		m.gates[2].Update(m.gates[0].Output(), m.gates[3].Output())
		m.gates[3].Update(m.gates[2].Output(), m.gates[1].Output())
		m.wireO.Update(m.gates[2].Output())
	}
}

type Byte struct {
	inputs  [8]circuit.Wire
	bits    [8]Bit
	outputs [8]circuit.Wire
	next    ByteComponent
}

func NewByte() *Byte {
	b := new(Byte)
	for i, _ := range b.bits {
		b.bits[i] = *NewBit()
	}
	return b
}

func (e *Byte) ConnectOutput(b ByteComponent) {
	e.next = b
}

func (e *Byte) GetOutputWire(index int) bool {
	return e.outputs[index].Get()
}

func (e *Byte) SetInputWire(index int, value bool) {
	e.inputs[index].Update(value)
}

func (e *Byte) Update(set bool) {
	for i, input := range e.inputs {
		e.bits[i].Update(input.Get(), set)
		e.outputs[i].Update(e.bits[i].Get())
	}

	if e.next != nil {
		for i, w := range e.outputs {
			e.next.SetInputWire(i, w.Get())
		}
	}
}
