package components

import "github.com/djhworld/simple-computer/circuit"

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
	inputA circuit.Wire
	inputB circuit.Wire
	inputC circuit.Wire
	inputD circuit.Wire

	notGates [4]circuit.NOTGate
	andGates [16]ANDGate4
	outputs  [16]circuit.Wire
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
	for i := range d.outputs {
		if d.outputs[i].Get() {
			return i
		}
	}

	return 0
}

func (d *Decoder4x16) Update(inputA, inputB, inputC, inputD bool) {
	d.inputA.Update(inputA)
	d.inputB.Update(inputB)
	d.inputC.Update(inputC)
	d.inputD.Update(inputD)

	// https://www.elprocus.com/designing-4-to-16-decoder-using-3-to-8-decoder/
	d.notGates[0].Update(d.inputA.Get())
	d.notGates[1].Update(d.inputB.Get())
	d.notGates[2].Update(d.inputC.Get())
	d.notGates[3].Update(d.inputD.Get())

	d.andGates[0].Update(d.notGates[0].Output(), d.notGates[1].Output(), d.notGates[2].Output(), d.notGates[3].Output())
	d.andGates[1].Update(d.notGates[0].Output(), d.notGates[1].Output(), d.notGates[2].Output(), d.inputD.Get())
	d.andGates[2].Update(d.notGates[0].Output(), d.notGates[1].Output(), d.inputC.Get(), d.notGates[3].Output())
	d.andGates[3].Update(d.notGates[0].Output(), d.notGates[1].Output(), d.inputC.Get(), d.inputD.Get())

	d.andGates[4].Update(d.notGates[0].Output(), d.inputB.Get(), d.notGates[2].Output(), d.notGates[3].Output())
	d.andGates[5].Update(d.notGates[0].Output(), d.inputB.Get(), d.notGates[2].Output(), d.inputD.Get())
	d.andGates[6].Update(d.notGates[0].Output(), d.inputB.Get(), d.inputC.Get(), d.notGates[3].Output())
	d.andGates[7].Update(d.notGates[0].Output(), d.inputB.Get(), d.inputC.Get(), d.inputD.Get())

	d.andGates[8].Update(d.inputA.Get(), d.notGates[1].Output(), d.notGates[2].Output(), d.notGates[3].Output())
	d.andGates[9].Update(d.inputA.Get(), d.notGates[1].Output(), d.notGates[2].Output(), d.inputD.Get())
	d.andGates[10].Update(d.inputA.Get(), d.notGates[1].Output(), d.inputC.Get(), d.notGates[3].Output())
	d.andGates[11].Update(d.inputA.Get(), d.notGates[1].Output(), d.inputC.Get(), d.inputD.Get())

	d.andGates[12].Update(d.inputA.Get(), d.inputB.Get(), d.notGates[2].Output(), d.notGates[3].Output())
	d.andGates[13].Update(d.inputA.Get(), d.inputB.Get(), d.notGates[2].Output(), d.inputD.Get())
	d.andGates[14].Update(d.inputA.Get(), d.inputB.Get(), d.inputC.Get(), d.notGates[3].Output())
	d.andGates[15].Update(d.inputA.Get(), d.inputB.Get(), d.inputC.Get(), d.inputD.Get())

	for i, _ := range d.outputs {
		d.outputs[i].Update(d.andGates[i].output.Get())
	}
}
