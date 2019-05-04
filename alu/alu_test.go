package alu

import (
	"testing"

	"github.com/djhworld/simple-computer/components"
)

var inputABus *components.Bus = components.NewBus(BUS_WIDTH)
var inputBBus *components.Bus = components.NewBus(BUS_WIDTH)
var outputBus *components.Bus = components.NewBus(BUS_WIDTH)
var flagsBus *components.Bus = components.NewBus(BUS_WIDTH)

func TestAluAdd(t *testing.T) {
	alu := NewALU(inputABus, inputBBus, outputBus, flagsBus)
	testOp(alu, ADD, 0x0000, 0x0000, false, 0x0000, true, false, false, true, t)
	testOp(alu, ADD, 0x0000, 0x0001, true, 0x0002, false, false, false, false, t)
	testOp(alu, ADD, 0x0001, 0x0002, false, 0x0003, false, false, false, false, t)
	testOp(alu, ADD, 0x0002, 0x0001, false, 0x0003, false, true, false, false, t)
	testOp(alu, ADD, 0x0001, 0x00FE, false, 0x00FF, false, false, false, false, t)
	testOp(alu, ADD, 0x00FE, 0x0001, false, 0x00FF, false, true, false, false, t)
	testOp(alu, ADD, 0xFFFF, 0x0004, false, 0x0003, false, true, true, false, t)
	testOp(alu, ADD, 0xFFFF, 0x0001, false, 0x0000, false, true, true, true, t)

	//carry in situations
	testOp(alu, ADD, 0x0000, 0x0000, true, 0x0001, true, false, false, false, t)
	testOp(alu, ADD, 0x0002, 0x0001, true, 0x0004, false, true, false, false, t)
	testOp(alu, ADD, 0xFFFF, 0x0000, true, 0x0000, false, true, true, true, t)

	j := uint16(0x8000)
	for i := uint16(1); i < 32768; i++ {
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
		testOp(alu, ADD, i, j, false, i+j, IsEqual, isLarger, false, isZero, t)

		isLarger = false
		if j > i {
			isLarger = true
		}

		isZero = false
		if j+i == 0 {
			isZero = true
		}

		testOp(alu, ADD, j, i, false, j+i, IsEqual, isLarger, false, isZero, t)
		j--
	}
}

func TestAluSHR(t *testing.T) {
	alu := NewALU(inputABus, inputBBus, outputBus, flagsBus)
	for i := uint16(32768); i > 1; i /= 2 {
		testOp(alu, SHR, i, i, false, i/2, true, false, false, false, t)
		testOp(alu, SHR, i, 0x00, false, i/2, false, true, false, false, t)
	}

	testOp(alu, SHR, 0x0000, 0x0000, false, 0x0000, true, false, false, true, t)
	testOp(alu, SHR, 0x0056, 0x0056, false, 0x002B, true, false, false, false, t)
	testOp(alu, SHR, 0x0004, 0x0001, false, 0x0002, false, true, false, false, t)
	testOp(alu, SHR, 0x0072, 0x0000, false, 0x0039, false, true, false, false, t)

	//should carry out
	testOp(alu, SHR, 0x00A1, 0x0001, false, 0x0050, false, true, true, false, t)

	//should carry in
	testOp(alu, SHR, 0x4A00, 0x01, true, 0xA500, false, true, false, false, t)

	//should do nothing with inputB
	testOp(alu, SHR, 0x0000, 0x0005, false, 0x0000, false, false, false, true, t)
}

func TestAluSHL(t *testing.T) {
	alu := NewALU(inputABus, inputBBus, outputBus, flagsBus)
	for i := uint16(1); i < 32767; i *= 2 {
		testOp(alu, SHL, i, i, false, i*2, true, false, false, false, t)
		testOp(alu, SHL, i, 0x0000, false, i*2, false, true, false, false, t)
	}

	testOp(alu, SHL, 0x0000, 0x0000, false, 0x0000, true, false, false, true, t)
	testOp(alu, SHL, 0x0059, 0x0059, false, 0x00B2, true, false, false, false, t)
	testOp(alu, SHL, 0x0004, 0x0001, false, 0x0008, false, true, false, false, t)

	testOp(alu, SHL, 0x0073, 0x0000, false, 0x00E6, false, true, false, false, t)

	//should carry out
	testOp(alu, SHL, 0xAA00, 0x0001, false, 0x5400, false, true, true, false, t)

	//should carry in
	testOp(alu, SHL, 0x004A, 0x0001, true, 0x0095, false, true, false, false, t)

	//should do nothing with inputB
	testOp(alu, SHL, 0x0000, 0x0005, false, 0x0000, false, false, false, true, t)
}

func TestNOT(t *testing.T) {
	alu := NewALU(inputABus, inputBBus, outputBus, flagsBus)
	testOp(alu, NOT, 0x0000, 0x0000, false, 0xFFFF, true, false, false, false, t)
	testOp(alu, NOT, 0x0001, 0x0000, false, 0xFFFE, false, true, false, false, t)
	testOp(alu, NOT, 0x00AA, 0x0000, false, 0xFF55, false, true, false, false, t)
	testOp(alu, NOT, 0xFFFF, 0x0000, false, 0x0000, false, true, false, true, t)

	//should ignore  inputB
	testOp(alu, NOT, 0xFFFF, 0x0084, false, 0x0000, false, true, false, true, t)
	//should ignore  carry in bit
	testOp(alu, NOT, 0xB9B9, 0x0000, true, 0x4646, false, true, false, false, t)
}

func TestAND(t *testing.T) {
	alu := NewALU(inputABus, inputBBus, outputBus, flagsBus)
	testOp(alu, AND, 0x0000, 0x0000, false, 0x0000, true, false, false, true, t)
	testOp(alu, AND, 0x0001, 0x0001, false, 0x0001, true, false, false, false, t)
	testOp(alu, AND, 0x00FF, 0x00FF, false, 0x00FF, true, false, false, false, t)
	testOp(alu, AND, 0x00FF, 0x0069, false, 0x0069, false, true, false, false, t)
	testOp(alu, AND, 0x0069, 0x00FF, false, 0x0069, false, false, false, false, t)
	testOp(alu, AND, 0x008A, 0x0000, false, 0x0000, false, true, false, true, t)
	testOp(alu, AND, 0x0000, 0x008A, false, 0x0000, false, false, false, true, t)

	//should ignore  carry in bit
	testOp(alu, AND, 0x009A, 0x009A, true, 0x009A, true, false, false, false, t)

	j := uint16(0xFFFF)
	for i := uint16(1); i < 65535; i++ {
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
		testOp(alu, AND, i, j, false, i&j, IsEqual, isLarger, false, isZero, t)

		isLarger = false
		if j > i {
			isLarger = true
		}

		isZero = false
		if j&i == 0 {
			isZero = true
		}

		testOp(alu, AND, j, i, false, j&i, IsEqual, isLarger, false, isZero, t)
		j--
	}
}

func TestOR(t *testing.T) {
	alu := NewALU(inputABus, inputBBus, outputBus, flagsBus)
	testOp(alu, OR, 0x0000, 0x0000, false, 0x0000, true, false, false, true, t)
	testOp(alu, OR, 0x0001, 0x0001, false, 0x0001, true, false, false, false, t)
	testOp(alu, OR, 0x00FF, 0x00FF, false, 0x00FF, true, false, false, false, t)

	testOp(alu, OR, 0x00FF, 0x0069, false, 0x00FF, false, true, false, false, t)
	testOp(alu, OR, 0x0069, 0x00FF, false, 0x00FF, false, false, false, false, t)
	testOp(alu, OR, 0x008A, 0x0000, false, 0x008A, false, true, false, false, t)
	testOp(alu, OR, 0x0000, 0x008A, false, 0x008A, false, false, false, false, t)

	//should ignore  carry in bit
	testOp(alu, OR, 0x009A, 0x009A, true, 0x009A, true, false, false, false, t)

	j := uint16(0xFFFF)
	for i := uint16(1); i < 65535; i++ {
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
		testOp(alu, OR, i, j, false, i|j, IsEqual, isLarger, false, isZero, t)

		isLarger = false
		if j > i {
			isLarger = true
		}

		isZero = false
		if j|i == 0 {
			isZero = true
		}

		testOp(alu, OR, j, i, false, j|i, IsEqual, isLarger, false, isZero, t)
		j--
	}
}

func TestXOR(t *testing.T) {
	alu := NewALU(inputABus, inputBBus, outputBus, flagsBus)
	testOp(alu, XOR, 0x0000, 0x0000, false, 0x0000, true, false, false, true, t)
	testOp(alu, XOR, 0x0001, 0x0001, false, 0x0000, true, false, false, true, t)
	testOp(alu, XOR, 0x00FF, 0x00FF, false, 0x0000, true, false, false, true, t)

	testOp(alu, XOR, 0x00FF, 0x0069, false, 0x0096, false, true, false, false, t)
	testOp(alu, XOR, 0x0069, 0x00FF, false, 0x0096, false, false, false, false, t)
	testOp(alu, XOR, 0x008A, 0x0000, false, 0x008A, false, true, false, false, t)
	testOp(alu, XOR, 0x0000, 0x008A, false, 0x008A, false, false, false, false, t)

	testOp(alu, XOR, 0x009A, 0x009A, true, 0x0000, true, false, false, true, t)

	j := uint16(0xFFFF)
	for i := uint16(1); i < 65535; i++ {
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
		testOp(alu, XOR, i, j, false, i^j, IsEqual, isLarger, false, isZero, t)

		isLarger = false
		if j > i {
			isLarger = true
		}

		isZero = false
		if j^i == 0 {
			isZero = true
		}

		testOp(alu, XOR, j, i, false, j^i, IsEqual, isLarger, false, isZero, t)
		j--
	}
}

func TestCMP(t *testing.T) {
	alu := NewALU(inputABus, inputBBus, outputBus, flagsBus)
	testOp(alu, CMP, 0x0000, 0x0000, false, 0x0000, true, false, false, false, t)
	testOp(alu, CMP, 0x0001, 0x0001, false, 0x0000, true, false, false, false, t)
	testOp(alu, CMP, 0x00FF, 0x00FF, false, 0x0000, true, false, false, false, t)

	j := uint16(0xFFFF)
	for i := uint16(1); i < 65535; i++ {
		var IsEqual bool = false
		if i == j {
			IsEqual = true
		}

		var isLarger bool = false
		if i > j {
			isLarger = true
		}

		testOp(alu, CMP, i, j, false, 0x0000, IsEqual, isLarger, false, false, t)

		isLarger = false
		if j > i {
			isLarger = true
		}

		testOp(alu, CMP, j, i, false, 0x0000, IsEqual, isLarger, false, false, t)
		j--
	}

}

func testOp(alu *ALU, op uint16, inputA, inputB uint16, CarryIn bool, expectedOutput uint16, expectedEqual, expectedIsLarger, expectedCarry, expectedZero bool, t *testing.T) {
	inputABus.SetValue(inputA)
	inputBBus.SetValue(inputB)
	setOp(alu, op)
	alu.CarryIn.Update(CarryIn)
	alu.Update()

	output := getValueOfBus(outputBus)
	if output != expectedOutput {
		t.Logf("Expected output of 0x%X but got 0x%X", expectedOutput, output)
		t.FailNow()
	}

	if carryFlagSet := alu.flagsOutputBus.GetOutputWire(0); carryFlagSet != expectedCarry {
		t.Logf("Expected is carry out flag to be %v but got %v", expectedCarry, carryFlagSet)
		t.FailNow()
	}
	if isLargerFlagSet := alu.flagsOutputBus.GetOutputWire(1); isLargerFlagSet != expectedIsLarger {
		t.Logf("Expected is larger flag to be %v but got %v", expectedIsLarger, isLargerFlagSet)
		t.FailNow()
	}
	if equalFlagSet := alu.flagsOutputBus.GetOutputWire(2); equalFlagSet != expectedEqual {
		t.Logf("Expected equal flag to be %v but got %v", expectedEqual, equalFlagSet)
		t.FailNow()
	}
	if zeroFlagSet := alu.flagsOutputBus.GetOutputWire(3); zeroFlagSet != expectedZero {
		t.Logf("Expected zero flag to be %v but got %v", expectedZero, zeroFlagSet)
		t.FailNow()
	}
}

func setOp(a *ALU, value uint16) {
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

func getValueOfBus(bus *components.Bus) uint16 {
	var x uint16 = 0
	var result uint16
	for i := BUS_WIDTH - 1; i >= 0; i-- {
		if bus.GetOutputWire(i) {
			result = result | (1 << x)
		} else {
			result = result & ^(1 << x)
		}
		x++
	}
	return result
}
