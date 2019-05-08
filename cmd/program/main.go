package main

import (
	"fmt"

	"github.com/djhworld/simple-computer/asm"
)

/*
11111000 = 0x0078
10000110 = 0x0086
10000010 = 0x0082
10000010 = 0x0082
10000010 = 0x0082
10000110 = 0x0086
11111000 = 0x0078
00000000
*/
var CHARACTERS map[rune][8]uint16 = map[rune][8]uint16{
	'A': [8]uint16{0x007C, 0x00C6, 0x0082, 0x00FE, 0x0082, 0x0082, 0x0082, 0x0000},
	'B': [8]uint16{0x00FC, 0x0086, 0x0082, 0x00FE, 0x0082, 0x0086, 0x00FC, 0x0000},
	'C': [8]uint16{0x007E, 0x00C0, 0x0080, 0x0080, 0x0080, 0x00C0, 0x007E, 0x0000},
	'D': [8]uint16{0x00F8, 0x0086, 0x0082, 0x0082, 0x0082, 0x0086, 0x00F8, 0x0000},
	'E': [8]uint16{0x007E, 0x00C0, 0x0080, 0x00FE, 0x0080, 0x00C0, 0x007E, 0x0000},
}

// important RAM areas
// 0x0000 - 0x03FF ASCII table
// 0x0400 - 0x0400 pen position
// 0x0401 - 0x0401 keycode register
// 0x0500 - 0xFEFD user code + memory
// 0xFEFE - 0xFEFF used to jump back to user code
// 0xFF00 - 0xFFFF temporary variables

const PEN_POSITION_ADDR uint16 = 0x0400
const KEYCODE_REGISTER uint16 = 0x0401
const ONE uint16 = 0x0001

const (
	DISPLAY_ADAPTER_ADDR  = uint16(0x0007)
	KEYBOARD_ADAPTER_ADDR = uint16(0x000F)
	USER_CODE_AREA        = uint16(0x0500)
)

func main() {
	instructions := asm.Instructions{}

	instructions.AddBlocks(
		[]asm.Instruction{asm.LABEL{"FONTDESCS"}},
		loadFontCharacterIntoFontRegion('A'),
		loadFontCharacterIntoFontRegion('B'),
		loadFontCharacterIntoFontRegion('C'),
		loadFontCharacterIntoFontRegion('D'),
		loadFontCharacterIntoFontRegion('E'),
	)
	// more chars...

	instructions.AddBlocks(updatePenPosition(0x02BF))

	// loop:
	// 	poll keyboard
	// 	if keycode != 0
	//   store keycode in memory
	//   jump to display character
	//   jump to loop:
	//  else:
	//	  jump to loop:

	instructions.AddBlocks(
		selectDisplayAdapter(asm.REG3),
		loadCharIntoKeycodeRegister('A'),
		loadFontCharacterIntoDisplayRAM("drawFontA"),
		loadCharIntoKeycodeRegister('B'),
		loadFontCharacterIntoDisplayRAM("drawFontB"),
		loadCharIntoKeycodeRegister('C'),
		loadFontCharacterIntoDisplayRAM("drawFontC"),
		loadCharIntoKeycodeRegister('D'),
		loadFontCharacterIntoDisplayRAM("drawFontD"),
		loadCharIntoKeycodeRegister('E'),
		loadFontCharacterIntoDisplayRAM("drawFontE"),
		deselectIO(asm.REG3),
	)

	fmt.Println(instructions.String())
}

func updatePenPosition(position uint16) []asm.Instruction {
	return []asm.Instruction{
		asm.DATA{asm.REG0, PEN_POSITION_ADDR},
		asm.DATA{asm.REG1, position},
		asm.STORE{asm.REG0, asm.REG1},
	}
}

func loadCharIntoKeycodeRegister(char rune) []asm.Instruction {
	return []asm.Instruction{
		asm.DATA{asm.REG0, KEYCODE_REGISTER},
		asm.DATA{asm.REG1, uint16(char)},
		asm.STORE{asm.REG0, asm.REG1},
	}
}

func selectDisplayAdapter(useRegister asm.REGISTER) []asm.Instruction {
	return []asm.Instruction{
		asm.DATA{useRegister, DISPLAY_ADAPTER_ADDR},
		asm.OUT{asm.ADDRESS_MODE, useRegister},
	}
}

func deselectIO(useRegister asm.REGISTER) []asm.Instruction {
	return []asm.Instruction{
		asm.XOR{useRegister, useRegister},
		asm.OUT{asm.ADDRESS_MODE, useRegister},
	}
}

func loadFontCharacterIntoDisplayRAM(labelPrefix string) []asm.Instruction {
	fontYAddr := uint16(0xFF00)
	lineWidth := uint16(0x001E)

	instructions := asm.Instructions{}

	instructions.Add(
		asm.DATA{asm.REG2, PEN_POSITION_ADDR},
		asm.LOAD{asm.REG2, asm.REG2},
	)

	// we can keep this value in reg2 to track where in display RAM we are writing
	penPositionRegister := asm.REG2

	// counter for what line of the font are we rendering
	instructions.Add(
		asm.DATA{asm.REG0, fontYAddr},
		asm.DATA{asm.REG1, 0x0000},
		asm.STORE{asm.REG0, asm.REG1},
	)

	// calculate memory position of font line
	// start of loop:
	instructions.Add(
		asm.LABEL{labelPrefix + "-STARTLOOP"},
		asm.DATA{asm.REG3, KEYCODE_REGISTER}, // load keycode
		asm.LOAD{asm.REG3, asm.REG3},
		asm.SHL{asm.REG3},
		asm.SHL{asm.REG3},
		asm.SHL{asm.REG3},             // memory address in RAM for start of font
		asm.DATA{asm.REG0, fontYAddr}, // fontY address
		asm.LOAD{asm.REG0, asm.REG0},  // load fontY
		asm.ADD{asm.REG0, asm.REG3},   // calculate memory position of fontstart+fontYinstructions = append(instructions, ADD{asm.REG0, asm.REG3})       // calculate memory position of fontstart+fontY

		//increment fontY by 1
		asm.DATA{asm.REG1, ONE},       // one
		asm.ADD{asm.REG1, asm.REG0},   // increment fontY by 1
		asm.DATA{asm.REG1, fontYAddr}, // fontY address
		asm.STORE{asm.REG1, asm.REG0}, // store new value of fontY in memory

		// load font line from memory
		asm.LOAD{asm.REG3, asm.REG0}, // load value from memory into reg0

		// write to display ram
		asm.OUT{asm.DATA_MODE, penPositionRegister}, // display RAM address
		asm.OUT{asm.DATA_MODE, asm.REG0},            // display RAM value
		asm.DATA{asm.REG1, lineWidth},
		asm.ADD{asm.REG1, penPositionRegister}, // move pen down by 1 line

		// check if we have rendered all 8 lines
		asm.DATA{asm.REG0, fontYAddr}, // fontY addr
		asm.LOAD{asm.REG0, asm.REG0},  //load fontY into reg0
		asm.DATA{asm.REG1, 0x0007},
		asm.CMP{asm.REG0, asm.REG1}, // if fontY == 0x0007 then we have rendered the last line

		// if all 8 lines rendered, jump out of loop, we're done
		asm.JMPF{[]string{"E"}, asm.LABEL{labelPrefix + "-ENDLOOP"}},

		// otherwise jump back to start of loop and render next line of font
		asm.JMP{asm.LABEL{labelPrefix + "-STARTLOOP"}},
	)

	//increment pen position by 1 as we are moving to the next character
	instructions.Add(
		asm.LABEL{labelPrefix + "-ENDLOOP"},
		asm.DATA{asm.REG0, PEN_POSITION_ADDR},
		asm.LOAD{asm.REG0, asm.REG0},
		asm.DATA{asm.REG1, ONE},               // one
		asm.ADD{asm.REG1, asm.REG0},           // increment pen position by 1
		asm.DATA{asm.REG1, PEN_POSITION_ADDR}, // pen position address
		asm.STORE{asm.REG1, asm.REG0},         // store new value of pen position in memory
	)

	return instructions.Get()
}

func loadFontCharacterIntoFontRegion(char rune) []asm.Instruction {
	fontDescription := CHARACTERS[char]

	instructions := []asm.Instruction{}

	for i := uint16(0); i < 8; i++ {
		line := fontDescription[i]
		instructions = append(instructions, asm.DATA{asm.REG0, (uint16(char) << uint16(3)) + i})
		instructions = append(instructions, asm.DATA{asm.REG1, line})
		instructions = append(instructions, asm.STORE{asm.REG0, asm.REG1})
	}

	return instructions
}
