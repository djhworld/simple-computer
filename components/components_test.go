package components

import (
	"fmt"
	"testing"

	"github.com/djhworld/simple-computer/circuit"
)

type DummyComponent struct {
	wires [BUS_WIDTH]circuit.Wire
	next  Component
}

func (d *DummyComponent) ConnectOutput(b Component) {
	d.next = b
}

func (d *DummyComponent) SetInputWire(index int, value bool) {
	d.wires[index].Update(value)
}

func (d *DummyComponent) GetOutputWire(index int) bool {
	return d.wires[index].Get()
}

func (d *DummyComponent) Update() {
	for i, w := range d.wires {
		d.next.SetInputWire(i, w.Get())
	}
}

func TestEnabler(t *testing.T) {
	d := new(DummyComponent)

	for i, _ := range d.wires {
		d.SetInputWire(i, true)
	}

	enabler := NewEnabler()
	d.ConnectOutput(enabler)
	d.Update()
	enabler.Update(false)

	for _, w := range enabler.outputs {
		if w.Get() {
			t.FailNow()
		}
	}
}

func TestEnablerWithEnableOn(t *testing.T) {
	d := new(DummyComponent)
	d.SetInputWire(0, false)
	d.SetInputWire(1, true)
	d.SetInputWire(2, false)
	d.SetInputWire(3, true)
	d.SetInputWire(4, false)
	d.SetInputWire(5, true)
	d.SetInputWire(6, false)
	d.SetInputWire(7, true)
	d.SetInputWire(8, false)
	d.SetInputWire(9, true)
	d.SetInputWire(10, false)
	d.SetInputWire(11, true)
	d.SetInputWire(12, false)
	d.SetInputWire(13, true)
	d.SetInputWire(14, false)
	d.SetInputWire(15, true)

	enabler := NewEnabler()
	d.ConnectOutput(enabler)
	d.Update()
	enabler.Update(true)

	results := [BUS_WIDTH]bool{}
	for i, w := range enabler.outputs {
		results[i] = w.Get()
	}

	if results != [BUS_WIDTH]bool{false, true, false, true, false, true, false, true, false, true, false, true, false, true, false, true} {
		t.FailNow()
	}
}

func TestANDGate3(t *testing.T) {
	combinations := [][]bool{
		[]bool{false, false, false, false},
		[]bool{false, false, true, false},
		[]bool{false, true, false, false},
		[]bool{false, true, true, false},
		[]bool{true, false, false, false},
		[]bool{true, false, true, false},
		[]bool{true, true, false, false},
		[]bool{true, true, true, true},
	}

	for _, combination := range combinations {
		wireA := circuit.NewWire("A", combination[0])
		wireB := circuit.NewWire("B", combination[1])
		wireC := circuit.NewWire("C", combination[2])

		gate1 := NewANDGate3()
		gate1.Update(wireA.Get(), wireB.Get(), wireC.Get())

		if gate1.output.Get() != combination[3] {
			t.Fail()
		}
	}
}

func TestORGate3(t *testing.T) {
	combinations := [][]bool{
		[]bool{false, false, false, false},
		[]bool{false, false, true, true},
		[]bool{false, true, false, true},
		[]bool{false, true, true, true},
		[]bool{true, false, false, true},
		[]bool{true, false, true, true},
		[]bool{true, true, false, true},
		[]bool{true, true, true, true},
	}

	for _, combination := range combinations {
		wireA := circuit.NewWire("A", combination[0])
		wireB := circuit.NewWire("B", combination[1])
		wireC := circuit.NewWire("C", combination[2])

		gate1 := NewORGate3()
		gate1.Update(wireA.Get(), wireB.Get(), wireC.Get())

		if gate1.output.Get() != combination[3] {
			t.Fail()
		}
	}
}

func TestANDGate4(t *testing.T) {
	combinations := [][]bool{
		[]bool{false, false, false, false, false},
		[]bool{true, false, false, false, false},
		[]bool{false, true, false, false, false},
		[]bool{false, false, true, false, false},
		[]bool{false, false, false, true, false},
		[]bool{true, true, false, false, false},
		[]bool{false, true, true, false, false},
		[]bool{false, false, true, true, false},
		[]bool{true, true, true, false, false},
		[]bool{false, true, true, true, false},
		[]bool{true, false, false, true, false},
		[]bool{true, false, true, true, false},
		[]bool{true, true, false, true, false},
		[]bool{true, false, true, false, false},
		[]bool{false, true, false, true, false},
		[]bool{true, true, true, true, true},
	}

	for _, combination := range combinations {
		wireA := circuit.NewWire("A", combination[0])
		wireB := circuit.NewWire("B", combination[1])
		wireC := circuit.NewWire("C", combination[2])
		wireD := circuit.NewWire("D", combination[3])

		gate1 := NewANDGate4()
		gate1.Update(wireA.Get(), wireB.Get(), wireC.Get(), wireD.Get())

		if gate1.output.Get() != combination[4] {
			t.Fail()
		}
	}
}

func TestORGate4(t *testing.T) {
	combinations := [][]bool{
		[]bool{false, false, false, false, false},
		[]bool{true, false, false, false, true},
		[]bool{false, true, false, false, true},
		[]bool{false, false, true, false, true},
		[]bool{false, false, false, true, true},
		[]bool{true, true, false, false, true},
		[]bool{false, true, true, false, true},
		[]bool{false, false, true, true, true},
		[]bool{true, true, true, false, true},
		[]bool{false, true, true, true, true},
		[]bool{true, false, false, true, true},
		[]bool{true, false, true, true, true},
		[]bool{true, true, false, true, true},
		[]bool{true, false, true, false, true},
		[]bool{false, true, false, true, true},
		[]bool{true, true, true, true, true},
	}

	for _, combination := range combinations {
		wireA := circuit.NewWire("A", combination[0])
		wireB := circuit.NewWire("B", combination[1])
		wireC := circuit.NewWire("C", combination[2])
		wireD := circuit.NewWire("D", combination[3])

		gate1 := NewORGate4()
		gate1.Update(wireA.Get(), wireB.Get(), wireC.Get(), wireD.Get())

		if gate1.output.Get() != combination[4] {
			t.Fail()
		}
	}
}

func TestORGate5(t *testing.T) {
	combinations := [][]bool{
		[]bool{false, false, false, false, false, false},
		[]bool{false, false, false, false, true, true},
		[]bool{false, false, false, true, false, true},
		[]bool{false, false, false, true, true, true},
		[]bool{false, false, true, false, false, true},
		[]bool{false, false, true, false, true, true},
		[]bool{false, false, true, true, false, true},
		[]bool{false, false, true, true, true, true},
		[]bool{false, true, false, false, false, true},
		[]bool{false, true, false, false, true, true},
		[]bool{false, true, false, true, false, true},
		[]bool{false, true, false, true, true, true},
		[]bool{false, true, true, false, false, true},
		[]bool{false, true, true, false, true, true},
		[]bool{false, true, true, true, false, true},
		[]bool{false, true, true, true, true, true},
		[]bool{true, false, false, false, false, true},
		[]bool{true, false, false, false, true, true},
		[]bool{true, false, false, true, false, true},
		[]bool{true, false, false, true, true, true},
		[]bool{true, false, true, false, false, true},
		[]bool{true, false, true, false, true, true},
		[]bool{true, false, true, true, false, true},
		[]bool{true, false, true, true, true, true},
		[]bool{true, true, false, false, false, true},
		[]bool{true, true, false, false, true, true},
		[]bool{true, true, false, true, false, true},
		[]bool{true, true, false, true, true, true},
		[]bool{true, true, true, false, false, true},
		[]bool{true, true, true, false, true, true},
		[]bool{true, true, true, true, false, true},
		[]bool{true, true, true, true, true, true},
	}

	for _, combination := range combinations {
		wireA := circuit.NewWire("A", combination[0])
		wireB := circuit.NewWire("B", combination[1])
		wireC := circuit.NewWire("C", combination[2])
		wireD := circuit.NewWire("D", combination[3])
		wireE := circuit.NewWire("E", combination[4])

		gate1 := NewORGate5()
		gate1.Update(wireA.Get(), wireB.Get(), wireC.Get(), wireD.Get(), wireE.Get())

		if gate1.output.Get() != combination[5] {
			t.Fail()
		}
	}
}

func TestORGate6(t *testing.T) {
	combinations := [][]bool{
		[]bool{false, false, false, false, false, false, false},
		[]bool{false, false, false, false, false, true, true},
		[]bool{false, false, false, false, true, false, true},
		[]bool{false, false, false, false, true, true, true},
		[]bool{false, false, false, true, false, false, true},
		[]bool{false, false, false, true, false, true, true},
		[]bool{false, false, false, true, true, false, true},
		[]bool{false, false, false, true, true, true, true},
		[]bool{false, false, true, false, false, false, true},
		[]bool{false, false, true, false, false, true, true},
		[]bool{false, false, true, false, true, false, true},
		[]bool{false, false, true, false, true, true, true},
		[]bool{false, false, true, true, false, false, true},
		[]bool{false, false, true, true, false, true, true},
		[]bool{false, false, true, true, true, false, true},
		[]bool{false, false, true, true, true, true, true},
		[]bool{false, true, false, false, false, false, true},
		[]bool{false, true, false, false, false, true, true},
		[]bool{false, true, false, false, true, false, true},
		[]bool{false, true, false, false, true, true, true},
		[]bool{false, true, false, true, false, false, true},
		[]bool{false, true, false, true, false, true, true},
		[]bool{false, true, false, true, true, false, true},
		[]bool{false, true, false, true, true, true, true},
		[]bool{false, true, true, false, false, false, true},
		[]bool{false, true, true, false, false, true, true},
		[]bool{false, true, true, false, true, false, true},
		[]bool{false, true, true, false, true, true, true},
		[]bool{false, true, true, true, false, false, true},
		[]bool{false, true, true, true, false, true, true},
		[]bool{false, true, true, true, true, false, true},
		[]bool{false, true, true, true, true, true, true},
		[]bool{true, false, false, false, false, false, true},
		[]bool{true, false, false, false, false, true, true},
		[]bool{true, false, false, false, true, false, true},
		[]bool{true, false, false, false, true, true, true},
		[]bool{true, false, false, true, false, false, true},
		[]bool{true, false, false, true, false, true, true},
		[]bool{true, false, false, true, true, false, true},
		[]bool{true, false, false, true, true, true, true},
		[]bool{true, false, true, false, false, false, true},
		[]bool{true, false, true, false, false, true, true},
		[]bool{true, false, true, false, true, false, true},
		[]bool{true, false, true, false, true, true, true},
		[]bool{true, false, true, true, false, false, true},
		[]bool{true, false, true, true, false, true, true},
		[]bool{true, false, true, true, true, false, true},
		[]bool{true, false, true, true, true, true, true},
		[]bool{true, true, false, false, false, false, true},
		[]bool{true, true, false, false, false, true, true},
		[]bool{true, true, false, false, true, false, true},
		[]bool{true, true, false, false, true, true, true},
		[]bool{true, true, false, true, false, false, true},
		[]bool{true, true, false, true, false, true, true},
		[]bool{true, true, false, true, true, false, true},
		[]bool{true, true, false, true, true, true, true},
		[]bool{true, true, true, false, false, false, true},
		[]bool{true, true, true, false, false, true, true},
		[]bool{true, true, true, false, true, false, true},
		[]bool{true, true, true, false, true, true, true},
		[]bool{true, true, true, true, false, false, true},
		[]bool{true, true, true, true, false, true, true},
		[]bool{true, true, true, true, true, false, true},
		[]bool{true, true, true, true, true, true, true},
	}

	for _, combination := range combinations {
		wireA := circuit.NewWire("A", combination[0])
		wireB := circuit.NewWire("B", combination[1])
		wireC := circuit.NewWire("C", combination[2])
		wireD := circuit.NewWire("D", combination[3])
		wireE := circuit.NewWire("E", combination[4])
		wireF := circuit.NewWire("F", combination[5])

		gate1 := NewORGate6()
		gate1.Update(wireA.Get(), wireB.Get(), wireC.Get(), wireD.Get(), wireE.Get(), wireF.Get())

		if gate1.output.Get() != combination[6] {
			t.Fail()
		}
	}
}

func TestLeftShifter(t *testing.T) {
	//left shifter always applies the same output
	for i := 1; i < 32767; i *= 2 {
		testLeftShifter(i, false, i*2, false, t)
	}

	testLeftShifter(0, false, 0, false, t)
	testLeftShifter(0x8000, false, 0, true, t)
	testLeftShifter(0xEEEF, false, 0xDDDE, true, t)
	testLeftShifter(0xFFFF, false, 0xFFFE, true, t)

	// shift in?
	testLeftShifter(0x0000, true, 0x0001, false, t)
	testLeftShifter(0x8000, true, 0x0001, true, t)
}

func testLeftShifter(input int, shiftIn bool, expectedOutput int, expectedShiftOut bool, t *testing.T) {
	l := NewLeftShifter()
	setWireOnComponent16(l, input)
	l.Update(shiftIn)
	if output := getValueOfOutput(l, BUS_WIDTH); output != expectedOutput {
		t.Logf("expected 0x%X got 0x%X", expectedOutput, output)
		t.FailNow()
	}
	if l.shiftOut.Get() != expectedShiftOut {
		t.FailNow()
	}
}

func TestRightShifter(t *testing.T) {
	//right shifter always applies the same output
	for i := 32768; i > 1; i /= 2 {
		testRightShifter(i, false, i/2, false, t)
	}

	testRightShifter(0, false, 0, false, t)
	testRightShifter(0x0001, false, 0, true, t)
	testRightShifter(0x8000, false, 0x4000, false, t)
	testRightShifter(0xEEEF, false, 0x7777, true, t)
	testRightShifter(0xFFFF, false, 0x7FFF, true, t)

	// shift in?
	testRightShifter(0x0000, true, 0x8000, false, t)
	testRightShifter(0x8000, true, 0xC000, false, t)
	testRightShifter(0x4AAA, true, 0xA555, false, t)
}

func testRightShifter(input int, shiftIn bool, expectedOutput int, expectedShiftOut bool, t *testing.T) {
	r := NewRightShifter()
	setWireOnComponent16(r, input)
	r.Update(shiftIn)
	output := getValueOfOutput(r, BUS_WIDTH)
	if output != expectedOutput {
		t.Logf("expected 0x%X got 0x%X", expectedOutput, output)
		t.FailNow()
	}
	if r.shiftOut.Get() != expectedShiftOut {
		t.FailNow()
	}
}

func TestNOTer(t *testing.T) {
	n := NewNOTer()

	n.SetInputWire(0, false)
	n.SetInputWire(1, true)
	n.SetInputWire(2, false)
	n.SetInputWire(3, true)
	n.SetInputWire(4, false)
	n.SetInputWire(5, true)
	n.SetInputWire(6, false)
	n.SetInputWire(7, true)
	n.SetInputWire(8, false)
	n.SetInputWire(9, true)
	n.SetInputWire(10, false)
	n.SetInputWire(11, true)
	n.SetInputWire(12, false)
	n.SetInputWire(13, true)
	n.SetInputWire(14, false)
	n.SetInputWire(15, true)

	n.Update()

	results := [BUS_WIDTH]bool{}
	for i, w := range n.outputs {
		results[i] = w.Get()
	}

	if results != [BUS_WIDTH]bool{true, false, true, false, true, false, true, false, true, false, true, false, true, false, true, false} {
		t.FailNow()
	}

	// check all bits invert
	for i, _ := range n.inputs {
		n.SetInputWire(i, false)
	}

	n.Update()

	for _, w := range n.outputs {
		if !w.Get() {
			t.FailNow()
		}
	}
}

func TestANDer(t *testing.T) {
	a := NewANDer()

	a.SetInputWire(0, true)
	a.SetInputWire(1, true)
	a.SetInputWire(2, true)
	a.SetInputWire(3, true)
	a.SetInputWire(4, true)
	a.SetInputWire(5, true)
	a.SetInputWire(6, true)
	a.SetInputWire(7, true)
	a.SetInputWire(8, true)
	a.SetInputWire(9, true)
	a.SetInputWire(10, true)
	a.SetInputWire(11, true)
	a.SetInputWire(12, true)
	a.SetInputWire(13, true)
	a.SetInputWire(14, true)
	a.SetInputWire(15, true)

	a.SetInputWire(16, false)
	a.SetInputWire(17, true)
	a.SetInputWire(18, false)
	a.SetInputWire(19, true)
	a.SetInputWire(20, false)
	a.SetInputWire(21, true)
	a.SetInputWire(22, false)
	a.SetInputWire(23, true)
	a.SetInputWire(24, false)
	a.SetInputWire(25, true)
	a.SetInputWire(26, false)
	a.SetInputWire(27, true)
	a.SetInputWire(28, false)
	a.SetInputWire(29, true)
	a.SetInputWire(30, false)
	a.SetInputWire(31, true)

	a.Update()

	results := [BUS_WIDTH]bool{}
	for i, w := range a.outputs {
		results[i] = w.Get()
	}

	if results != [BUS_WIDTH]bool{false, true, false, true, false, true, false, true, false, true, false, true, false, true, false, true} {
		t.FailNow()
	}

	// check all bits invert
	for i, _ := range a.inputs {
		a.SetInputWire(i, true)
	}

	a.Update()

	for _, w := range a.outputs {
		if !w.Get() {
			t.FailNow()
		}
	}
}

func TestORer(t *testing.T) {
	o := NewORer()

	o.SetInputWire(0, false)
	o.SetInputWire(1, true)
	o.SetInputWire(2, true)
	o.SetInputWire(3, true)
	o.SetInputWire(4, true)
	o.SetInputWire(5, true)
	o.SetInputWire(6, true)
	o.SetInputWire(7, true)
	o.SetInputWire(8, true)
	o.SetInputWire(9, true)
	o.SetInputWire(10, true)
	o.SetInputWire(11, true)
	o.SetInputWire(12, true)
	o.SetInputWire(13, true)
	o.SetInputWire(14, true)
	o.SetInputWire(15, true)

	o.SetInputWire(16, false)
	o.SetInputWire(17, true)
	o.SetInputWire(18, true)
	o.SetInputWire(19, true)
	o.SetInputWire(20, true)
	o.SetInputWire(21, true)
	o.SetInputWire(22, true)
	o.SetInputWire(23, true)
	o.SetInputWire(24, true)
	o.SetInputWire(25, true)
	o.SetInputWire(26, true)
	o.SetInputWire(27, true)
	o.SetInputWire(28, true)
	o.SetInputWire(29, true)
	o.SetInputWire(30, true)
	o.SetInputWire(31, false)

	o.Update()

	results := [BUS_WIDTH]bool{}
	for i, w := range o.outputs {
		results[i] = w.Get()
	}

	if results != [BUS_WIDTH]bool{false, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true} {
		t.FailNow()
	}

	// check all bits  don't change
	for i, _ := range o.inputs {
		o.SetInputWire(i, false)
	}

	o.Update()

	for _, w := range o.outputs {
		if w.Get() {
			t.FailNow()
		}
	}
}

func TestXORer(t *testing.T) {
	o := NewXORer()

	o.SetInputWire(0, true)
	o.SetInputWire(1, true)
	o.SetInputWire(2, false)
	o.SetInputWire(3, false)
	o.SetInputWire(4, true)
	o.SetInputWire(5, true)
	o.SetInputWire(6, false)
	o.SetInputWire(7, false)
	o.SetInputWire(8, true)
	o.SetInputWire(9, true)
	o.SetInputWire(10, false)
	o.SetInputWire(11, false)
	o.SetInputWire(12, true)
	o.SetInputWire(13, true)
	o.SetInputWire(14, false)
	o.SetInputWire(15, false)

	o.SetInputWire(16, true)
	o.SetInputWire(17, false)
	o.SetInputWire(18, true)
	o.SetInputWire(19, false)
	o.SetInputWire(20, true)
	o.SetInputWire(21, false)
	o.SetInputWire(22, true)
	o.SetInputWire(23, false)
	o.SetInputWire(24, true)
	o.SetInputWire(25, false)
	o.SetInputWire(26, true)
	o.SetInputWire(27, false)
	o.SetInputWire(28, true)
	o.SetInputWire(29, false)
	o.SetInputWire(30, true)
	o.SetInputWire(31, false)

	o.Update()

	results := [BUS_WIDTH]bool{}
	for i, w := range o.outputs {
		results[i] = w.Get()
	}

	if results != [BUS_WIDTH]bool{false, true, true, false, false, true, true, false, false, true, true, false, false, true, true, false} {
		t.FailNow()
	}

	// check all bits  don't change
	for i, _ := range o.inputs {
		o.SetInputWire(i, false)
	}

	o.Update()

	for _, w := range o.outputs {
		if w.Get() {
			t.FailNow()
		}
	}
}

func setWireOnComponent32(b Component, inputA int, inputB int) {
	var x uint16 = 0
	for i := 16 - 1; i >= 0; i-- {
		r := (inputA & (1 << x))
		if r != 0 {
			b.SetInputWire(i, true)
		} else {
			b.SetInputWire(i, false)
		}
		x++
	}

	x = 0
	for i := 32 - 1; i >= 16; i-- {
		r := (inputB & (1 << x))
		if r != 0 {
			b.SetInputWire(i, true)
		} else {
			b.SetInputWire(i, false)
		}
		x++
	}
}

func setWireOnComponent16(b Component, input int) {
	var x uint16 = 0
	for i := 15; i >= 0; i-- {
		r := (input & (1 << x))
		if r != 0 {
			b.SetInputWire(i, true)
		} else {
			b.SetInputWire(i, false)
		}
		x++
	}
}

func getValueOfOutput(b Component, outputBits int) int {
	var x int = 0
	var result int
	for i := (outputBits - 1); i >= 0; i-- {
		if b.GetOutputWire(i) {
			result = result | (1 << uint16(x))
		} else {
			result = result & ^(1 << uint16(x))
		}
		x++
	}
	return result
}

func TestCompare2(t *testing.T) {
	c := NewCompare2()

	c.Update(true, true, true, false)
}

func TestComparator(t *testing.T) {
	testComparatorReturnsCorrectResult(0, 0, true, false, t)
	testComparatorReturnsCorrectResult(1, 0, false, true, t)
	testComparatorReturnsCorrectResult(0, 1, false, false, t)
	testComparatorReturnsCorrectResult(50, 50, true, false, t)
	testComparatorReturnsCorrectResult(65535, 0, false, true, t)
	testComparatorReturnsCorrectResult(0, 65535, false, false, t)
	testComparatorReturnsCorrectResult(115, 2, false, true, t)
	testComparatorReturnsCorrectResult(4, 1, false, true, t)
}

func testComparatorReturnsCorrectResult(inputA int, inputB int, expectedIsEqual bool, expectedIsLarger bool, t *testing.T) {
	c := NewComparator()
	setWireOnComponent32(c, inputA, inputB)

	c.Update()

	if c.Equal() != expectedIsEqual {
		t.Log(fmt.Sprintf("Expected equals to be of %v but got %v", expectedIsEqual, c.Equal()))
		t.FailNow()
	}

	if c.Larger() != expectedIsLarger {
		t.Log(fmt.Sprintf("Expected is larger to be of %v but got %v", expectedIsLarger, c.Larger()))
		t.FailNow()
	}
}

func TestIsZero(t *testing.T) {
	z := NewIsZero()
	for i := 0; i < BUS_WIDTH; i++ {
		z.SetInputWire(i, false)
	}

	z.Update()

	if z.output.Get() != true {
		t.FailNow()
	}

	z = NewIsZero()
	for i := 0; i < BUS_WIDTH; i++ {
		z.SetInputWire(i, true)
	}

	z.Update()

	if z.output.Get() != false {
		t.FailNow()
	}

	for i := 0; i < BUS_WIDTH; i++ {
		z := NewIsZero()

		z.SetInputWire(i, true)

		z.Update()

		result := z.output.Get()

		if result != false {
			t.FailNow()
		}
	}

}

func TestBusOne(t *testing.T) {
	inputBus := NewBus(BUS_WIDTH)
	outputBus := NewBus(BUS_WIDTH)
	b := NewBusOne(inputBus, outputBus)
	for i := 0; i < 65536; i++ {
		testBusOneReturnsCorrectResult(b, inputBus, outputBus, i, false, i, t)
		//should always return 1 regardless of input
		testBusOneReturnsCorrectResult(b, inputBus, outputBus, i, true, 1, t)
	}
}

func testBusOneReturnsCorrectResult(b *BusOne, inputBus, outputBus *Bus, input int, enableBus1 bool, expectedOutput int, t *testing.T) {
	setWireOnComponent16(inputBus, input)

	if enableBus1 {
		b.Enable()
	} else {
		b.Disable()
	}

	b.Update()
	output := getValueOfOutput(outputBus, BUS_WIDTH)

	if output != expectedOutput {
		t.Log(fmt.Sprintf("Expected output to be to be of 0x%X but got 0x%X", expectedOutput, output))
		t.FailNow()
	}
}
