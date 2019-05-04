package components

import "github.com/djhworld/simple-computer/circuit"

type ORGate3 struct {
	inputA circuit.Wire
	inputB circuit.Wire
	inputC circuit.Wire
	orA    circuit.ORGate
	orB    circuit.ORGate
	output circuit.Wire
}

func NewORGate3() *ORGate3 {
	a := new(ORGate3)

	a.inputA = *circuit.NewWire("a", false)
	a.inputB = *circuit.NewWire("b", false)
	a.inputC = *circuit.NewWire("c", false)
	a.output = *circuit.NewWire("d", false)

	a.orA = *circuit.NewORGate()
	a.orB = *circuit.NewORGate()

	return a
}

func (g *ORGate3) Output() bool {
	return g.output.Get()
}

func (g *ORGate3) Update(inputA bool, inputB bool, inputC bool) {
	g.orA.Update(inputA, inputB)
	g.orB.Update(g.orA.Output(), inputC)

	g.output.Update(g.orB.Output())
}

type ORGate4 struct {
	inputA circuit.Wire
	inputB circuit.Wire
	inputC circuit.Wire
	inputD circuit.Wire
	orA    circuit.ORGate
	orB    circuit.ORGate
	orC    circuit.ORGate
	output circuit.Wire
}

func NewORGate4() *ORGate4 {
	a := new(ORGate4)

	a.inputA = *circuit.NewWire("a", false)
	a.inputB = *circuit.NewWire("b", false)
	a.inputC = *circuit.NewWire("c", false)
	a.inputD = *circuit.NewWire("d", false)
	a.output = *circuit.NewWire("o", false)

	a.orA = *circuit.NewORGate()
	a.orB = *circuit.NewORGate()
	a.orC = *circuit.NewORGate()

	return a
}

func (g *ORGate4) Update(inputA, inputB, inputC, inputD bool) {
	g.orA.Update(inputA, inputB)
	g.orB.Update(g.orA.Output(), inputC)
	g.orC.Update(g.orB.Output(), inputD)
	g.output.Update(g.orC.Output())
}

func (g *ORGate4) Output() bool {
	return g.output.Get()
}

type ORGate5 struct {
	inputA circuit.Wire
	inputB circuit.Wire
	inputC circuit.Wire
	inputD circuit.Wire
	inputE circuit.Wire
	orA    circuit.ORGate
	orB    circuit.ORGate
	orC    circuit.ORGate
	orD    circuit.ORGate
	output circuit.Wire
}

func NewORGate5() *ORGate5 {
	a := new(ORGate5)

	a.inputA = *circuit.NewWire("a", false)
	a.inputB = *circuit.NewWire("b", false)
	a.inputC = *circuit.NewWire("c", false)
	a.inputD = *circuit.NewWire("d", false)
	a.inputE = *circuit.NewWire("e", false)
	a.output = *circuit.NewWire("o", false)

	a.orA = *circuit.NewORGate()
	a.orB = *circuit.NewORGate()
	a.orC = *circuit.NewORGate()
	a.orD = *circuit.NewORGate()

	return a
}

func (g *ORGate5) Update(inputA, inputB, inputC, inputD, inputE bool) {
	g.orA.Update(inputA, inputB)
	g.orB.Update(g.orA.Output(), inputC)
	g.orC.Update(g.orB.Output(), inputD)
	g.orD.Update(g.orC.Output(), inputE)
	g.output.Update(g.orD.Output())
}

func (g *ORGate5) Output() bool {
	return g.output.Get()
}

type ORGate6 struct {
	inputA circuit.Wire
	inputB circuit.Wire
	inputC circuit.Wire
	inputD circuit.Wire
	inputE circuit.Wire
	inputF circuit.Wire
	orA    circuit.ORGate
	orB    circuit.ORGate
	orC    circuit.ORGate
	orD    circuit.ORGate
	orE    circuit.ORGate
	output circuit.Wire
}

func NewORGate6() *ORGate6 {
	a := new(ORGate6)

	a.inputA = *circuit.NewWire("a", false)
	a.inputB = *circuit.NewWire("b", false)
	a.inputC = *circuit.NewWire("c", false)
	a.inputD = *circuit.NewWire("d", false)
	a.inputE = *circuit.NewWire("e", false)
	a.inputF = *circuit.NewWire("f", false)
	a.output = *circuit.NewWire("o", false)

	a.orA = *circuit.NewORGate()
	a.orB = *circuit.NewORGate()
	a.orC = *circuit.NewORGate()
	a.orD = *circuit.NewORGate()
	a.orE = *circuit.NewORGate()

	return a
}

func (g *ORGate6) Update(inputA, inputB, inputC, inputD, inputE, inputF bool) {
	g.orA.Update(inputA, inputB)
	g.orB.Update(g.orA.Output(), inputC)
	g.orC.Update(g.orB.Output(), inputD)
	g.orD.Update(g.orC.Output(), inputE)
	g.orE.Update(g.orD.Output(), inputF)
	g.output.Update(g.orE.Output())
}

func (g *ORGate6) Output() bool {
	return g.output.Get()
}

type ANDGate8 struct {
	andA   circuit.ANDGate
	andB   circuit.ANDGate
	andC   circuit.ANDGate
	andD   circuit.ANDGate
	andE   circuit.ANDGate
	andF   circuit.ANDGate
	andG   circuit.ANDGate
	output circuit.Wire
}

func NewANDGate8() *ANDGate8 {
	a := new(ANDGate8)

	a.output = *circuit.NewWire("o", false)

	a.andA = *circuit.NewANDGate()
	a.andB = *circuit.NewANDGate()
	a.andC = *circuit.NewANDGate()
	a.andD = *circuit.NewANDGate()
	a.andE = *circuit.NewANDGate()
	a.andF = *circuit.NewANDGate()
	a.andG = *circuit.NewANDGate()

	return a
}

func (gt *ANDGate8) Update(a, b, c, d, e, f, g, h bool) {
	gt.andA.Update(a, b)
	gt.andB.Update(gt.andA.Output(), c)
	gt.andC.Update(gt.andB.Output(), d)
	gt.andD.Update(gt.andC.Output(), e)
	gt.andE.Update(gt.andD.Output(), f)
	gt.andF.Update(gt.andE.Output(), g)
	gt.andG.Update(gt.andF.Output(), h)
	gt.output.Update(gt.andG.Output())
}

func (g *ANDGate8) Output() bool {
	return g.output.Get()
}

type ANDGate4 struct {
	inputA circuit.Wire
	inputB circuit.Wire
	inputC circuit.Wire
	inputD circuit.Wire
	andA   circuit.ANDGate
	andB   circuit.ANDGate
	andC   circuit.ANDGate
	output circuit.Wire
}

func NewANDGate4() *ANDGate4 {
	a := new(ANDGate4)

	a.inputA = *circuit.NewWire("a", false)
	a.inputB = *circuit.NewWire("b", false)
	a.inputC = *circuit.NewWire("c", false)
	a.inputD = *circuit.NewWire("d", false)
	a.output = *circuit.NewWire("o", false)

	a.andA = *circuit.NewANDGate()
	a.andB = *circuit.NewANDGate()
	a.andC = *circuit.NewANDGate()

	return a
}

func (g *ANDGate4) Output() bool {
	return g.output.Get()
}

func (g *ANDGate4) Update(inputA, inputB, inputC, inputD bool) {
	g.andA.Update(inputA, inputB)
	g.andB.Update(g.andA.Output(), inputC)
	g.andC.Update(g.andB.Output(), inputD)
	g.output.Update(g.andC.Output())
}

type ANDGate5 struct {
	andA   circuit.ANDGate
	andB   circuit.ANDGate
	andC   circuit.ANDGate
	andD   circuit.ANDGate
	output circuit.Wire
}

func NewANDGate5() *ANDGate5 {
	a := new(ANDGate5)

	a.output = *circuit.NewWire("o", false)

	a.andA = *circuit.NewANDGate()
	a.andB = *circuit.NewANDGate()
	a.andC = *circuit.NewANDGate()
	a.andD = *circuit.NewANDGate()

	return a
}

func (g *ANDGate5) Output() bool {
	return g.output.Get()
}

func (g *ANDGate5) Update(inputA, inputB, inputC, inputD, inputE bool) {
	g.andA.Update(inputA, inputB)
	g.andB.Update(g.andA.Output(), inputC)
	g.andC.Update(g.andB.Output(), inputD)
	g.andD.Update(g.andC.Output(), inputE)
	g.output.Update(g.andD.Output())
}

type ANDGate3 struct {
	inputA circuit.Wire
	inputB circuit.Wire
	inputC circuit.Wire
	andA   circuit.ANDGate
	andB   circuit.ANDGate
	output circuit.Wire
}

func NewANDGate3() *ANDGate3 {
	a := new(ANDGate3)

	a.inputA = *circuit.NewWire("a", false)
	a.inputB = *circuit.NewWire("b", false)
	a.inputC = *circuit.NewWire("c", false)
	a.output = *circuit.NewWire("d", false)

	a.andA = *circuit.NewANDGate()
	a.andB = *circuit.NewANDGate()

	return a
}

func (g *ANDGate3) Output() bool {
	return g.output.Get()
}

func (g *ANDGate3) Update(inputA bool, inputB bool, inputC bool) {
	g.andA.Update(inputA, inputB)
	g.andB.Update(g.andA.Output(), inputC)

	g.output.Update(g.andB.Output())
}
