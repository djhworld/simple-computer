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
