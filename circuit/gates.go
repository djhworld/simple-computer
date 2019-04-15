package circuit

type NANDGate struct {
	inputA Wire
	inputB Wire
	output Wire
}

func NewNANDGate() *NANDGate {
	return &NANDGate{
		inputA: *NewWire("I", false),
		inputB: *NewWire("S", false),
		output: *NewWire("O", false),
	}
}

func (g *NANDGate) Output() bool {
	return g.output.Get()
}

func (g *NANDGate) Update(inputA, inputB bool) {
	g.inputA.Update(inputA)
	g.inputB.Update(inputB)
	if g.inputA.Get() && g.inputB.Get() {
		g.output.Update(false)
	} else {
		g.output.Update(true)
	}
}

type ANDGate struct {
	inputA Wire
	inputB Wire
	output Wire
}

func NewANDGate() *ANDGate {
	return &ANDGate{
		inputA: *NewWire("I", false),
		inputB: *NewWire("S", false),
		output: *NewWire("O", false),
	}
}

func (g *ANDGate) Update(inputA bool, inputB bool) {
	g.inputA.Update(inputA)
	g.inputB.Update(inputB)
	if g.inputA.Get() && g.inputB.Get() {
		g.output.Update(true)
	} else {
		g.output.Update(false)
	}
}

func (g *ANDGate) Output() bool {
	return g.output.Get()
}

type NOTGate struct {
	inputA Wire
	output Wire
}

func NewNOTGate() *NOTGate {
	return &NOTGate{
		inputA: *NewWire("I", false),
		output: *NewWire("O", false),
	}
}

func (g *NOTGate) Update(input bool) {
	g.inputA.Update(input)
	if g.inputA.Get() {
		g.output.Update(false)
	} else {
		g.output.Update(true)
	}
}

func (g *NOTGate) Output() bool {
	return g.output.Get()
}

type ORGate struct {
	inputA Wire
	inputB Wire
	output Wire
}

func NewORGate() *ORGate {
	return &ORGate{
		inputA: *NewWire("I", false),
		inputB: *NewWire("S", false),
		output: *NewWire("O", false),
	}
}

func (g *ORGate) Output() bool {
	return g.output.Get()
}

func (g *ORGate) Update(inputA, inputB bool) {
	g.inputA.Update(inputA)
	g.inputB.Update(inputB)
	if !g.inputA.Get() && !g.inputB.Get() {
		g.output.Update(false)
	} else {
		g.output.Update(true)
	}

}

type XORGate struct {
	inputA Wire
	inputB Wire
	output Wire
}

func NewXORGate() *XORGate {
	return &XORGate{
		inputA: *NewWire("I", false),
		inputB: *NewWire("S", false),
		output: *NewWire("O", false),
	}
}

func (g *XORGate) Output() bool {
	return g.output.Get()
}

func (g *XORGate) Update(inputA, inputB bool) {
	g.inputA.Update(inputA)
	g.inputB.Update(inputB)
	if !g.inputA.Get() && !g.inputB.Get() {
		g.output.Update(false)
	} else if g.inputA.Get() && g.inputB.Get() {
		g.output.Update(false)
	} else {
		g.output.Update(true)
	}
}

type NORGate struct {
	inputA Wire
	inputB Wire
	output Wire
}

func NewNORGate() *NORGate {
	return &NORGate{
		inputA: *NewWire("I", false),
		inputB: *NewWire("S", false),
		output: *NewWire("O", false),
	}
}

func (g *NORGate) Update(inputA, inputB bool) {
	g.inputA.Update(inputA)
	g.inputB.Update(inputB)
	if !g.inputA.Get() && !g.inputB.Get() {
		g.output.Update(true)
	} else {
		g.output.Update(false)
	}
}
