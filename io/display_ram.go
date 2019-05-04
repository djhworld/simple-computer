package io

import (
	"github.com/djhworld/simple-computer/circuit"
	"github.com/djhworld/simple-computer/components"
	"github.com/djhworld/simple-computer/memory"
)

// Display RAM is special as the writes (inputs) and reads (outputs) are two separate
// units that operate independently.
type displayRAM struct {
	InputAddressRegister components.Register
	inputRowDecoder      components.Decoder8x256
	inputColDecoder      components.Decoder8x256

	OutputAddressRegister components.Register
	outputRowDecoder      components.Decoder8x256
	outputColDecoder      components.Decoder8x256

	data      [256][256]memory.Cell
	set       circuit.Wire
	enable    circuit.Wire
	inputBus  *components.Bus
	outputBus *components.Bus
}

func newDisplayRAM(inputBus, outputBus *components.Bus) *displayRAM {
	m := new(displayRAM)
	m.InputAddressRegister = *components.NewRegister("IMAR", inputBus, inputBus)
	m.inputRowDecoder = *components.NewDecoder8x256()
	m.inputColDecoder = *components.NewDecoder8x256()

	m.OutputAddressRegister = *components.NewRegister("OMAR", outputBus, outputBus)
	m.outputRowDecoder = *components.NewDecoder8x256()
	m.outputColDecoder = *components.NewDecoder8x256()

	m.inputBus = inputBus
	m.outputBus = outputBus

	// 0xF0 x 0xA0
	for i := 0; i < 256; i++ {
		for j := 0; j < 256; j++ {
			c := memory.NewCell(inputBus, outputBus)
			inputBus.SetValue(0x0000)
			c.Update(true, false)
			c.Update(false, false)
			m.data[i][j] = *c
		}
	}

	return m
}

func (m *displayRAM) Enable() {
	m.enable.Update(true)
}

func (m *displayRAM) Disable() {
	m.enable.Update(false)
}

func (m *displayRAM) Set() {
	m.set.Update(true)
}

func (m *displayRAM) Unset() {
	m.set.Update(false)
}

func (m *displayRAM) UpdateIncoming() {
	m.InputAddressRegister.Update()
	m.inputRowDecoder.Update(
		m.InputAddressRegister.Bit(0),
		m.InputAddressRegister.Bit(1),
		m.InputAddressRegister.Bit(2),
		m.InputAddressRegister.Bit(3),
		m.InputAddressRegister.Bit(4),
		m.InputAddressRegister.Bit(5),
		m.InputAddressRegister.Bit(6),
		m.InputAddressRegister.Bit(7),
	)
	m.inputColDecoder.Update(
		m.InputAddressRegister.Bit(8),
		m.InputAddressRegister.Bit(9),
		m.InputAddressRegister.Bit(10),
		m.InputAddressRegister.Bit(11),
		m.InputAddressRegister.Bit(12),
		m.InputAddressRegister.Bit(13),
		m.InputAddressRegister.Bit(14),
		m.InputAddressRegister.Bit(15),
	)

	var row int = m.inputRowDecoder.Index()
	var col int = m.inputColDecoder.Index()

	m.data[row][col].Update(m.set.Get(), false)
}

func (m *displayRAM) UpdateOutgoing() {
	m.OutputAddressRegister.Update()
	m.outputRowDecoder.Update(
		m.OutputAddressRegister.Bit(0),
		m.OutputAddressRegister.Bit(1),
		m.OutputAddressRegister.Bit(2),
		m.OutputAddressRegister.Bit(3),
		m.OutputAddressRegister.Bit(4),
		m.OutputAddressRegister.Bit(5),
		m.OutputAddressRegister.Bit(6),
		m.OutputAddressRegister.Bit(7),
	)
	m.outputColDecoder.Update(
		m.OutputAddressRegister.Bit(8),
		m.OutputAddressRegister.Bit(9),
		m.OutputAddressRegister.Bit(10),
		m.OutputAddressRegister.Bit(11),
		m.OutputAddressRegister.Bit(12),
		m.OutputAddressRegister.Bit(13),
		m.OutputAddressRegister.Bit(14),
		m.OutputAddressRegister.Bit(15),
	)

	var row int = m.outputRowDecoder.Index()
	var col int = m.outputColDecoder.Index()

	m.data[row][col].Update(false, m.enable.Get())
}
