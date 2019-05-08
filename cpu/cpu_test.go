package cpu

import (
	"testing"

	"github.com/djhworld/simple-computer/components"
	"github.com/djhworld/simple-computer/memory"
)

var BUS *components.Bus = components.NewBus(BUS_WIDTH)
var MEMORY *memory.Memory64K = memory.NewMemory64K(BUS)

func SetUpCPU() *CPU {
	return NewCPU(BUS, MEMORY)
}

func ClearMem() {
	for i := uint16(0); i < 65535; i++ {
		setMemoryLocation2(MEMORY, i, 0x0000)
	}
	setMemoryLocation2(MEMORY, 0xFFFF, 0x0000)
}

func TestIARIncrementedOnEveryCycle(t *testing.T) {
	ClearMem()
	c := SetUpCPU()

	c.SetIAR(0x0000)

	var q uint16
	for q = 0; q < 1000; q++ {
		doFetchDecodeExecute(c)
		checkIAR(c, q+1, t)
	}
}

func TestInstructionReceivedFromMemory(t *testing.T) {
	ClearMem()
	c := SetUpCPU()

	instructions := []uint16{0x008A, 0x0082, 0x0088, 0x0094, 0x00B1}

	var addr uint16 = 0xF00F
	for _, b := range instructions {
		setMemoryLocation(c, addr, b)
		addr++
	}

	c.SetIAR(0xF00F)

	for _, b := range instructions {
		doFetchDecodeExecute(c)
		checkIR(c, b, t)
	}
}

func TestFlagsRegisterAllFalse(t *testing.T) {
	ClearMem()
	c := SetUpCPU()

	setMemoryLocation(c, 0x0000, 0x0081)
	setRegisters(c, [4]uint16{0x0009, 0x000A, 0x0002, 0x0003})
	c.SetIAR(0x0000)

	doFetchDecodeExecute(c)

	checkFlagsRegister(c, false, false, false, false, t)
}

func TestFlagsRegisterCarryFlagEnabled(t *testing.T) {
	ClearMem()
	c := SetUpCPU()

	setMemoryLocation(c, 0x0000, 0x0081)
	setRegisters(c, [4]uint16{0x0020, 0xFFFF, 0x0002, 0x0003})
	c.SetIAR(0x0000)

	doFetchDecodeExecute(c)

	checkFlagsRegister(c, true, false, false, false, t)
}

func TestFlagsRegisterIsLargerFlagEnabled(t *testing.T) {
	ClearMem()
	c := SetUpCPU()

	setMemoryLocation(c, 0x0000, 0x0081)
	setRegisters(c, [4]uint16{0x0021, 0x0020, 0x0002, 0x0003})
	c.SetIAR(0x0000)

	doFetchDecodeExecute(c)

	checkFlagsRegister(c, false, true, false, false, t)
}

func TestFlagsRegisterIsEqualsFlagEnabled(t *testing.T) {
	ClearMem()
	c := SetUpCPU()

	setMemoryLocation(c, 0x0000, 0x0081)
	setRegisters(c, [4]uint16{0x0021, 0x0021, 0x0002, 0x0003})
	c.SetIAR(0x0000)

	doFetchDecodeExecute(c)

	checkFlagsRegister(c, false, false, true, false, t)
}

func TestFlagsRegisterIsZeroFlagEnabled(t *testing.T) {
	ClearMem()
	c := SetUpCPU()

	setMemoryLocation(c, 0x0000, 0x0081)
	setRegisters(c, [4]uint16{0x0001, 0xFFFF, 0x0002, 0x0003})
	c.SetIAR(0x0000)

	doFetchDecodeExecute(c)

	checkFlagsRegister(c, true, false, false, true, t)
}

func TestFlagsRegisterMultipleEnabled(t *testing.T) {
	ClearMem()
	c := SetUpCPU()

	setMemoryLocation(c, 0x0000, 0x0081)
	setRegisters(c, [4]uint16{0xFFFF, 0x0001, 0x0002, 0x0003})
	c.SetIAR(0x0000)

	doFetchDecodeExecute(c)

	checkFlagsRegister(c, true, true, false, true, t)
}
func TestSTThenLD(t *testing.T) {
	ClearMem()
	c := SetUpCPU()

	for i := uint16(0); i < 512; i++ {
		//  ST R2, R3
		setMemoryLocation(c, i, 0x001B)
	}
	c.SetIAR(0x0000)

	//storing values into memory
	var value uint16 = 0x0400
	for i := 0x0200; i < 0x0400; i++ {
		setRegisters(c, [4]uint16{0x0001, 0x0001, uint16(i), value})
		doFetchDecodeExecute(c)
		value--
	}

	for i := uint16(0); i < 512; i++ {
		// LD R2, R3
		setMemoryLocation(c, i, 0x000B)
	}
	c.SetIAR(0x0000)

	//retrieving them back from memory
	value = 0x0400
	for i := 0x0200; i < 0x0400; i++ {
		setRegisters(c, [4]uint16{0x0001, 0x0001, uint16(i), 0x0001})
		doFetchDecodeExecute(c)
		checkRegisters(c, 0x0001, 0x0001, uint16(i), value, t)
		value--
	}
}

func TestLD4Times(t *testing.T) {
	ClearMem()
	c := SetUpCPU()

	var addr uint16 = 0x00A2
	values := []uint16{0x0088, 0x0090, 0x0092, 0x00AB}

	for i := uint16(0); i < uint16(len(values)); i++ {
		setMemoryLocation(c, i, 0x0001)
		setMemoryLocation(c, addr, values[i])
		addr++
	}

	c.SetIAR(0x0000)

	addr = 0x00A2
	for _, v := range values {
		setRegister(c, 0, addr)
		setRegister(c, 1, 0x0001)
		setRegister(c, 2, 0x0001)
		setRegister(c, 3, 0x0001)

		doFetchDecodeExecute(c)
		checkRegisters(c, addr, v, 0x0001, 0x0001, t)
		addr++
	}
}

func TestLD(t *testing.T) {
	ClearMem()

	testLD(0x0000, 0x0080, 0x0023, []uint16{0x0080, 0x0081, 0x0082, 0x0083}, []uint16{0x0023, 0x0081, 0x0082, 0x0083}, t)
	testLD(0x0001, 0x0084, 0x00F2, []uint16{0x0084, 0x0085, 0x0086, 0x0087}, []uint16{0x0084, 0x00F2, 0x0086, 0x0087}, t)
	testLD(0x0002, 0x0088, 0x0001, []uint16{0x0088, 0x0089, 0x008A, 0x008B}, []uint16{0x0088, 0x0089, 0x0001, 0x008B}, t)
	testLD(0x0003, 0x008C, 0x005A, []uint16{0x008C, 0x008D, 0x008E, 0x008F}, []uint16{0x008C, 0x008D, 0x008E, 0x005A}, t)

	testLD(0x0004, 0x0091, 0x0023, []uint16{0x0090, 0x0091, 0x0092, 0x0093}, []uint16{0x0023, 0x0091, 0x0092, 0x0093}, t)
	testLD(0x0005, 0x0095, 0x00F2, []uint16{0x0094, 0x0095, 0x0096, 0x0097}, []uint16{0x0094, 0x00F2, 0x0096, 0x0097}, t)
	testLD(0x0006, 0x0099, 0x0001, []uint16{0x0098, 0x0099, 0x009A, 0x009B}, []uint16{0x0098, 0x0099, 0x0001, 0x009B}, t)
	testLD(0x0007, 0x009D, 0x005A, []uint16{0x009C, 0x009D, 0x009E, 0x009F}, []uint16{0x009C, 0x009D, 0x009E, 0x005A}, t)

	testLD(0x0008, 0x00A2, 0x0023, []uint16{0x00A0, 0x00A1, 0x00A2, 0x00A3}, []uint16{0x0023, 0x00A1, 0x00A2, 0x00A3}, t)
	testLD(0x0009, 0x00A6, 0x00F2, []uint16{0x00A4, 0x00A5, 0x00A6, 0x00A7}, []uint16{0x00A4, 0x00F2, 0x00A6, 0x00A7}, t)
	testLD(0x000A, 0x00AA, 0x0001, []uint16{0x00A8, 0x00A9, 0x00AA, 0x00AB}, []uint16{0x00A8, 0x00A9, 0x0001, 0x00AB}, t)
	testLD(0x000B, 0x00AE, 0x005A, []uint16{0x00AC, 0x00AD, 0x00AE, 0x00AF}, []uint16{0x00AC, 0x00AD, 0x00AE, 0x005A}, t)

	testLD(0x000C, 0x00B3, 0x0023, []uint16{0x00B0, 0x00B1, 0x00B2, 0x00B3}, []uint16{0x0023, 0x00B1, 0x00B2, 0x00B3}, t)
	testLD(0x000D, 0x00B7, 0x00F2, []uint16{0x00B4, 0x00B5, 0x00B6, 0x00B7}, []uint16{0x00B4, 0x00F2, 0x00B6, 0x00B7}, t)
	testLD(0x000E, 0x22BB, 0xAB01, []uint16{0x00B8, 0x00B9, 0x00BA, 0x22BB}, []uint16{0x00B8, 0x00B9, 0xAB01, 0x22BB}, t)
	testLD(0x000F, 0x00BF, 0x005A, []uint16{0x00BC, 0x00BD, 0x00BE, 0x00BF}, []uint16{0x00BC, 0x00BD, 0x00BE, 0x005A}, t)
}

func testLD(instruction uint16, memAddress, memValue uint16, inputRegisters, expectedOutputRegisters []uint16, t *testing.T) {
	c := SetUpCPU()
	var insAddr uint16 = 0x0000
	setMemoryLocation(c, insAddr, instruction)
	c.SetIAR(insAddr)

	setMemoryLocation(c, memAddress, memValue)

	for i, v := range inputRegisters {
		setRegister(c, i, v)
	}

	doFetchDecodeExecute(c)
	checkRegisters(c, expectedOutputRegisters[0], expectedOutputRegisters[1], expectedOutputRegisters[2], expectedOutputRegisters[3], t)
}

func TestST(t *testing.T) {
	ClearMem()
	testST(0x0010, [4]uint16{0x00A0, 0x0001, 0x0001, 0x0001}, 0x00A0, 0x00A0, t)
	testST(0x0011, [4]uint16{0x00A1, 0x0029, 0x0001, 0x0001}, 0x00A1, 0x0029, t)
	testST(0x0012, [4]uint16{0x00A2, 0x0001, 0x007F, 0x0001}, 0x00A2, 0x007F, t)
	testST(0x0013, [4]uint16{0x00A3, 0x0001, 0x0001, 0x001B}, 0x00A3, 0x001B, t)

	testST(0x0014, [4]uint16{0x00A0, 0x00B4, 0x0001, 0x0001}, 0x00B4, 0x00A0, t)
	testST(0x0015, [4]uint16{0x0001, 0x00B5, 0x0001, 0x0001}, 0x00B5, 0x00B5, t)
	testST(0x0016, [4]uint16{0x0001, 0x00B6, 0x007F, 0x0001}, 0x00B6, 0x007F, t)
	testST(0x0017, [4]uint16{0x0001, 0x00B7, 0x0001, 0x001B}, 0x00B7, 0x001B, t)

	testST(0x0018, [4]uint16{0x00A0, 0x0001, 0x00C8, 0x0001}, 0x00C8, 0x00A0, t)
	testST(0x0019, [4]uint16{0x0001, 0x0029, 0x00C9, 0x0001}, 0x00C9, 0x0029, t)
	testST(0x001A, [4]uint16{0x0001, 0x0001, 0x00CA, 0x0001}, 0x00CA, 0x00CA, t)
	testST(0x001B, [4]uint16{0x0001, 0x0001, 0x00CB, 0x001B}, 0x00CB, 0x001B, t)

	testST(0x001C, [4]uint16{0x00A0, 0x0001, 0x0001, 0x00DC}, 0x00DC, 0x00A0, t)
	testST(0x001D, [4]uint16{0x0001, 0x0029, 0x0001, 0x00DD}, 0x00DD, 0x0029, t)
	testST(0x001E, [4]uint16{0x0001, 0x0001, 0x1A7F, 0xFCDE}, 0xFCDE, 0x1A7F, t)
	testST(0x001F, [4]uint16{0x0001, 0x0001, 0x0001, 0x00DF}, 0x00DF, 0x00DF, t)
}

func testST(instruction uint16, inputRegisters [4]uint16, expectedValueAddress, expectedValue uint16, t *testing.T) {
	c := SetUpCPU()

	// ST value into memory
	var insAddr uint16 = 0x0000
	setMemoryLocation(c, insAddr, instruction)
	c.SetIAR(insAddr)

	setRegisters(c, inputRegisters)

	doFetchDecodeExecute(c)

	//LD value into register zero
	setMemoryLocation(c, insAddr+1, 0x0000)
	c.SetIAR(insAddr + 1)

	setRegisters(c, [4]uint16{expectedValueAddress, inputRegisters[1], inputRegisters[2], inputRegisters[3]})

	doFetchDecodeExecute(c)

	if c.gpReg0.Value() != expectedValue {
		t.Logf("Expected register 0 to have value of: %X but got %X", expectedValue, c.gpReg0.Value())
		t.FailNow()
	}
}

func TestDATA(t *testing.T) {
	ClearMem()
	c := SetUpCPU()

	var insAddr uint16 = 0x0000

	// DATA R0
	setMemoryLocation(c, insAddr, 0x0020)
	setMemoryLocation(c, insAddr+1, 0xF071)

	// DATA R1
	setMemoryLocation(c, insAddr+2, 0x0021)
	setMemoryLocation(c, insAddr+3, 0xF172)

	// DATA R2
	setMemoryLocation(c, insAddr+4, 0x0022)
	setMemoryLocation(c, insAddr+5, 0xF273)

	// DATA R3
	setMemoryLocation(c, insAddr+6, 0x0023)
	setMemoryLocation(c, insAddr+7, 0xF374)

	c.SetIAR(insAddr)

	setRegisters(c, [4]uint16{0x0001, 0x0001, 0x0001, 0x0001})

	for i := 0; i < 4; i++ {
		doFetchDecodeExecute(c)
	}

	checkRegisters(c, 0xF071, 0xF172, 0xF273, 0xF374, t)

	// check IAR has incremented 2 each time
	checkIAR(c, 0x0008, t)
}

func TestJMPR(t *testing.T) {
	ClearMem()
	testJMPR(0x0030, [4]uint16{0x0083, 0x0001, 0x0001, 0x0001}, 0x0083, t)
	testJMPR(0x0031, [4]uint16{0x0001, 0x00F1, 0x0001, 0x0001}, 0x00F1, t)
	testJMPR(0x0032, [4]uint16{0x0001, 0x0001, 0x00BB, 0x0001}, 0x00BB, t)
	testJMPR(0x0033, [4]uint16{0x0001, 0x0001, 0x0001, 0xFF19}, 0xFF19, t)
}

func testJMPR(instruction uint16, inputRegisters [4]uint16, expectedIAR uint16, t *testing.T) {
	c := SetUpCPU()
	var insAddr uint16 = 0x0000

	// JMPR
	setMemoryLocation(c, insAddr, instruction)

	c.SetIAR(insAddr)

	setRegisters(c, inputRegisters)

	doFetchDecodeExecute(c)

	// registers shouldn't change
	checkRegisters(c, inputRegisters[0], inputRegisters[1], inputRegisters[2], inputRegisters[3], t)

	// check IAR has jumped to the new location
	checkIAR(c, expectedIAR, t)
}

func TestJMP(t *testing.T) {
	ClearMem()
	for i := 0; i < 0x0400; i++ {
		testJMP(uint16(i), t)
	}
}

func testJMP(expectedIAR uint16, t *testing.T) {
	c := SetUpCPU()
	var insAddr uint16 = 0x0000

	// JMP
	setMemoryLocation(c, insAddr, 0x0040)
	setMemoryLocation(c, insAddr+1, expectedIAR)

	c.SetIAR(insAddr)

	inputRegisters := [4]uint16{0x0001, 0x0001, 0x0001, 0x0001}
	setRegisters(c, inputRegisters)

	doFetchDecodeExecute(c)

	// registers shouldn't change
	checkRegisters(c, inputRegisters[0], inputRegisters[1], inputRegisters[2], inputRegisters[3], t)

	// check IAR has jumped to the new location
	checkIAR(c, expectedIAR, t)
}

func TestJMPC(t *testing.T) {
	ClearMem()
	testJMPConditional(0x0058, 0x0090, 0x0081, [4]uint16{0x0004, 0xFFFF, 0x0001, 0x002}, 0x0090, t)
	// should not jump in false case
	testJMPConditional(0x0058, 0x0091, 0x0081, [4]uint16{0x0005, 0x0006, 0x0001, 0x0001}, 0x0003, t)
}

func TestJMPA(t *testing.T) {
	ClearMem()
	testJMPConditional(0x0054, 0x0020, 0x00F1, [4]uint16{0x0002, 0x0001, 0x0001, 0x002}, 0x0020, t)
	// should not jump in false case
	testJMPConditional(0x0054, 0x0021, 0x00F1, [4]uint16{0x0001, 0x0003, 0x0001, 0x0001}, 0x0003, t)
}

func TestJMPE(t *testing.T) {
	ClearMem()
	testJMPConditional(0x0052, 0x00AE, 0x00F1, [4]uint16{0x0000, 0x0000, 0x0001, 0x002}, 0x00AE, t)
	// should not jump in false case
	testJMPConditional(0x0052, 0x00AF, 0x00F1, [4]uint16{0x0010, 0x0011, 0x0001, 0x0001}, 0x0003, t)
}
func TestJMPZ(t *testing.T) {
	ClearMem()
	// perform NOT on R0 (0x00FF) to trigger zero flag
	testJMPConditional(0x0051, 0x00AE, 0x00B0, [4]uint16{0xFFFF, 0x0001, 0x0001, 0x001}, 0x00AE, t)

	// should not jump in false case
	testJMPConditional(0x0051, 0x00AF, 0x00B0, [4]uint16{0x0000, 0x0011, 0x0001, 0x0001}, 0x0003, t)
}

func TestJMPCA(t *testing.T) {
	ClearMem()
	// Jump If Carry or A larger
	// carry condition
	testJMPConditional(0x005C, 0x0090, 0x0081, [4]uint16{0x0004, 0xFFFF, 0x0001, 0x002}, 0x0090, t)
	// a is larger
	testJMPConditional(0x005C, 0x0090, 0x0081, [4]uint16{0x000A, 0x0001, 0x0001, 0x002}, 0x0090, t)
	// should not jump in false case
	testJMPConditional(0x005C, 0x0091, 0x0081, [4]uint16{0x0001, 0x0001, 0x0001, 0x0001}, 0x0003, t)
}

func TestJMPCE(t *testing.T) {
	ClearMem()
	// Jump If Carry or A = B
	// carry condition
	testJMPConditional(0x005A, 0x0090, 0x0081, [4]uint16{0x0004, 0xFFFF, 0x0001, 0x002}, 0x0090, t)
	// a = b
	testJMPConditional(0x005A, 0x0090, 0x0081, [4]uint16{0x0002, 0x0002, 0x0001, 0x002}, 0x0090, t)
	// should not jump in false case
	testJMPConditional(0x005A, 0x0091, 0x0081, [4]uint16{0x0001, 0x0002, 0x0001, 0x0001}, 0x0003, t)
}

func TestJMPCZ(t *testing.T) {
	ClearMem()
	// Jump If Carry or zero flag
	// carry condition
	testJMPConditional(0x0059, 0x0090, 0x0081, [4]uint16{0x0004, 0xFFFF, 0x0001, 0x002}, 0x0090, t)
	// zero flag
	testJMPConditional(0x0059, 0x0090, 0x00B0, [4]uint16{0xFFFF, 0x00FE, 0x00FE, 0x00FE}, 0x0090, t)
	// should not jump in false case
	testJMPConditional(0x0059, 0x0091, 0x0081, [4]uint16{0x0001, 0x0002, 0x0001, 0x0001}, 0x0003, t)
}

func TestJMPAE(t *testing.T) {
	ClearMem()
	// Jump is A is larger or A = B
	// a larger
	testJMPConditional(0x0056, 0x0020, 0x00F1, [4]uint16{0x0002, 0x0001, 0x0001, 0x002}, 0x0020, t)
	//a = b
	testJMPConditional(0x0056, 0x0020, 0x00F1, [4]uint16{0x0002, 0x0002, 0x0001, 0x002}, 0x0020, t)
	// should not jump in false case
	testJMPConditional(0x0056, 0x0021, 0x00F1, [4]uint16{0x0001, 0x0003, 0x0001, 0x0001}, 0x0003, t)
}

func TestJMPAZ(t *testing.T) {
	ClearMem()
	// Jump is A is larger or zero flag
	// a larger
	testJMPConditional(0x0055, 0x0020, 0x00F1, [4]uint16{0x0002, 0x0001, 0x0001, 0x002}, 0x0020, t)
	// zero flag (using and)
	testJMPConditional(0x0055, 0x0020, 0x00C1, [4]uint16{0x0001, 0x00FE, 0x0002, 0x0002}, 0x0020, t)

	// should not jump in false case
	testJMPConditional(0x0055, 0x0021, 0x00F1, [4]uint16{0x0001, 0x0003, 0x0001, 0x0001}, 0x0003, t)
}

func TestJMPEZ(t *testing.T) {
	ClearMem()
	// Jump if A = B or zero flag
	// a = b
	testJMPConditional(0x0053, 0x0020, 0x00F1, [4]uint16{0x0002, 0x0002, 0x0001, 0x002}, 0x0020, t)
	// zero flag (using and)
	testJMPConditional(0x0053, 0x0020, 0x00C1, [4]uint16{0x0001, 0x00FE, 0x0002, 0x0002}, 0x0020, t)

	// should not jump in false case
	testJMPConditional(0x0053, 0x0021, 0x00F1, [4]uint16{0x0001, 0x0003, 0x0001, 0x0001}, 0x0003, t)
}

func TestJMPCAE(t *testing.T) {
	ClearMem()
	// Jump if Carry OR A is Larger OR A = B

	// carry condition
	testJMPConditional(0x005E, 0x0090, 0x0081, [4]uint16{0x0004, 0xFFFF, 0x0001, 0x002}, 0x0090, t)

	// a larger
	testJMPConditional(0x005E, 0x0020, 0x00F1, [4]uint16{0x0002, 0x0001, 0x0001, 0x002}, 0x0020, t)

	// a = b
	testJMPConditional(0x005E, 0x0020, 0x00F1, [4]uint16{0x0002, 0x0002, 0x0001, 0x002}, 0x0020, t)

	// should not jump in false case
	testJMPConditional(0x005E, 0x0021, 0x00F1, [4]uint16{0x0001, 0x0003, 0x0001, 0x0001}, 0x0003, t)
}

func TestJMPCAZ(t *testing.T) {
	ClearMem()
	// Jump if Carry OR A is Larger OR zero flag

	// carry condition
	testJMPConditional(0x005D, 0x0090, 0x0081, [4]uint16{0x0004, 0xFFFF, 0x0001, 0x002}, 0x0090, t)

	// a larger
	testJMPConditional(0x005D, 0x0020, 0x00F1, [4]uint16{0x0002, 0x0001, 0x0001, 0x002}, 0x0020, t)

	// zero flag (using and)
	testJMPConditional(0x005D, 0x0020, 0x00C1, [4]uint16{0x0001, 0x00FE, 0x0002, 0x0002}, 0x0020, t)

	// should not jump in false case
	testJMPConditional(0x005D, 0x0021, 0x00F1, [4]uint16{0x0001, 0x0003, 0x0001, 0x0001}, 0x0003, t)
}

func TestJMPCEZ(t *testing.T) {
	ClearMem()
	// Jump if Carry OR a = b OR zero flag

	// carry condition
	testJMPConditional(0x005B, 0x0090, 0x0081, [4]uint16{0x0004, 0xFFFF, 0x0001, 0x002}, 0x0090, t)

	// a = b
	testJMPConditional(0x005B, 0x0020, 0x00F1, [4]uint16{0x0002, 0x0002, 0x0001, 0x002}, 0x0020, t)

	// zero flag (using and)
	testJMPConditional(0x005B, 0x0020, 0x00C1, [4]uint16{0x0001, 0x00FE, 0x0002, 0x0002}, 0x0020, t)

	// should not jump in false case
	testJMPConditional(0x005B, 0x0021, 0x00F1, [4]uint16{0x0001, 0x0003, 0x0001, 0x0001}, 0x0003, t)
}

func TestJMPAEZ(t *testing.T) {
	ClearMem()
	// Jump if a is larger OR a = b OR zero flag

	// a larger
	testJMPConditional(0x0057, 0x0020, 0x00F1, [4]uint16{0x0002, 0x0001, 0x0001, 0x002}, 0x0020, t)

	// a = b
	testJMPConditional(0x0057, 0x0020, 0x00F1, [4]uint16{0x0002, 0x0002, 0x0001, 0x002}, 0x0020, t)

	// zero flag (using and)
	testJMPConditional(0x0057, 0x0020, 0x00C1, [4]uint16{0x0001, 0x00FE, 0x0002, 0x0002}, 0x0020, t)

	// should not jump in false case
	testJMPConditional(0x0057, 0x0021, 0x00F1, [4]uint16{0x0001, 0x0003, 0x0001, 0x0001}, 0x0003, t)
}

func TestJMPCAEZ(t *testing.T) {
	ClearMem()
	// Jump if Carry OR a is larger OR a = b OR zero flag

	// carry condition
	testJMPConditional(0x005F, 0x0090, 0x0081, [4]uint16{0x0004, 0xFFFF, 0x0001, 0x002}, 0x0090, t)

	// a larger
	testJMPConditional(0x005F, 0x0020, 0x00F1, [4]uint16{0x0002, 0x0001, 0x0001, 0x002}, 0x0020, t)

	// a = b
	testJMPConditional(0x005F, 0x0020, 0x00F1, [4]uint16{0x0002, 0x0002, 0x0001, 0x002}, 0x0020, t)

	// zero flag (using and)
	testJMPConditional(0x005F, 0x0020, 0x00C1, [4]uint16{0x0001, 0x00FE, 0x0002, 0x0002}, 0x0020, t)

	// should not jump in false case
	testJMPConditional(0x005F, 0x0021, 0x00F1, [4]uint16{0x0001, 0x0003, 0x0001, 0x0001}, 0x0003, t)
}

func testJMPConditional(jmpConditionInstr, destination, initialInstr uint16, inputRegisters [4]uint16, expectedIAR uint16, t *testing.T) {
	c := SetUpCPU()

	var insAddr uint16 = 0x0000

	setMemoryLocation(c, insAddr, initialInstr)
	setMemoryLocation(c, insAddr+1, jmpConditionInstr)
	setMemoryLocation(c, insAddr+2, destination)

	c.SetIAR(insAddr)

	setRegisters(c, inputRegisters)

	doFetchDecodeExecute(c)
	doFetchDecodeExecute(c)

	// check IAR
	checkIAR(c, expectedIAR, t)
}

func TestCLF(t *testing.T) {
	ClearMem()
	// carry + zero + greater
	testCLF(0x0081, [4]uint16{0xFFFF, 0x0001, 0x0000, 0x0000}, t)

	// equal flag
	testCLF(0x0081, [4]uint16{0x0001, 0x0001, 0x0000, 0x0000}, t)

	// all flags should be false anyway
	testCLF(0x0081, [4]uint16{0x0001, 0x0002, 0x0000, 0x0000}, t)
}

func testCLF(initialInstruction uint16, initialRegisters [4]uint16, t *testing.T) {
	var insAddr uint16 = 0x0000

	c := SetUpCPU()

	setMemoryLocation(c, insAddr, initialInstruction)
	setMemoryLocation(c, insAddr+1, 0x0060)

	c.SetIAR(insAddr)

	setRegisters(c, initialRegisters)

	doFetchDecodeExecute(c)
	doFetchDecodeExecute(c)

	checkFlagsRegister(c, false, false, false, false, t)
}

func TestALUAdd(t *testing.T) {
	ClearMem()

	var inputs [4]uint16 = [4]uint16{0x0002, 0x0003, 0xFD04, 0x0005}
	testInstruction(0x0080, inputs, [4]uint16{inputs[0] + inputs[0], inputs[1], inputs[2], inputs[3]}, t)
	testInstruction(0x0081, inputs, [4]uint16{inputs[0], inputs[1] + inputs[0], inputs[2], inputs[3]}, t)
	testInstruction(0x0082, inputs, [4]uint16{inputs[0], inputs[1], inputs[2] + inputs[0], inputs[3]}, t)
	testInstruction(0x0083, inputs, [4]uint16{inputs[0], inputs[1], inputs[2], inputs[3] + inputs[0]}, t)

	testInstruction(0x0084, inputs, [4]uint16{inputs[0] + inputs[1], inputs[1], inputs[2], inputs[3]}, t)
	testInstruction(0x0085, inputs, [4]uint16{inputs[0], inputs[1] + inputs[1], inputs[2], inputs[3]}, t)
	testInstruction(0x0086, inputs, [4]uint16{inputs[0], inputs[1], inputs[2] + inputs[1], inputs[3]}, t)
	testInstruction(0x0087, inputs, [4]uint16{inputs[0], inputs[1], inputs[2], inputs[3] + inputs[1]}, t)

	testInstruction(0x0088, inputs, [4]uint16{inputs[0] + inputs[2], inputs[1], inputs[2], inputs[3]}, t)
	testInstruction(0x0089, inputs, [4]uint16{inputs[0], inputs[1] + inputs[2], inputs[2], inputs[3]}, t)
	testInstruction(0x008A, inputs, [4]uint16{inputs[0], inputs[1], inputs[2] + inputs[2], inputs[3]}, t)
	testInstruction(0x008B, inputs, [4]uint16{inputs[0], inputs[1], inputs[2], inputs[3] + inputs[2]}, t)

	testInstruction(0x008C, inputs, [4]uint16{inputs[0] + inputs[3], inputs[1], inputs[2], inputs[3]}, t)
	testInstruction(0x008D, inputs, [4]uint16{inputs[0], inputs[1] + inputs[3], inputs[2], inputs[3]}, t)
	testInstruction(0x008E, inputs, [4]uint16{inputs[0], inputs[1], inputs[2] + inputs[3], inputs[3]}, t)
	testInstruction(0x008F, inputs, [4]uint16{inputs[0], inputs[1], inputs[2], inputs[3] + inputs[3]}, t)
}

func TestALUAddWithCarry(t *testing.T) {
	ClearMem()
	testALUAddWithCarry(
		0x0080,
		[4]uint16{0xFFFE, 0x0000, 0x0000, 0x0000},
		[4]uint16{0x0001, 0x0000, 0x0000, 0x0000},
		t,
	)

	testALUAddWithCarry(
		0x0081,
		[4]uint16{0xFFFE, 0x0005, 0x0000, 0x0000},
		[4]uint16{0x0000, 0x0001, 0x0000, 0x0000},
		t,
	)

	testALUAddWithCarry(
		0x0082,
		[4]uint16{0xFFFE, 0x0000, 0x0005, 0x0000},
		[4]uint16{0x0000, 0x0000, 0x0001, 0x0000},
		t,
	)

	testALUAddWithCarry(
		0x0083,
		[4]uint16{0xFFFE, 0x0000, 0x0000, 0x0005},
		[4]uint16{0x0000, 0x0000, 0x0000, 0x0001},
		t,
	)
}

func testALUAddWithCarry(instruction uint16, inputRegisters, expectedOutputRegisters [4]uint16, t *testing.T) {
	c := SetUpCPU()

	setMemoryLocation(c, 0x0000, instruction)
	setMemoryLocation(c, 0x0001, instruction)

	c.SetIAR(0x0000)

	setRegisters(c, inputRegisters)
	doFetchDecodeExecute(c)
	setRegisters(c, [4]uint16{0x0000, 0x0000, 0x0000, 0x0000}) // zeros so that we can see the carry flag cause a change
	doFetchDecodeExecute(c)

	checkRegisters(c, expectedOutputRegisters[0], expectedOutputRegisters[1], expectedOutputRegisters[2], expectedOutputRegisters[3], t)
}

func TestALUNOT(t *testing.T) {
	ClearMem()
	var inputs [4]uint16 = [4]uint16{0xFFFF, 0xFEFE, 0xFDFD, 0xFCFC}
	testInstruction(0x00B0, inputs, [4]uint16{0x0000, 0xFEFE, 0xFDFD, 0xFCFC}, t)
	testInstruction(0x00B5, inputs, [4]uint16{0xFFFF, 0x0101, 0xFDFD, 0xFCFC}, t)
	testInstruction(0x00BA, inputs, [4]uint16{0xFFFF, 0xFEFE, 0x0202, 0xFCFC}, t)
	testInstruction(0x00BF, inputs, [4]uint16{0xFFFF, 0xFEFE, 0xFDFD, 0x0303}, t)
}

func TestALUAND(t *testing.T) {
	ClearMem()
	var inputs [4]uint16 = [4]uint16{0x0002, 0xABC3, 0x0004, 0x0005}
	testInstruction(0x00C0, inputs, [4]uint16{inputs[0] & inputs[0], inputs[1], inputs[2], inputs[3]}, t)
	testInstruction(0x00C1, inputs, [4]uint16{inputs[0], inputs[1] & inputs[0], inputs[2], inputs[3]}, t)
	testInstruction(0x00C2, inputs, [4]uint16{inputs[0], inputs[1], inputs[2] & inputs[0], inputs[3]}, t)
	testInstruction(0x00C3, inputs, [4]uint16{inputs[0], inputs[1], inputs[2], inputs[3] & inputs[0]}, t)

	testInstruction(0x00C4, inputs, [4]uint16{inputs[0] & inputs[1], inputs[1], inputs[2], inputs[3]}, t)
	testInstruction(0x00C5, inputs, [4]uint16{inputs[0], inputs[1] & inputs[1], inputs[2], inputs[3]}, t)
	testInstruction(0x00C6, inputs, [4]uint16{inputs[0], inputs[1], inputs[2] & inputs[1], inputs[3]}, t)
	testInstruction(0x00C7, inputs, [4]uint16{inputs[0], inputs[1], inputs[2], inputs[3] & inputs[1]}, t)

	testInstruction(0x00C8, inputs, [4]uint16{inputs[0] & inputs[2], inputs[1], inputs[2], inputs[3]}, t)
	testInstruction(0x00C9, inputs, [4]uint16{inputs[0], inputs[1] & inputs[2], inputs[2], inputs[3]}, t)
	testInstruction(0x00CA, inputs, [4]uint16{inputs[0], inputs[1], inputs[2] & inputs[2], inputs[3]}, t)
	testInstruction(0x00CB, inputs, [4]uint16{inputs[0], inputs[1], inputs[2], inputs[3] & inputs[2]}, t)

	testInstruction(0x00CC, inputs, [4]uint16{inputs[0] & inputs[3], inputs[1], inputs[2], inputs[3]}, t)
	testInstruction(0x00CD, inputs, [4]uint16{inputs[0], inputs[1] & inputs[3], inputs[2], inputs[3]}, t)
	testInstruction(0x00CE, inputs, [4]uint16{inputs[0], inputs[1], inputs[2] & inputs[3], inputs[3]}, t)
	testInstruction(0x00CF, inputs, [4]uint16{inputs[0], inputs[1], inputs[2], inputs[3] & inputs[3]}, t)
}

func TestALUOR(t *testing.T) {
	ClearMem()
	var inputs [4]uint16 = [4]uint16{0x2092, 0x0091, 0xCF45, 0x00AF}

	testInstruction(0x00D0, inputs, [4]uint16{inputs[0] | inputs[0], inputs[1], inputs[2], inputs[3]}, t)
	testInstruction(0x00D1, inputs, [4]uint16{inputs[0], inputs[1] | inputs[0], inputs[2], inputs[3]}, t)
	testInstruction(0x00D2, inputs, [4]uint16{inputs[0], inputs[1], inputs[2] | inputs[0], inputs[3]}, t)
	testInstruction(0x00D3, inputs, [4]uint16{inputs[0], inputs[1], inputs[2], inputs[3] | inputs[0]}, t)

	testInstruction(0x00D4, inputs, [4]uint16{inputs[0] | inputs[1], inputs[1], inputs[2], inputs[3]}, t)
	testInstruction(0x00D5, inputs, [4]uint16{inputs[0], inputs[1] | inputs[1], inputs[2], inputs[3]}, t)
	testInstruction(0x00D6, inputs, [4]uint16{inputs[0], inputs[1], inputs[2] | inputs[1], inputs[3]}, t)
	testInstruction(0x00D7, inputs, [4]uint16{inputs[0], inputs[1], inputs[2], inputs[3] | inputs[1]}, t)

	testInstruction(0x00D8, inputs, [4]uint16{inputs[0] | inputs[2], inputs[1], inputs[2], inputs[3]}, t)
	testInstruction(0x00D9, inputs, [4]uint16{inputs[0], inputs[1] | inputs[2], inputs[2], inputs[3]}, t)
	testInstruction(0x00DA, inputs, [4]uint16{inputs[0], inputs[1], inputs[2] | inputs[2], inputs[3]}, t)
	testInstruction(0x00DB, inputs, [4]uint16{inputs[0], inputs[1], inputs[2], inputs[3] | inputs[2]}, t)

	testInstruction(0x00DC, inputs, [4]uint16{inputs[0] | inputs[3], inputs[1], inputs[2], inputs[3]}, t)
	testInstruction(0x00DD, inputs, [4]uint16{inputs[0], inputs[1] | inputs[3], inputs[2], inputs[3]}, t)
	testInstruction(0x00DE, inputs, [4]uint16{inputs[0], inputs[1], inputs[2] | inputs[3], inputs[3]}, t)
	testInstruction(0x00DF, inputs, [4]uint16{inputs[0], inputs[1], inputs[2], inputs[3] | inputs[3]}, t)
}

func TestALUXOR(t *testing.T) {
	ClearMem()
	var inputs [4]uint16 = [4]uint16{0x0092, 0x8791, 0x0045, 0xD1AF}

	testInstruction(0x00E0, inputs, [4]uint16{inputs[0] ^ inputs[0], inputs[1], inputs[2], inputs[3]}, t)
	testInstruction(0x00E1, inputs, [4]uint16{inputs[0], inputs[1] ^ inputs[0], inputs[2], inputs[3]}, t)
	testInstruction(0x00E2, inputs, [4]uint16{inputs[0], inputs[1], inputs[2] ^ inputs[0], inputs[3]}, t)
	testInstruction(0x00E3, inputs, [4]uint16{inputs[0], inputs[1], inputs[2], inputs[3] ^ inputs[0]}, t)

	testInstruction(0x00E4, inputs, [4]uint16{inputs[0] ^ inputs[1], inputs[1], inputs[2], inputs[3]}, t)
	testInstruction(0x00E5, inputs, [4]uint16{inputs[0], inputs[1] ^ inputs[1], inputs[2], inputs[3]}, t)
	testInstruction(0x00E6, inputs, [4]uint16{inputs[0], inputs[1], inputs[2] ^ inputs[1], inputs[3]}, t)
	testInstruction(0x00E7, inputs, [4]uint16{inputs[0], inputs[1], inputs[2], inputs[3] ^ inputs[1]}, t)

	testInstruction(0x00E8, inputs, [4]uint16{inputs[0] ^ inputs[2], inputs[1], inputs[2], inputs[3]}, t)
	testInstruction(0x00E9, inputs, [4]uint16{inputs[0], inputs[1] ^ inputs[2], inputs[2], inputs[3]}, t)
	testInstruction(0x00EA, inputs, [4]uint16{inputs[0], inputs[1], inputs[2] ^ inputs[2], inputs[3]}, t)
	testInstruction(0x00EB, inputs, [4]uint16{inputs[0], inputs[1], inputs[2], inputs[3] ^ inputs[2]}, t)

	testInstruction(0x00EC, inputs, [4]uint16{inputs[0] ^ inputs[3], inputs[1], inputs[2], inputs[3]}, t)
	testInstruction(0x00ED, inputs, [4]uint16{inputs[0], inputs[1] ^ inputs[3], inputs[2], inputs[3]}, t)
	testInstruction(0x00EE, inputs, [4]uint16{inputs[0], inputs[1], inputs[2] ^ inputs[3], inputs[3]}, t)
	testInstruction(0x00EF, inputs, [4]uint16{inputs[0], inputs[1], inputs[2], inputs[3] ^ inputs[3]}, t)
}

func TestCMP(t *testing.T) {
	ClearMem()
	var inputs [4]uint16 = [4]uint16{0xAB92, 0x0091, 0x0045, 0x00AF}

	var instruction uint16 = 0x00F0
	for a := 0; a < 4; a++ {
		for b := 0; b < 4; b++ {
			testCMP(instruction, inputs, inputs, a, b, t)
			instruction++
		}
	}

	var zeroes [4]uint16 = [4]uint16{0x0000, 0x0000, 0x0000, 0x0000}
	instruction = 0x00F0
	for a := 0; a < 4; a++ {
		for b := 0; b < 4; b++ {
			testCMP(instruction, zeroes, zeroes, a, b, t)
			instruction++
		}
	}
}

func testCMP(instruction uint16, inputRegisters [4]uint16, expectedOutputRegisters [4]uint16, compareA, compareB int, t *testing.T) {
	c := SetUpCPU()

	setMemoryLocation(c, 0x0000, instruction)

	for i, r := range inputRegisters {
		setRegister(c, i, r)
	}

	c.SetIAR(0x0000)

	doFetchDecodeExecute(c)

	checkRegisters(c, expectedOutputRegisters[0], expectedOutputRegisters[1], expectedOutputRegisters[2], expectedOutputRegisters[3], t)
	checkFlagsRegister(c, false, inputRegisters[compareA] > inputRegisters[compareB], inputRegisters[compareA] == inputRegisters[compareB], false, t)
}

func testInstruction(instruction uint16, inputRegisters [4]uint16, expectedOutputRegisters [4]uint16, t *testing.T) {
	c := SetUpCPU()
	setMemoryLocation(c, 0x0000, instruction)

	for i, r := range inputRegisters {
		setRegister(c, i, r)
	}

	c.SetIAR(0x0000)

	doFetchDecodeExecute(c)

	checkRegisters(c, expectedOutputRegisters[0], expectedOutputRegisters[1], expectedOutputRegisters[2], expectedOutputRegisters[3], t)
}

func TestALUShiftLeft(t *testing.T) {
	ClearMem()
	var ones [4]uint16 = [4]uint16{0x0001, 0x0001, 0x0001, 0x0001}
	var shifts uint16
	for shifts = 0; shifts < 16; shifts++ {
		testShift(0x00A0, ones, [4]uint16{1 << shifts, 0x0001, 0x0001, 0x0001}, shifts, t)
		testShift(0x00A5, ones, [4]uint16{0x0001, 1 << shifts, 0x0001, 0x0001}, shifts, t)
		testShift(0x00AA, ones, [4]uint16{0x0001, 0x0001, 1 << shifts, 0x0001}, shifts, t)
		testShift(0x00AF, ones, [4]uint16{0x0001, 0x0001, 0x0001, 1 << shifts}, shifts, t)
	}
}

func TestALUShiftRight(t *testing.T) {
	var input [4]uint16 = [4]uint16{0x8000, 0x8000, 0x8000, 0x8000}
	var shifts uint16
	for shifts = 0; shifts < 16; shifts++ {
		testShift(0x0090, input, [4]uint16{0x8000 >> shifts, 0x8000, 0x8000, 0x8000}, shifts, t)
		testShift(0x0095, input, [4]uint16{0x8000, 0x8000 >> shifts, 0x8000, 0x8000}, shifts, t)
		testShift(0x009A, input, [4]uint16{0x8000, 0x8000, 0x8000 >> shifts, 0x8000}, shifts, t)
		testShift(0x009F, input, [4]uint16{0x8000, 0x8000, 0x8000, 0x8000 >> shifts}, shifts, t)
	}
}

func testShift(instruction uint16, inputRegisters [4]uint16, expectedOutputRegisters [4]uint16, shifts uint16, t *testing.T) {
	c := SetUpCPU()

	var i uint16
	for i = 0; i < shifts; i++ {
		setMemoryLocation(c, i, instruction)
	}

	for i, r := range inputRegisters {
		setRegister(c, i, r)
	}

	c.SetIAR(0x0000)

	for i = 0; i < shifts; i++ {
		doFetchDecodeExecute(c)
	}

	checkRegisters(c, expectedOutputRegisters[0], expectedOutputRegisters[1], expectedOutputRegisters[2], expectedOutputRegisters[3], t)
}

func TestSubtract(t *testing.T) {
	ClearMem()
	testSubtract(0, 0, t)
	testSubtract(1, 0, t)
	testSubtract(37, 21, t)
	testSubtract(0x00FF, 0x00FF, t)
	testSubtract(10, 3, t)
	testSubtract(100, 99, t)
}

func testSubtract(inputA, inputB uint16, t *testing.T) {
	c := SetUpCPU()

	setRegisters(c, [4]uint16{inputA, inputB, 1, 0})
	setMemoryLocation(c, 0x0000, 0x00B5) // NOT
	setMemoryLocation(c, 0x0001, 0x0089) // ADD R2, R1
	setMemoryLocation(c, 0x0002, 0x0060) // CLF
	setMemoryLocation(c, 0x0003, 0x0081) // ADD R0, R1

	c.SetIAR(0x0000)

	doFetchDecodeExecute(c)
	doFetchDecodeExecute(c)
	doFetchDecodeExecute(c)
	doFetchDecodeExecute(c)

	checkRegister(c, 1, inputA-inputB, t)
}

func TestMultiply(t *testing.T) {
	ClearMem()
	testMultiply(0, 0, t)
	testMultiply(1, 1, t)
	testMultiply(1, 2, t)
	testMultiply(2, 1, t)
	testMultiply(5, 5, t)
	testMultiply(8, 12, t)
	testMultiply(19, 13, t)
}

func testMultiply(inputA, inputB uint16, t *testing.T) {
	c := SetUpCPU()

	setMemoryLocation(c, 50, 0x0023) // DATA R3
	setMemoryLocation(c, 51, 0x0001) // .. 1
	setMemoryLocation(c, 52, 0x00EA) // XOR R2, R2
	setMemoryLocation(c, 53, 0x0060) // CLF
	setMemoryLocation(c, 54, 0x0090) // SHR R0
	setMemoryLocation(c, 55, 0x0058) // JC
	setMemoryLocation(c, 56, 59)     // ...addr 59
	setMemoryLocation(c, 57, 0x0040) // JMP
	setMemoryLocation(c, 58, 61)     // ...addr 61
	setMemoryLocation(c, 59, 0x0060) // CLF
	setMemoryLocation(c, 60, 0x0086) // ADD R1, R2
	setMemoryLocation(c, 61, 0x0060) // CLF
	setMemoryLocation(c, 62, 0x00A5) // SHL R1
	setMemoryLocation(c, 63, 0x00AF) // SHL R3
	setMemoryLocation(c, 64, 0x0058) // JC
	setMemoryLocation(c, 65, 68)     // ...addr 68
	setMemoryLocation(c, 66, 0x0040) // JMP
	setMemoryLocation(c, 67, 53)     // ...addr 53

	setRegisters(c, [4]uint16{inputA, inputB, 0, 0})

	c.SetIAR(50)

	for {
		doFetchDecodeExecute(c)
		if c.iar.Value() >= 68 {
			break
		}
	}

	checkRegister(c, 2, inputA*inputB, t)
}

func TestIOInputInstruction(t *testing.T) {
	ClearMem()
	// IN Data, RB
	zeros := [4]uint16{0, 0, 0, 0}
	testIOInputInstruction(0x0070, zeros, [4]uint16{0x00DA, 0x0000, 0x0000, 0x0000}, t)
	testIOInputInstruction(0x0071, zeros, [4]uint16{0x0000, 0x00DA, 0x0000, 0x0000}, t)
	testIOInputInstruction(0x0072, zeros, [4]uint16{0x0000, 0x0000, 0x00DA, 0x0000}, t)
	testIOInputInstruction(0x0073, zeros, [4]uint16{0x0000, 0x0000, 0x0000, 0x00DA}, t)
	// IN Addr, RB
	testIOInputInstruction(0x0074, zeros, [4]uint16{0x00AD, 0x0000, 0x0000, 0x0000}, t)
	testIOInputInstruction(0x0075, zeros, [4]uint16{0x0000, 0x00AD, 0x0000, 0x0000}, t)
	testIOInputInstruction(0x0076, zeros, [4]uint16{0x0000, 0x0000, 0x00AD, 0x0000}, t)
	testIOInputInstruction(0x0077, zeros, [4]uint16{0x0000, 0x0000, 0x0000, 0x00AD}, t)
}

func testIOInputInstruction(instruction uint16, inputRegisters, expectedRegisters [4]uint16, t *testing.T) {
	c := SetUpCPU()
	c.ConnectPeripheral(NewDumbPeripheral())

	setMemoryLocation(c, 0x0000, instruction)
	setRegisters(c, inputRegisters)

	c.SetIAR(0x0000)

	doFetchDecodeExecute(c)

	checkRegisters(c, expectedRegisters[0], expectedRegisters[1], expectedRegisters[2], expectedRegisters[3], t)
}

// should put a value in the peripheral register, and set the address/data flag on the peripheral
func TestIOOutputInstruction(t *testing.T) {
	ClearMem()
	// OUT Data, RB
	testIOOutputInstruction(0x0078, [4]uint16{0x00DD, 0x0008, 0x0007, 0x0006}, 0x00DD, true, false, t)
	testIOOutputInstruction(0x0079, [4]uint16{0x0009, 0x00DD, 0x0007, 0x0006}, 0x00DD, true, false, t)
	testIOOutputInstruction(0x007A, [4]uint16{0x0009, 0x0008, 0x00DD, 0x0006}, 0x00DD, true, false, t)
	testIOOutputInstruction(0x007B, [4]uint16{0x0009, 0x0008, 0x0007, 0x00DD}, 0x00DD, true, false, t)

	// OUT Addr, RB
	testIOOutputInstruction(0x007C, [4]uint16{0x00AA, 0x0008, 0x0007, 0x0006}, 0x00AA, false, true, t)
	testIOOutputInstruction(0x007D, [4]uint16{0x0009, 0x00AA, 0x0007, 0x0006}, 0x00AA, false, true, t)
	testIOOutputInstruction(0x007E, [4]uint16{0x0009, 0x0008, 0x00AA, 0x0006}, 0x00AA, false, true, t)
	testIOOutputInstruction(0x007F, [4]uint16{0x0009, 0x0008, 0x0007, 0x00AA}, 0x00AA, false, true, t)
}

func testIOOutputInstruction(instruction uint16, inputRegisters [4]uint16, expectedPeripheralValue uint16, expectedDataMode, expectedAddressMode bool, t *testing.T) {
	c := SetUpCPU()
	peripheral := NewDumbPeripheral()
	c.ConnectPeripheral(peripheral)

	setMemoryLocation(c, 0x0000, instruction)
	setRegisters(c, inputRegisters)

	c.SetIAR(0x0000)

	doFetchDecodeExecute(c)

	if peripheral.value.Value() != expectedPeripheralValue {
		t.FailNow()
	}

	if peripheral.outputDataMode != expectedDataMode {
		t.FailNow()
	}

	if peripheral.outputAddressMode != expectedAddressMode {
		t.FailNow()
	}
}

// Dumb peripheral that just contains a register
// In 'ENABLE' state:
//     emits the values 0x00DA for data mode, and 0x00AD for address mode
// In SET state:
//    sets the flags outputDataMode and outputAddressMode based on the data/address value on the bus
type DumbPeripheral struct {
	ioBus   *components.IOBus
	mainBus *components.Bus

	value             components.Register
	outputDataMode    bool
	outputAddressMode bool
}

func NewDumbPeripheral() *DumbPeripheral {
	p := new(DumbPeripheral)
	return p
}

func (p *DumbPeripheral) Connect(ioBus *components.IOBus, mainBus *components.Bus) {
	p.ioBus = ioBus
	p.mainBus = mainBus
	p.value = *components.NewRegister("P", p.mainBus, p.mainBus)

}

func (p *DumbPeripheral) Update() {
	p.updateEnabled()
	p.value.Update()
	p.updateSet()
	p.value.Update()
}

func (p *DumbPeripheral) refreshValue(addressValue, dataValue uint16) {
	p.value.Set()
	p.value.Update()

	if p.ioBus.GetOutputWire(components.DATA_OR_ADDRESS) {
		//address mode
		p.mainBus.SetValue(addressValue)
	} else {
		p.mainBus.SetValue(dataValue)
	}
	p.value.Update()
	p.value.Unset()
	p.value.Update()
}

func (p *DumbPeripheral) updateEnabled() {
	if p.ioBus.GetOutputWire(components.CLOCK_ENABLE) {
		p.refreshValue(0x00AD, 0x00DA)
		p.value.Enable()
	} else {
		p.value.Disable()
	}
}

func (p *DumbPeripheral) updateSet() {
	if p.ioBus.GetOutputWire(components.CLOCK_SET) {
		p.outputDataMode = (p.ioBus.GetOutputWire(components.DATA_OR_ADDRESS) == false)
		p.outputAddressMode = (p.ioBus.GetOutputWire(components.DATA_OR_ADDRESS) == true)
		p.value.Set()
	} else {
		p.value.Unset()
	}
}

func doFetchDecodeExecute(c *CPU) {
	for i := 0; i < 6; i++ {
		c.Step()
	}
}

func setMemoryLocation(c *CPU, address uint16, value uint16) {
	c.memory.AddressRegister.Set()
	c.mainBus.SetValue(address)
	c.memory.Update()

	c.memory.AddressRegister.Unset()
	c.memory.Update()

	c.mainBus.SetValue(value)
	c.memory.Set()
	c.memory.Update()

	c.memory.Unset()
	c.memory.Update()
}

func setMemoryLocation2(m *memory.Memory64K, address uint16, value uint16) {
	m.AddressRegister.Set()
	BUS.SetValue(address)
	m.Update()

	m.AddressRegister.Unset()
	m.Update()

	BUS.SetValue(value)
	m.Set()
	m.Update()

	m.Unset()
	m.Update()
}

func setRegister(c *CPU, register int, value uint16) {
	switch register {
	case 0:
		c.gpReg0.Set()
		c.gpReg0.Update()
		c.mainBus.SetValue(value)
		c.gpReg0.Update()
		c.gpReg0.Unset()
		c.gpReg0.Update()
	case 1:
		c.gpReg1.Set()
		c.gpReg1.Update()
		c.mainBus.SetValue(value)
		c.gpReg1.Update()
		c.gpReg1.Unset()
		c.gpReg1.Update()
	case 2:
		c.gpReg2.Set()
		c.gpReg2.Update()
		c.mainBus.SetValue(value)
		c.gpReg2.Update()
		c.gpReg2.Unset()
		c.gpReg2.Update()
	case 3:
		c.gpReg3.Set()
		c.gpReg3.Update()
		c.mainBus.SetValue(value)
		c.gpReg3.Update()
		c.gpReg3.Unset()
		c.gpReg3.Update()
	}
}

func checkIAR(c *CPU, expValue uint16, t *testing.T) {
	if c.iar.Value() != expValue {
		t.Logf("Expected IAR to have value of: %X but got %X", expValue, c.iar.Value())
		t.FailNow()
	}
}

func checkIR(c *CPU, expValue uint16, t *testing.T) {
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

func checkRegister(c *CPU, register int, expectedValue uint16, t *testing.T) {
	var regValue uint16
	switch register {
	case 0:
		regValue = c.gpReg0.Value()
	case 1:
		regValue = c.gpReg1.Value()
	case 2:
		regValue = c.gpReg2.Value()
	case 3:
		regValue = c.gpReg3.Value()
	default:
		t.Logf("Unknown register %d", register)
		t.FailNow()
	}

	if regValue != expectedValue {
		t.Logf("Expected register %d to have value of: %X but got %X", register, expectedValue, regValue)
		t.FailNow()
	}
}

func checkRegisters(c *CPU, expReg0, expReg1, expReg2, expReg3 uint16, t *testing.T) {
	checkRegister(c, 0, expReg0, t)
	checkRegister(c, 1, expReg1, t)
	checkRegister(c, 2, expReg2, t)
	checkRegister(c, 3, expReg3, t)
}

func setRegisters(c *CPU, values [4]uint16) {
	for i, v := range values {
		setRegister(c, i, v)
	}
}
