package memory

import (
	"fmt"
	"strings"

	"github.com/djhworld/simple-computer/circuit"
	"github.com/djhworld/simple-computer/components"
)

type Cell struct {
	value components.Register
	gates [3]circuit.ANDGate
}

func NewCell(bus *components.Bus) *Cell {
	c := new(Cell)
	c.value = *components.NewRegister("", bus, bus)
	c.gates[0] = *circuit.NewANDGate()
	c.gates[1] = *circuit.NewANDGate()
	c.gates[2] = *circuit.NewANDGate()
	return c
}

func (c *Cell) Update(set bool, enable bool) {
	c.gates[0].Update(true, true)
	c.gates[1].Update(c.gates[0].Output(), set)
	c.gates[2].Update(c.gates[0].Output(), enable)

	if c.gates[1].Output() {
		c.value.Set()
	} else {
		c.value.Unset()
	}

	if c.gates[2].Output() {
		c.value.Enable()
	} else {
		c.value.Disable()
	}
	c.value.Update()
}

type Memory256 struct {
	AddressRegister components.Register
	rowDecoder      components.Decoder4x16
	colDecoder      components.Decoder4x16
	data            [16][16]Cell
	set             circuit.Wire
	enable          circuit.Wire
	bus             *components.Bus
}

func NewMemory256(bus *components.Bus) *Memory256 {
	m := new(Memory256)
	m.AddressRegister = *components.NewRegister("MAR", bus, bus)
	m.rowDecoder = *components.NewDecoder4x16()
	m.colDecoder = *components.NewDecoder4x16()
	m.bus = bus

	for i := 0; i < 16; i++ {
		for j := 0; j < 16; j++ {
			m.data[i][j] = *NewCell(bus)
		}
	}

	return m
}

func (m *Memory256) Enable() {
	m.enable.Update(true)
}

func (m *Memory256) Disable() {
	m.enable.Update(false)
}

func (m *Memory256) Set() {
	m.set.Update(true)
}

func (m *Memory256) Unset() {
	m.set.Update(false)
}

func (m *Memory256) Update() {
	m.AddressRegister.Update()
	m.rowDecoder.Update(m.AddressRegister.Bit(0), m.AddressRegister.Bit(1), m.AddressRegister.Bit(2), m.AddressRegister.Bit(3))
	m.colDecoder.Update(m.AddressRegister.Bit(4), m.AddressRegister.Bit(5), m.AddressRegister.Bit(6), m.AddressRegister.Bit(7))

	var row int = m.rowDecoder.Index()
	var col int = m.colDecoder.Index()

	m.data[row][col].Update(m.set.Get(), m.enable.Get())
}

func (m *Memory256) String() string {
	var row int = m.rowDecoder.Index()
	var col int = m.colDecoder.Index()

	var builder strings.Builder
	builder.WriteString(fmt.Sprint("Memory\n--------------------------------------\n"))
	builder.WriteString(fmt.Sprintf("RD: %d\tCD: %d\tS: %v\tE: %v\t%s\n", row, col, m.set.Get(), m.enable.Get(), m.AddressRegister.String()))

	for i := 0; i < 16; i++ {
		for j := 0; j < 16; j++ {
			val := m.data[i][j].value.Value()
			if val <= 0x0F {
				builder.WriteString(fmt.Sprintf("0x0%X\t", val))
			} else {
				builder.WriteString(fmt.Sprintf("0x%X\t", val))
			}
		}
		builder.WriteString(fmt.Sprint("\n"))

	}
	return builder.String()

}
