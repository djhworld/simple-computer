package memory

import (
	"fmt"
	"strings"

	"github.com/djhworld/simple-computer/circuit"
	"github.com/djhworld/simple-computer/components"
)

const BUS_WIDTH = 16

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

type Memory64K struct {
	AddressRegister components.Register
	rowDecoder      components.Decoder8x256
	colDecoder      components.Decoder8x256
	data            [256][256]Cell
	set             circuit.Wire
	enable          circuit.Wire
	bus             *components.Bus
}

func NewMemory64K(bus *components.Bus) *Memory64K {
	m := new(Memory64K)
	m.AddressRegister = *components.NewRegister("MAR", bus, bus)
	m.rowDecoder = *components.NewDecoder8x256()
	m.colDecoder = *components.NewDecoder8x256()
	m.bus = bus

	for i := 0; i < 256; i++ {
		for j := 0; j < 256; j++ {
			m.data[i][j] = *NewCell(bus)
		}
	}

	return m
}

func (m *Memory64K) Enable() {
	m.enable.Update(true)
}

func (m *Memory64K) Disable() {
	m.enable.Update(false)
}

func (m *Memory64K) Set() {
	m.set.Update(true)
}

func (m *Memory64K) Unset() {
	m.set.Update(false)
}

func (m *Memory64K) Update() {
	m.AddressRegister.Update()
	m.rowDecoder.Update(
		m.AddressRegister.Bit(0),
		m.AddressRegister.Bit(1),
		m.AddressRegister.Bit(2),
		m.AddressRegister.Bit(3),
		m.AddressRegister.Bit(4),
		m.AddressRegister.Bit(5),
		m.AddressRegister.Bit(6),
		m.AddressRegister.Bit(7),
	)
	m.colDecoder.Update(
		m.AddressRegister.Bit(8),
		m.AddressRegister.Bit(9),
		m.AddressRegister.Bit(10),
		m.AddressRegister.Bit(11),
		m.AddressRegister.Bit(12),
		m.AddressRegister.Bit(13),
		m.AddressRegister.Bit(14),
		m.AddressRegister.Bit(15),
	)

	var row int = m.rowDecoder.Index()
	var col int = m.colDecoder.Index()

	m.data[row][col].Update(m.set.Get(), m.enable.Get())
}

func (m *Memory64K) String() string {
	var row int = m.rowDecoder.Index()
	var col int = m.colDecoder.Index()

	var builder strings.Builder
	builder.WriteString(fmt.Sprint("Memory\n--------------------------------------\n"))
	builder.WriteString(fmt.Sprintf("RD: %d\tCD: %d\tS: %v\tE: %v\t%s\n", row, col, m.set.Get(), m.enable.Get(), m.AddressRegister.String()))

	for i := 0; i < 256; i++ {
		for j := 0; j < 256; j++ {
			val := m.data[i][j].value.Value()
			if val <= 0x000F {
				builder.WriteString(fmt.Sprintf("0x000%X\t", val))
			} else if val <= 0x00FF {
				builder.WriteString(fmt.Sprintf("0x00%X\t", val))
			} else if val <= 0x0FFF {
				builder.WriteString(fmt.Sprintf("0x0%X\t", val))
			} else {
				builder.WriteString(fmt.Sprintf("0x%X\t", val))
			}
		}
		builder.WriteString(fmt.Sprint("\n"))

	}
	return builder.String()

}
