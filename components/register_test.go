package components

import (
	"testing"
)

func TestRegisterByteIsSet(t *testing.T) {
	b := NewBus()
	setBus(b, 0x57)

	r := NewRegister("r", b, b)

	r.Set()
	r.Disable()
	r.Update()

	if !checkValueIs(r.byter, 0x57) {
		t.Fail()
	}

	r.Unset()
	setBus(b, 0x21)
	r.Update()

	// value should not change
	if !checkValueIs(r.byter, 0x57) {
		t.Fail()
	}

	r.Set()
	setBus(b, 0x39)
	r.Update()

	// value should change
	if !checkValueIs(r.byter, 0x39) {
		t.Fail()
	}
}

func TestRegisterOutputIsZeroWhenDisabled(t *testing.T) {
	b := NewBus()
	setBus(b, 0xF1)

	r := NewRegister("r", b, b)

	r.Set()
	r.Disable()
	r.Update()

	//set bus to new value
	setBus(b, 0x21)

	if !checkBus(b, 0x21) {
		t.Fail()
	}

	if !checkRegisterOutput(r, 0x00) {
		t.Fail()
	}
}

func TestRegisterOutputIsByteValueWhenEnable(t *testing.T) {
	b := NewBus()
	setBus(b, 0xF1)

	r := NewRegister("r", b, b)

	r.Set()
	r.Enable()
	r.Update()

	if !checkBus(b, 0xF1) {
		t.Fail()
	}
}

func TestBusIsUpdatedOnRegisterEnable(t *testing.T) {
	b := NewBus()
	r := NewRegister("r", b, b)
	setBus(b, 0xF1)

	r.Disable()
	r.Set()

	r.Update()
	r.Unset()
	setBus(b, 0x28)
	r.Enable()
	r.Update()

	if !checkBus(b, 0xF1) {
		t.Fail()
	}
}

func TestOutputBusIsUsedWhenDifferentFromInputBus(t *testing.T) {
	inputBus := NewBus()
	outputBus := NewBus()
	r := NewRegister("r", inputBus, outputBus)
	setBus(inputBus, 0xF1)
	setBus(outputBus, 0x93)

	r.Disable()
	r.Set()

	r.Update()
	r.Unset()
	setBus(inputBus, 0x28)
	r.Enable()
	r.Update()

	if !checkBus(inputBus, 0x28) {
		t.Fail()
	}

	// output bus should contain the register value
	if !checkBus(outputBus, 0xF1) {
		t.Fail()
	}
}

func TestBusIsNOTUpdatedOnRegisterDisabled(t *testing.T) {
	b := NewBus()
	r := NewRegister("r", b, b)
	setBus(b, 0xF1)

	r.Set()
	r.Disable()
	r.Update()

	r.Unset()
	setBus(b, 0x33)
	r.Disable()
	r.Update()

	if !checkBus(b, 0x33) {
		t.Fail()
	}
}

func checkValueIs(b ByteComponent, expected byte) bool {
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

func checkRegisterOutput(r *Register, expected byte) bool {
	var result byte
	for i := 7; i >= 0; i-- {
		if r.outputs[i].Get() {
			result = result | (1 << byte(i))
		} else {
			result = result & ^(1 << byte(i))
		}
	}
	return result == expected
}

func setBus(b *Bus, value byte) {
	for i := 7; i >= 0; i-- {
		r := (value & (1 << byte(i)))
		if r != 0 {
			b.SetInputWire(i, true)
		} else {
			b.SetInputWire(i, false)
		}
	}
}

func checkBus(b *Bus, expected byte) bool {
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
