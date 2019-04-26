package memory

import (
	"testing"

	"github.com/djhworld/simple-computer/components"
)

func TestMemory64KWrite(t *testing.T) {
	bus := components.NewBus(BUS_WIDTH)
	m := NewMemory64K(bus)

	var i uint16
	var q uint16 = 0xFFFF
	for i = 0x0000; i < 0xFFFF; i++ {
		m.AddressRegister.Set()
		setBusValue(bus, i)
		m.Update()

		m.AddressRegister.Unset()
		m.Update()

		setBusValue(bus, q)
		m.Set()
		m.Update()

		m.Unset()
		m.Update()

		q--
	}

	var expected uint16 = 0xFFFF
	for i = 0x0000; i < 0xFFFF; i++ {
		m.AddressRegister.Set()
		setBusValue(bus, i)
		m.Update()

		m.AddressRegister.Unset()
		m.Update()

		m.Enable()
		m.Update()

		m.Disable()
		m.Update()

		checkBus(bus, expected)
		expected--
	}
}

func TestMemory64KDoesNotUpdateWhenSetFlagIsOff(t *testing.T) {
	bus := components.NewBus(BUS_WIDTH)
	m := NewMemory64K(bus)

	var i uint16
	var q uint16 = 0xFFFF
	for i = 0x0000; i < 0xFFFF; i++ {
		m.AddressRegister.Set()
		setBusValue(bus, i)
		m.Update()

		m.AddressRegister.Unset()
		m.Update()

		setBusValue(bus, q)

		m.Unset()
		m.Update()

		q--
	}

	var expected uint16 = 0xFFFF
	for i = 0x0000; i < 0xFFFF; i++ {
		m.AddressRegister.Set()
		setBusValue(bus, i)
		m.Update()

		m.AddressRegister.Unset()
		m.Update()

		m.Enable()
		m.Update()

		m.Disable()
		m.Update()

		checkBus(bus, expected)
	}
}

func setBusValue(b *components.Bus, value uint16) {
	for i := BUS_WIDTH - 1; i >= 0; i-- {
		r := (value & (1 << uint16(i)))
		if r != 0 {
			b.SetInputWire(i, true)
		} else {
			b.SetInputWire(i, false)
		}
	}
}

func checkBus(b *components.Bus, expected uint16) bool {
	var result uint16
	for i := BUS_WIDTH - 1; i >= 0; i-- {
		if b.GetOutputWire(i) {
			result = result | (1 << uint16(i))
		} else {
			result = result & ^(1 << uint16(i))
		}
	}
	return result == expected
}
