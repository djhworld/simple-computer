package main

import (
	"fmt"
	"math"

	"github.com/djhworld/simple-computer/asm"
)

const (
	sineTableBase = uint16(0xFF60) // 30 y-pixel values (one per character column)
	sineYScratch  = uint16(0xFF7E) // scratch: y*2 during draw; saved head during rotate
	sineYScratch2 = uint16(0xFF7F) // scratch: original y during y*30 computation
)

// sinePixelValue is the pixel pattern written to each display RAM cell.
// The display reads bits 7–0 of the 16-bit value: bit 7 = leftmost pixel.
// 0x003C = 0011 1100 → pixels 2,3,4,5 lit (four centre pixels per column).
const sinePixelValue = uint16(0x003C)

// precomputed y*30 values for one full cycle of sin over 30 columns.
// Storing y*30 avoids the multiplication in assembly — each column just loads
// and adds its x offset to get the display RAM address directly.
// Centre at y=80, amplitude 60px → y range 20–140, y*30 range 600–4200.
var sineTable = func() [30]uint16 {
	var t [30]uint16
	for i := range t {
		y := 80.0 + 60.0*math.Sin(2*math.Pi*float64(i)/30.0)
		t[i] = uint16(math.Round(y)) * 30
	}
	return t
}()

func sineWave(instructions asm.Instructions) {
	instructions.Add(asm.DEFLABEL{"main"})

	// write precomputed sine y-values into display RAM table
	for i, y := range sineTable {
		instructions.Add(
			asm.DATA{asm.REG0, asm.NUMBER{sineTableBase + uint16(i)}},
			asm.DATA{asm.REG1, asm.NUMBER{y}},
			asm.STORE{asm.REG0, asm.REG1},
		)
	}

	// select display adapter (stays selected for the entire program)
	instructions.AddBlocks(selectDisplayAdapter(asm.REG0))

	// animation loop: erase → advance wave → draw
	instructions.Add(asm.DEFLABEL{"sine-loop"})

	for x := 0; x < 30; x++ {
		instructions.AddBlocks(sineColumnWrite(uint16(x), 0x0000))
	}
	instructions.AddBlocks(sineAdvanceTable())
	for x := 0; x < 30; x++ {
		instructions.AddBlocks(sineColumnWrite(uint16(x), sinePixelValue))
	}

	instructions.Add(asm.JMP{asm.LABEL{"sine-loop"}})

	fmt.Println(instructions.String())
}

// sineColumnWrite reads the precomputed y*30 value from table[x], adds the
// column index x to get the display RAM address, then writes pixelValue to
// that row and the rows immediately above and below (3 rows total).
func sineColumnWrite(x, pixelValue uint16) []asm.Instruction {
	ins := asm.Instructions{}

	// R0 = y*30 (precomputed at generation time, stored in table)
	ins.Add(
		asm.DATA{asm.REG0, asm.NUMBER{sineTableBase + x}},
		asm.LOAD{asm.REG0, asm.REG0},
	)

	// R0 = y*30 + x = base display address
	ins.Add(
		asm.DATA{asm.REG1, asm.NUMBER{x}},
		asm.CLF{},
		asm.ADD{asm.REG1, asm.REG0},
	)

	// save base address — we need it for all three rows
	ins.Add(
		asm.DATA{asm.REG1, asm.NUMBER{sineYScratch2}},
		asm.STORE{asm.REG1, asm.REG0},
	)

	// row y-1: base_addr - 30
	ins.Add(
		asm.DATA{asm.REG1, asm.NUMBER{30}},
		asm.NOT{asm.REG1},
		asm.DATA{asm.REG2, asm.NUMBER{0x0001}},
		asm.CLF{},
		asm.ADD{asm.REG2, asm.REG1},          // R1 = -30
		asm.CLF{},
		asm.ADD{asm.REG1, asm.REG0},          // R0 = base_addr - 30
		asm.OUT{asm.DATA_MODE, asm.REG0},
		asm.DATA{asm.REG1, asm.NUMBER{pixelValue}},
		asm.OUT{asm.DATA_MODE, asm.REG1},
	)

	// row y: base_addr
	ins.Add(
		asm.DATA{asm.REG0, asm.NUMBER{sineYScratch2}},
		asm.LOAD{asm.REG0, asm.REG0},
		asm.OUT{asm.DATA_MODE, asm.REG0},
		asm.DATA{asm.REG1, asm.NUMBER{pixelValue}},
		asm.OUT{asm.DATA_MODE, asm.REG1},
	)

	// row y+1: base_addr + 30
	ins.Add(
		asm.DATA{asm.REG1, asm.NUMBER{30}},
		asm.CLF{},
		asm.ADD{asm.REG1, asm.REG0},          // R0 = base_addr + 30
		asm.OUT{asm.DATA_MODE, asm.REG0},
		asm.DATA{asm.REG1, asm.NUMBER{pixelValue}},
		asm.OUT{asm.DATA_MODE, asm.REG1},
	)

	return ins.Get()
}

// sineAdvanceTable rotates the sine table left by one entry, so each frame
// the wave appears to scroll one character column to the left.
func sineAdvanceTable() []asm.Instruction {
	ins := asm.Instructions{}

	// save table[0] to scratch before it gets overwritten
	ins.Add(
		asm.DATA{asm.REG0, asm.NUMBER{sineTableBase}},
		asm.LOAD{asm.REG0, asm.REG0},
		asm.DATA{asm.REG1, asm.NUMBER{sineYScratch}},
		asm.STORE{asm.REG1, asm.REG0},
	)

	// table[i] = table[i+1] for i = 0..28
	for i := 0; i < 29; i++ {
		ins.Add(
			asm.DATA{asm.REG0, asm.NUMBER{sineTableBase + uint16(i+1)}},
			asm.LOAD{asm.REG0, asm.REG0},
			asm.DATA{asm.REG1, asm.NUMBER{sineTableBase + uint16(i)}},
			asm.STORE{asm.REG1, asm.REG0},
		)
	}

	// table[29] = saved table[0]
	ins.Add(
		asm.DATA{asm.REG0, asm.NUMBER{sineYScratch}},
		asm.LOAD{asm.REG0, asm.REG0},
		asm.DATA{asm.REG1, asm.NUMBER{sineTableBase + 29}},
		asm.STORE{asm.REG1, asm.REG0},
	)

	return ins.Get()
}
