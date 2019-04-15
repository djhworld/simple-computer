package circuit

import (
	"testing"
)

func TestNANDGate(t *testing.T) {
	combinations := [][]bool{
		[]bool{false, false, true},
		[]bool{true, false, true},
		[]bool{false, true, true},
		[]bool{true, true, false},
	}

	for _, combination := range combinations {
		gate1 := NewNANDGate()
		gate1.Update(combination[0], combination[1])

		if gate1.output.Get() != combination[2] {
			t.Fail()
		}
	}
}

func TestANDGate(t *testing.T) {
	combinations := [][]bool{
		[]bool{false, false, false},
		[]bool{true, false, false},
		[]bool{false, true, false},
		[]bool{true, true, true},
	}

	for _, combination := range combinations {
		gate1 := NewANDGate()
		gate1.Update(combination[0], combination[1])

		if gate1.output.Get() != combination[2] {
			t.Fail()
		}
	}
}

func TestORGate(t *testing.T) {
	combinations := [][]bool{
		[]bool{false, false, false},
		[]bool{true, false, true},
		[]bool{false, true, true},
		[]bool{true, true, true},
	}

	for _, combination := range combinations {
		gate1 := NewORGate()
		gate1.Update(combination[0], combination[1])

		if gate1.output.Get() != combination[2] {
			t.Fail()
		}
	}
}

func TestNOTGate(t *testing.T) {
	combinations := [][]bool{
		[]bool{false, true},
		[]bool{true, false},
	}

	for _, combination := range combinations {
		gate1 := NewNOTGate()
		gate1.Update(combination[0])

		if gate1.output.Get() != combination[1] {
			t.Fail()
		}
	}
}

func TestXORGate(t *testing.T) {
	combinations := [][]bool{
		[]bool{false, false, false},
		[]bool{true, false, true},
		[]bool{false, true, true},
		[]bool{true, true, false},
	}

	for _, combination := range combinations {
		gate1 := NewXORGate()
		gate1.Update(combination[0], combination[1])

		if gate1.output.Get() != combination[2] {
			t.Fail()
		}
	}
}

func TestNORGate(t *testing.T) {
	combinations := [][]bool{
		[]bool{false, false, true},
		[]bool{true, false, false},
		[]bool{false, true, false},
		[]bool{true, true, false},
	}

	for _, combination := range combinations {
		gate1 := NewNORGate()
		gate1.Update(combination[0], combination[1])

		if gate1.output.Get() != combination[2] {
			t.Fail()
		}
	}
}
