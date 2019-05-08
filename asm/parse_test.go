package asm

import (
	"reflect"
	"strings"
	"testing"
)

func TestParseLabel(t *testing.T) {
	p := Parser{}

	result, err := p.Parse(strings.NewReader(`
	start:
	end:
	`))

	expected := []Instruction{LABEL{"start"}, LABEL{"end"}}

	if err != nil {
		t.FailNow()
	}

	if result[0] != expected[0] {
		t.FailNow()
	}

	if result[1] != expected[1] {
		t.FailNow()
	}
}

func TestParseADD(t *testing.T) {
	input := `
		ADD R0, R1
		ADD R1,R0
		ADD   R2,   R3
	`

	expected := []Instruction{ADD{REG0, REG1}, ADD{REG1, REG0}, ADD{REG2, REG3}}

	testParseInstructions(input, expected, t)

}

func TestParseLD(t *testing.T) {
	input := `
		LD R0, R1
		LD R1,R0
		LD    R2,   R3
	`

	expected := []Instruction{LOAD{REG0, REG1}, LOAD{REG1, REG0}, LOAD{REG2, REG3}}

	testParseInstructions(input, expected, t)
}

func TestParseST(t *testing.T) {
	input := `
		ST R0, R1
		ST R1,R0
		ST    R2,   R3
	`

	expected := []Instruction{STORE{REG0, REG1}, STORE{REG1, REG0}, STORE{REG2, REG3}}

	testParseInstructions(input, expected, t)
}

func TestParseXOR(t *testing.T) {
	input := `
		XOR R0, R1
		XOR R1,R0
		XOR    R2,   R3
	`

	expected := []Instruction{XOR{REG0, REG1}, XOR{REG1, REG0}, XOR{REG2, REG3}}

	testParseInstructions(input, expected, t)
}

func TestParseCMP(t *testing.T) {
	input := `
		CMP R0, R1
		CMP R1,R0
		CMP    R2,   R3
	`

	expected := []Instruction{CMP{REG0, REG1}, CMP{REG1, REG0}, CMP{REG2, REG3}}

	testParseInstructions(input, expected, t)
}

func TestParseOR(t *testing.T) {
	input := `
		OR R0, R1
		OR R1,R0
		OR    R2,   R3
	`

	expected := []Instruction{OR{REG0, REG1}, OR{REG1, REG0}, OR{REG2, REG3}}

	testParseInstructions(input, expected, t)
}

func TestParseAND(t *testing.T) {
	input := `
		AND R0, R1
		AND R1,R0
		AND    R2,   R3
	`

	expected := []Instruction{AND{REG0, REG1}, AND{REG1, REG0}, AND{REG2, REG3}}

	testParseInstructions(input, expected, t)
}

func TestParseCLF(t *testing.T) {
	input := `
	CLF           
	`

	expected := []Instruction{CLF{}}

	testParseInstructions(input, expected, t)
}

func TestParseShifts(t *testing.T) {
	input := `
	SHL R0
	SHL   R1
	SHL R2
	SHL  R3
	SHR R0
	SHR   R1
	SHR R2
	SHR  R3
	`

	expected := []Instruction{
		SHL{REG0},
		SHL{REG1},
		SHL{REG2},
		SHL{REG3},
		SHR{REG0},
		SHR{REG1},
		SHR{REG2},
		SHR{REG3},
	}

	testParseInstructions(input, expected, t)
}

func TestParseNOT(t *testing.T) {
	input := `
	NOT R0
	NOT   R1
	NOT R2
	NOT  R3
	`

	expected := []Instruction{
		NOT{REG0},
		NOT{REG1},
		NOT{REG2},
		NOT{REG3},
	}

	testParseInstructions(input, expected, t)
}

func TestParseJR(t *testing.T) {
	input := `
	JR R0
	JR  R1
	JR    R2
	JR R3
	`

	expected := []Instruction{
		JR{REG0},
		JR{REG1},
		JR{REG2},
		JR{REG3},
	}

	testParseInstructions(input, expected, t)
}

func TestParseDATA(t *testing.T) {
	input := `
		DATA R0, 0xFE85
		DATA R1,0x2963
		DATA R2,    0x000A
		DATA R3,19
		DATA R0, 0xf9fe
	`

	expected := []Instruction{DATA{REG0, 0xFE85}, DATA{REG1, 0x2963}, DATA{REG2, 0x000A}, DATA{REG3, 19}, DATA{REG0, 0xF9FE}}

	testParseInstructions(input, expected, t)
}

func TestParseJMP(t *testing.T) {
	input := `
		JMP mylabel
		JMP   foo
		JMP bar
	`

	expected := []Instruction{
		JMP{LABEL{"mylabel"}},
		JMP{LABEL{"foo"}},
		JMP{LABEL{"bar"}},
	}

	testParseInstructions(input, expected, t)
}

// JMP(CAEZ)
// set instruction address register to next byte (2 byte instruction)
// jump if <flag(s)> are true
// ----------------------
// 0x0051 = JMPZ <value>
// 0x0052 = JMPE <value>
// 0x0053 = JMPEZ <value>
// 0x0054 = JMPA <value>
// 0x0055 = JMPAZ <value>
// 0x0056 = JMPAE <value>
// 0x0057 = JMPAEZ <value>
// 0x0058 = JMPC <value>
// 0x0059 = JMPCZ <value>
// 0x005A = JMPCE <value>
// 0x005B = JMPCEZ <value>
// 0x005C = JMPCA <value>
// 0x005D = JMPCAZ <value>
// 0x005E = JMPCAE <value>
// 0x005F = JMPCAEZ <value>

func TestParseFlagJMP(t *testing.T) {
	input := `
		JMPZ  flagZ
		JMPE  flagE
		JMPEZ   flagEz
		JMPA  flagA
		JMPAZ   flagAz
		JMPAE   flagAe
		JMPAEZ   flagAez
		JMPC  flagC
		JMPCZ  flagCz
		JMPCE  flagCe
		JMPCEZ  flagCez
		JMPCA  flagCa
		JMPCAZ  flagCaz
		JMPCAE  flagCae
		JMPCAEZ  flagCaez
	`

	expected := []Instruction{
		JMPF{[]string{"Z"}, LABEL{"flagZ"}},
		JMPF{[]string{"E"}, LABEL{"flagE"}},
		JMPF{[]string{"E", "Z"}, LABEL{"flagEz"}},
		JMPF{[]string{"A"}, LABEL{"flagA"}},
		JMPF{[]string{"A", "Z"}, LABEL{"flagAz"}},
		JMPF{[]string{"A", "E"}, LABEL{"flagAe"}},
		JMPF{[]string{"A", "E", "Z"}, LABEL{"flagAez"}},
		JMPF{[]string{"C"}, LABEL{"flagC"}},
		JMPF{[]string{"C", "Z"}, LABEL{"flagCz"}},
		JMPF{[]string{"C", "E"}, LABEL{"flagCe"}},
		JMPF{[]string{"C", "E", "Z"}, LABEL{"flagCez"}},
		JMPF{[]string{"C", "A"}, LABEL{"flagCa"}},
		JMPF{[]string{"C", "A", "Z"}, LABEL{"flagCaz"}},
		JMPF{[]string{"C", "A", "E"}, LABEL{"flagCae"}},
		JMPF{[]string{"C", "A", "E", "Z"}, LABEL{"flagCaez"}},
	}

	testParseInstructions(input, expected, t)
}

func TestParseIO(t *testing.T) {
	input := `
		OUT Addr, R0
		OUT Addr,   R1
		OUT Addr,R2
		OUT  Addr, R3
		OUT Data, R0
		OUT Data,   R1
		OUT Data,R2
		OUT Data, R3
		IN Addr, R0
		IN Addr,   R1
		IN Addr,R2
		IN  Addr, R3
		IN Data, R0
		IN Data,   R1
		IN Data,R2
		IN Data, R3
	`

	expected := []Instruction{
		OUT{ADDRESS_MODE, REG0},
		OUT{ADDRESS_MODE, REG1},
		OUT{ADDRESS_MODE, REG2},
		OUT{ADDRESS_MODE, REG3},
		OUT{DATA_MODE, REG0},
		OUT{DATA_MODE, REG1},
		OUT{DATA_MODE, REG2},
		OUT{DATA_MODE, REG3},
		IN{ADDRESS_MODE, REG0},
		IN{ADDRESS_MODE, REG1},
		IN{ADDRESS_MODE, REG2},
		IN{ADDRESS_MODE, REG3},
		IN{DATA_MODE, REG0},
		IN{DATA_MODE, REG1},
		IN{DATA_MODE, REG2},
		IN{DATA_MODE, REG3},
	}

	testParseInstructions(input, expected, t)
}

func TestParseSimpleProgram(t *testing.T) {
	input := `
	loop:
		DATA R0, 0x1235
		ADD R1, R0
		CMP R0, R1
	JMPE foo 
	JMP loop
	foo:
		AND R1, R0
		OR R1, R2
	`

	expected := []Instruction{
		LABEL{"loop"},
		DATA{REG0, 0x1235},
		ADD{REG1, REG0},
		CMP{REG0, REG1},
		JMPF{[]string{"E"}, LABEL{"foo"}},
		JMP{LABEL{"loop"}},
		LABEL{"foo"},
		AND{REG1, REG0},
		OR{REG1, REG2},
	}

	testParseInstructions(input, expected, t)
}

func testParseInstructions(input string, expected []Instruction, t *testing.T) {
	p := Parser{}

	result, err := p.Parse(strings.NewReader(input))

	if err != nil {
		t.Logf("encountered error %v", err)
		t.FailNow()
	}

	if len(result) != len(expected) {
		t.Logf("expected %d instructions but got %d", len(expected), len(result))
		t.FailNow()
	}

	for i := range result {
		if reflect.DeepEqual(result[i], expected[i]) == false {
			t.Logf("expected instruction %d to be %v but got %v", i, expected[i], result[i])
			t.FailNow()
		}
	}
}
