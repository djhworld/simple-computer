package components

import "testing"

func TestByte(t *testing.T) {
	d := new(DummyComponent)

	for i, _ := range d.wires {
		d.SetInputWire(i, true)
	}

	b := NewByte()
	d.ConnectOutput(b)
	d.Update()
	b.Update(false)

	for _, w := range b.outputs {
		if !w.Get() {
			t.FailNow()
		}
	}
}

func TestByteWithSetOn(t *testing.T) {
	d := new(DummyComponent)

	for i, _ := range d.wires {
		d.SetInputWire(i, true)
	}

	b := NewByte()
	d.ConnectOutput(b)
	d.Update()
	b.Update(true)

	results := [8]bool{}
	for i, w := range b.outputs {
		results[i] = w.Get()
	}

	if results != [8]bool{true, true, true, true, true, true, true, true} {
		t.FailNow()
	}
}
