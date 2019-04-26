package components

import (
	"testing"
)

func TestRegisterWordIsSet(t *testing.T) {
	b := NewBus(BUS_WIDTH)
	setBus(b, 0x3957)

	r := NewRegister("r", b, b)

	r.Set()
	r.Disable()
	r.Update()

	if !checkValueIs(r.word, 0x3957) {
		t.Fail()
	}

	r.Unset()
	setBus(b, 0xFE21)
	r.Update()

	// value should not change
	if !checkValueIs(r.word, 0x3957) {
		t.Fail()
	}

	r.Set()
	setBus(b, 0x0039)
	r.Update()

	// value should change
	if !checkValueIs(r.word, 0x0039) {
		t.Fail()
	}
}

func TestRegisterOutputIsZeroWhenDisabled(t *testing.T) {
	b := NewBus(BUS_WIDTH)
	setBus(b, 0xABF1)

	r := NewRegister("r", b, b)

	r.Set()
	r.Disable()
	r.Update()

	//set bus to new value
	setBus(b, 0x0021)

	if !checkBus(b, 0x0021) {
		t.Fail()
	}

	if !checkRegisterOutput(r, 0x0000) {
		t.Fail()
	}
}

func TestRegisterOutputIsWordValueWhenEnable(t *testing.T) {
	b := NewBus(BUS_WIDTH)
	setBus(b, 0x54F1)

	r := NewRegister("r", b, b)

	r.Set()
	r.Enable()
	r.Update()

	if !checkBus(b, 0x54F1) {
		t.Fail()
	}
}

func TestBusIsUpdatedOnRegisterEnable(t *testing.T) {
	b := NewBus(BUS_WIDTH)
	r := NewRegister("r", b, b)
	setBus(b, 0x32F1)

	r.Disable()
	r.Set()

	r.Update()
	r.Unset()
	setBus(b, 0x9028)
	r.Enable()
	r.Update()

	if !checkBus(b, 0x32F1) {
		t.Fail()
	}
}

func TestOutputBusIsUsedWhenDifferentFromInputBus(t *testing.T) {
	inputBus := NewBus(BUS_WIDTH)
	outputBus := NewBus(BUS_WIDTH)
	r := NewRegister("r", inputBus, outputBus)
	setBus(inputBus, 0x00F1)
	setBus(outputBus, 0x0093)

	r.Disable()
	r.Set()

	r.Update()
	r.Unset()
	setBus(inputBus, 0x0028)
	r.Enable()
	r.Update()

	if !checkBus(inputBus, 0x0028) {
		t.Fail()
	}

	// output bus should contain the register value
	if !checkBus(outputBus, 0x00F1) {
		t.Fail()
	}
}

func TestBusIsNOTUpdatedOnRegisterDisabled(t *testing.T) {
	b := NewBus(BUS_WIDTH)
	r := NewRegister("r", b, b)
	setBus(b, 0x00F1)

	r.Set()
	r.Disable()
	r.Update()

	r.Unset()
	setBus(b, 0x0033)
	r.Disable()
	r.Update()

	if !checkBus(b, 0x0033) {
		t.Fail()
	}
}

func checkValueIs(b Component, expected uint16) bool {
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

func checkRegisterOutput(r *Register, expected uint16) bool {
	var result uint16
	for i := BUS_WIDTH - 1; i >= 0; i-- {
		if r.outputs[i].Get() {
			result = result | (1 << uint16(i))
		} else {
			result = result & ^(1 << uint16(i))
		}
	}
	return result == expected
}

func setBus(b *Bus, value uint16) {
	for i := BUS_WIDTH - 1; i >= 0; i-- {
		r := (value & (1 << uint16(i)))
		if r != 0 {
			b.SetInputWire(i, true)
		} else {
			b.SetInputWire(i, false)
		}
	}
}

func checkBus(b *Bus, expected uint16) bool {
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
