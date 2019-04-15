package components

import (
	"fmt"

	"github.com/djhworld/simple-computer/circuit"
)

type Register struct {
	name      string
	set       circuit.Wire
	enable    circuit.Wire
	inputs    [8]circuit.Wire
	byter     *Byte
	enabler   *Enabler
	outputs   [8]circuit.Wire
	inputBus  *Bus
	outputBus *Bus
}

func NewRegister(name string, inputBus *Bus, outputBus *Bus) *Register {
	r := new(Register)
	r.name = name
	r.byter = NewByte()
	r.enabler = NewEnabler()
	r.enable = *circuit.NewWire("E", false)
	r.set = *circuit.NewWire("S", false)
	r.inputBus = inputBus
	r.outputBus = outputBus
	r.byter.ConnectOutput(r.enabler)
	return r
}

func (r *Register) Bit(index int) bool {
	return r.byter.GetOutputWire(index)
}

func (r *Register) Enable() {
	r.enable.Update(true)
}

func (r *Register) Disable() {
	r.enable.Update(false)
}

func (r *Register) Set() {
	r.set.Update(true)
}

func (r *Register) Unset() {
	r.set.Update(false)
}

func (r *Register) Update() {
	if r.inputBus != nil {
		//for i := 0; i < 8; i++ {
		for i := 8 - 1; i >= 0; i-- {
			r.inputs[i].Update(r.inputBus.GetOutputWire(i))
		}
	}

	for i, w := range r.inputs {
		r.byter.SetInputWire(i, w.Get())
	}

	r.byter.Update(r.set.Get())
	r.enabler.Update(r.enable.Get())

	for i, w := range r.enabler.outputs {
		r.outputs[i].Update(w.Get())
	}

	if r.enable.Get() && r.outputBus != nil {
		//for i, w := range r.outputs {
		for i := 8 - 1; i >= 0; i-- {
			r.outputBus.SetInputWire(i, r.outputs[i].Get())
		}
	}
}

func (r *Register) Value() byte {
	var value byte
	var x int = 0
	for i := 7; i >= 0; i-- {
		if r.byter.GetOutputWire(i) {
			value = value | (1 << byte(x))
		} else {
			value = value & ^(1 << byte(x))
		}
		x++
	}

	return value
}

func (r *Register) String() string {
	var output byte
	var x int = 0
	for i := 7; i >= 0; i-- {
		if r.outputs[i].Get() {
			output = output | (1 << byte(x))
		} else {
			output = output & ^(1 << byte(x))
		}
		x++
	}

	return fmt.Sprintf("%s: 0x%X (output = 0x%X) E: %v S: %v", r.name, r.Value(), output, r.enable.Get(), r.set.Get())
}
