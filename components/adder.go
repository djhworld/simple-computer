package components

import "github.com/djhworld/simple-computer/circuit"

type Adder struct {
	inputs   [16]circuit.Wire
	carryIn  circuit.Wire
	adds     [8]Add2
	carryOut circuit.Wire
	outputs  [8]circuit.Wire
	next     ByteComponent
}

func NewAdder() *Adder {
	a := new(Adder)

	for i, _ := range a.adds {
		a.adds[i] = *NewAdd2()
	}

	return a
}

func (a *Adder) ConnectOutput(b ByteComponent) {
	a.next = b
}

func (a *Adder) GetOutputWire(index int) bool {
	return a.outputs[index].Get()
}

func (a *Adder) SetInputWire(index int, value bool) {
	a.inputs[index].Update(value)
}

func (a *Adder) Carry() bool {
	return a.carryOut.Get()
}

func (a *Adder) Update(carryIn bool) {
	a.carryIn.Update(carryIn)

	awire := 15
	bwire := 7
	//	for i, _ := range a.adds {
	for i := len(a.adds) - 1; i >= 0; i-- {
		aval := a.inputs[awire].Get()
		bval := a.inputs[bwire].Get()

		a.adds[i].Update(aval, bval, a.carryIn.Get())
		a.outputs[i].Update(a.adds[i].Sum())
		a.carryOut.Update(a.adds[i].Carry())
		a.carryIn.Update(a.adds[i].Carry())
		awire--
		bwire--
	}
}

type Add2 struct {
	inputA   circuit.Wire
	inputB   circuit.Wire
	carryIn  circuit.Wire
	xor1     circuit.XORGate
	xor2     circuit.XORGate
	and1     circuit.ANDGate
	and2     circuit.ANDGate
	or1      circuit.ORGate
	carryOut circuit.Wire
	sumOut   circuit.Wire
}

func NewAdd2() *Add2 {
	a := new(Add2)

	a.inputA = *circuit.NewWire("a", false)
	a.inputB = *circuit.NewWire("b", false)
	a.carryIn = *circuit.NewWire("c", false)
	a.carryOut = *circuit.NewWire("co", false)
	a.sumOut = *circuit.NewWire("so", false)

	a.xor1 = *circuit.NewXORGate()
	a.xor2 = *circuit.NewXORGate()
	a.and1 = *circuit.NewANDGate()
	a.and2 = *circuit.NewANDGate()
	a.or1 = *circuit.NewORGate()

	return a
}

func (g *Add2) Update(inputA, inputB, carryIn bool) {
	g.inputA.Update(inputA)
	g.inputB.Update(inputB)
	g.carryIn.Update(carryIn)

	g.xor1.Update(g.inputA.Get(), g.inputB.Get())
	g.xor2.Update(g.xor1.Output(), g.carryIn.Get())

	g.sumOut.Update(g.xor2.Output())

	g.and1.Update(g.carryIn.Get(), g.xor1.Output())
	g.and2.Update(g.inputA.Get(), g.inputB.Get())

	g.or1.Update(g.and1.Output(), g.and2.Output())

	g.carryOut.Update(g.or1.Output())
}

func (g *Add2) Sum() bool {
	return g.sumOut.Get()
}

func (g *Add2) Carry() bool {
	return g.carryOut.Get()
}
