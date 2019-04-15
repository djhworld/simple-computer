package components

import (
	"fmt"
	"testing"
)

func TestStepper(t *testing.T) {
	testStepperStateAfter(0, -1, t)
	testStepperStateAfter(1, 1, t)
	testStepperStateAfter(2, 2, t)
	testStepperStateAfter(3, 3, t)
	testStepperStateAfter(4, 4, t)
	testStepperStateAfter(5, 5, t)
	testStepperStateAfter(6, 6, t)
}
func testStepperStateAfter(cycles, expectedOutput int, t *testing.T) {
	stepper := NewStepper()
	for i := 0; i < cycles; i++ {
		stepper.Update(true)
		stepper.Update(false)
	}

	output := getOutput(stepper)
	if expectedOutput != output {
		t.Logf("expected %X but got %X", expectedOutput, output)
		t.FailNow()
	}
}

func getOutput(stepper *Stepper) int {
	fmt.Println(stepper)
	for i := 0; i < 7; i++ {
		if stepper.GetOutputWire(i) {
			return i + 1
		}
	}
	return -1
}
