package cpu

import (
	"testing"

	"github.com/djhworld/simple-computer/components"
	"github.com/djhworld/simple-computer/memory"
)

func TestIARIncrementedOnEveryCycle(t *testing.T) {
	b := components.NewBus()
	m := memory.NewMemory256(b)
	c := NewCPU(b, m)

	setIAR(c, 0x00)

	var q byte
	for q = 0; q < 0xFF; q++ {
		doFetchDecodeExecute(c)
		checkIAR(c, q+1, t)
	}
}

func TestInstructionReceivedFromMemory(t *testing.T) {
	b := components.NewBus()
	m := memory.NewMemory256(b)
	c := NewCPU(b, m)

	instructions := []byte{0x8A, 0x82, 0x88, 0x94, 0xB1}

	var addr byte = 0x0F
	for _, b := range instructions {
		setMemoryLocation(c, addr, b)
		addr++
	}

	setIAR(c, 0x0F)

	for _, b := range instructions {
		doFetchDecodeExecute(c)
		checkIR(c, b, t)
	}
}

func TestFlagsRegisterAllFalse(t *testing.T) {
	b := components.NewBus()
	m := memory.NewMemory256(b)
	c := NewCPU(b, m)

	setMemoryLocation(c, 0x00, 0x81)
	setRegisters(c, [4]byte{0x09, 0x0A, 0x02, 0x03})
	setIAR(c, 0x00)

	doFetchDecodeExecute(c)

	checkFlagsRegister(c, false, false, false, false, t)
}

func TestFlagsRegisterCarryFlagEnabled(t *testing.T) {
	b := components.NewBus()
	m := memory.NewMemory256(b)
	c := NewCPU(b, m)

	setMemoryLocation(c, 0x00, 0x81)
	setRegisters(c, [4]byte{0x20, 0xFF, 0x02, 0x03})
	setIAR(c, 0x00)

	doFetchDecodeExecute(c)

	checkFlagsRegister(c, true, false, false, false, t)
}

func TestFlagsRegisterIsLargerFlagEnabled(t *testing.T) {
	b := components.NewBus()
	m := memory.NewMemory256(b)
	c := NewCPU(b, m)

	setMemoryLocation(c, 0x00, 0x81)
	setRegisters(c, [4]byte{0x21, 0x20, 0x02, 0x03})
	setIAR(c, 0x00)

	doFetchDecodeExecute(c)

	checkFlagsRegister(c, false, true, false, false, t)
}

func TestFlagsRegisterIsEqualsFlagEnabled(t *testing.T) {
	b := components.NewBus()
	m := memory.NewMemory256(b)
	c := NewCPU(b, m)

	setMemoryLocation(c, 0x00, 0x81)
	setRegisters(c, [4]byte{0x21, 0x21, 0x02, 0x03})
	setIAR(c, 0x00)

	doFetchDecodeExecute(c)

	checkFlagsRegister(c, false, false, true, false, t)
}

func TestFlagsRegisterIsZeroFlagEnabled(t *testing.T) {
	b := components.NewBus()
	m := memory.NewMemory256(b)
	c := NewCPU(b, m)

	setMemoryLocation(c, 0x00, 0x81)
	setRegisters(c, [4]byte{0x01, 0xFF, 0x02, 0x03})
	setIAR(c, 0x00)

	doFetchDecodeExecute(c)

	checkFlagsRegister(c, true, false, false, true, t)
}

func TestFlagsRegisterMultipleEnabled(t *testing.T) {
	b := components.NewBus()
	m := memory.NewMemory256(b)
	c := NewCPU(b, m)

	setMemoryLocation(c, 0x00, 0x81)
	setRegisters(c, [4]byte{0xFF, 0x01, 0x02, 0x03})
	setIAR(c, 0x00)

	doFetchDecodeExecute(c)

	checkFlagsRegister(c, true, true, false, true, t)
}
func TestSTThenLD(t *testing.T) {
	b := components.NewBus()
	m := memory.NewMemory256(b)
	c := NewCPU(b, m)

	for i := byte(0); i < 128; i++ {
		setMemoryLocation(c, i, 0x1B)
	}
	setIAR(c, 0x00)

	var value byte = 0xFE
	for i := 0x80; i <= 0xFF; i++ {
		setRegisters(c, [4]byte{0x01, 0x01, byte(i), value})
		doFetchDecodeExecute(c)
		value--
	}

	for i := byte(0); i < 128; i++ {
		setMemoryLocation(c, i, 0x0B)
	}
	setIAR(c, 0x00)

	value = 0xFE
	for i := 0x80; i <= 0xFF; i++ {
		setRegisters(c, [4]byte{0x01, 0x01, byte(i), 0x01})
		doFetchDecodeExecute(c)
		checkRegisters(c, 0x01, 0x01, byte(i), value, t)
		value--
	}
}

func TestLD4Times(t *testing.T) {
	b := components.NewBus()
	m := memory.NewMemory256(b)
	c := NewCPU(b, m)

	var addr byte = 0xA2
	values := []byte{0x88, 0x90, 0x92, 0xAB}

	for i := byte(0); i < byte(len(values)); i++ {
		setMemoryLocation(c, i, 0x01)
		setMemoryLocation(c, addr, values[i])
		addr++
	}

	setIAR(c, 0x00)

	addr = 0xA2
	for _, v := range values {
		setRegister(c, 0, addr)
		setRegister(c, 1, 0x01)
		setRegister(c, 2, 0x01)
		setRegister(c, 3, 0x01)

		doFetchDecodeExecute(c)
		checkRegisters(c, addr, v, 0x01, 0x01, t)
		addr++
	}
}

func TestLD(t *testing.T) {
	testLD(0x00, 0x80, 0x23, []byte{0x80, 0x81, 0x82, 0x83}, []byte{0x23, 0x81, 0x82, 0x83}, t)
	testLD(0x01, 0x84, 0xF2, []byte{0x84, 0x85, 0x86, 0x87}, []byte{0x84, 0xF2, 0x86, 0x87}, t)
	testLD(0x02, 0x88, 0x01, []byte{0x88, 0x89, 0x8A, 0x8B}, []byte{0x88, 0x89, 0x01, 0x8B}, t)
	testLD(0x03, 0x8C, 0x5A, []byte{0x8C, 0x8D, 0x8E, 0x8F}, []byte{0x8C, 0x8D, 0x8E, 0x5A}, t)

	testLD(0x04, 0x91, 0x23, []byte{0x90, 0x91, 0x92, 0x93}, []byte{0x23, 0x91, 0x92, 0x93}, t)
	testLD(0x05, 0x95, 0xF2, []byte{0x94, 0x95, 0x96, 0x97}, []byte{0x94, 0xF2, 0x96, 0x97}, t)
	testLD(0x06, 0x99, 0x01, []byte{0x98, 0x99, 0x9A, 0x9B}, []byte{0x98, 0x99, 0x01, 0x9B}, t)
	testLD(0x07, 0x9D, 0x5A, []byte{0x9C, 0x9D, 0x9E, 0x9F}, []byte{0x9C, 0x9D, 0x9E, 0x5A}, t)

	testLD(0x08, 0xA2, 0x23, []byte{0xA0, 0xA1, 0xA2, 0xA3}, []byte{0x23, 0xA1, 0xA2, 0xA3}, t)
	testLD(0x09, 0xA6, 0xF2, []byte{0xA4, 0xA5, 0xA6, 0xA7}, []byte{0xA4, 0xF2, 0xA6, 0xA7}, t)
	testLD(0x0A, 0xAA, 0x01, []byte{0xA8, 0xA9, 0xAA, 0xAB}, []byte{0xA8, 0xA9, 0x01, 0xAB}, t)
	testLD(0x0B, 0xAE, 0x5A, []byte{0xAC, 0xAD, 0xAE, 0xAF}, []byte{0xAC, 0xAD, 0xAE, 0x5A}, t)

	testLD(0x0C, 0xB3, 0x23, []byte{0xB0, 0xB1, 0xB2, 0xB3}, []byte{0x23, 0xB1, 0xB2, 0xB3}, t)
	testLD(0x0D, 0xB7, 0xF2, []byte{0xB4, 0xB5, 0xB6, 0xB7}, []byte{0xB4, 0xF2, 0xB6, 0xB7}, t)
	testLD(0x0E, 0xBB, 0x01, []byte{0xB8, 0xB9, 0xBA, 0xBB}, []byte{0xB8, 0xB9, 0x01, 0xBB}, t)
	testLD(0x0F, 0xBF, 0x5A, []byte{0xBC, 0xBD, 0xBE, 0xBF}, []byte{0xBC, 0xBD, 0xBE, 0x5A}, t)
}

func testLD(instruction byte, memAddress, memValue byte, inputRegisters, expectedOutputRegisters []byte, t *testing.T) {
	b := components.NewBus()
	m := memory.NewMemory256(b)
	c := NewCPU(b, m)

	var insAddr byte = 0x00
	setMemoryLocation(c, insAddr, instruction)
	setIAR(c, insAddr)

	setMemoryLocation(c, memAddress, memValue)

	for i, v := range inputRegisters {
		setRegister(c, i, v)
	}

	doFetchDecodeExecute(c)
	checkRegisters(c, expectedOutputRegisters[0], expectedOutputRegisters[1], expectedOutputRegisters[2], expectedOutputRegisters[3], t)
}

func TestST(t *testing.T) {
	testST(0x10, [4]byte{0xA0, 0x01, 0x01, 0x01}, 0xA0, 0xA0, t)
	testST(0x11, [4]byte{0xA1, 0x29, 0x01, 0x01}, 0xA1, 0x29, t)
	testST(0x12, [4]byte{0xA2, 0x01, 0x7F, 0x01}, 0xA2, 0x7F, t)
	testST(0x13, [4]byte{0xA3, 0x01, 0x01, 0x1B}, 0xA3, 0x1B, t)

	testST(0x14, [4]byte{0xA0, 0xB4, 0x01, 0x01}, 0xB4, 0xA0, t)
	testST(0x15, [4]byte{0x01, 0xB5, 0x01, 0x01}, 0xB5, 0xB5, t)
	testST(0x16, [4]byte{0x01, 0xB6, 0x7F, 0x01}, 0xB6, 0x7F, t)
	testST(0x17, [4]byte{0x01, 0xB7, 0x01, 0x1B}, 0xB7, 0x1B, t)

	testST(0x18, [4]byte{0xA0, 0x01, 0xC8, 0x01}, 0xC8, 0xA0, t)
	testST(0x19, [4]byte{0x01, 0x29, 0xC9, 0x01}, 0xC9, 0x29, t)
	testST(0x1A, [4]byte{0x01, 0x01, 0xCA, 0x01}, 0xCA, 0xCA, t)
	testST(0x1B, [4]byte{0x01, 0x01, 0xCB, 0x1B}, 0xCB, 0x1B, t)

	testST(0x1C, [4]byte{0xA0, 0x01, 0x01, 0xDC}, 0xDC, 0xA0, t)
	testST(0x1D, [4]byte{0x01, 0x29, 0x01, 0xDD}, 0xDD, 0x29, t)
	testST(0x1E, [4]byte{0x01, 0x01, 0x7F, 0xDE}, 0xDE, 0x7F, t)
	testST(0x1F, [4]byte{0x01, 0x01, 0x01, 0xDF}, 0xDF, 0xDF, t)
}

func testST(instruction byte, inputRegisters [4]byte, expectedValueAddress, expectedValue byte, t *testing.T) {
	b := components.NewBus()
	m := memory.NewMemory256(b)
	c := NewCPU(b, m)

	// ST value into memory
	var insAddr byte = 0x00
	setMemoryLocation(c, insAddr, instruction)
	setIAR(c, insAddr)

	setRegisters(c, inputRegisters)

	doFetchDecodeExecute(c)

	//LD value into register zero
	setMemoryLocation(c, insAddr+1, 0x00)
	setIAR(c, insAddr+1)

	setRegisters(c, [4]byte{expectedValueAddress, inputRegisters[1], inputRegisters[2], inputRegisters[3]})

	doFetchDecodeExecute(c)

	if c.gpReg0.Value() != expectedValue {
		t.Logf("Expected register 0 to have value of: %X but got %X", expectedValue, c.gpReg0.Value())
		t.FailNow()
	}
}

func TestDATA(t *testing.T) {
	b := components.NewBus()
	m := memory.NewMemory256(b)
	c := NewCPU(b, m)

	var insAddr byte = 0x00

	// DATA R0
	setMemoryLocation(c, insAddr, 0x20)
	setMemoryLocation(c, insAddr+1, 0x71)

	// DATA R1
	setMemoryLocation(c, insAddr+2, 0x21)
	setMemoryLocation(c, insAddr+3, 0x72)

	// DATA R2
	setMemoryLocation(c, insAddr+4, 0x22)
	setMemoryLocation(c, insAddr+5, 0x73)

	// DATA R3
	setMemoryLocation(c, insAddr+6, 0x23)
	setMemoryLocation(c, insAddr+7, 0x74)

	setIAR(c, insAddr)

	setRegisters(c, [4]byte{0x01, 0x01, 0x01, 0x01})

	for i := 0; i < 4; i++ {
		doFetchDecodeExecute(c)
	}

	checkRegisters(c, 0x71, 0x72, 0x73, 0x74, t)

	// check IAR has incremented 2 each time
	checkIAR(c, 0x08, t)
}

func TestJMPR(t *testing.T) {
	testJMPR(0x30, [4]byte{0x83, 0x01, 0x01, 0x01}, 0x83, t)
	testJMPR(0x31, [4]byte{0x01, 0xF1, 0x01, 0x01}, 0xF1, t)
	testJMPR(0x32, [4]byte{0x01, 0x01, 0xBB, 0x01}, 0xBB, t)
	testJMPR(0x33, [4]byte{0x01, 0x01, 0x01, 0x19}, 0x19, t)
}

func testJMPR(instruction byte, inputRegisters [4]byte, expectedIAR byte, t *testing.T) {
	b := components.NewBus()
	m := memory.NewMemory256(b)
	c := NewCPU(b, m)

	var insAddr byte = 0x00

	// JMPR
	setMemoryLocation(c, insAddr, instruction)

	setIAR(c, insAddr)

	setRegisters(c, inputRegisters)

	doFetchDecodeExecute(c)

	// registers shouldn't change
	checkRegisters(c, inputRegisters[0], inputRegisters[1], inputRegisters[2], inputRegisters[3], t)

	// check IAR has jumped to the new location
	checkIAR(c, expectedIAR, t)
}

func TestJMP(t *testing.T) {
	for i := 0; i < 0xFF; i++ {
		testJMP(byte(i), t)
	}
}

func testJMP(expectedIAR byte, t *testing.T) {
	b := components.NewBus()
	m := memory.NewMemory256(b)
	c := NewCPU(b, m)

	var insAddr byte = 0x00

	// JMP
	setMemoryLocation(c, insAddr, 0x40)
	setMemoryLocation(c, insAddr+1, expectedIAR)

	setIAR(c, insAddr)

	inputRegisters := [4]byte{0x01, 0x01, 0x01, 0x01}
	setRegisters(c, inputRegisters)

	doFetchDecodeExecute(c)

	// registers shouldn't change
	checkRegisters(c, inputRegisters[0], inputRegisters[1], inputRegisters[2], inputRegisters[3], t)

	// check IAR has jumped to the new location
	checkIAR(c, expectedIAR, t)
}

func TestJMPC(t *testing.T) {
	testJMPConditional(0x58, 0x90, 0x81, [4]byte{0x04, 0xFF, 0x01, 0x2}, 0x90, t)
	// should not jump in false case
	testJMPConditional(0x58, 0x91, 0x81, [4]byte{0x05, 0x06, 0x01, 0x01}, 0x03, t)
}

func TestJMPA(t *testing.T) {
	testJMPConditional(0x54, 0x20, 0xF1, [4]byte{0x02, 0x01, 0x01, 0x2}, 0x20, t)
	// should not jump in false case
	testJMPConditional(0x54, 0x21, 0xF1, [4]byte{0x01, 0x03, 0x01, 0x01}, 0x03, t)
}

func TestJMPE(t *testing.T) {
	testJMPConditional(0x52, 0xAE, 0xF1, [4]byte{0x00, 0x00, 0x01, 0x2}, 0xAE, t)
	// should not jump in false case
	testJMPConditional(0x52, 0xAF, 0xF1, [4]byte{0x10, 0x11, 0x01, 0x01}, 0x03, t)
}
func TestJMPZ(t *testing.T) {
	// perform NOT on R0 (0xFF) to trigger zero flag
	testJMPConditional(0x51, 0xAE, 0xB0, [4]byte{0xFF, 0x01, 0x01, 0x1}, 0xAE, t)

	// should not jump in false case
	testJMPConditional(0x51, 0xAF, 0xB0, [4]byte{0x00, 0x11, 0x01, 0x01}, 0x03, t)
}

func TestJMPCA(t *testing.T) {
	// Jump If Carry or A larger
	// carry condition
	testJMPConditional(0x5C, 0x90, 0x81, [4]byte{0x04, 0xFF, 0x01, 0x2}, 0x90, t)
	// a is larger
	testJMPConditional(0x5C, 0x90, 0x81, [4]byte{0x0A, 0x01, 0x01, 0x2}, 0x90, t)
	// should not jump in false case
	testJMPConditional(0x5C, 0x91, 0x81, [4]byte{0x01, 0x01, 0x01, 0x01}, 0x03, t)
}

func TestJMPCE(t *testing.T) {
	// Jump If Carry or A = B
	// carry condition
	testJMPConditional(0x5A, 0x90, 0x81, [4]byte{0x04, 0xFF, 0x01, 0x2}, 0x90, t)
	// a = b
	testJMPConditional(0x5A, 0x90, 0x81, [4]byte{0x02, 0x02, 0x01, 0x2}, 0x90, t)
	// should not jump in false case
	testJMPConditional(0x5A, 0x91, 0x81, [4]byte{0x01, 0x02, 0x01, 0x01}, 0x03, t)
}

func TestJMPCZ(t *testing.T) {
	// Jump If Carry or zero flag
	// carry condition
	testJMPConditional(0x59, 0x90, 0x81, [4]byte{0x04, 0xFF, 0x01, 0x2}, 0x90, t)
	// zero flag
	testJMPConditional(0x59, 0x90, 0xB0, [4]byte{0xFF, 0xFE, 0xFE, 0xFE}, 0x90, t)
	// should not jump in false case
	testJMPConditional(0x59, 0x91, 0x81, [4]byte{0x01, 0x02, 0x01, 0x01}, 0x03, t)
}

func TestJMPAE(t *testing.T) {
	// Jump is A is larger or A = B
	// a larger
	testJMPConditional(0x56, 0x20, 0xF1, [4]byte{0x02, 0x01, 0x01, 0x2}, 0x20, t)
	//a = b
	testJMPConditional(0x56, 0x20, 0xF1, [4]byte{0x02, 0x02, 0x01, 0x2}, 0x20, t)
	// should not jump in false case
	testJMPConditional(0x56, 0x21, 0xF1, [4]byte{0x01, 0x03, 0x01, 0x01}, 0x03, t)
}

func TestJMPAZ(t *testing.T) {
	// Jump is A is larger or zero flag
	// a larger
	testJMPConditional(0x55, 0x20, 0xF1, [4]byte{0x02, 0x01, 0x01, 0x2}, 0x20, t)
	// zero flag (using and)
	testJMPConditional(0x55, 0x20, 0xC1, [4]byte{0x01, 0xFE, 0x02, 0x02}, 0x20, t)

	// should not jump in false case
	testJMPConditional(0x55, 0x21, 0xF1, [4]byte{0x01, 0x03, 0x01, 0x01}, 0x03, t)
}

func TestJMPEZ(t *testing.T) {
	// Jump if A = B or zero flag
	// a = b
	testJMPConditional(0x53, 0x20, 0xF1, [4]byte{0x02, 0x02, 0x01, 0x2}, 0x20, t)
	// zero flag (using and)
	testJMPConditional(0x53, 0x20, 0xC1, [4]byte{0x01, 0xFE, 0x02, 0x02}, 0x20, t)

	// should not jump in false case
	testJMPConditional(0x53, 0x21, 0xF1, [4]byte{0x01, 0x03, 0x01, 0x01}, 0x03, t)
}

func TestJMPCAE(t *testing.T) {
	// Jump if Carry OR A is Larger OR A = B

	// carry condition
	testJMPConditional(0x5E, 0x90, 0x81, [4]byte{0x04, 0xFF, 0x01, 0x2}, 0x90, t)

	// a larger
	testJMPConditional(0x5E, 0x20, 0xF1, [4]byte{0x02, 0x01, 0x01, 0x2}, 0x20, t)

	// a = b
	testJMPConditional(0x5E, 0x20, 0xF1, [4]byte{0x02, 0x02, 0x01, 0x2}, 0x20, t)

	// should not jump in false case
	testJMPConditional(0x5E, 0x21, 0xF1, [4]byte{0x01, 0x03, 0x01, 0x01}, 0x03, t)
}

func TestJMPCAZ(t *testing.T) {
	// Jump if Carry OR A is Larger OR zero flag

	// carry condition
	testJMPConditional(0x5D, 0x90, 0x81, [4]byte{0x04, 0xFF, 0x01, 0x2}, 0x90, t)

	// a larger
	testJMPConditional(0x5D, 0x20, 0xF1, [4]byte{0x02, 0x01, 0x01, 0x2}, 0x20, t)

	// zero flag (using and)
	testJMPConditional(0x5D, 0x20, 0xC1, [4]byte{0x01, 0xFE, 0x02, 0x02}, 0x20, t)

	// should not jump in false case
	testJMPConditional(0x5D, 0x21, 0xF1, [4]byte{0x01, 0x03, 0x01, 0x01}, 0x03, t)
}

func TestJMPCEZ(t *testing.T) {
	// Jump if Carry OR a = b OR zero flag

	// carry condition
	testJMPConditional(0x5B, 0x90, 0x81, [4]byte{0x04, 0xFF, 0x01, 0x2}, 0x90, t)

	// a = b
	testJMPConditional(0x5B, 0x20, 0xF1, [4]byte{0x02, 0x02, 0x01, 0x2}, 0x20, t)

	// zero flag (using and)
	testJMPConditional(0x5B, 0x20, 0xC1, [4]byte{0x01, 0xFE, 0x02, 0x02}, 0x20, t)

	// should not jump in false case
	testJMPConditional(0x5B, 0x21, 0xF1, [4]byte{0x01, 0x03, 0x01, 0x01}, 0x03, t)
}

func TestJMPAEZ(t *testing.T) {
	// Jump if a is larger OR a = b OR zero flag

	// a larger
	testJMPConditional(0x57, 0x20, 0xF1, [4]byte{0x02, 0x01, 0x01, 0x2}, 0x20, t)

	// a = b
	testJMPConditional(0x57, 0x20, 0xF1, [4]byte{0x02, 0x02, 0x01, 0x2}, 0x20, t)

	// zero flag (using and)
	testJMPConditional(0x57, 0x20, 0xC1, [4]byte{0x01, 0xFE, 0x02, 0x02}, 0x20, t)

	// should not jump in false case
	testJMPConditional(0x57, 0x21, 0xF1, [4]byte{0x01, 0x03, 0x01, 0x01}, 0x03, t)
}

func TestJMPCAEZ(t *testing.T) {
	// Jump if Carry OR a is larger OR a = b OR zero flag

	// carry condition
	testJMPConditional(0x5F, 0x90, 0x81, [4]byte{0x04, 0xFF, 0x01, 0x2}, 0x90, t)

	// a larger
	testJMPConditional(0x5F, 0x20, 0xF1, [4]byte{0x02, 0x01, 0x01, 0x2}, 0x20, t)

	// a = b
	testJMPConditional(0x5F, 0x20, 0xF1, [4]byte{0x02, 0x02, 0x01, 0x2}, 0x20, t)

	// zero flag (using and)
	testJMPConditional(0x5F, 0x20, 0xC1, [4]byte{0x01, 0xFE, 0x02, 0x02}, 0x20, t)

	// should not jump in false case
	testJMPConditional(0x5F, 0x21, 0xF1, [4]byte{0x01, 0x03, 0x01, 0x01}, 0x03, t)
}

func testJMPConditional(jmpConditionInstr, destination, initialInstr byte, inputRegisters [4]byte, expectedIAR byte, t *testing.T) {
	b := components.NewBus()
	m := memory.NewMemory256(b)
	c := NewCPU(b, m)

	var insAddr byte = 0x00

	setMemoryLocation(c, insAddr, initialInstr)
	setMemoryLocation(c, insAddr+1, jmpConditionInstr)
	setMemoryLocation(c, insAddr+2, destination)

	setIAR(c, insAddr)

	setRegisters(c, inputRegisters)

	doFetchDecodeExecute(c)
	doFetchDecodeExecute(c)

	// check IAR
	checkIAR(c, expectedIAR, t)
}

func TestALUAdd(t *testing.T) {
	var inputs [4]byte = [4]byte{0x02, 0x03, 0x04, 0x05}

	testInstruction(0x80, inputs, [4]byte{inputs[0] + inputs[0], inputs[1], inputs[2], inputs[3]}, t)
	testInstruction(0x81, inputs, [4]byte{inputs[0], inputs[1] + inputs[0], inputs[2], inputs[3]}, t)
	testInstruction(0x82, inputs, [4]byte{inputs[0], inputs[1], inputs[2] + inputs[0], inputs[3]}, t)
	testInstruction(0x83, inputs, [4]byte{inputs[0], inputs[1], inputs[2], inputs[3] + inputs[0]}, t)

	testInstruction(0x84, inputs, [4]byte{inputs[0] + inputs[1], inputs[1], inputs[2], inputs[3]}, t)
	testInstruction(0x85, inputs, [4]byte{inputs[0], inputs[1] + inputs[1], inputs[2], inputs[3]}, t)
	testInstruction(0x86, inputs, [4]byte{inputs[0], inputs[1], inputs[2] + inputs[1], inputs[3]}, t)
	testInstruction(0x87, inputs, [4]byte{inputs[0], inputs[1], inputs[2], inputs[3] + inputs[1]}, t)

	testInstruction(0x88, inputs, [4]byte{inputs[0] + inputs[2], inputs[1], inputs[2], inputs[3]}, t)
	testInstruction(0x89, inputs, [4]byte{inputs[0], inputs[1] + inputs[2], inputs[2], inputs[3]}, t)
	testInstruction(0x8A, inputs, [4]byte{inputs[0], inputs[1], inputs[2] + inputs[2], inputs[3]}, t)
	testInstruction(0x8B, inputs, [4]byte{inputs[0], inputs[1], inputs[2], inputs[3] + inputs[2]}, t)

	testInstruction(0x8C, inputs, [4]byte{inputs[0] + inputs[3], inputs[1], inputs[2], inputs[3]}, t)
	testInstruction(0x8D, inputs, [4]byte{inputs[0], inputs[1] + inputs[3], inputs[2], inputs[3]}, t)
	testInstruction(0x8E, inputs, [4]byte{inputs[0], inputs[1], inputs[2] + inputs[3], inputs[3]}, t)
	testInstruction(0x8F, inputs, [4]byte{inputs[0], inputs[1], inputs[2], inputs[3] + inputs[3]}, t)

	// TODO work with carry flags etc?

}

func TestALUNOT(t *testing.T) {
	var inputs [4]byte = [4]byte{0xFF, 0xFE, 0xFD, 0xFC}

	testInstruction(0xB0, inputs, [4]byte{0x00, 0xFE, 0xFD, 0xFC}, t)
	testInstruction(0xB5, inputs, [4]byte{0xFF, 0x01, 0xFD, 0xFC}, t)
	testInstruction(0xBA, inputs, [4]byte{0xFF, 0xFE, 0x02, 0xFC}, t)
	testInstruction(0xBF, inputs, [4]byte{0xFF, 0xFE, 0xFD, 0x03}, t)
}

func TestALUAND(t *testing.T) {
	var inputs [4]byte = [4]byte{0x02, 0x03, 0x04, 0x05}

	testInstruction(0xC0, inputs, [4]byte{inputs[0] & inputs[0], inputs[1], inputs[2], inputs[3]}, t)
	testInstruction(0xC1, inputs, [4]byte{inputs[0], inputs[1] & inputs[0], inputs[2], inputs[3]}, t)
	testInstruction(0xC2, inputs, [4]byte{inputs[0], inputs[1], inputs[2] & inputs[0], inputs[3]}, t)
	testInstruction(0xC3, inputs, [4]byte{inputs[0], inputs[1], inputs[2], inputs[3] & inputs[0]}, t)

	testInstruction(0xC4, inputs, [4]byte{inputs[0] & inputs[1], inputs[1], inputs[2], inputs[3]}, t)
	testInstruction(0xC5, inputs, [4]byte{inputs[0], inputs[1] & inputs[1], inputs[2], inputs[3]}, t)
	testInstruction(0xC6, inputs, [4]byte{inputs[0], inputs[1], inputs[2] & inputs[1], inputs[3]}, t)
	testInstruction(0xC7, inputs, [4]byte{inputs[0], inputs[1], inputs[2], inputs[3] & inputs[1]}, t)

	testInstruction(0xC8, inputs, [4]byte{inputs[0] & inputs[2], inputs[1], inputs[2], inputs[3]}, t)
	testInstruction(0xC9, inputs, [4]byte{inputs[0], inputs[1] & inputs[2], inputs[2], inputs[3]}, t)
	testInstruction(0xCA, inputs, [4]byte{inputs[0], inputs[1], inputs[2] & inputs[2], inputs[3]}, t)
	testInstruction(0xCB, inputs, [4]byte{inputs[0], inputs[1], inputs[2], inputs[3] & inputs[2]}, t)

	testInstruction(0xCC, inputs, [4]byte{inputs[0] & inputs[3], inputs[1], inputs[2], inputs[3]}, t)
	testInstruction(0xCD, inputs, [4]byte{inputs[0], inputs[1] & inputs[3], inputs[2], inputs[3]}, t)
	testInstruction(0xCE, inputs, [4]byte{inputs[0], inputs[1], inputs[2] & inputs[3], inputs[3]}, t)
	testInstruction(0xCF, inputs, [4]byte{inputs[0], inputs[1], inputs[2], inputs[3] & inputs[3]}, t)
}

func TestALUOR(t *testing.T) {
	var inputs [4]byte = [4]byte{0x92, 0x91, 0x45, 0xAF}

	testInstruction(0xD0, inputs, [4]byte{inputs[0] | inputs[0], inputs[1], inputs[2], inputs[3]}, t)
	testInstruction(0xD1, inputs, [4]byte{inputs[0], inputs[1] | inputs[0], inputs[2], inputs[3]}, t)
	testInstruction(0xD2, inputs, [4]byte{inputs[0], inputs[1], inputs[2] | inputs[0], inputs[3]}, t)
	testInstruction(0xD3, inputs, [4]byte{inputs[0], inputs[1], inputs[2], inputs[3] | inputs[0]}, t)

	testInstruction(0xD4, inputs, [4]byte{inputs[0] | inputs[1], inputs[1], inputs[2], inputs[3]}, t)
	testInstruction(0xD5, inputs, [4]byte{inputs[0], inputs[1] | inputs[1], inputs[2], inputs[3]}, t)
	testInstruction(0xD6, inputs, [4]byte{inputs[0], inputs[1], inputs[2] | inputs[1], inputs[3]}, t)
	testInstruction(0xD7, inputs, [4]byte{inputs[0], inputs[1], inputs[2], inputs[3] | inputs[1]}, t)

	testInstruction(0xD8, inputs, [4]byte{inputs[0] | inputs[2], inputs[1], inputs[2], inputs[3]}, t)
	testInstruction(0xD9, inputs, [4]byte{inputs[0], inputs[1] | inputs[2], inputs[2], inputs[3]}, t)
	testInstruction(0xDA, inputs, [4]byte{inputs[0], inputs[1], inputs[2] | inputs[2], inputs[3]}, t)
	testInstruction(0xDB, inputs, [4]byte{inputs[0], inputs[1], inputs[2], inputs[3] | inputs[2]}, t)

	testInstruction(0xDC, inputs, [4]byte{inputs[0] | inputs[3], inputs[1], inputs[2], inputs[3]}, t)
	testInstruction(0xDD, inputs, [4]byte{inputs[0], inputs[1] | inputs[3], inputs[2], inputs[3]}, t)
	testInstruction(0xDE, inputs, [4]byte{inputs[0], inputs[1], inputs[2] | inputs[3], inputs[3]}, t)
	testInstruction(0xDF, inputs, [4]byte{inputs[0], inputs[1], inputs[2], inputs[3] | inputs[3]}, t)
}

func TestALUXOR(t *testing.T) {
	var inputs [4]byte = [4]byte{0x92, 0x91, 0x45, 0xAF}
	testInstruction(0xE0, inputs, [4]byte{inputs[0] ^ inputs[0], inputs[1], inputs[2], inputs[3]}, t)
	testInstruction(0xE1, inputs, [4]byte{inputs[0], inputs[1] ^ inputs[0], inputs[2], inputs[3]}, t)
	testInstruction(0xE2, inputs, [4]byte{inputs[0], inputs[1], inputs[2] ^ inputs[0], inputs[3]}, t)
	testInstruction(0xE3, inputs, [4]byte{inputs[0], inputs[1], inputs[2], inputs[3] ^ inputs[0]}, t)

	testInstruction(0xE4, inputs, [4]byte{inputs[0] ^ inputs[1], inputs[1], inputs[2], inputs[3]}, t)
	testInstruction(0xE5, inputs, [4]byte{inputs[0], inputs[1] ^ inputs[1], inputs[2], inputs[3]}, t)
	testInstruction(0xE6, inputs, [4]byte{inputs[0], inputs[1], inputs[2] ^ inputs[1], inputs[3]}, t)
	testInstruction(0xE7, inputs, [4]byte{inputs[0], inputs[1], inputs[2], inputs[3] ^ inputs[1]}, t)

	testInstruction(0xE8, inputs, [4]byte{inputs[0] ^ inputs[2], inputs[1], inputs[2], inputs[3]}, t)
	testInstruction(0xE9, inputs, [4]byte{inputs[0], inputs[1] ^ inputs[2], inputs[2], inputs[3]}, t)
	testInstruction(0xEA, inputs, [4]byte{inputs[0], inputs[1], inputs[2] ^ inputs[2], inputs[3]}, t)
	testInstruction(0xEB, inputs, [4]byte{inputs[0], inputs[1], inputs[2], inputs[3] ^ inputs[2]}, t)

	testInstruction(0xEC, inputs, [4]byte{inputs[0] ^ inputs[3], inputs[1], inputs[2], inputs[3]}, t)
	testInstruction(0xED, inputs, [4]byte{inputs[0], inputs[1] ^ inputs[3], inputs[2], inputs[3]}, t)
	testInstruction(0xEE, inputs, [4]byte{inputs[0], inputs[1], inputs[2] ^ inputs[3], inputs[3]}, t)
	testInstruction(0xEF, inputs, [4]byte{inputs[0], inputs[1], inputs[2], inputs[3] ^ inputs[3]}, t)
}

func TestCMP(t *testing.T) {
	var inputs [4]byte = [4]byte{0x92, 0x91, 0x45, 0xAF}

	// outputs should remain the same
	var i byte
	for i = 0; i <= 0x0F; i++ {
		testInstruction(0xF0+i, inputs, inputs, t)
	}

	//TODO work with carry flags etc?
}

func testInstruction(instruction byte, inputRegisters [4]byte, expectedOutputRegisters [4]byte, t *testing.T) {
	b := components.NewBus()
	m := memory.NewMemory256(b)
	c := NewCPU(b, m)

	setMemoryLocation(c, 0x00, instruction)

	for i, r := range inputRegisters {
		setRegister(c, i, r)
	}

	setIAR(c, 0x00)

	doFetchDecodeExecute(c)

	checkRegisters(c, expectedOutputRegisters[0], expectedOutputRegisters[1], expectedOutputRegisters[2], expectedOutputRegisters[3], t)
}

func TestALUShiftLeft(t *testing.T) {
	var ones [4]byte = [4]byte{0x01, 0x01, 0x01, 0x01}
	var shifts byte
	for shifts = 0; shifts < 8; shifts++ {
		testShift(0xA0, ones, [4]byte{1 << shifts, 0x01, 0x01, 0x01}, shifts, t)
		testShift(0xA5, ones, [4]byte{0x01, 1 << shifts, 0x01, 0x01}, shifts, t)
		testShift(0xAA, ones, [4]byte{0x01, 0x01, 1 << shifts, 0x01}, shifts, t)
		testShift(0xAF, ones, [4]byte{0x01, 0x01, 0x01, 1 << shifts}, shifts, t)
	}
}

func TestALUShiftRight(t *testing.T) {
	var input [4]byte = [4]byte{0x80, 0x80, 0x80, 0x80}
	var shifts byte
	for shifts = 0; shifts < 8; shifts++ {
		testShift(0x90, input, [4]byte{0x80 >> shifts, 0x80, 0x80, 0x80}, shifts, t)
		testShift(0x95, input, [4]byte{0x80, 0x80 >> shifts, 0x80, 0x80}, shifts, t)
		testShift(0x9A, input, [4]byte{0x80, 0x80, 0x80 >> shifts, 0x80}, shifts, t)
		testShift(0x9F, input, [4]byte{0x80, 0x80, 0x80, 0x80 >> shifts}, shifts, t)
	}
}

func testShift(instruction byte, inputRegisters [4]byte, expectedOutputRegisters [4]byte, shifts byte, t *testing.T) {
	b := components.NewBus()
	m := memory.NewMemory256(b)
	c := NewCPU(b, m)

	var i byte
	for i = 0; i < shifts; i++ {
		setMemoryLocation(c, i, instruction)
	}

	for i, r := range inputRegisters {
		setRegister(c, i, r)
	}

	setIAR(c, 0x00)

	for i = 0; i < shifts; i++ {
		doFetchDecodeExecute(c)
	}

	checkRegisters(c, expectedOutputRegisters[0], expectedOutputRegisters[1], expectedOutputRegisters[2], expectedOutputRegisters[3], t)
}

func doFetchDecodeExecute(c *CPU) {
	for i := 0; i < 6; i++ {
		c.step(true)
		c.step(false)
	}

}

func setIAR(c *CPU, value byte) {
	setBus(c.mainBus, value)

	c.iar.Set()
	c.iar.Update()
	c.iar.Unset()
	c.iar.Update()
}

func setMemoryLocation(c *CPU, address byte, value byte) {
	c.memory.AddressRegister.Set()
	setBus(c.mainBus, address)
	c.memory.Update()

	c.memory.AddressRegister.Unset()
	c.memory.Update()

	setBus(c.mainBus, value)
	c.memory.Set()
	c.memory.Update()

	c.memory.Unset()
	c.memory.Update()
}

func setRegister(c *CPU, register int, value byte) {
	switch register {
	case 0:
		c.gpReg0.Set()
		c.gpReg0.Update()
		setBus(c.mainBus, value)
		c.gpReg0.Update()
		c.gpReg0.Unset()
		c.gpReg0.Update()
	case 1:
		c.gpReg1.Set()
		c.gpReg1.Update()
		setBus(c.mainBus, value)
		c.gpReg1.Update()
		c.gpReg1.Unset()
		c.gpReg1.Update()
	case 2:
		c.gpReg2.Set()
		c.gpReg2.Update()
		setBus(c.mainBus, value)
		c.gpReg2.Update()
		c.gpReg2.Unset()
		c.gpReg2.Update()
	case 3:
		c.gpReg3.Set()
		c.gpReg3.Update()
		setBus(c.mainBus, value)
		c.gpReg3.Update()
		c.gpReg3.Unset()
		c.gpReg3.Update()
	}
}

func checkIAR(c *CPU, expValue byte, t *testing.T) {
	if c.iar.Value() != expValue {
		t.Logf("Expected IAR to have value of: %X but got %X", expValue, c.iar.Value())
		t.FailNow()
	}
}

func checkIR(c *CPU, expValue byte, t *testing.T) {
	if c.ir.Value() != expValue {
		t.Logf("Expected IR to have value of: %X but got %X", expValue, c.ir.Value())
		t.FailNow()
	}
}

func checkFlagsRegister(c *CPU, expectedCarry, expectedIsLarger, expectedIsEqual, expectedIsZero bool, t *testing.T) {
	if carryFlagSet := c.flagsBus.GetOutputWire(0); carryFlagSet != expectedCarry {
		t.Logf("Expected is carry out flag to be %v but got %v", expectedCarry, carryFlagSet)
		t.FailNow()
	}
	if isLargerFlagSet := c.flagsBus.GetOutputWire(1); isLargerFlagSet != expectedIsLarger {
		t.Logf("Expected is larger flag to be %v but got %v", expectedIsLarger, isLargerFlagSet)
		t.FailNow()
	}
	if equalFlagSet := c.flagsBus.GetOutputWire(2); equalFlagSet != expectedIsEqual {
		t.Logf("Expected equal flag to be %v but got %v", expectedIsEqual, equalFlagSet)
		t.FailNow()
	}
	if zeroFlagSet := c.flagsBus.GetOutputWire(3); zeroFlagSet != expectedIsZero {
		t.Logf("Expected zero flag to be %v but got %v", expectedIsZero, zeroFlagSet)
		t.FailNow()
	}
}

func checkRegisters(c *CPU, expReg0, expReg1, expReg2, expReg3 byte, t *testing.T) {
	if c.gpReg0.Value() != expReg0 {
		t.Logf("Expected register 0 to have value of: %X but got %X", expReg0, c.gpReg0.Value())
		t.FailNow()
	}

	if c.gpReg1.Value() != expReg1 {
		t.Logf("Expected register 1 to have value of: %X but got %X", expReg1, c.gpReg1.Value())
		t.FailNow()
	}

	if c.gpReg2.Value() != expReg2 {
		t.Logf("Expected register 2 to have value of: %X but got %X", expReg2, c.gpReg2.Value())
		t.FailNow()
	}

	if c.gpReg3.Value() != expReg3 {
		t.Logf("Expected register 3 to have value of: %X but got %X", expReg3, c.gpReg3.Value())
		t.FailNow()
	}

}

func setRegisters(c *CPU, values [4]byte) {
	for i, v := range values {
		setRegister(c, i, v)
	}
}
