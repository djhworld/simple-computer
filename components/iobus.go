package components

import (
	"github.com/djhworld/simple-computer/circuit"
)

const (
	CLOCK_SET       = 0
	CLOCK_ENABLE    = 1
	MODE            = 2
	DATA_OR_ADDRESS = 3
)

type IOBus struct {
	wires [4]circuit.Wire
}

func NewIOBus() *IOBus {
	b := new(IOBus)
	b.wires[CLOCK_SET] = *circuit.NewWire("", false)
	b.wires[CLOCK_ENABLE] = *circuit.NewWire("", false)
	b.wires[MODE] = *circuit.NewWire("", false)
	b.wires[DATA_OR_ADDRESS] = *circuit.NewWire("", false)
	return b
}

func (i *IOBus) Set() {
	i.wires[CLOCK_SET].Update(true)
}

func (i *IOBus) Unset() {
	i.wires[CLOCK_SET].Update(false)
}

func (i *IOBus) Enable() {
	i.wires[CLOCK_ENABLE].Update(true)
}

func (i *IOBus) Disable() {
	i.wires[CLOCK_ENABLE].Update(false)
}

func (i *IOBus) IsSet() bool {
	return i.wires[CLOCK_SET].Get()
}

func (i *IOBus) IsEnable() bool {
	return i.wires[CLOCK_ENABLE].Get()
}

func (i *IOBus) IsInputMode() bool {
	return i.wires[MODE].Get() == false
}

func (i *IOBus) IsOutputMode() bool {
	return i.wires[MODE].Get() == true
}

func (i *IOBus) IsDataMode() bool {
	return i.wires[DATA_OR_ADDRESS].Get() == false
}

func (i *IOBus) IsAddressMode() bool {
	return i.wires[DATA_OR_ADDRESS].Get() == true
}

func (i *IOBus) Update(mode, dataOrAddress bool) {
	i.wires[MODE].Update(mode)
	i.wires[DATA_OR_ADDRESS].Update(dataOrAddress)
}

func (i *IOBus) GetOutputWire(index int) bool {
	return i.wires[index].Get()
}
