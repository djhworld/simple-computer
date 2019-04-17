package components

import (
	"fmt"

	"github.com/djhworld/simple-computer/circuit"
)

type ByteComponent interface {
	ConnectOutput(ByteComponent)
	SetInputWire(int, bool)
	GetOutputWire(int) bool
}

type Enabler struct {
	inputs  [8]circuit.Wire
	gates   [8]circuit.ANDGate
	outputs [8]circuit.Wire
	next    ByteComponent
}

func NewEnabler() *Enabler {
	e := new(Enabler)

	for i, _ := range e.gates {
		e.gates[i] = *circuit.NewANDGate()
	}
	return e
}

func (e *Enabler) ConnectOutput(b ByteComponent) {
	e.next = b
}

func (e *Enabler) GetOutputWire(index int) bool {
	return e.outputs[index].Get()
}

func (e *Enabler) SetInputWire(index int, value bool) {
	e.inputs[index].Update(value)
}

func (e *Enabler) Update(enable bool) {
	for i, g := range e.gates {
		g.Update(e.inputs[i].Get(), enable)
		e.outputs[i].Update(g.Output())
	}

	if e.next != nil {
		for i, w := range e.outputs {
			e.next.SetInputWire(i, w.Get())
		}
	}
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

func (g *ANDGate4) Update(inputA, inputB, inputC, inputD bool) {
	g.andA.Update(inputA, inputB)
	g.andB.Update(g.andA.Output(), inputC)
	g.andC.Update(g.andB.Output(), inputD)
	g.output.Update(g.andC.Output())
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

// TODO not sure if this is exactly how this should look...
type LeftShifter struct {
	inputs   [8]circuit.Wire
	outputs  [8]circuit.Wire
	shiftIn  circuit.Wire
	shiftOut circuit.Wire
	next     ByteComponent
}

func NewLeftShifter() *LeftShifter {
	return new(LeftShifter)
}

func (l *LeftShifter) ConnectOutput(b ByteComponent) {
	l.next = b
}

func (l *LeftShifter) GetOutputWire(index int) bool {
	return l.outputs[index].Get()
}

func (l *LeftShifter) SetInputWire(index int, value bool) {
	l.inputs[index].Update(value)
}

func (l *LeftShifter) ShiftOut() bool {
	return l.shiftOut.Get()
}

func (l *LeftShifter) Update(shiftIn bool) {
	l.shiftIn.Update(shiftIn)
	l.shiftOut.Update(l.inputs[0].Get())
	l.outputs[0].Update(l.inputs[1].Get())
	l.outputs[1].Update(l.inputs[2].Get())
	l.outputs[2].Update(l.inputs[3].Get())
	l.outputs[3].Update(l.inputs[4].Get())
	l.outputs[4].Update(l.inputs[5].Get())
	l.outputs[5].Update(l.inputs[6].Get())
	l.outputs[6].Update(l.inputs[7].Get())
	l.outputs[7].Update(l.shiftIn.Get())
}

type RightShifter struct {
	inputs   [8]circuit.Wire
	shiftIn  circuit.Wire
	shiftOut circuit.Wire
	outputs  [8]circuit.Wire
	next     ByteComponent
}

func NewRightShifter() *RightShifter {
	return new(RightShifter)
}

func (r *RightShifter) ConnectOutput(b ByteComponent) {
	r.next = b
}

func (r *RightShifter) GetOutputWire(index int) bool {
	return r.outputs[index].Get()
}

func (r *RightShifter) SetInputWire(index int, value bool) {
	r.inputs[index].Update(value)
}

func (r *RightShifter) ShiftOut() bool {
	return r.shiftOut.Get()
}

func (r *RightShifter) Update(shiftIn bool) {
	r.shiftIn.Update(shiftIn)
	r.outputs[0].Update(r.shiftIn.Get())
	r.outputs[1].Update(r.inputs[0].Get())
	r.outputs[2].Update(r.inputs[1].Get())
	r.outputs[3].Update(r.inputs[2].Get())
	r.outputs[4].Update(r.inputs[3].Get())
	r.outputs[5].Update(r.inputs[4].Get())
	r.outputs[6].Update(r.inputs[5].Get())
	r.outputs[7].Update(r.inputs[6].Get())
	r.shiftOut.Update(r.inputs[7].Get())
}

type IsZero struct {
	inputs  [8]circuit.Wire
	orer    ORer
	notGate circuit.NOTGate
	output  circuit.Wire
}

func NewIsZero() *IsZero {
	z := new(IsZero)
	z.orer = *NewORer()
	z.notGate = *circuit.NewNOTGate()

	return z
}

func (z *IsZero) Reset() {
	z.output.Update(false)
}

func (z *IsZero) ConnectOutput(b ByteComponent) {
	// noop
}

func (z *IsZero) GetOutputWire(index int) bool {
	// only 1 wire
	return z.output.Get()
}

func (z *IsZero) SetInputWire(index int, value bool) {
	z.inputs[index].Update(value)
}

func (z *IsZero) Update() {
	for i, _ := range z.inputs {
		z.orer.SetInputWire(i, z.inputs[i].Get())
		z.orer.SetInputWire(i+8, z.inputs[i].Get())
	}
	z.orer.Update()

	for i, _ := range z.orer.outputs {
		if z.orer.outputs[i].Get() {
			z.notGate.Update(true)
			z.output.Update(z.notGate.Output())
			return
		} else {
			z.notGate.Update(false)
		}
	}

	z.output.Update(z.notGate.Output())

}

type NOTer struct {
	inputs  [8]circuit.Wire
	gates   [8]circuit.NOTGate
	outputs [8]circuit.Wire
	next    ByteComponent
}

func NewNOTer() *NOTer {
	n := new(NOTer)

	for i, _ := range n.gates {
		n.gates[i] = *circuit.NewNOTGate()
	}

	return n
}

func (n *NOTer) ConnectOutput(b ByteComponent) {
	n.next = b
}

func (n *NOTer) GetOutputWire(index int) bool {
	return n.outputs[index].Get()
}

func (n *NOTer) SetInputWire(index int, value bool) {
	n.inputs[index].Update(value)
}

func (n *NOTer) Update() {
	for i, _ := range n.gates {
		n.gates[i].Update(n.inputs[i].Get())
		n.outputs[i].Update(n.gates[i].Output())
	}
}

type ANDer struct {
	inputs  [16]circuit.Wire
	gates   [8]circuit.ANDGate
	outputs [8]circuit.Wire
	next    ByteComponent
}

func NewANDer() *ANDer {
	a := new(ANDer)

	for i, _ := range a.gates {
		a.gates[i] = *circuit.NewANDGate()
	}

	return a
}

func (a *ANDer) ConnectOutput(b ByteComponent) {
	a.next = b
}

func (a *ANDer) GetOutputWire(index int) bool {
	return a.outputs[index].Get()
}

func (a *ANDer) SetInputWire(index int, value bool) {
	a.inputs[index].Update(value)
}

func (a *ANDer) Update() {
	awire := 8
	bwire := 0
	for i, _ := range a.gates {
		a.gates[i].Update(a.inputs[awire].Get(), a.inputs[bwire].Get())
		a.outputs[i].Update(a.gates[i].Output())
		awire++
		bwire++
	}
}

type ORer struct {
	inputs  [16]circuit.Wire
	gates   [8]circuit.ORGate
	outputs [8]circuit.Wire
	next    ByteComponent
}

func NewORer() *ORer {
	o := new(ORer)

	for i, _ := range o.gates {
		o.gates[i] = *circuit.NewORGate()
	}

	return o
}

func (o *ORer) ConnectOutput(b ByteComponent) {
	o.next = b
}

func (o *ORer) GetOutputWire(index int) bool {
	return o.outputs[index].Get()
}

func (o *ORer) SetInputWire(index int, value bool) {
	o.inputs[index].Update(value)
}

func (o *ORer) Update() {
	awire := 8
	bwire := 0
	for i, _ := range o.gates {
		o.gates[i].Update(o.inputs[awire].Get(), o.inputs[bwire].Get())
		o.outputs[i].Update(o.gates[i].Output())
		awire++
		bwire++
	}
}

type XORer struct {
	inputs  [16]circuit.Wire
	gates   [8]circuit.XORGate
	outputs [8]circuit.Wire
	next    ByteComponent
}

func NewXORer() *XORer {
	o := new(XORer)

	for i, _ := range o.gates {
		o.gates[i] = *circuit.NewXORGate()
	}

	return o
}

func (o *XORer) ConnectOutput(b ByteComponent) {
	o.next = b
}

func (o *XORer) GetOutputWire(index int) bool {
	return o.outputs[index].Get()
}

func (o *XORer) SetInputWire(index int, value bool) {
	o.inputs[index].Update(value)
}

func (o *XORer) Update() {
	awire := 8
	bwire := 0
	for i, _ := range o.gates {
		o.gates[i].Update(o.inputs[awire].Get(), o.inputs[bwire].Get())
		o.outputs[i].Update(o.gates[i].Output())
		awire++
		bwire++
	}
}

type Compare2 struct {
	inputA      circuit.Wire
	inputB      circuit.Wire
	xor1        circuit.XORGate
	not1        circuit.NOTGate
	and1        circuit.ANDGate
	andgate3    ANDGate3
	or1         circuit.ORGate
	out         circuit.Wire
	equalIn     circuit.Wire
	equalOut    circuit.Wire
	isLargerIn  circuit.Wire
	isLargerOut circuit.Wire
}

func NewCompare2() *Compare2 {
	c := new(Compare2)

	c.inputA = *circuit.NewWire("a", false)
	c.inputB = *circuit.NewWire("b", false)
	c.equalIn = *circuit.NewWire("eqin", false)
	c.equalOut = *circuit.NewWire("eqin", false)
	c.isLargerIn = *circuit.NewWire("largerin", false)
	c.isLargerOut = *circuit.NewWire("largerin", false)
	c.out = *circuit.NewWire("out", false)

	c.xor1 = *circuit.NewXORGate()
	c.not1 = *circuit.NewNOTGate()
	c.and1 = *circuit.NewANDGate()
	c.andgate3 = *NewANDGate3()
	c.or1 = *circuit.NewORGate()

	return c
}

func (g *Compare2) Equal() bool {
	return g.equalOut.Get()
}

func (g *Compare2) Larger() bool {
	return g.isLargerOut.Get()
}

func (g *Compare2) Output() bool {
	return g.out.Get()
}

func (g *Compare2) Update(inputA, inputB, equalIn, isLargerIn bool) {
	g.inputA.Update(inputA)
	g.inputB.Update(inputB)
	g.equalIn.Update(equalIn)
	g.isLargerIn.Update(isLargerIn)

	g.xor1.Update(g.inputA.Get(), g.inputB.Get())
	g.not1.Update(g.xor1.Output())
	g.and1.Update(g.not1.Output(), g.equalIn.Get())
	g.equalOut.Update(g.and1.Output())

	g.andgate3.Update(g.equalIn.Get(), g.inputA.Get(), g.xor1.Output())
	g.or1.Update(g.andgate3.Output(), g.isLargerIn.Get())
	g.isLargerOut.Update(g.or1.Output())

	g.out.Update(g.xor1.Output())
}

type Comparator struct {
	inputs       [16]circuit.Wire
	equalIn      circuit.Wire
	aIsLargerIn  circuit.Wire
	compares     [8]Compare2
	outputs      [8]circuit.Wire
	equalOut     circuit.Wire
	aIsLargerOut circuit.Wire
	next         ByteComponent
}

func NewComparator() *Comparator {
	c := new(Comparator)

	for i, _ := range c.compares {
		c.compares[i] = *NewCompare2()
	}

	return c
}

func (c *Comparator) ConnectOutput(b ByteComponent) {
	c.next = b
}

func (c *Comparator) GetOutputWire(index int) bool {
	return c.outputs[index].Get()
}

func (c *Comparator) SetInputWire(index int, value bool) {
	c.inputs[index].Update(value)
}

func (c *Comparator) PrintInput() {
	for i := range c.inputs {
		if c.inputs[i].Get() {
			fmt.Print(1)
		} else {
			fmt.Print(0)

		}
		if i == 7 {
			fmt.Print("    ")
		}
	}
	fmt.Println()
}

func (c *Comparator) Update() {
	// these start out as 1 and 0 respectively
	c.equalIn.Update(true)
	c.aIsLargerIn.Update(false)

	// top 8 bits are <b>, bottom 8 bits are <a>
	awire := 0
	bwire := 8

	for i := range c.compares {
		c.compares[i].Update(c.inputs[awire].Get(), c.inputs[bwire].Get(), c.equalIn.Get(), c.aIsLargerIn.Get())
		c.outputs[i].Update(c.compares[i].Output())
		c.equalOut.Update(c.compares[i].Equal())
		c.aIsLargerOut.Update(c.compares[i].Larger())

		c.equalIn.Update(c.compares[i].Equal())
		c.aIsLargerIn.Update(c.compares[i].Larger())
		awire++
		bwire++
	}
}

func (g *Comparator) Equal() bool {
	return g.equalOut.Get()
}

func (g *Comparator) Larger() bool {
	return g.aIsLargerOut.Get()
}

type BusOne struct {
	inputBus  *Bus
	outputBus *Bus
	inputs    [8]circuit.Wire
	bus1      circuit.Wire
	andGates  [7]circuit.ANDGate
	notGate   circuit.NOTGate
	orGate    circuit.ORGate
	outputs   [8]circuit.Wire
	next      ByteComponent
}

func NewBusOne(inputBus, outputBus *Bus) *BusOne {
	b := new(BusOne)
	b.inputBus = inputBus
	b.outputBus = outputBus

	for i, _ := range b.andGates {
		b.andGates[i] = *circuit.NewANDGate()
	}

	b.notGate = *circuit.NewNOTGate()
	b.orGate = *circuit.NewORGate()

	return b
}

func (b *BusOne) ConnectOutput(bc ByteComponent) {
	b.next = bc
}

func (b *BusOne) GetOutputWire(index int) bool {
	return b.outputs[index].Get()
}

func (b *BusOne) SetInputWire(index int, value bool) {
	b.inputs[index].Update(value)
}

func (b *BusOne) Enable() {
	b.bus1.Update(true)
}

func (b *BusOne) Disable() {
	b.bus1.Update(false)
}

func (b *BusOne) Update() {
	for i := 8 - 1; i >= 0; i-- {
		b.inputs[i].Update(b.inputBus.GetOutputWire(i))
	}

	b.notGate.Update(b.bus1.Get())

	b.andGates[0].Update(b.inputs[0].Get(), b.notGate.Output())
	b.andGates[1].Update(b.inputs[1].Get(), b.notGate.Output())
	b.andGates[2].Update(b.inputs[2].Get(), b.notGate.Output())
	b.andGates[3].Update(b.inputs[3].Get(), b.notGate.Output())
	b.andGates[4].Update(b.inputs[4].Get(), b.notGate.Output())
	b.andGates[5].Update(b.inputs[5].Get(), b.notGate.Output())
	b.andGates[6].Update(b.inputs[6].Get(), b.notGate.Output())
	b.orGate.Update(b.inputs[7].Get(), b.bus1.Get())

	b.outputs[0].Update(b.andGates[0].Output())
	b.outputs[1].Update(b.andGates[1].Output())
	b.outputs[2].Update(b.andGates[2].Output())
	b.outputs[3].Update(b.andGates[3].Output())
	b.outputs[4].Update(b.andGates[4].Output())
	b.outputs[5].Update(b.andGates[5].Output())
	b.outputs[6].Update(b.andGates[6].Output())
	b.outputs[7].Update(b.orGate.Output())

	for i := 8 - 1; i >= 0; i-- {
		b.outputBus.SetInputWire(i, b.outputs[i].Get())
	}
}

func (b *BusOne) String() string {
	var output byte
	var x int = 0
	for i := 7; i >= 0; i-- {
		if b.outputs[i].Get() {
			output = output | (1 << byte(x))
		} else {
			output = output & ^(1 << byte(x))
		}
		x++
	}
	return fmt.Sprintf("0x%X", output)
}
