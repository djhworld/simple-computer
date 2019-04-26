package components

import (
	"fmt"
	"testing"
)

func TestAdder(t *testing.T) {
	testAdderReturnsCorrectResult(0, 0, false, 0, false, t)
	testAdderReturnsCorrectResult(1, 0, false, 1, false, t)
	testAdderReturnsCorrectResult(0, 1, false, 1, false, t)
	testAdderReturnsCorrectResult(1, 1, false, 2, false, t)
	testAdderReturnsCorrectResult(2, 1, false, 3, false, t)
	testAdderReturnsCorrectResult(1, 2, false, 3, false, t)
	testAdderReturnsCorrectResult(4, 9, false, 13, false, t)
	testAdderReturnsCorrectResult(64, 64, false, 128, false, t)
	testAdderReturnsCorrectResult(127, 128, false, 255, false, t)

	//carry in
	testAdderReturnsCorrectResult(0, 0, true, 1, false, t)
	testAdderReturnsCorrectResult(20, 2, true, 23, false, t)
	testAdderReturnsCorrectResult(0xFFFF, 0, true, 0, true, t)

	testAdderReturnsCorrectResult(32768, 32768, false, 0, true, t)
	testAdderReturnsCorrectResult(32769, 32768, false, 1, true, t)
	testAdderReturnsCorrectResult(65535, 2, false, 1, true, t)
}

func testAdderReturnsCorrectResult(inputA int, inputB int, carryIn bool, expectedResult int, expectedCarry bool, t *testing.T) {
	a := NewAdder()
	setWireOnComponent32(a, inputA, inputB)

	a.Update(carryIn)

	result := getValueOfOutput(a, BUS_WIDTH)

	if result != expectedResult {
		t.Log(fmt.Sprintf("Expected output of %d but got %d", expectedResult, result))
		t.FailNow()
	}

	if a.Carry() != expectedCarry {
		t.Log(fmt.Sprintf("Expected carry of %v but got %v", expectedCarry, a.Carry()))
		t.FailNow()
	}
}

func TestAdd2(t *testing.T) {
	a := NewAdd2()

	a.Update(true, true, false)

	if a.sumOut.Get() != false {
		t.FailNow()
	}

	if a.carryOut.Get() != true {
		t.FailNow()
	}

	a.Update(true, false, false)

	if a.sumOut.Get() != true {
		t.FailNow()
	}

	if a.carryOut.Get() != false {
		t.FailNow()
	}

	a.Update(false, false, true)

	if a.sumOut.Get() != true {
		t.FailNow()
	}

	if a.carryOut.Get() != false {
		t.FailNow()
	}

	a.Update(true, true, true)

	if a.sumOut.Get() != true {
		t.FailNow()
	}

	if a.carryOut.Get() != true {
		t.FailNow()
	}

	a.Update(true, false, true)

	if a.sumOut.Get() != false {
		t.FailNow()
	}

	if a.carryOut.Get() != true {
		t.FailNow()
	}
}
