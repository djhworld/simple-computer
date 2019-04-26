package circuit

type NANDGate struct {
	output Wire
}

func NewNANDGate() *NANDGate {
	return &NANDGate{
		output: *NewWire("O", false),
	}
}

func (g *NANDGate) Output() bool {
	return g.output.Get()
}

func (g *NANDGate) Update(inputA, inputB bool) {
	g.output.Update(!(inputA && inputB))
}

type ANDGate struct {
	output Wire
}

func NewANDGate() *ANDGate {
	return &ANDGate{
		output: *NewWire("O", false),
	}
}

func (g *ANDGate) Update(inputA bool, inputB bool) {
	g.output.Update((inputA && inputB))
}

func (g *ANDGate) Output() bool {
	return g.output.Get()
}

type NOTGate struct {
	output Wire
}

func NewNOTGate() *NOTGate {
	return &NOTGate{
		output: *NewWire("O", false),
	}
}

func (g *NOTGate) Update(input bool) {
	g.output.Update(!input)
}

func (g *NOTGate) Output() bool {
	return g.output.Get()
}

type ORGate struct {
	output Wire
}

func NewORGate() *ORGate {
	return &ORGate{
		output: *NewWire("O", false),
	}
}

func (g *ORGate) Output() bool {
	return g.output.Get()
}

func (g *ORGate) Update(inputA, inputB bool) {
	g.output.Update(!(!inputA && !inputB))
}

type XORGate struct {
	output Wire
}

func NewXORGate() *XORGate {
	return &XORGate{
		output: *NewWire("O", false),
	}
}

func (g *XORGate) Output() bool {
	return g.output.Get()
}

func (g *XORGate) Update(inputA, inputB bool) {
	g.output.Update(!((!inputA && !inputB) || (inputA && inputB)))
}

type NORGate struct {
	output Wire
}

func NewNORGate() *NORGate {
	return &NORGate{
		output: *NewWire("O", false),
	}
}

func (g *NORGate) Update(inputA, inputB bool) {
	g.output.Update(!inputA && !inputB)
}
