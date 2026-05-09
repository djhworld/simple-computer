package asm

import (
	"fmt"
	"strings"

	"github.com/djhworld/simple-computer/utils"
)

type LabelResolver func(LABEL) (uint16, error)
type SymbolResolver func(SYMBOL) (uint16, error)

type REGISTER int
type IO_MODE string

const (
	REG0 = REGISTER(iota)
	REG1
	REG2
	REG3
)

const (
	ADDRESS_MODE = IO_MODE("Addr")
	DATA_MODE    = IO_MODE("Data")
)

type Instruction interface {
	String() string
	Emit(LabelResolver, SymbolResolver) ([]uint16, error)
	Size() int
}

// DATA
// put value in memory into register (2 byte instruction)
// ----------------------
// 0x0020 = DATA R0
// 0x0021 = DATA R1
// 0x0022 = DATA R2
// 0x0023 = DATA R3
type DATA struct {
	ToRegister REGISTER
	Data       marker
}

func (d DATA) Size() int {
	return 2
}

func (d DATA) Emit(labelResolver LabelResolver, symbolResolver SymbolResolver) ([]uint16, error) {
	instruction := dataOpcodes[d.ToRegister]

	if v, ok := d.Data.(SYMBOL); ok {
		//TODO get value from symbols map....
		resolvedSymbol, err := symbolResolver(v)
		if err != nil {
			return nil, err
		}
		return []uint16{instruction, resolvedSymbol}, nil
	} else if v, ok := d.Data.(NUMBER); ok {
		return []uint16{instruction, v.Value}, nil
	} else {
		return nil, fmt.Errorf("Unsupported operand for Data %v", d.Data)
	}
}

func (d DATA) String() string {
	if v, ok := d.Data.(NUMBER); ok {
		return fmt.Sprintf("DATA R%d, %s", d.ToRegister, utils.ValueToString(v.Value))
	} else if v, ok := d.Data.(SYMBOL); ok {
		return fmt.Sprintf("DATA R%d, %s", d.ToRegister, v.String())
	}

	return fmt.Sprintf("DATA R%d, %v", d.ToRegister, d.Data)
}

// SHL
// ----------------------
// 0x00A0 = SHL R0
// 0x00A5 = SHL R1
// 0x00AA = SHL R2
// 0x00AF = SHL R3
type SHL struct {
	Register REGISTER
}

func (s SHL) Size() int {
	return 1
}

func (s SHL) Emit(labelResolver LabelResolver, symbolResolver SymbolResolver) ([]uint16, error) {
	return []uint16{shlOpcodes[s.Register]}, nil
}

func (s SHL) String() string {
	result := fmt.Sprintf("SHL R%d", s.Register)
	return result
}

// SHR
// ----------------------
// 0x0090 = SHL R0
// 0x0095 = SHL R1
// 0x009A = SHL R2
// 0x009F = SHL R3
type SHR struct {
	Register REGISTER
}

func (s SHR) Size() int {
	return 1
}

func (s SHR) Emit(labelResolver LabelResolver, symbolResolver SymbolResolver) ([]uint16, error) {
	return []uint16{shrOpcodes[s.Register]}, nil
}

func (s SHR) String() string {
	result := fmt.Sprintf("SHR R%d", s.Register)
	return result
}

// JR
// set instruction address register to value in register
// ----------------------
// 0x0030 = JR R0
// 0x0031 = JR R1
// 0x0032 = JR R2
// 0x0033 = JR R3
type JR struct {
	Register REGISTER
}

func (j JR) Size() int {
	return 1
}

func (j JR) Emit(labelResolver LabelResolver, symbolResolver SymbolResolver) ([]uint16, error) {
	return []uint16{jrOpcodes[j.Register]}, nil
}

func (j JR) String() string {
	result := fmt.Sprintf("JR R%d", j.Register)
	return result
}

// NOT
// ----------------------
// 0x00B0 = NOT R0
// 0x00B5 = NOT R1
// 0x00BA = NOT R2
// 0x00BF = NOT R3
type NOT struct {
	Register REGISTER
}

func (n NOT) Size() int {
	return 1
}

func (n NOT) Emit(labelResolver LabelResolver, symbolResolver SymbolResolver) ([]uint16, error) {
	return []uint16{notOpcodes[n.Register]}, nil
}

func (n NOT) String() string {
	result := fmt.Sprintf("NOT R%d", n.Register)
	return result
}

// STORES
// ----------------------
// arg A = memory address for value
// arg B = value to store in memory
// 0x0010 = ST R0, R0
// 0x0011 = ST R0, R1
// 0x0012 = ST R0, R2
// 0x0013 = ST R0, R3

// 0x0014 = ST R1, R0
// 0x0015 = ST R1, R1
// 0x0016 = ST R1, R2
// 0x0017 = ST R1, R3

// 0x0018 = ST R2, R0
// 0x0019 = ST R2, R1
// 0x001A = ST R2, R2
// 0x001B = ST R2, R3

// 0x001C = ST R3, R0
// 0x001D = ST R3, R1
// 0x001E = ST R3, R2
// 0x001F = ST R3, R3

type STORE struct {
	FromRegister REGISTER
	ToRegister   REGISTER
}

func (s STORE) Size() int {
	return 1
}

func (s STORE) Emit(labelResolver LabelResolver, symbolResolver SymbolResolver) ([]uint16, error) {
	return []uint16{stBases[s.FromRegister] + uint16(s.ToRegister)}, nil
}

func (s STORE) String() string {
	result := fmt.Sprintf("ST R%d, R%d", s.FromRegister, s.ToRegister)
	return result
}

// LOADS
// ----------------------
// arg A = memory address to load from
// arg B = register to store value in
// 0x0000 = LD R0, R0
// 0x0001 = LD R0, R1
// 0x0002 = LD R0, R2
// 0x0003 = LD R0, R3

// 0x0004 = LD R1, R0
// 0x0005 = LD R1, R1
// 0x0006 = LD R1, R2
// 0x0007 = LD R1, R3

// 0x0008 = LD R2, R0
// 0x0009 = LD R2, R1
// 0x000A = LD R2, R2
// 0x000B = LD R2, R3

// 0x000C = LD R3, R0
// 0x000D = LD R3, R1
// 0x000E = LD R3, R2
// 0x000F = LD R3, R3
type LOAD struct {
	MemoryAddressReg REGISTER
	ToRegister       REGISTER
}

func (l LOAD) Size() int {
	return 1
}

func (l LOAD) Emit(labelResolver LabelResolver, symbolResolver SymbolResolver) ([]uint16, error) {
	return []uint16{ldBases[l.MemoryAddressReg] + uint16(l.ToRegister)}, nil
}

func (l LOAD) String() string {
	result := fmt.Sprintf("LD R%d, R%d", l.MemoryAddressReg, l.ToRegister)
	return result
}

// OUT
// ----------------------
// 0x0078 = OUT Data, R0
// 0x0079 = OUT Data, R1
// 0x007A = OUT Data, R2
// 0x007B = OUT Data, R3

// 0x007C = OUT Addr, R0
// 0x007D = OUT Addr, R1
// 0x007E = OUT Addr, R2
// 0x007F = OUT Addr, R3
type OUT struct {
	IoMode       IO_MODE
	FromRegister REGISTER
}

func (o OUT) Size() int {
	return 1
}

func (o OUT) Emit(labelResolver LabelResolver, symbolResolver SymbolResolver) ([]uint16, error) {
	switch o.IoMode {
	case DATA_MODE:
		return []uint16{opOUTData + uint16(o.FromRegister)}, nil
	case ADDRESS_MODE:
		return []uint16{opOUTAddr + uint16(o.FromRegister)}, nil
	default:
		return nil, fmt.Errorf("unsupported io mode for OUT instruction")
	}
}

func (o OUT) String() string {
	result := fmt.Sprintf("OUT %s, R%d", o.IoMode, o.FromRegister)
	return result
}

// IN
// ----------------------
// 0x0070 = IN Data, R0
// 0x0071 = IN Data, R1
// 0x0072 = IN Data, R2
// 0x0073 = IN Data, R3

// 0x0074 = IN Addr, R0
// 0x0075 = IN Addr, R1
// 0x0076 = IN Addr, R2
// 0x0077 = IN Addr, R3
type IN struct {
	IoMode     IO_MODE
	ToRegister REGISTER
}

func (i IN) Size() int {
	return 1
}

func (i IN) Emit(labelResolver LabelResolver, symbolResolver SymbolResolver) ([]uint16, error) {
	switch i.IoMode {
	case DATA_MODE:
		return []uint16{opINData + uint16(i.ToRegister)}, nil
	case ADDRESS_MODE:
		return []uint16{opINAddr + uint16(i.ToRegister)}, nil
	default:
		return nil, fmt.Errorf("unsupported io mode for IN instruction")
	}
}

func (i IN) String() string {
	result := fmt.Sprintf("IN %s, R%d", i.IoMode, i.ToRegister)
	return result
}

// XORS
// ----------------------
// 0x00E0 = XOR R0, R0
// 0x00E1 = XOR R0, R1
// 0x00E2 = XOR R0, R2
// 0x00E3 = XOR R0, R3

// 0x00E4 = XOR R1, R0
// 0x00E5 = XOR R1, R1
// 0x00E6 = XOR R1, R2
// 0x00E7 = XOR R1, R3

// 0x00E8 = XOR R2, R0
// 0x00E9 = XOR R2, R1
// 0x00EA = XOR R2, R2
// 0x00EB = XOR R2, R3

// 0x00EC = XOR R3, R0
// 0x00ED = XOR R3, R1
// 0x00EE = XOR R3, R2
// 0x00EF = XOR R3, R3
type XOR struct {
	ARegister REGISTER
	BRegister REGISTER
}

func (x XOR) Size() int {
	return 1
}

func (x XOR) Emit(labelResolver LabelResolver, symbolResolver SymbolResolver) ([]uint16, error) {
	return []uint16{xorBases[x.ARegister] + uint16(x.BRegister)}, nil
}

func (x XOR) String() string {
	result := fmt.Sprintf("XOR R%d, R%d", x.ARegister, x.BRegister)
	return result
}

// ORS
// ----------------------
// 0x00D0 = OR R0, R0
// 0x00D1 = OR R0, R1
// 0x00D2 = OR R0, R2
// 0x00D3 = OR R0, R3

// 0x00D4 = OR R1, R0
// 0x00D5 = OR R1, R1
// 0x00D6 = OR R1, R2
// 0x00D7 = OR R1, R3

// 0x00D8 = OR R2, R0
// 0x00D9 = OR R2, R1
// 0x00DA = OR R2, R2
// 0x00DB = OR R2, R3

// 0x00DC = OR R3, R0
// 0x00DD = OR R3, R1
// 0x00DE = OR R3, R2
// 0x00DF = OR R3, R3
type OR struct {
	ARegister REGISTER
	BRegister REGISTER
}

func (o OR) Size() int {
	return 1
}

func (o OR) Emit(labelResolver LabelResolver, symbolResolver SymbolResolver) ([]uint16, error) {
	return []uint16{orBases[o.ARegister] + uint16(o.BRegister)}, nil
}

func (o OR) String() string {
	result := fmt.Sprintf("OR R%d, R%d", o.ARegister, o.BRegister)
	return result
}

// ANDS
// ----------------------
// 0x00C0 = AND R0, R0
// 0x00C1 = AND R0, R1
// 0x00C2 = AND R0, R2
// 0x00C3 = AND R0, R3

// 0x00C4 = AND R1, R0
// 0x00C5 = AND R1, R1
// 0x00C6 = AND R1, R2
// 0x00C7 = AND R1, R3

// 0x00C8 = AND R2, R0
// 0x00C9 = AND R2, R1
// 0x00CA = AND R2, R2
// 0x00CB = AND R2, R3

// 0x00CC = AND R3, R0
// 0x00CD = AND R3, R1
// 0x00CE = AND R3, R2
// 0x00CF = AND R3, R3
type AND struct {
	ARegister REGISTER
	BRegister REGISTER
}

func (a AND) Size() int {
	return 1
}

func (a AND) Emit(labelResolver LabelResolver, symbolResolver SymbolResolver) ([]uint16, error) {
	return []uint16{andBases[a.ARegister] + uint16(a.BRegister)}, nil
}

func (a AND) String() string {
	result := fmt.Sprintf("AND R%d, R%d", a.ARegister, a.BRegister)
	return result
}

// CMP
// ----------------------
// 0x00F0 = CMP R0, R0
// 0x00F1 = CMP R0, R1
// 0x00F2 = CMP R0, R2
// 0x00F3 = CMP R0, R3

// 0x00F4 = CMP R1, R0
// 0x00F5 = CMP R1, R1
// 0x00F6 = CMP R1, R2
// 0x00F7 = CMP R1, R3

// 0x00F8 = CMP R2, R0
// 0x00F9 = CMP R2, R1
// 0x00FA = CMP R2, R2
// 0x00FB = CMP R2, R3

// 0x00FC = CMP R3, R0
// 0x00FD = CMP R3, R1
// 0x00FE = CMP R3, R2
// 0x00FF = CMP R3, R3
type CMP struct {
	ARegister REGISTER
	BRegister REGISTER
}

func (c CMP) Size() int {
	return 1
}

func (c CMP) Emit(labelResolver LabelResolver, symbolResolver SymbolResolver) ([]uint16, error) {
	return []uint16{cmpBases[c.ARegister] + uint16(c.BRegister)}, nil
}

func (c CMP) String() string {
	result := fmt.Sprintf("CMP R%d, R%d", c.ARegister, c.BRegister)
	return result
}

// CLF (CLEAR FLAGS)
// ----------------------
// 0x0060 CLF
type CLF struct {
}

func (c CLF) Size() int {
	return 1
}

func (c CLF) Emit(labelResolver LabelResolver, symbolResolver SymbolResolver) ([]uint16, error) {
	return []uint16{opCLF}, nil
}

func (c CLF) String() string {
	result := fmt.Sprintf("CLF")
	return result
}

type JMPF struct {
	Flags   []string
	JumpLoc LABEL
}

func (j JMPF) Size() int {
	return 2
}

func (j JMPF) Emit(labelResolver LabelResolver, symbolResolver SymbolResolver) ([]uint16, error) {
	flags := strings.Join(j.Flags, "")

	instruction, ok := jmpfOpcodes[flags]
	if !ok {
		return nil, fmt.Errorf("unsupported flag combination '%s' for JMPF", flags)
	}

	resolvedAddress, err := labelResolver(j.JumpLoc)
	if err != nil {
		return nil, err
	}
	return []uint16{instruction, resolvedAddress}, nil
}

func (j JMPF) String() string {
	flags := strings.Join(j.Flags, "")
	result := fmt.Sprintf("JMP%s %s", flags, j.JumpLoc)
	return result
}

func (j JMP) Size() int {
	return 2
}

type JMP struct {
	JumpLoc LABEL
}

func (j JMP) Emit(labelResolver LabelResolver, symbolResolver SymbolResolver) ([]uint16, error) {
	resolvedAddress, err := labelResolver(j.JumpLoc)
	if err != nil {
		return nil, err
	}
	return []uint16{opJMP, resolvedAddress}, nil
}

func (j JMP) String() string {
	result := fmt.Sprintf("JMP %s", j.JumpLoc)
	return result
}

// ADDS
// ----------------------
// 0x0080 = ADD R0, R0
// 0x0081 = ADD R0, R1
// 0x0082 = ADD R0, R2
// 0x0083 = ADD R0, R3

// 0x0084 = ADD R1, R0
// 0x0085 = ADD R1, R1
// 0x0086 = ADD R1, R2
// 0x0087 = ADD R1, R3

// 0x0088 = ADD R2, R0
// 0x0089 = ADD R2, R1
// 0x008A = ADD R2, R2
// 0x008B = ADD R2, R3

// 0x008C = ADD R3, R0
// 0x008D = ADD R3, R1
// 0x008E = ADD R3, R2
// 0x008F = ADD R3, R3
type ADD struct {
	ARegister REGISTER
	BRegister REGISTER
}

func (a ADD) Size() int {
	return 1
}

func (a ADD) Emit(labelResolver LabelResolver, symbolResolver SymbolResolver) ([]uint16, error) {
	return []uint16{addBases[a.ARegister] + uint16(a.BRegister)}, nil
}

func (a ADD) String() string {
	result := fmt.Sprintf("ADD R%d, R%d", a.ARegister, a.BRegister)
	return result
}

// PLACEHOLDER INSTRUCTIONS - these are used by the assembler
type DEFLABEL struct {
	Name string
}

func (l DEFLABEL) Size() int {
	return 0
}

func (l DEFLABEL) Emit(labelResolver LabelResolver, symbolResolver SymbolResolver) ([]uint16, error) {
	// noop
	return nil, nil
}

func (l DEFLABEL) String() string {
	return l.Name
}

type DEFSYMBOL struct {
	Name  string
	Value uint16
}

func (s DEFSYMBOL) Size() int {
	return 0
}

func (s DEFSYMBOL) Emit(labelResolver LabelResolver, symbolResolver SymbolResolver) ([]uint16, error) {
	// noop
	return nil, nil
}

func (s DEFSYMBOL) String() string {
	return fmt.Sprintf("%%%s = 0x%X", s.Name, s.Value)
}

// PSUEDO INSTRUCTIONS - these are  composite instructions that may map to multiple opcodes

type CALL struct {
	Routine LABEL
}

func (c CALL) Size() int {
	return 4
}

func (c CALL) Emit(labelResolver LabelResolver, symbolResolver SymbolResolver) ([]uint16, error) {
	nextInsAddress, err := symbolResolver(SYMBOL{NEXTINSTRUCTION})
	if err != nil {
		return nil, err
	}

	compositeInstructions := []Instruction{
		DATA{REG3, NUMBER{nextInsAddress}},
		JMP{c.Routine},
	}

	emitted := []uint16{}
	for _, ins := range compositeInstructions {
		if e, err := ins.Emit(labelResolver, symbolResolver); err == nil {
			emitted = append(emitted, e...)
		} else {
			return nil, err
		}
	}

	return emitted, nil
}

func (c CALL) String() string {
	return fmt.Sprintf("CALL %s", c.Routine)
}

// Instructions - useful list data structure for convienience
type Instructions struct {
	instructions []Instruction
}

func (s *Instructions) Add(ins ...Instruction) {
	for _, i := range ins {
		s.instructions = append(s.instructions, i)
	}
}

func (s *Instructions) AddBlocks(blocks ...[]Instruction) {
	for _, block := range blocks {
		s.Add(block...)
	}
}

func (s *Instructions) Get() []Instruction {
	return s.instructions
}

func (s *Instructions) String() string {
	result := strings.Builder{}

	for _, ins := range s.instructions {
		if _, ok := ins.(DEFLABEL); ok {
			l := ins.(DEFLABEL)
			result.WriteString("\n")
			result.WriteString(l.Name)
			result.WriteString(":\n")
			continue
		}

		if _, ok := ins.(DEFSYMBOL); ok {
			s := ins.(DEFSYMBOL)
			result.WriteString("\n")
			result.WriteString(s.String())
			result.WriteString("\n")
			continue
		}

		result.WriteString("\t")
		result.WriteString(ins.String())
		result.WriteString("\n")
	}

	return result.String()
}
