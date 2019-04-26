package components

import "testing"

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

	results := [16]bool{}
	for i, w := range w.outputs {
		results[i] = w.Get()
	}

	if results != [16]bool{true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true} {
		t.FailNow()
	}
}
