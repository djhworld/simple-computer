package components

import (
	"github.com/djhworld/simple-computer/circuit"
)

type Decoder2x4 struct {
	inputA circuit.Wire
	inputB circuit.Wire

	notGates [2]circuit.NOTGate
	andGates [4]circuit.ANDGate
	outputs  [4]circuit.Wire
}

func NewDecoder2x4() *Decoder2x4 {
	d := new(Decoder2x4)

	for i, _ := range d.notGates {
		d.notGates[i] = *circuit.NewNOTGate()
	}

	for i, _ := range d.andGates {
		d.andGates[i] = *circuit.NewANDGate()
	}

	return d
}

func (d *Decoder2x4) GetOutputWire(index int) bool {
	return d.outputs[index].Get()
}

func (d *Decoder2x4) Update(inputA bool, inputB bool) {
	d.inputA.Update(inputA)
	d.inputB.Update(inputB)

	d.notGates[0].Update(d.inputA.Get())
	d.notGates[1].Update(d.inputB.Get())

	d.andGates[0].Update(d.notGates[0].Output(), d.notGates[1].Output())
	d.andGates[1].Update(d.notGates[0].Output(), d.inputB.Get())
	d.andGates[2].Update(d.inputA.Get(), d.notGates[1].Output())
	d.andGates[3].Update(d.inputA.Get(), d.inputB.Get())

	for i, _ := range d.outputs {
		d.outputs[i].Update(d.andGates[i].Output())
	}
}

type Decoder3x8 struct {
	inputA circuit.Wire
	inputB circuit.Wire
	inputC circuit.Wire

	notGates [3]circuit.NOTGate
	andGates [8]ANDGate3
	outputs  [8]circuit.Wire
}

func NewDecoder3x8() *Decoder3x8 {
	d := new(Decoder3x8)

	for i, _ := range d.notGates {
		d.notGates[i] = *circuit.NewNOTGate()
	}

	for i, _ := range d.andGates {
		d.andGates[i] = *NewANDGate3()
	}

	return d
}

func (d *Decoder3x8) GetOutputWire(index int) bool {
	return d.outputs[index].Get()
}

// Returns the index which is enabled
func (d *Decoder3x8) Index() int {
	for i := range d.outputs {
		if d.outputs[i].Get() {
			return i
		}
	}

	return 0
}

func (d *Decoder3x8) Update(inputA, inputB, inputC bool) {
	d.inputA.Update(inputA)
	d.inputB.Update(inputB)
	d.inputC.Update(inputC)

	d.notGates[0].Update(d.inputA.Get())
	d.notGates[1].Update(d.inputB.Get())
	d.notGates[2].Update(d.inputC.Get())

	d.andGates[0].Update(d.notGates[0].Output(), d.notGates[1].Output(), d.notGates[2].Output())
	d.andGates[1].Update(d.notGates[0].Output(), d.notGates[1].Output(), d.inputC.Get())
	d.andGates[2].Update(d.notGates[0].Output(), d.inputB.Get(), d.notGates[2].Output())
	d.andGates[3].Update(d.notGates[0].Output(), d.inputB.Get(), d.inputC.Get())

	d.andGates[4].Update(d.inputA.Get(), d.notGates[1].Output(), d.notGates[2].Output())
	d.andGates[5].Update(d.inputA.Get(), d.notGates[1].Output(), d.inputC.Get())
	d.andGates[6].Update(d.inputA.Get(), d.inputB.Get(), d.notGates[2].Output())
	d.andGates[7].Update(d.inputA.Get(), d.inputB.Get(), d.inputC.Get())

	for i, _ := range d.outputs {
		d.outputs[i].Update(d.andGates[i].Output())
	}
}

type Decoder4x16 struct {
	notGates [4]circuit.NOTGate
	andGates [16]ANDGate4
	outputs  [16]circuit.Wire
	index    int
}

func NewDecoder4x16() *Decoder4x16 {
	d := new(Decoder4x16)

	for i, _ := range d.notGates {
		d.notGates[i] = *circuit.NewNOTGate()
	}

	for i, _ := range d.andGates {
		d.andGates[i] = *NewANDGate4()
	}

	return d
}

// Returns the index which is enabled
func (d *Decoder4x16) Index() int {
	return d.index
}

func (d *Decoder4x16) GetOutputWire(index int) bool {
	return d.outputs[index].Get()
}

func (d *Decoder4x16) Update(inputA, inputB, inputC, inputD bool) {
	// https://www.elprocus.com/designing-4-to-16-decoder-using-3-to-8-decoder/
	d.notGates[0].Update(inputA)
	d.notGates[1].Update(inputB)
	d.notGates[2].Update(inputC)
	d.notGates[3].Update(inputD)

	d.andGates[0].Update(d.notGates[0].Output(), d.notGates[1].Output(), d.notGates[2].Output(), d.notGates[3].Output())
	d.andGates[1].Update(d.notGates[0].Output(), d.notGates[1].Output(), d.notGates[2].Output(), inputD)
	d.andGates[2].Update(d.notGates[0].Output(), d.notGates[1].Output(), inputC, d.notGates[3].Output())
	d.andGates[3].Update(d.notGates[0].Output(), d.notGates[1].Output(), inputC, inputD)

	d.andGates[4].Update(d.notGates[0].Output(), inputB, d.notGates[2].Output(), d.notGates[3].Output())
	d.andGates[5].Update(d.notGates[0].Output(), inputB, d.notGates[2].Output(), inputD)
	d.andGates[6].Update(d.notGates[0].Output(), inputB, inputC, d.notGates[3].Output())
	d.andGates[7].Update(d.notGates[0].Output(), inputB, inputC, inputD)

	d.andGates[8].Update(inputA, d.notGates[1].Output(), d.notGates[2].Output(), d.notGates[3].Output())
	d.andGates[9].Update(inputA, d.notGates[1].Output(), d.notGates[2].Output(), inputD)
	d.andGates[10].Update(inputA, d.notGates[1].Output(), inputC, d.notGates[3].Output())
	d.andGates[11].Update(inputA, d.notGates[1].Output(), inputC, inputD)

	d.andGates[12].Update(inputA, inputB, d.notGates[2].Output(), d.notGates[3].Output())
	d.andGates[13].Update(inputA, inputB, d.notGates[2].Output(), inputD)
	d.andGates[14].Update(inputA, inputB, inputC, d.notGates[3].Output())
	d.andGates[15].Update(inputA, inputB, inputC, inputD)

	d.index = 0
	for i := 0; i < len(d.outputs); i++ {
		d.outputs[i].Update(d.andGates[i].output.Get())
		if d.outputs[i].Get() {
			d.index += i
		}
	}
}

// Decoder is constructed from 16 4x16 decoders
// https://www.quora.com/How-can-I-construct-a-8X256-decoder-using-4X16-decoders-only
type Decoder8x256 struct {
	decoderSelector Decoder4x16
	decoders4x16    [16]Decoder4x16
	index           int
}

func NewDecoder8x256() *Decoder8x256 {
	d := new(Decoder8x256)

	d.decoderSelector = *NewDecoder4x16()

	for i := range d.decoders4x16 {
		d.decoders4x16[i] = *NewDecoder4x16()
	}

	return d
}

// Returns the index which is enabled
func (d *Decoder8x256) Index() int {
	return d.index
}

func (dc *Decoder8x256) Update(a, b, c, d, e, f, g, h bool) {
	dc.index = 0

	dc.decoderSelector.Update(e, f, g, h)
	for i := 0; i < 16; i++ {
		dc.updateDecoder(a, b, c, d, i, 16*i)
	}
}

func (dc *Decoder8x256) updateDecoder(a, b, c, d bool, decoderIndex int, outputWireStart int) {
	if dc.decoderSelector.GetOutputWire(decoderIndex) {
		dc.decoders4x16[decoderIndex].Update(a, b, c, d)

		for i := 0; i < 16; i++ {
			if dc.decoders4x16[decoderIndex].outputs[i].Get() {
				dc.index = outputWireStart + i
			}
		}
	}
}
