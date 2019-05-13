package asm

import (
	"fmt"
	"reflect"
	"testing"
)

func TestTwoRegInstructionsString(t *testing.T) {
	var TABLE map[Instruction]string = map[Instruction]string{
		LOAD{REG0, REG0}: "LD R0, R0",
		LOAD{REG0, REG1}: "LD R0, R1",
		LOAD{REG0, REG2}: "LD R0, R2",
		LOAD{REG0, REG3}: "LD R0, R3",
		LOAD{REG1, REG0}: "LD R1, R0",
		LOAD{REG1, REG1}: "LD R1, R1",
		LOAD{REG1, REG2}: "LD R1, R2",
		LOAD{REG1, REG3}: "LD R1, R3",
		LOAD{REG2, REG0}: "LD R2, R0",
		LOAD{REG2, REG1}: "LD R2, R1",
		LOAD{REG2, REG2}: "LD R2, R2",
		LOAD{REG2, REG3}: "LD R2, R3",
		LOAD{REG3, REG0}: "LD R3, R0",
		LOAD{REG3, REG1}: "LD R3, R1",
		LOAD{REG3, REG2}: "LD R3, R2",
		LOAD{REG3, REG3}: "LD R3, R3",

		STORE{REG0, REG0}: "ST R0, R0",
		STORE{REG0, REG1}: "ST R0, R1",
		STORE{REG0, REG2}: "ST R0, R2",
		STORE{REG0, REG3}: "ST R0, R3",
		STORE{REG1, REG0}: "ST R1, R0",
		STORE{REG1, REG1}: "ST R1, R1",
		STORE{REG1, REG2}: "ST R1, R2",
		STORE{REG1, REG3}: "ST R1, R3",
		STORE{REG2, REG0}: "ST R2, R0",
		STORE{REG2, REG1}: "ST R2, R1",
		STORE{REG2, REG2}: "ST R2, R2",
		STORE{REG2, REG3}: "ST R2, R3",
		STORE{REG3, REG0}: "ST R3, R0",
		STORE{REG3, REG1}: "ST R3, R1",
		STORE{REG3, REG2}: "ST R3, R2",
		STORE{REG3, REG3}: "ST R3, R3",

		ADD{REG0, REG0}: "ADD R0, R0",
		ADD{REG0, REG1}: "ADD R0, R1",
		ADD{REG0, REG2}: "ADD R0, R2",
		ADD{REG0, REG3}: "ADD R0, R3",
		ADD{REG1, REG0}: "ADD R1, R0",
		ADD{REG1, REG1}: "ADD R1, R1",
		ADD{REG1, REG2}: "ADD R1, R2",
		ADD{REG1, REG3}: "ADD R1, R3",
		ADD{REG2, REG0}: "ADD R2, R0",
		ADD{REG2, REG1}: "ADD R2, R1",
		ADD{REG2, REG2}: "ADD R2, R2",
		ADD{REG2, REG3}: "ADD R2, R3",
		ADD{REG3, REG0}: "ADD R3, R0",
		ADD{REG3, REG1}: "ADD R3, R1",
		ADD{REG3, REG2}: "ADD R3, R2",
		ADD{REG3, REG3}: "ADD R3, R3",

		AND{REG0, REG0}: "AND R0, R0",
		AND{REG0, REG1}: "AND R0, R1",
		AND{REG0, REG2}: "AND R0, R2",
		AND{REG0, REG3}: "AND R0, R3",
		AND{REG1, REG0}: "AND R1, R0",
		AND{REG1, REG1}: "AND R1, R1",
		AND{REG1, REG2}: "AND R1, R2",
		AND{REG1, REG3}: "AND R1, R3",
		AND{REG2, REG0}: "AND R2, R0",
		AND{REG2, REG1}: "AND R2, R1",
		AND{REG2, REG2}: "AND R2, R2",
		AND{REG2, REG3}: "AND R2, R3",
		AND{REG3, REG0}: "AND R3, R0",
		AND{REG3, REG1}: "AND R3, R1",
		AND{REG3, REG2}: "AND R3, R2",
		AND{REG3, REG3}: "AND R3, R3",

		OR{REG0, REG0}: "OR R0, R0",
		OR{REG0, REG1}: "OR R0, R1",
		OR{REG0, REG2}: "OR R0, R2",
		OR{REG0, REG3}: "OR R0, R3",
		OR{REG1, REG0}: "OR R1, R0",
		OR{REG1, REG1}: "OR R1, R1",
		OR{REG1, REG2}: "OR R1, R2",
		OR{REG1, REG3}: "OR R1, R3",
		OR{REG2, REG0}: "OR R2, R0",
		OR{REG2, REG1}: "OR R2, R1",
		OR{REG2, REG2}: "OR R2, R2",
		OR{REG2, REG3}: "OR R2, R3",
		OR{REG3, REG0}: "OR R3, R0",
		OR{REG3, REG1}: "OR R3, R1",
		OR{REG3, REG2}: "OR R3, R2",
		OR{REG3, REG3}: "OR R3, R3",

		XOR{REG0, REG0}: "XOR R0, R0",
		XOR{REG0, REG1}: "XOR R0, R1",
		XOR{REG0, REG2}: "XOR R0, R2",
		XOR{REG0, REG3}: "XOR R0, R3",
		XOR{REG1, REG0}: "XOR R1, R0",
		XOR{REG1, REG1}: "XOR R1, R1",
		XOR{REG1, REG2}: "XOR R1, R2",
		XOR{REG1, REG3}: "XOR R1, R3",
		XOR{REG2, REG0}: "XOR R2, R0",
		XOR{REG2, REG1}: "XOR R2, R1",
		XOR{REG2, REG2}: "XOR R2, R2",
		XOR{REG2, REG3}: "XOR R2, R3",
		XOR{REG3, REG0}: "XOR R3, R0",
		XOR{REG3, REG1}: "XOR R3, R1",
		XOR{REG3, REG2}: "XOR R3, R2",
		XOR{REG3, REG3}: "XOR R3, R3",

		CMP{REG0, REG0}: "CMP R0, R0",
		CMP{REG0, REG1}: "CMP R0, R1",
		CMP{REG0, REG2}: "CMP R0, R2",
		CMP{REG0, REG3}: "CMP R0, R3",
		CMP{REG1, REG0}: "CMP R1, R0",
		CMP{REG1, REG1}: "CMP R1, R1",
		CMP{REG1, REG2}: "CMP R1, R2",
		CMP{REG1, REG3}: "CMP R1, R3",
		CMP{REG2, REG0}: "CMP R2, R0",
		CMP{REG2, REG1}: "CMP R2, R1",
		CMP{REG2, REG2}: "CMP R2, R2",
		CMP{REG2, REG3}: "CMP R2, R3",
		CMP{REG3, REG0}: "CMP R3, R0",
		CMP{REG3, REG1}: "CMP R3, R1",
		CMP{REG3, REG2}: "CMP R3, R2",
		CMP{REG3, REG3}: "CMP R3, R3",
	}

	for ins, expected := range TABLE {
		if ins.String() != expected {
			t.Logf("Expected %s got %s when testing %s", expected, ins.String(), ins)
			t.FailNow()
		}
	}
}

func TestTwoRegInstructions(t *testing.T) {
	var TABLE map[Instruction][]uint16 = map[Instruction][]uint16{
		LOAD{REG0, REG0}: []uint16{0x00},
		LOAD{REG0, REG1}: []uint16{0x01},
		LOAD{REG0, REG2}: []uint16{0x02},
		LOAD{REG0, REG3}: []uint16{0x03},
		LOAD{REG1, REG0}: []uint16{0x04},
		LOAD{REG1, REG1}: []uint16{0x05},
		LOAD{REG1, REG2}: []uint16{0x06},
		LOAD{REG1, REG3}: []uint16{0x07},
		LOAD{REG2, REG0}: []uint16{0x08},
		LOAD{REG2, REG1}: []uint16{0x09},
		LOAD{REG2, REG2}: []uint16{0x0A},
		LOAD{REG2, REG3}: []uint16{0x0B},
		LOAD{REG3, REG0}: []uint16{0x0C},
		LOAD{REG3, REG1}: []uint16{0x0D},
		LOAD{REG3, REG2}: []uint16{0x0E},
		LOAD{REG3, REG3}: []uint16{0x0F},

		STORE{REG0, REG0}: []uint16{0x10},
		STORE{REG0, REG1}: []uint16{0x11},
		STORE{REG0, REG2}: []uint16{0x12},
		STORE{REG0, REG3}: []uint16{0x13},
		STORE{REG1, REG0}: []uint16{0x14},
		STORE{REG1, REG1}: []uint16{0x15},
		STORE{REG1, REG2}: []uint16{0x16},
		STORE{REG1, REG3}: []uint16{0x17},
		STORE{REG2, REG0}: []uint16{0x18},
		STORE{REG2, REG1}: []uint16{0x19},
		STORE{REG2, REG2}: []uint16{0x1A},
		STORE{REG2, REG3}: []uint16{0x1B},
		STORE{REG3, REG0}: []uint16{0x1C},
		STORE{REG3, REG1}: []uint16{0x1D},
		STORE{REG3, REG2}: []uint16{0x1E},
		STORE{REG3, REG3}: []uint16{0x1F},

		ADD{REG0, REG0}: []uint16{0x80},
		ADD{REG0, REG1}: []uint16{0x81},
		ADD{REG0, REG2}: []uint16{0x82},
		ADD{REG0, REG3}: []uint16{0x83},
		ADD{REG1, REG0}: []uint16{0x84},
		ADD{REG1, REG1}: []uint16{0x85},
		ADD{REG1, REG2}: []uint16{0x86},
		ADD{REG1, REG3}: []uint16{0x87},
		ADD{REG2, REG0}: []uint16{0x88},
		ADD{REG2, REG1}: []uint16{0x89},
		ADD{REG2, REG2}: []uint16{0x8A},
		ADD{REG2, REG3}: []uint16{0x8B},
		ADD{REG3, REG0}: []uint16{0x8C},
		ADD{REG3, REG1}: []uint16{0x8D},
		ADD{REG3, REG2}: []uint16{0x8E},
		ADD{REG3, REG3}: []uint16{0x8F},

		AND{REG0, REG0}: []uint16{0xC0},
		AND{REG0, REG1}: []uint16{0xC1},
		AND{REG0, REG2}: []uint16{0xC2},
		AND{REG0, REG3}: []uint16{0xC3},
		AND{REG1, REG0}: []uint16{0xC4},
		AND{REG1, REG1}: []uint16{0xC5},
		AND{REG1, REG2}: []uint16{0xC6},
		AND{REG1, REG3}: []uint16{0xC7},
		AND{REG2, REG0}: []uint16{0xC8},
		AND{REG2, REG1}: []uint16{0xC9},
		AND{REG2, REG2}: []uint16{0xCA},
		AND{REG2, REG3}: []uint16{0xCB},
		AND{REG3, REG0}: []uint16{0xCC},
		AND{REG3, REG1}: []uint16{0xCD},
		AND{REG3, REG2}: []uint16{0xCE},
		AND{REG3, REG3}: []uint16{0xCF},

		OR{REG0, REG0}: []uint16{0xD0},
		OR{REG0, REG1}: []uint16{0xD1},
		OR{REG0, REG2}: []uint16{0xD2},
		OR{REG0, REG3}: []uint16{0xD3},
		OR{REG1, REG0}: []uint16{0xD4},
		OR{REG1, REG1}: []uint16{0xD5},
		OR{REG1, REG2}: []uint16{0xD6},
		OR{REG1, REG3}: []uint16{0xD7},
		OR{REG2, REG0}: []uint16{0xD8},
		OR{REG2, REG1}: []uint16{0xD9},
		OR{REG2, REG2}: []uint16{0xDA},
		OR{REG2, REG3}: []uint16{0xDB},
		OR{REG3, REG0}: []uint16{0xDC},
		OR{REG3, REG1}: []uint16{0xDD},
		OR{REG3, REG2}: []uint16{0xDE},
		OR{REG3, REG3}: []uint16{0xDF},

		XOR{REG0, REG0}: []uint16{0xE0},
		XOR{REG0, REG1}: []uint16{0xE1},
		XOR{REG0, REG2}: []uint16{0xE2},
		XOR{REG0, REG3}: []uint16{0xE3},
		XOR{REG1, REG0}: []uint16{0xE4},
		XOR{REG1, REG1}: []uint16{0xE5},
		XOR{REG1, REG2}: []uint16{0xE6},
		XOR{REG1, REG3}: []uint16{0xE7},
		XOR{REG2, REG0}: []uint16{0xE8},
		XOR{REG2, REG1}: []uint16{0xE9},
		XOR{REG2, REG2}: []uint16{0xEA},
		XOR{REG2, REG3}: []uint16{0xEB},
		XOR{REG3, REG0}: []uint16{0xEC},
		XOR{REG3, REG1}: []uint16{0xED},
		XOR{REG3, REG2}: []uint16{0xEE},
		XOR{REG3, REG3}: []uint16{0xEF},

		CMP{REG0, REG0}: []uint16{0xF0},
		CMP{REG0, REG1}: []uint16{0xF1},
		CMP{REG0, REG2}: []uint16{0xF2},
		CMP{REG0, REG3}: []uint16{0xF3},
		CMP{REG1, REG0}: []uint16{0xF4},
		CMP{REG1, REG1}: []uint16{0xF5},
		CMP{REG1, REG2}: []uint16{0xF6},
		CMP{REG1, REG3}: []uint16{0xF7},
		CMP{REG2, REG0}: []uint16{0xF8},
		CMP{REG2, REG1}: []uint16{0xF9},
		CMP{REG2, REG2}: []uint16{0xFA},
		CMP{REG2, REG3}: []uint16{0xFB},
		CMP{REG3, REG0}: []uint16{0xFC},
		CMP{REG3, REG1}: []uint16{0xFD},
		CMP{REG3, REG2}: []uint16{0xFE},
		CMP{REG3, REG3}: []uint16{0xFF},
	}

	for ins, expected := range TABLE {
		if emit, err := ins.Emit(nil, nil); err == nil {
			if reflect.DeepEqual(emit, expected) == false {
				t.Logf("Expected %v got %v when testing %s", expected, emit, ins)
				t.FailNow()
			}
		} else {
			t.Logf("Got error %v when testing %s", err, ins)
			t.FailNow()
		}
	}
}

func TestOneRegInstructionsString(t *testing.T) {
	var TABLE map[Instruction]string = map[Instruction]string{
		SHR{REG0}: "SHR R0",
		SHR{REG1}: "SHR R1",
		SHR{REG2}: "SHR R2",
		SHR{REG3}: "SHR R3",

		SHL{REG0}: "SHL R0",
		SHL{REG1}: "SHL R1",
		SHL{REG2}: "SHL R2",
		SHL{REG3}: "SHL R3",

		NOT{REG0}: "NOT R0",
		NOT{REG1}: "NOT R1",
		NOT{REG2}: "NOT R2",
		NOT{REG3}: "NOT R3",

		JR{REG0}: "JR R0",
		JR{REG1}: "JR R1",
		JR{REG2}: "JR R2",
		JR{REG3}: "JR R3",
	}

	for ins, expected := range TABLE {
		if ins.String() != expected {
			t.Logf("Expected %s got %s when testing %s", expected, ins.String(), ins)
			t.FailNow()
		}
	}
}

func TestOneRegInstructions(t *testing.T) {
	var TABLE map[Instruction][]uint16 = map[Instruction][]uint16{
		SHR{REG0}: []uint16{0x90},
		SHR{REG1}: []uint16{0x95},
		SHR{REG2}: []uint16{0x9A},
		SHR{REG3}: []uint16{0x9F},

		SHL{REG0}: []uint16{0xA0},
		SHL{REG1}: []uint16{0xA5},
		SHL{REG2}: []uint16{0xAA},
		SHL{REG3}: []uint16{0xAF},

		NOT{REG0}: []uint16{0xB0},
		NOT{REG1}: []uint16{0xB5},
		NOT{REG2}: []uint16{0xBA},
		NOT{REG3}: []uint16{0xBF},

		JR{REG0}: []uint16{0x30},
		JR{REG1}: []uint16{0x31},
		JR{REG2}: []uint16{0x32},
		JR{REG3}: []uint16{0x33},
	}

	for ins, expected := range TABLE {
		if emit, err := ins.Emit(nil, nil); err == nil {
			if reflect.DeepEqual(emit, expected) == false {
				t.Logf("Expected %v got %v when testing %s", expected, emit, ins)
				t.FailNow()
			}
		} else {
			t.Logf("Got error %v when testing %s", err, ins)
			t.FailNow()
		}
	}
}

func TestDATAInstructionsString(t *testing.T) {
	var TABLE map[Instruction]string = map[Instruction]string{
		DATA{REG0, NUMBER{0x0001}}: "DATA R0, 0x0001",
		DATA{REG1, NUMBER{0x0002}}: "DATA R1, 0x0002",
		DATA{REG2, NUMBER{0x0003}}: "DATA R2, 0x0003",
		DATA{REG3, NUMBER{0x0004}}: "DATA R3, 0x0004",
		DATA{REG0, SYMBOL{"aaa"}}:  "DATA R0, %aaa",
		DATA{REG1, SYMBOL{"bbb"}}:  "DATA R1, %bbb",
		DATA{REG2, SYMBOL{"ccc"}}:  "DATA R2, %ccc",
		DATA{REG3, SYMBOL{"ddd"}}:  "DATA R3, %ddd",
	}

	for ins, expected := range TABLE {
		if ins.String() != expected {
			t.Logf("Expected %s got %s when testing %s", expected, ins.String(), ins)
			t.FailNow()
		}
	}
}

func TestDATAInstruction(t *testing.T) {
	var TABLE map[Instruction][]uint16 = map[Instruction][]uint16{
		DATA{REG0, NUMBER{0x0001}}: []uint16{0x20, 0x0001},
		DATA{REG1, NUMBER{0x0002}}: []uint16{0x21, 0x0002},
		DATA{REG2, NUMBER{0x0003}}: []uint16{0x22, 0x0003},
		DATA{REG3, NUMBER{0x0004}}: []uint16{0x23, 0x0004},
		DATA{REG0, SYMBOL{"foo"}}:  []uint16{0x20, 0xA000},
		DATA{REG1, SYMBOL{"bar"}}:  []uint16{0x21, 0xB000},
		DATA{REG2, SYMBOL{"baz"}}:  []uint16{0x22, 0xC000},
		DATA{REG3, SYMBOL{"bee"}}:  []uint16{0x23, 0xD000},
	}

	dummySymbolResolver := func(s SYMBOL) (uint16, error) {
		switch s.Name {
		case "foo":
			return 0xA000, nil
		case "bar":
			return 0xB000, nil
		case "baz":
			return 0xC000, nil
		case "bee":
			return 0xD000, nil
		default:
			return 0x0000, fmt.Errorf("received unknown symbol")
		}
	}

	for ins, expected := range TABLE {
		if emit, err := ins.Emit(nil, dummySymbolResolver); err == nil {
			if reflect.DeepEqual(emit, expected) == false {
				t.Logf("Expected %v got %v when testing %s", expected, emit, ins)
				t.FailNow()
			}
		} else {
			t.Logf("Got error %v when testing %s", err, ins)
			t.FailNow()
		}
	}
}

func TestIOInstructionsString(t *testing.T) {
	var TABLE map[Instruction]string = map[Instruction]string{
		IN{DATA_MODE, REG0}: "IN Data, R0",
		IN{DATA_MODE, REG1}: "IN Data, R1",
		IN{DATA_MODE, REG2}: "IN Data, R2",
		IN{DATA_MODE, REG3}: "IN Data, R3",

		OUT{DATA_MODE, REG0}: "OUT Data, R0",
		OUT{DATA_MODE, REG1}: "OUT Data, R1",
		OUT{DATA_MODE, REG2}: "OUT Data, R2",
		OUT{DATA_MODE, REG3}: "OUT Data, R3",

		IN{ADDRESS_MODE, REG0}: "IN Addr, R0",
		IN{ADDRESS_MODE, REG1}: "IN Addr, R1",
		IN{ADDRESS_MODE, REG2}: "IN Addr, R2",
		IN{ADDRESS_MODE, REG3}: "IN Addr, R3",

		OUT{ADDRESS_MODE, REG0}: "OUT Addr, R0",
		OUT{ADDRESS_MODE, REG1}: "OUT Addr, R1",
		OUT{ADDRESS_MODE, REG2}: "OUT Addr, R2",
		OUT{ADDRESS_MODE, REG3}: "OUT Addr, R3",
	}

	for ins, expected := range TABLE {
		if ins.String() != expected {
			t.Logf("Expected %s got %s when testing %s", expected, ins.String(), ins)
			t.FailNow()
		}
	}
}
func TestIOInstructions(t *testing.T) {
	var TABLE map[Instruction][]uint16 = map[Instruction][]uint16{
		IN{DATA_MODE, REG0}:    []uint16{0x70},
		IN{DATA_MODE, REG1}:    []uint16{0x71},
		IN{DATA_MODE, REG2}:    []uint16{0x72},
		IN{DATA_MODE, REG3}:    []uint16{0x73},
		IN{ADDRESS_MODE, REG0}: []uint16{0x74},
		IN{ADDRESS_MODE, REG1}: []uint16{0x75},
		IN{ADDRESS_MODE, REG2}: []uint16{0x76},
		IN{ADDRESS_MODE, REG3}: []uint16{0x77},

		OUT{DATA_MODE, REG0}:    []uint16{0x78},
		OUT{DATA_MODE, REG1}:    []uint16{0x79},
		OUT{DATA_MODE, REG2}:    []uint16{0x7A},
		OUT{DATA_MODE, REG3}:    []uint16{0x7B},
		OUT{ADDRESS_MODE, REG0}: []uint16{0x7C},
		OUT{ADDRESS_MODE, REG1}: []uint16{0x7D},
		OUT{ADDRESS_MODE, REG2}: []uint16{0x7E},
		OUT{ADDRESS_MODE, REG3}: []uint16{0x7F},
	}

	for ins, expected := range TABLE {
		if emit, err := ins.Emit(nil, nil); err == nil {
			if reflect.DeepEqual(emit, expected) == false {
				t.Logf("Expected %v got %v when testing %s", expected, emit, ins)
				t.FailNow()
			}
		} else {
			t.Logf("Got error %v when testing %s", err, ins)
			t.FailNow()
		}
	}
}

func TestJMPInstructionString(t *testing.T) {
	var TABLE map[Instruction]string = map[Instruction]string{
		JMP{LABEL{"foo"}}: "JMP foo",
		JMP{LABEL{"bar"}}: "JMP bar",
	}

	for ins, expected := range TABLE {
		if ins.String() != expected {
			t.Logf("Expected %s got %s when testing %s", expected, ins.String(), ins)
			t.FailNow()
		}
	}
}

func TestJMPInstruction(t *testing.T) {
	var TABLE map[Instruction][]uint16 = map[Instruction][]uint16{
		JMP{LABEL{"foo"}}: []uint16{0x40, 0x0001},
		JMP{LABEL{"bar"}}: []uint16{0x40, 0x0002},
	}

	dummyLabelResolver := func(l LABEL) (uint16, error) {
		if l.Name == "foo" {
			return 0x0001, nil
		} else if l.Name == "bar" {
			return 0x0002, nil
		}

		return 0x0000, fmt.Errorf("received unknown label")
	}

	for ins, expected := range TABLE {
		if emit, err := ins.Emit(dummyLabelResolver, nil); err == nil {
			if reflect.DeepEqual(emit, expected) == false {
				t.Logf("Expected %v got %v when testing %s", expected, emit, ins)
				t.FailNow()
			}
		} else {
			t.Logf("Got error %v when testing %s", err, ins)
			t.FailNow()
		}
	}
}

func TestCLFInstructionString(t *testing.T) {
	var TABLE map[Instruction]string = map[Instruction]string{
		CLF{}: "CLF",
	}

	for ins, expected := range TABLE {
		if ins.String() != expected {
			t.Logf("Expected %s got %s when testing %s", expected, ins.String(), ins)
			t.FailNow()
		}
	}
}

func TestCALLInstructionString(t *testing.T) {
	var TABLE map[Instruction]string = map[Instruction]string{
		CALL{LABEL{"foo"}}: "CALL foo",
		CALL{LABEL{"bar"}}: "CALL bar",
	}

	for ins, expected := range TABLE {
		if ins.String() != expected {
			t.Logf("Expected %s got %s when testing %s", expected, ins.String(), ins)
			t.FailNow()
		}
	}
}
func TestCALLInstruction(t *testing.T) {
	var TABLE map[Instruction][]uint16 = map[Instruction][]uint16{
		CALL{LABEL{"foo"}}: []uint16{0x23, 0x1234, 0x40, 0x0001},
		CALL{LABEL{"bar"}}: []uint16{0x23, 0x1234, 0x40, 0x0002},
	}
	dummyLabelResolver := func(l LABEL) (uint16, error) {
		if l.Name == "foo" {
			return 0x0001, nil
		} else if l.Name == "bar" {
			return 0x0002, nil
		}

		return 0x0000, fmt.Errorf("received unknown label")
	}

	dummySymbolResolver := func(s SYMBOL) (uint16, error) {
		if s.Name == "NEXTINSTRUCTION" {
			return 0x1234, nil
		}

		return 0x0000, fmt.Errorf("received unknown label")
	}


	for ins, expected := range TABLE {
		if emit, err := ins.Emit(dummyLabelResolver, dummySymbolResolver); err == nil {
			if reflect.DeepEqual(emit, expected) == false {
				t.Logf("Expected %v got %v when testing %s", expected, emit, ins)
				t.FailNow()
			}
		} else {
			t.Logf("Got error %v when testing %s", err, ins)
			t.FailNow()
		}
	}
}
func TestCLFInstruction(t *testing.T) {
	var TABLE map[Instruction][]uint16 = map[Instruction][]uint16{
		CLF{}: []uint16{0x60},
	}

	for ins, expected := range TABLE {
		if emit, err := ins.Emit(nil, nil); err == nil {
			if reflect.DeepEqual(emit, expected) == false {
				t.Logf("Expected %v got %v when testing %s", expected, emit, ins)
				t.FailNow()
			}
		} else {
			t.Logf("Got error %v when testing %s", err, ins)
			t.FailNow()
		}
	}
}

func TestJMPFlagInstructionsString(t *testing.T) {
	var instructions []Instruction = []Instruction{
		JMPF{[]string{"Z"}, LABEL{"foo1"}},
		JMPF{[]string{"E"}, LABEL{"foo2"}},
		JMPF{[]string{"E", "Z"}, LABEL{"foo3"}},
		JMPF{[]string{"A"}, LABEL{"foo4"}},
		JMPF{[]string{"A", "Z"}, LABEL{"foo5"}},
		JMPF{[]string{"A", "E"}, LABEL{"foo6"}},
		JMPF{[]string{"A", "E", "Z"}, LABEL{"foo7"}},
		JMPF{[]string{"C"}, LABEL{"foo8"}},
		JMPF{[]string{"C", "Z"}, LABEL{"foo9"}},
		JMPF{[]string{"C", "E"}, LABEL{"foo10"}},
		JMPF{[]string{"C", "E", "Z"}, LABEL{"foo11"}},
		JMPF{[]string{"C", "A"}, LABEL{"foo12"}},
		JMPF{[]string{"C", "A", "Z"}, LABEL{"foo13"}},
		JMPF{[]string{"C", "A", "E"}, LABEL{"foo14"}},
		JMPF{[]string{"C", "A", "E", "Z"}, LABEL{"foo15"}},
	}

	var expected []string = []string{
		"JMPZ foo1",
		"JMPE foo2",
		"JMPEZ foo3",
		"JMPA foo4",
		"JMPAZ foo5",
		"JMPAE foo6",
		"JMPAEZ foo7",
		"JMPC foo8",
		"JMPCZ foo9",
		"JMPCE foo10",
		"JMPCEZ foo11",
		"JMPCA foo12",
		"JMPCAZ foo13",
		"JMPCAE foo14",
		"JMPCAEZ foo15",
	}

	for i, ins := range instructions {
		if ins.String() != expected[i] {
			t.Logf("Expected %s got %s when testing %s", expected[i], ins.String(), ins)
			t.FailNow()
		}
	}
}

func TestJMPFlagInstructions(t *testing.T) {
	var instructions []Instruction = []Instruction{
		JMPF{[]string{"Z"}, LABEL{"foo"}},
		JMPF{[]string{"E"}, LABEL{"foo"}},
		JMPF{[]string{"E", "Z"}, LABEL{"foo"}},
		JMPF{[]string{"A"}, LABEL{"foo"}},
		JMPF{[]string{"A", "Z"}, LABEL{"foo"}},
		JMPF{[]string{"A", "E"}, LABEL{"foo"}},
		JMPF{[]string{"A", "E", "Z"}, LABEL{"foo"}},
		JMPF{[]string{"C"}, LABEL{"foo"}},
		JMPF{[]string{"C", "Z"}, LABEL{"foo"}},
		JMPF{[]string{"C", "E"}, LABEL{"foo"}},
		JMPF{[]string{"C", "E", "Z"}, LABEL{"foo"}},
		JMPF{[]string{"C", "A"}, LABEL{"foo"}},
		JMPF{[]string{"C", "A", "Z"}, LABEL{"foo"}},
		JMPF{[]string{"C", "A", "E"}, LABEL{"foo"}},
		JMPF{[]string{"C", "A", "E", "Z"}, LABEL{"foo"}},
	}

	var expected [][]uint16 = [][]uint16{
		[]uint16{0x0051, 0x1234},
		[]uint16{0x0052, 0x1234},
		[]uint16{0x0053, 0x1234},
		[]uint16{0x0054, 0x1234},
		[]uint16{0x0055, 0x1234},
		[]uint16{0x0056, 0x1234},
		[]uint16{0x0057, 0x1234},
		[]uint16{0x0058, 0x1234},
		[]uint16{0x0059, 0x1234},
		[]uint16{0x005A, 0x1234},
		[]uint16{0x005B, 0x1234},
		[]uint16{0x005C, 0x1234},
		[]uint16{0x005D, 0x1234},
		[]uint16{0x005E, 0x1234},
		[]uint16{0x005F, 0x1234},
	}

	dummyLabelResolver := func(l LABEL) (uint16, error) {
		if l.Name == "foo" {
			return 0x1234, nil
		}
		return 0x0000, fmt.Errorf("received unknown label")
	}

	for i, ins := range instructions {
		if emit, err := ins.Emit(dummyLabelResolver, nil); err == nil {
			if reflect.DeepEqual(emit, expected[i]) == false {
				t.Logf("Expected %v got %v when testing %s", expected[i], emit, ins)
				t.FailNow()
			}
		} else {
			t.Logf("Got error %v when testing %s", err, ins)
			t.FailNow()
		}
	}
}
