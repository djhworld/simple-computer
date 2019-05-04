package io

import (
	"testing"

	"github.com/djhworld/simple-computer/components"
)

func TestAdapterPutsKeycodeRegisterOnBusAndClearsIt(t *testing.T) {
	ioBus := components.NewIOBus()
	mainBus := components.NewBus(BUS_WIDTH)

	adapter := NewKeyboardAdapter()
	adapter.Connect(ioBus, mainBus)

	mainBus.SetValue(0x000F)
	adapter.KeyboardInBus.SetValue(0x1234)

	adapter.Update()

	ioBus.Set()
	ioBus.Update(true, true)
	adapter.Update()
	ioBus.Unset()
	adapter.Update()

	ioBus.Enable()
	ioBus.Update(false, false)
	adapter.Update()

	adapter.Update()

	if !checkBus(mainBus, 0x1234) {
		t.FailNow()
	}

	if adapter.keycodeRegister.Value() != 0x0000 {
		t.FailNow()
	}
}

func checkBus(b *components.Bus, expected uint16) bool {
	var x int = 0
	var result uint16
	for i := BUS_WIDTH - 1; i >= 0; i-- {
		if b.GetOutputWire(i) {
			result = result | (1 << uint16(x))
		} else {
			result = result & ^(1 << uint16(x))
		}
		x++
	}
	return result == expected
}
