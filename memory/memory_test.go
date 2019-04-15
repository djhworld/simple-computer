package memory

import (
	"fmt"
	"testing"

	"github.com/djhworld/simple-computer/components"
)

func TestMemory256Write(t *testing.T) {
	bus := components.NewBus()
	m := NewMemory256(bus)

	var i byte
	var q byte = 0xFF
	for i = 0x00; i < 0xFF; i++ {
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

	var expected byte = 0xFF
	for i = 0x00; i < 0xFF; i++ {
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

	fmt.Println(m)
}

func TestMemory256DoesNotUpdateWhenSetFlagIsOff(t *testing.T) {
	bus := components.NewBus()
	m := NewMemory256(bus)

	var i byte
	var q byte = 0xFF
	for i = 0x00; i < 0xFF; i++ {
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

	var expected byte = 0xFF
	for i = 0x00; i < 0xFF; i++ {
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

func setBusValue(b *components.Bus, value byte) {
	for i := 7; i >= 0; i-- {
		r := (value & (1 << byte(i)))
		if r != 0 {
			b.SetInputWire(i, true)
		} else {
			b.SetInputWire(i, false)
		}
	}
}

func checkBus(b *components.Bus, expected byte) bool {
	var result byte
	for i := 7; i >= 0; i-- {
		if b.GetOutputWire(i) {
			result = result | (1 << byte(i))
		} else {
			result = result & ^(1 << byte(i))
		}
	}
	return result == expected
}
