package components

import (
	"testing"

	"github.com/djhworld/simple-computer/arch"
)

func TestBit(t *testing.T) {
	b := NewBit()
	b.Update(false, true)
	if b.Get() != false {
		t.FailNow()
	}

	b.Update(false, false)
	if b.Get() != false {
		t.FailNow()
	}

	b.Update(true, true)
	if b.Get() != true {
		t.FailNow()
	}

	b.Update(false, false)
	if b.Get() != true {
		t.FailNow()
	}
}
func TestWord(t *testing.T) {
	d := new(DummyComponent)

	for i, _ := range d.wires {
		d.SetInputWire(i, true)
	}

	w := NewWord()
	d.ConnectOutput(w)
	d.Update()
	w.Update(false)

	for _, w := range w.outputs {
		if !w.Get() {
			t.FailNow()
		}
	}
}

func TestWordWithSetOn(t *testing.T) {
	d := new(DummyComponent)

	for i, _ := range d.wires {
		d.SetInputWire(i, true)
	}

	w := NewWord()
	d.ConnectOutput(w)
	d.Update()
	w.Update(true)

	results := [arch.BUS_WIDTH]bool{}
	for i, w := range w.outputs {
		results[i] = w.Get()
	}

	if results != [arch.BUS_WIDTH]bool{true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true} {
		t.FailNow()
	}
}
