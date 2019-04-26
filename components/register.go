package components

import (
	"fmt"

	"github.com/djhworld/simple-computer/circuit"
	"github.com/djhworld/simple-computer/utils"
)

type Register struct {
	name      string
	set       circuit.Wire
	enable    circuit.Wire
	word      *Word
	enabler   *Enabler
	outputs   [BUS_WIDTH]circuit.Wire
	inputBus  *Bus
	outputBus *Bus
}

func NewRegister(name string, inputBus *Bus, outputBus *Bus) *Register {
	r := new(Register)
	r.name = name
	r.word = NewWord()
	r.enabler = NewEnabler()
	r.enable = *circuit.NewWire("E", false)
	r.set = *circuit.NewWire("S", false)
	r.inputBus = inputBus
	r.outputBus = outputBus
	r.word.ConnectOutput(r.enabler)
	return r
}

func (r *Register) Bit(index int) bool {
	return r.word.GetOutputWire(index)
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
	for i := BUS_WIDTH - 1; i >= 0; i-- {
		r.word.SetInputWire(i, r.inputBus.GetOutputWire(i))
	}

	r.word.Update(r.set.Get())
	r.enabler.Update(r.enable.Get())

	for i := 0; i < len(r.enabler.outputs); i++ {
		r.outputs[i].Update(r.enabler.outputs[i].Get())
	}

	if r.enable.Get() {
		for i := BUS_WIDTH - 1; i >= 0; i-- {
			r.outputBus.SetInputWire(i, r.outputs[i].Get())
		}
	}
}

func (r *Register) Value() uint16 {
	var value uint16
	var x uint16 = 0
	for i := BUS_WIDTH - 1; i >= 0; i-- {
		if r.word.GetOutputWire(i) {
			value = value | (1 << x)
		} else {
			value = value & ^(1 << x)
		}
		x++
	}

	return value
}

func (r *Register) String() string {
	var output uint16
	var x uint16 = 0
	for i := BUS_WIDTH - 1; i >= 0; i-- {
		if r.outputs[i].Get() {
			output = output | (1 << x)
		} else {
			output = output & ^(1 << x)
		}
		x++
	}

	return fmt.Sprintf("%s: %s (output = %s) E: %v S: %v", r.name, utils.ValueToString(r.Value()), utils.ValueToString(output), r.enable.Get(), r.set.Get())
}

