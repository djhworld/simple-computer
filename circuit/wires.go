package circuit

type Wire struct {
	Name  string
	value bool
}

func NewWire(name string, value bool) *Wire {
	return &Wire{
		Name:  name,
		value: value,
	}
}

func (w *Wire) Update(value bool) {
	w.value = value
}

func (w *Wire) Get() bool {
	return w.value
}
