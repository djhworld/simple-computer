package main

import (
	"github.com/djhworld/simple-computer/components"
	"github.com/djhworld/simple-computer/cpu"
	"github.com/djhworld/simple-computer/memory"
	"time"
)

func main() {
	b := components.NewBus(16)
	m := memory.NewMemory64K(b)
	c := cpu.NewCPU(b, m)

	setMemoryLocation(m, b, 46, 0x20) // DATA R0
	setMemoryLocation(m, b, 47, 0x07) // 7
	setMemoryLocation(m, b, 48, 0x21) // DATA R1
	setMemoryLocation(m, b, 49, 0x0A) // 10
	setMemoryLocation(m, b, 50, 0x23) // DATA R3
	setMemoryLocation(m, b, 51, 0x01) // .. 1
	setMemoryLocation(m, b, 52, 0xEA) // XOR R2, R2
	setMemoryLocation(m, b, 53, 0x60) // CLF
	setMemoryLocation(m, b, 54, 0x90) // SHR R0
	setMemoryLocation(m, b, 55, 0x58) // JC
	setMemoryLocation(m, b, 56, 59)   // ...addr 59
	setMemoryLocation(m, b, 57, 0x40) // JMP
	setMemoryLocation(m, b, 58, 61)   // ...addr 61
	setMemoryLocation(m, b, 59, 0x60) // CLF
	setMemoryLocation(m, b, 60, 0x86) // ADD R1, R2
	setMemoryLocation(m, b, 61, 0x60) // CLF
	setMemoryLocation(m, b, 62, 0xA5) // SHL R1
	setMemoryLocation(m, b, 63, 0xAF) // SHL R3
	setMemoryLocation(m, b, 64, 0x58) // JC
	setMemoryLocation(m, b, 65, 68)   // ...addr 68
	setMemoryLocation(m, b, 66, 0x40) // JMP
	setMemoryLocation(m, b, 67, 53)   // ...addr 53
	setMemoryLocation(m, b, 68, 0x23) // DATA R3
	setMemoryLocation(m, b, 69, 0xFF) // 255
	setMemoryLocation(m, b, 70, 0x1E) // ST R3, R2
	setMemoryLocation(m, b, 71, 0x40) // JUMP
	setMemoryLocation(m, b, 72, 46) // JUMP

	c.Run(500 * time.Microsecond)

}

func setMemoryLocation(m *memory.Memory64K, mainBus *components.Bus, address uint16, value uint16) {
	m.AddressRegister.Set()
	setBus(mainBus, address)
	m.Update()

	m.AddressRegister.Unset()
	m.Update()

	setBus(mainBus, value)
	m.Set()
	m.Update()

	m.Unset()
	m.Update()
}

func setBus(b *components.Bus, value uint16) {
	var x = 0
	for i := 16-1; i >= 0; i-- {
		r := (value & (1 << uint16(x)))
		if r != 0 {
			b.SetInputWire(i, true)
		} else {
			b.SetInputWire(i, false)
		}

		x++
	}
}
