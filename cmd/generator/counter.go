package main

import (
	"fmt"

	"github.com/djhworld/simple-computer/asm"
)

const (
	counterHexTableBase = uint16(0xFF40) // 16-entry ASCII hex digit lookup table
	counterValueAddr    = uint16(0xFF50) // 16-bit counter
	counterDelayAddr    = uint16(0xFF51) // delay loop scratch
)

func counter(instructions asm.Instructions) {
	instructions.Add(asm.DEFLABEL{"main"})

	// populate hex digit ASCII lookup table at 0xFF40–0xFF4F
	// '0'–'9' = 0x30–0x39
	for i := uint16(0); i < 10; i++ {
		instructions.Add(
			asm.DATA{asm.REG0, asm.NUMBER{counterHexTableBase + i}},
			asm.DATA{asm.REG1, asm.NUMBER{0x0030 + i}},
			asm.STORE{asm.REG0, asm.REG1},
		)
	}
	// 'A'–'F' = 0x41–0x46
	for i := uint16(0); i < 6; i++ {
		instructions.Add(
			asm.DATA{asm.REG0, asm.NUMBER{counterHexTableBase + 10 + i}},
			asm.DATA{asm.REG1, asm.NUMBER{0x0041 + i}},
			asm.STORE{asm.REG0, asm.REG1},
		)
	}

	// initialise counter to 0
	instructions.Add(
		asm.DATA{asm.REG0, asm.NUMBER{counterValueAddr}},
		asm.DATA{asm.REG1, asm.NUMBER{0x0000}},
		asm.STORE{asm.REG0, asm.REG1},
	)

	// main display loop: reset pen, render label + 4 hex digits, delay, increment
	instructions.Add(asm.DEFLABEL{"counter-loop"})
	instructions.AddBlocks(
		updatePenPosition(0x00F0),
		resetLinex(),
		renderString("COUNT: 0x"),
	)
	instructions.AddBlocks(counterRenderNibble(12)) // bits 15–12
	instructions.AddBlocks(counterRenderNibble(8))  // bits 11–8
	instructions.AddBlocks(counterRenderNibble(4))  // bits 7–4
	instructions.AddBlocks(counterRenderNibble(0))  // bits 3–0
	instructions.AddBlocks(counterDelay())
	instructions.AddBlocks(counterIncrement())
	instructions.Add(asm.JMP{asm.LABEL{"counter-loop"}})

	fmt.Println(instructions.String())
}

// counterRenderNibble shifts the counter right by `shift` bits (using CLF before each
// SHR so the carry-in is always 0), masks the lowest 4 bits, looks up the ASCII hex
// character in the table, and renders it via the font routine.
func counterRenderNibble(shift int) []asm.Instruction {
	ins := asm.Instructions{}

	ins.Add(
		asm.DATA{asm.REG0, asm.NUMBER{counterValueAddr}},
		asm.LOAD{asm.REG0, asm.REG0},
	)

	for i := 0; i < shift; i++ {
		ins.Add(asm.CLF{})
		ins.Add(asm.SHR{asm.REG0})
	}
	ins.Add(asm.CLF{})

	// R0 = R0 & 0x000F  (isolate nibble)
	ins.Add(
		asm.DATA{asm.REG1, asm.NUMBER{0x000F}},
		asm.AND{asm.REG1, asm.REG0},
	)

	// R1 = hexTableBase + nibble,  R0 = mem[R1]  (ASCII char)
	ins.Add(
		asm.DATA{asm.REG1, asm.NUMBER{counterHexTableBase}},
		asm.ADD{asm.REG0, asm.REG1},
		asm.LOAD{asm.REG1, asm.REG0},
	)

	// store ASCII in KEYCODE-REGISTER and render
	ins.Add(
		asm.DATA{asm.REG1, asm.SYMBOL{"KEYCODE-REGISTER"}},
		asm.STORE{asm.REG1, asm.REG0},
	)
	ins.AddBlocks(callRoutine("ROUTINE-io-drawFontCharacter"))

	return ins.Get()
}

// counterIncrement adds 1 to the value stored at counterValueAddr, wrapping at 0xFFFF.
func counterIncrement() []asm.Instruction {
	return []asm.Instruction{
		asm.DATA{asm.REG0, asm.NUMBER{counterValueAddr}},
		asm.LOAD{asm.REG0, asm.REG0},
		asm.DATA{asm.REG1, asm.NUMBER{0x0001}},
		asm.ADD{asm.REG1, asm.REG0}, // R0 = counter + 1
		asm.CLF{},
		asm.DATA{asm.REG1, asm.NUMBER{counterValueAddr}},
		asm.STORE{asm.REG1, asm.REG0},
	}
}

// counterDelay counts down from 0x4000 to slow the counter to a human-readable rate.
func counterDelay() []asm.Instruction {
	ins := asm.Instructions{}
	ins.Add(
		asm.DATA{asm.REG0, asm.NUMBER{counterDelayAddr}},
		asm.DATA{asm.REG1, asm.NUMBER{0x0100}},
		asm.STORE{asm.REG0, asm.REG1},

		asm.DEFLABEL{"counter-delay-loop"},
		asm.DATA{asm.REG0, asm.NUMBER{counterDelayAddr}},
		asm.LOAD{asm.REG0, asm.REG0},
		asm.DATA{asm.REG1, asm.NUMBER{0x0000}},
		asm.NOT{asm.REG1},            // R1 = 0xFFFF = -1
		asm.CLF{},
		asm.ADD{asm.REG1, asm.REG0},  // R0 = R0 - 1
		asm.CLF{},
		asm.DATA{asm.REG1, asm.NUMBER{counterDelayAddr}},
		asm.STORE{asm.REG1, asm.REG0},
		asm.DATA{asm.REG1, asm.NUMBER{0x0000}},
		asm.CMP{asm.REG0, asm.REG1},
		asm.JMPF{[]string{"E"}, asm.LABEL{"counter-exit-delay"}},
		asm.JMP{asm.LABEL{"counter-delay-loop"}},

		asm.DEFLABEL{"counter-exit-delay"},
	)
	return ins.Get()
}
