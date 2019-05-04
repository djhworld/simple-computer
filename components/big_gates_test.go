package components

import (
	"testing"

	"github.com/djhworld/simple-computer/circuit"
)

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

func TestANDGate5(t *testing.T) {
	combinations := [][]bool{
		[]bool{false, false, false, false, false, false},
		[]bool{false, false, false, false, true, false},
		[]bool{false, false, false, true, false, false},
		[]bool{false, false, false, true, true, false},
		[]bool{false, false, true, false, false, false},
		[]bool{false, false, true, false, true, false},
		[]bool{false, false, true, true, false, false},
		[]bool{false, false, true, true, true, false},
		[]bool{false, true, false, false, false, false},
		[]bool{false, true, false, false, true, false},
		[]bool{false, true, false, true, false, false},
		[]bool{false, true, false, true, true, false},
		[]bool{false, true, true, false, false, false},
		[]bool{false, true, true, false, true, false},
		[]bool{false, true, true, true, false, false},
		[]bool{false, true, true, true, true, false},
		[]bool{true, false, false, false, false, false},
		[]bool{true, false, false, false, true, false},
		[]bool{true, false, false, true, false, false},
		[]bool{true, false, false, true, true, false},
		[]bool{true, false, true, false, false, false},
		[]bool{true, false, true, false, true, false},
		[]bool{true, false, true, true, false, false},
		[]bool{true, false, true, true, true, false},
		[]bool{true, true, false, false, false, false},
		[]bool{true, true, false, false, true, false},
		[]bool{true, true, false, true, false, false},
		[]bool{true, true, false, true, true, false},
		[]bool{true, true, true, false, false, false},
		[]bool{true, true, true, false, true, false},
		[]bool{true, true, true, true, false, false},
		[]bool{true, true, true, true, true, true},
	}

	for _, combination := range combinations {
		wireA := circuit.NewWire("A", combination[0])
		wireB := circuit.NewWire("B", combination[1])
		wireC := circuit.NewWire("C", combination[2])
		wireD := circuit.NewWire("D", combination[3])
		wireE := circuit.NewWire("E", combination[4])

		gate1 := NewANDGate5()
		gate1.Update(wireA.Get(), wireB.Get(), wireC.Get(), wireD.Get(), wireE.Get())

		if gate1.output.Get() != combination[5] {
			t.Fail()
		}
	}
}

func TestANDGate8(t *testing.T) {
	combinations := [][]bool{
		[]bool{false, false, false, false, false, false, false, false, false},
		[]bool{true, false, false, false, false, false, false, false, false},
		[]bool{false, true, false, false, false, false, false, false, false},
		[]bool{true, true, false, false, false, false, false, false, false},
		[]bool{false, false, true, false, false, false, false, false, false},
		[]bool{true, false, true, false, false, false, false, false, false},
		[]bool{false, true, true, false, false, false, false, false, false},
		[]bool{true, true, true, false, false, false, false, false, false},
		[]bool{false, false, false, true, false, false, false, false, false},
		[]bool{true, false, false, true, false, false, false, false, false},
		[]bool{false, true, false, true, false, false, false, false, false},
		[]bool{true, true, false, true, false, false, false, false, false},
		[]bool{false, false, true, true, false, false, false, false, false},
		[]bool{true, false, true, true, false, false, false, false, false},
		[]bool{false, true, true, true, false, false, false, false, false},
		[]bool{true, true, true, true, false, false, false, false, false},
		[]bool{false, false, false, false, true, false, false, false, false},
		[]bool{true, false, false, false, true, false, false, false, false},
		[]bool{false, true, false, false, true, false, false, false, false},
		[]bool{true, true, false, false, true, false, false, false, false},
		[]bool{false, false, true, false, true, false, false, false, false},
		[]bool{true, false, true, false, true, false, false, false, false},
		[]bool{false, true, true, false, true, false, false, false, false},
		[]bool{true, true, true, false, true, false, false, false, false},
		[]bool{false, false, false, true, true, false, false, false, false},
		[]bool{true, false, false, true, true, false, false, false, false},
		[]bool{false, true, false, true, true, false, false, false, false},
		[]bool{true, true, false, true, true, false, false, false, false},
		[]bool{false, false, true, true, true, false, false, false, false},
		[]bool{true, false, true, true, true, false, false, false, false},
		[]bool{false, true, true, true, true, false, false, false, false},
		[]bool{true, true, true, true, true, false, false, false, false},
		[]bool{false, false, false, false, false, true, false, false, false},
		[]bool{true, false, false, false, false, true, false, false, false},
		[]bool{false, true, false, false, false, true, false, false, false},
		[]bool{true, true, false, false, false, true, false, false, false},
		[]bool{false, false, true, false, false, true, false, false, false},
		[]bool{true, false, true, false, false, true, false, false, false},
		[]bool{false, true, true, false, false, true, false, false, false},
		[]bool{true, true, true, false, false, true, false, false, false},
		[]bool{false, false, false, true, false, true, false, false, false},
		[]bool{true, false, false, true, false, true, false, false, false},
		[]bool{false, true, false, true, false, true, false, false, false},
		[]bool{true, true, false, true, false, true, false, false, false},
		[]bool{false, false, true, true, false, true, false, false, false},
		[]bool{true, false, true, true, false, true, false, false, false},
		[]bool{false, true, true, true, false, true, false, false, false},
		[]bool{true, true, true, true, false, true, false, false, false},
		[]bool{false, false, false, false, true, true, false, false, false},
		[]bool{true, false, false, false, true, true, false, false, false},
		[]bool{false, true, false, false, true, true, false, false, false},
		[]bool{true, true, false, false, true, true, false, false, false},
		[]bool{false, false, true, false, true, true, false, false, false},
		[]bool{true, false, true, false, true, true, false, false, false},
		[]bool{false, true, true, false, true, true, false, false, false},
		[]bool{true, true, true, false, true, true, false, false, false},
		[]bool{false, false, false, true, true, true, false, false, false},
		[]bool{true, false, false, true, true, true, false, false, false},
		[]bool{false, true, false, true, true, true, false, false, false},
		[]bool{true, true, false, true, true, true, false, false, false},
		[]bool{false, false, true, true, true, true, false, false, false},
		[]bool{true, false, true, true, true, true, false, false, false},
		[]bool{false, true, true, true, true, true, false, false, false},
		[]bool{true, true, true, true, true, true, false, false, false},
		[]bool{false, false, false, false, false, false, true, false, false},
		[]bool{true, false, false, false, false, false, true, false, false},
		[]bool{false, true, false, false, false, false, true, false, false},
		[]bool{true, true, false, false, false, false, true, false, false},
		[]bool{false, false, true, false, false, false, true, false, false},
		[]bool{true, false, true, false, false, false, true, false, false},
		[]bool{false, true, true, false, false, false, true, false, false},
		[]bool{true, true, true, false, false, false, true, false, false},
		[]bool{false, false, false, true, false, false, true, false, false},
		[]bool{true, false, false, true, false, false, true, false, false},
		[]bool{false, true, false, true, false, false, true, false, false},
		[]bool{true, true, false, true, false, false, true, false, false},
		[]bool{false, false, true, true, false, false, true, false, false},
		[]bool{true, false, true, true, false, false, true, false, false},
		[]bool{false, true, true, true, false, false, true, false, false},
		[]bool{true, true, true, true, false, false, true, false, false},
		[]bool{false, false, false, false, true, false, true, false, false},
		[]bool{true, false, false, false, true, false, true, false, false},
		[]bool{false, true, false, false, true, false, true, false, false},
		[]bool{true, true, false, false, true, false, true, false, false},
		[]bool{false, false, true, false, true, false, true, false, false},
		[]bool{true, false, true, false, true, false, true, false, false},
		[]bool{false, true, true, false, true, false, true, false, false},
		[]bool{true, true, true, false, true, false, true, false, false},
		[]bool{false, false, false, true, true, false, true, false, false},
		[]bool{true, false, false, true, true, false, true, false, false},
		[]bool{false, true, false, true, true, false, true, false, false},
		[]bool{true, true, false, true, true, false, true, false, false},
		[]bool{false, false, true, true, true, false, true, false, false},
		[]bool{true, false, true, true, true, false, true, false, false},
		[]bool{false, true, true, true, true, false, true, false, false},
		[]bool{true, true, true, true, true, false, true, false, false},
		[]bool{false, false, false, false, false, true, true, false, false},
		[]bool{true, false, false, false, false, true, true, false, false},
		[]bool{false, true, false, false, false, true, true, false, false},
		[]bool{true, true, false, false, false, true, true, false, false},
		[]bool{false, false, true, false, false, true, true, false, false},
		[]bool{true, false, true, false, false, true, true, false, false},
		[]bool{false, true, true, false, false, true, true, false, false},
		[]bool{true, true, true, false, false, true, true, false, false},
		[]bool{false, false, false, true, false, true, true, false, false},
		[]bool{true, false, false, true, false, true, true, false, false},
		[]bool{false, true, false, true, false, true, true, false, false},
		[]bool{true, true, false, true, false, true, true, false, false},
		[]bool{false, false, true, true, false, true, true, false, false},
		[]bool{true, false, true, true, false, true, true, false, false},
		[]bool{false, true, true, true, false, true, true, false, false},
		[]bool{true, true, true, true, false, true, true, false, false},
		[]bool{false, false, false, false, true, true, true, false, false},
		[]bool{true, false, false, false, true, true, true, false, false},
		[]bool{false, true, false, false, true, true, true, false, false},
		[]bool{true, true, false, false, true, true, true, false, false},
		[]bool{false, false, true, false, true, true, true, false, false},
		[]bool{true, false, true, false, true, true, true, false, false},
		[]bool{false, true, true, false, true, true, true, false, false},
		[]bool{true, true, true, false, true, true, true, false, false},
		[]bool{false, false, false, true, true, true, true, false, false},
		[]bool{true, false, false, true, true, true, true, false, false},
		[]bool{false, true, false, true, true, true, true, false, false},
		[]bool{true, true, false, true, true, true, true, false, false},
		[]bool{false, false, true, true, true, true, true, false, false},
		[]bool{true, false, true, true, true, true, true, false, false},
		[]bool{false, true, true, true, true, true, true, false, false},
		[]bool{true, true, true, true, true, true, true, false, false},
		[]bool{false, false, false, false, false, false, false, true, false},
		[]bool{true, false, false, false, false, false, false, true, false},
		[]bool{false, true, false, false, false, false, false, true, false},
		[]bool{true, true, false, false, false, false, false, true, false},
		[]bool{false, false, true, false, false, false, false, true, false},
		[]bool{true, false, true, false, false, false, false, true, false},
		[]bool{false, true, true, false, false, false, false, true, false},
		[]bool{true, true, true, false, false, false, false, true, false},
		[]bool{false, false, false, true, false, false, false, true, false},
		[]bool{true, false, false, true, false, false, false, true, false},
		[]bool{false, true, false, true, false, false, false, true, false},
		[]bool{true, true, false, true, false, false, false, true, false},
		[]bool{false, false, true, true, false, false, false, true, false},
		[]bool{true, false, true, true, false, false, false, true, false},
		[]bool{false, true, true, true, false, false, false, true, false},
		[]bool{true, true, true, true, false, false, false, true, false},
		[]bool{false, false, false, false, true, false, false, true, false},
		[]bool{true, false, false, false, true, false, false, true, false},
		[]bool{false, true, false, false, true, false, false, true, false},
		[]bool{true, true, false, false, true, false, false, true, false},
		[]bool{false, false, true, false, true, false, false, true, false},
		[]bool{true, false, true, false, true, false, false, true, false},
		[]bool{false, true, true, false, true, false, false, true, false},
		[]bool{true, true, true, false, true, false, false, true, false},
		[]bool{false, false, false, true, true, false, false, true, false},
		[]bool{true, false, false, true, true, false, false, true, false},
		[]bool{false, true, false, true, true, false, false, true, false},
		[]bool{true, true, false, true, true, false, false, true, false},
		[]bool{false, false, true, true, true, false, false, true, false},
		[]bool{true, false, true, true, true, false, false, true, false},
		[]bool{false, true, true, true, true, false, false, true, false},
		[]bool{true, true, true, true, true, false, false, true, false},
		[]bool{false, false, false, false, false, true, false, true, false},
		[]bool{true, false, false, false, false, true, false, true, false},
		[]bool{false, true, false, false, false, true, false, true, false},
		[]bool{true, true, false, false, false, true, false, true, false},
		[]bool{false, false, true, false, false, true, false, true, false},
		[]bool{true, false, true, false, false, true, false, true, false},
		[]bool{false, true, true, false, false, true, false, true, false},
		[]bool{true, true, true, false, false, true, false, true, false},
		[]bool{false, false, false, true, false, true, false, true, false},
		[]bool{true, false, false, true, false, true, false, true, false},
		[]bool{false, true, false, true, false, true, false, true, false},
		[]bool{true, true, false, true, false, true, false, true, false},
		[]bool{false, false, true, true, false, true, false, true, false},
		[]bool{true, false, true, true, false, true, false, true, false},
		[]bool{false, true, true, true, false, true, false, true, false},
		[]bool{true, true, true, true, false, true, false, true, false},
		[]bool{false, false, false, false, true, true, false, true, false},
		[]bool{true, false, false, false, true, true, false, true, false},
		[]bool{false, true, false, false, true, true, false, true, false},
		[]bool{true, true, false, false, true, true, false, true, false},
		[]bool{false, false, true, false, true, true, false, true, false},
		[]bool{true, false, true, false, true, true, false, true, false},
		[]bool{false, true, true, false, true, true, false, true, false},
		[]bool{true, true, true, false, true, true, false, true, false},
		[]bool{false, false, false, true, true, true, false, true, false},
		[]bool{true, false, false, true, true, true, false, true, false},
		[]bool{false, true, false, true, true, true, false, true, false},
		[]bool{true, true, false, true, true, true, false, true, false},
		[]bool{false, false, true, true, true, true, false, true, false},
		[]bool{true, false, true, true, true, true, false, true, false},
		[]bool{false, true, true, true, true, true, false, true, false},
		[]bool{true, true, true, true, true, true, false, true, false},
		[]bool{false, false, false, false, false, false, true, true, false},
		[]bool{true, false, false, false, false, false, true, true, false},
		[]bool{false, true, false, false, false, false, true, true, false},
		[]bool{true, true, false, false, false, false, true, true, false},
		[]bool{false, false, true, false, false, false, true, true, false},
		[]bool{true, false, true, false, false, false, true, true, false},
		[]bool{false, true, true, false, false, false, true, true, false},
		[]bool{true, true, true, false, false, false, true, true, false},
		[]bool{false, false, false, true, false, false, true, true, false},
		[]bool{true, false, false, true, false, false, true, true, false},
		[]bool{false, true, false, true, false, false, true, true, false},
		[]bool{true, true, false, true, false, false, true, true, false},
		[]bool{false, false, true, true, false, false, true, true, false},
		[]bool{true, false, true, true, false, false, true, true, false},
		[]bool{false, true, true, true, false, false, true, true, false},
		[]bool{true, true, true, true, false, false, true, true, false},
		[]bool{false, false, false, false, true, false, true, true, false},
		[]bool{true, false, false, false, true, false, true, true, false},
		[]bool{false, true, false, false, true, false, true, true, false},
		[]bool{true, true, false, false, true, false, true, true, false},
		[]bool{false, false, true, false, true, false, true, true, false},
		[]bool{true, false, true, false, true, false, true, true, false},
		[]bool{false, true, true, false, true, false, true, true, false},
		[]bool{true, true, true, false, true, false, true, true, false},
		[]bool{false, false, false, true, true, false, true, true, false},
		[]bool{true, false, false, true, true, false, true, true, false},
		[]bool{false, true, false, true, true, false, true, true, false},
		[]bool{true, true, false, true, true, false, true, true, false},
		[]bool{false, false, true, true, true, false, true, true, false},
		[]bool{true, false, true, true, true, false, true, true, false},
		[]bool{false, true, true, true, true, false, true, true, false},
		[]bool{true, true, true, true, true, false, true, true, false},
		[]bool{false, false, false, false, false, true, true, true, false},
		[]bool{true, false, false, false, false, true, true, true, false},
		[]bool{false, true, false, false, false, true, true, true, false},
		[]bool{true, true, false, false, false, true, true, true, false},
		[]bool{false, false, true, false, false, true, true, true, false},
		[]bool{true, false, true, false, false, true, true, true, false},
		[]bool{false, true, true, false, false, true, true, true, false},
		[]bool{true, true, true, false, false, true, true, true, false},
		[]bool{false, false, false, true, false, true, true, true, false},
		[]bool{true, false, false, true, false, true, true, true, false},
		[]bool{false, true, false, true, false, true, true, true, false},
		[]bool{true, true, false, true, false, true, true, true, false},
		[]bool{false, false, true, true, false, true, true, true, false},
		[]bool{true, false, true, true, false, true, true, true, false},
		[]bool{false, true, true, true, false, true, true, true, false},
		[]bool{true, true, true, true, false, true, true, true, false},
		[]bool{false, false, false, false, true, true, true, true, false},
		[]bool{true, false, false, false, true, true, true, true, false},
		[]bool{false, true, false, false, true, true, true, true, false},
		[]bool{true, true, false, false, true, true, true, true, false},
		[]bool{false, false, true, false, true, true, true, true, false},
		[]bool{true, false, true, false, true, true, true, true, false},
		[]bool{false, true, true, false, true, true, true, true, false},
		[]bool{true, true, true, false, true, true, true, true, false},
		[]bool{false, false, false, true, true, true, true, true, false},
		[]bool{true, false, false, true, true, true, true, true, false},
		[]bool{false, true, false, true, true, true, true, true, false},
		[]bool{true, true, false, true, true, true, true, true, false},
		[]bool{false, false, true, true, true, true, true, true, false},
		[]bool{true, false, true, true, true, true, true, true, false},
		[]bool{false, true, true, true, true, true, true, true, false},
		[]bool{true, true, true, true, true, true, true, true, true},
	}

	for _, combination := range combinations {
		gate1 := NewANDGate8()
		gate1.Update(combination[0], combination[1], combination[2], combination[3], combination[4], combination[5], combination[6], combination[7])

		if gate1.Output() != combination[8] {
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
