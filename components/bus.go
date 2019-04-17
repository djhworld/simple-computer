package components

import "github.com/djhworld/simple-computer/circuit"

type Bus struct {
	wires [8]circuit.Wire
}

func NewBus() *Bus {
	b := new(Bus)
	for i, _ := range b.wires {
		b.wires[i] = *circuit.NewWire("", false)
	}

	return b
}

func (b *Bus) ConnectOutput(ByteComponent) {

}
func (b *Bus) SetInputWire(index int, value bool) {
	b.wires[index].Update(value)
}
func (b *Bus) GetOutputWire(index int) bool {
	return b.wires[index].Get()
}

func (b *Bus) String() string {
	result := ""
	for i := 0; i < 8; i++ {
		if b.GetOutputWire(i) {
			result += "1"
		} else {
			result += "0"
		}
	}
	return result
}
