package components

import "github.com/djhworld/simple-computer/circuit"

type Bus struct {
	wires []circuit.Wire
	width int
}

func NewBus(width int) *Bus {
	b := new(Bus)
	b.width = width
	b.wires = make([]circuit.Wire, b.width)
	for i, _ := range b.wires {
		b.wires[i] = *circuit.NewWire("", false)
	}

	return b
}

func (b *Bus) ConnectOutput(Component) {

}

func (b *Bus) SetInputWire(index int, value bool) {
	b.wires[index].Update(value)
}
func (b *Bus) GetOutputWire(index int) bool {
	return b.wires[index].Get()
}

func (b *Bus) SetValue(value uint16) {
	var x = 0
	for i := b.width - 1; i >= 0; i-- {
		r := (value & (1 << uint16(x)))
		if r != 0 {
			b.SetInputWire(i, true)
		} else {
			b.SetInputWire(i, false)
		}

		x++
	}
}

func (b *Bus) String() string {
	result := ""
	for i := 0; i < b.width; i++ {
		if b.GetOutputWire(i) {
			result += "1"
		} else {
			result += "0"
		}
	}
	return result
}
