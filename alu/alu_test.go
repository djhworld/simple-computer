package alu

import (
	"testing"

	"github.com/djhworld/simple-computer/components"
)

func TestAluAdd(t *testing.T) {
	testOp(ADD, 0x00, 0x00, false, 0x00, true, false, false, true, t)
	testOp(ADD, 0x00, 0x01, true, 0x02, false, false, false, false, t)
	testOp(ADD, 0x01, 0x02, false, 0x03, false, false, false, false, t)
	testOp(ADD, 0x02, 0x01, false, 0x03, false, true, false, false, t)
	testOp(ADD, 0x01, 0xFE, false, 0xFF, false, false, false, false, t)
	testOp(ADD, 0xFE, 0x01, false, 0xFF, false, true, false, false, t)
	testOp(ADD, 0xFF, 0x04, false, 0x03, false, true, true, false, t)
	testOp(ADD, 0xFF, 0x01, false, 0x00, false, true, true, true, t)

	//carry in situations
	testOp(ADD, 0x00, 0x00, true, 0x01, true, false, false, false, t)
	testOp(ADD, 0x02, 0x01, true, 0x04, false, true, false, false, t)
	testOp(ADD, 0xFF, 0x00, true, 0x00, false, true, true, true, t)

	j := 0x80
	for i := 1; i < 128; i++ {
		var IsEqual bool = false
		if i == j {
			IsEqual = true
		}

		var isLarger bool = false
		if i > j {
			isLarger = true
		}

		var isZero bool = false
		if i+j == 0 {
			isZero = true
		}
		testOp(ADD, i, j, false, i+j, IsEqual, isLarger, false, isZero, t)

		isLarger = false
		if j > i {
			isLarger = true
		}

		isZero = false
		if j+i == 0 {
			isZero = true
		}

		testOp(ADD, j, i, false, j+i, IsEqual, isLarger, false, isZero, t)
		j--
	}
}

func TestAluSHR(t *testing.T) {
	for i := 128; i > 1; i /= 2 {
		testOp(SHR, i, i, false, i/2, true, false, false, false, t)
		testOp(SHR, i, 0x00, false, i/2, false, true, false, false, t)
	}

	testOp(SHR, 0x00, 0x00, false, 0x00, true, false, false, true, t)
	testOp(SHR, 0x56, 0x56, false, 0x2B, true, false, false, false, t)
	testOp(SHR, 0x04, 0x01, false, 0x02, false, true, false, false, t)
	testOp(SHR, 0x72, 0x00, false, 0x39, false, true, false, false, t)

	//should carry out
	testOp(SHR, 0xA1, 0x01, false, 0x50, false, true, true, false, t)

	//should carry in
	testOp(SHR, 0x4A, 0x01, true, 0xA5, false, true, false, false, t)

	//should do nothing with inputB
	testOp(SHR, 0x00, 0x05, false, 0x00, false, false, false, true, t)
}

func TestAluSHL(t *testing.T) {
	for i := 1; i < 127; i *= 2 {
		testOp(SHL, i, i, false, i*2, true, false, false, false, t)
		testOp(SHL, i, 0x00, false, i*2, false, true, false, false, t)
	}

	testOp(SHL, 0x00, 0x00, false, 0x00, true, false, false, true, t)
	testOp(SHL, 0x59, 0x59, false, 0xB2, true, false, false, false, t)
	testOp(SHL, 0x04, 0x01, false, 0x08, false, true, false, false, t)

	testOp(SHL, 0x73, 0x00, false, 0xE6, false, true, false, false, t)

	//should carry out
	testOp(SHL, 0xAA, 0x01, false, 0x54, false, true, true, false, t)

	//should carry in
	testOp(SHL, 0x4A, 0x01, true, 0x95, false, true, false, false, t)

	//should do nothing with inputB
	testOp(SHL, 0x00, 0x05, false, 0x00, false, false, false, true, t)
}

func TestNOT(t *testing.T) {
	testOp(NOT, 0x00, 0x00, false, 0xFF, true, false, false, false, t)
	testOp(NOT, 0x01, 0x00, false, 0xFE, false, true, false, false, t)
	testOp(NOT, 0xAA, 0x00, false, 0x55, false, true, false, false, t)
	testOp(NOT, 0xFF, 0x00, false, 0x00, false, true, false, true, t)

	//should ignore  inputB
	testOp(NOT, 0xFF, 0x84, false, 0x00, false, true, false, true, t)
	//should ignore  carry in bit
	testOp(NOT, 0xB9, 0x00, true, 0x46, false, true, false, false, t)
}

func TestAND(t *testing.T) {
	testOp(AND, 0x00, 0x00, false, 0x00, true, false, false, true, t)
	testOp(AND, 0x01, 0x01, false, 0x01, true, false, false, false, t)
	testOp(AND, 0xFF, 0xFF, false, 0xFF, true, false, false, false, t)
	testOp(AND, 0xFF, 0x69, false, 0x69, false, true, false, false, t)
	testOp(AND, 0x69, 0xFF, false, 0x69, false, false, false, false, t)
	testOp(AND, 0x8A, 0x00, false, 0x00, false, true, false, true, t)
	testOp(AND, 0x00, 0x8A, false, 0x00, false, false, false, true, t)

	//should ignore  carry in bit
	testOp(AND, 0x9A, 0x9A, true, 0x9A, true, false, false, false, t)

	j := 0xFF
	for i := 1; i < 256; i++ {
		var IsEqual bool = false
		if i == j {
			IsEqual = true
		}

		var isLarger bool = false
		if i > j {
			isLarger = true
		}

		var isZero bool = false
		if i&j == 0 {
			isZero = true
		}
		testOp(AND, i, j, false, i&j, IsEqual, isLarger, false, isZero, t)

		isLarger = false
		if j > i {
			isLarger = true
		}

		isZero = false
		if j&i == 0 {
			isZero = true
		}

		testOp(AND, j, i, false, j&i, IsEqual, isLarger, false, isZero, t)
		j--
	}
}

func TestOR(t *testing.T) {
	testOp(OR, 0x00, 0x00, false, 0x00, true, false, false, true, t)
	testOp(OR, 0x01, 0x01, false, 0x01, true, false, false, false, t)
	testOp(OR, 0xFF, 0xFF, false, 0xFF, true, false, false, false, t)

	testOp(OR, 0xFF, 0x69, false, 0xFF, false, true, false, false, t)
	testOp(OR, 0x69, 0xFF, false, 0xFF, false, false, false, false, t)
	testOp(OR, 0x8A, 0x00, false, 0x8A, false, true, false, false, t)
	testOp(OR, 0x00, 0x8A, false, 0x8A, false, false, false, false, t)

	//should ignore  carry in bit
	testOp(OR, 0x9A, 0x9A, true, 0x9A, true, false, false, false, t)

	j := 0xFF
	for i := 1; i < 256; i++ {
		var IsEqual bool = false
		if i == j {
			IsEqual = true
		}

		var isLarger bool = false
		if i > j {
			isLarger = true
		}

		var isZero bool = false
		if i|j == 0 {
			isZero = true
		}
		testOp(OR, i, j, false, i|j, IsEqual, isLarger, false, isZero, t)

		isLarger = false
		if j > i {
			isLarger = true
		}

		isZero = false
		if j|i == 0 {
			isZero = true
		}

		testOp(OR, j, i, false, j|i, IsEqual, isLarger, false, isZero, t)
		j--
	}
}

func TestXOR(t *testing.T) {
	testOp(XOR, 0x00, 0x00, false, 0x00, true, false, false, true, t)
	testOp(XOR, 0x01, 0x01, false, 0x00, true, false, false, true, t)
	testOp(XOR, 0xFF, 0xFF, false, 0x00, true, false, false, true, t)

	testOp(XOR, 0xFF, 0x69, false, 0x96, false, true, false, false, t)
	testOp(XOR, 0x69, 0xFF, false, 0x96, false, false, false, false, t)
	testOp(XOR, 0x8A, 0x00, false, 0x8A, false, true, false, false, t)
	testOp(XOR, 0x00, 0x8A, false, 0x8A, false, false, false, false, t)

	testOp(XOR, 0x9A, 0x9A, true, 0x00, true, false, false, true, t)

	j := 0xFF
	for i := 1; i < 256; i++ {
		var IsEqual bool = false
		if i == j {
			IsEqual = true
		}

		var isLarger bool = false
		if i > j {
			isLarger = true
		}

		var isZero bool = false
		if i^j == 0 {
			isZero = true
		}
		testOp(XOR, i, j, false, i^j, IsEqual, isLarger, false, isZero, t)

		isLarger = false
		if j > i {
			isLarger = true
		}

		isZero = false
		if j^i == 0 {
			isZero = true
		}

		testOp(XOR, j, i, false, j^i, IsEqual, isLarger, false, isZero, t)
		j--
	}
}

func TestCMP(t *testing.T) {
	testOp(CMP, 0x00, 0x00, false, 0x00, true, false, false, false, t)
	testOp(CMP, 0x01, 0x01, false, 0x00, true, false, false, false, t)
	testOp(CMP, 0xFF, 0xFF, false, 0x00, true, false, false, false, t)

	j := 0xFF
	for i := 1; i < 256; i++ {
		var IsEqual bool = false
		if i == j {
			IsEqual = true
		}

		var isLarger bool = false
		if i > j {
			isLarger = true
		}

		testOp(CMP, i, j, false, 0x00, IsEqual, isLarger, false, false, t)

		isLarger = false
		if j > i {
			isLarger = true
		}

		testOp(CMP, j, i, false, 0x00, IsEqual, isLarger, false, false, t)
		j--
	}

}

func testOp(op byte, inputA, inputB int, CarryIn bool, expectedOutput int, expectedEqual, expectedIsLarger, expectedCarry, expectedZero bool, t *testing.T) {
	inputABus := components.NewBus()
	inputBBus := components.NewBus()
	outputBus := components.NewBus()
	alu := NewALU(inputABus, inputBBus, outputBus)
	setBusValue(inputABus, inputA)
	setBusValue(inputBBus, inputB)
	setOp(alu, op)
	alu.CarryIn.Update(CarryIn)
	alu.Update()

	output := getValueOfBus(outputBus)
	if output != expectedOutput {
		t.Logf("Expected output of 0x%X but got 0x%X", expectedOutput, output)
		t.FailNow()
	}

	if alu.IsEqual.Get() != expectedEqual {
		t.Logf("Expected equal flag to be %v but got %v", expectedEqual, alu.IsEqual.Get())
		t.FailNow()
	}

	if alu.AisLarger.Get() != expectedIsLarger {
		t.Logf("Expected is larger flag to be %v but got %v", expectedIsLarger, alu.AisLarger.Get())
		t.FailNow()
	}

	if alu.CarryOut.Get() != expectedCarry {
		t.Logf("Expected is carry out flag to be %v but got %v", expectedCarry, alu.CarryOut.Get())
		t.FailNow()
	}

	if alu.isZero.GetOutputWire(0) != expectedZero {
		t.Logf("Expected zero flag to be %v but got %v", expectedEqual, alu.isZero.GetOutputWire(0))
		t.FailNow()
	}
}

func setBusValue(bus *components.Bus, value int) {
	x := 0
	for i := 8 - 1; i >= 0; i-- {
		r := (value & (1 << byte(x)))
		if r != 0 {
			bus.SetInputWire(i, true)
		} else {
			bus.SetInputWire(i, false)
		}
		x++
	}
}

func setOp(a *ALU, value byte) {
	value = value & 0x07
	for i := 2; i >= 0; i-- {
		r := (value & (1 << byte(i)))
		if r != 0 {
			a.Op[i].Update(true)
		} else {
			a.Op[i].Update(false)
		}
	}
}

func getValueOfBus(bus *components.Bus) int {
	var x int = 0
	var result int
	//for i := 0; i < 8; i++ {
	for i := 8 - 1; i >= 0; i-- {
		if bus.GetOutputWire(i) {
			result = result | (1 << byte(x))
		} else {
			result = result & ^(1 << byte(x))
		}
		x++
	}
	return result
}
