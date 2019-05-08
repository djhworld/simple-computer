package asm

import "fmt"

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

func drawLetters() {
	instructions := Instructions{}

	instructions.AddBlocks(
		[]Instruction{LABEL{"FONTDESCS"}},
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
		selectDisplayAdapter(REG3),
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
		deselectIO(REG3),
	)

	fmt.Println(instructions.String())
}

func updatePenPosition(position uint16) []Instruction {
	return []Instruction{
		DATA{REG0, PEN_POSITION_ADDR},
		DATA{REG1, position},
		STORE{REG0, REG1},
	}
}

func loadCharIntoKeycodeRegister(char rune) []Instruction {
	return []Instruction{
		DATA{REG0, KEYCODE_REGISTER},
		DATA{REG1, uint16(char)},
		STORE{REG0, REG1},
	}
}

func selectDisplayAdapter(useRegister REGISTER) []Instruction {
	return []Instruction{
		DATA{useRegister, DISPLAY_ADAPTER_ADDR},
		OUT{ADDRESS_MODE, useRegister},
	}
}

func deselectIO(useRegister REGISTER) []Instruction {
	return []Instruction{
		XOR{useRegister, useRegister},
		OUT{ADDRESS_MODE, useRegister},
	}
}

func loadFontCharacterIntoDisplayRAM(labelPrefix string) []Instruction {
	fontYAddr := uint16(0xFF00)
	lineWidth := uint16(0x001E)

	instructions := Instructions{}

	instructions.Add(
		DATA{REG2, PEN_POSITION_ADDR},
		LOAD{REG2, REG2},
	)

	// we can keep this value in reg2 to track where in display RAM we are writing
	penPositionRegister := REG2

	// counter for what line of the font are we rendering
	instructions.Add(
		DATA{REG0, fontYAddr},
		DATA{REG1, 0x0000},
		STORE{REG0, REG1},
	)

	// calculate memory position of font line
	// start of loop:
	instructions.Add(
		LABEL{labelPrefix + "-STARTLOOP"},
		DATA{REG3, KEYCODE_REGISTER}, // load keycode
		LOAD{REG3, REG3},
		SHL{REG3},
		SHL{REG3},
		SHL{REG3},             // memory address in RAM for start of font
		DATA{REG0, fontYAddr}, // fontY address
		LOAD{REG0, REG0},      // load fontY
		ADD{REG0, REG3},       // calculate memory position of fontstart+fontYinstructions = append(instructions, ADD{REG0, REG3})       // calculate memory position of fontstart+fontY

		//increment fontY by 1
		DATA{REG1, ONE},       // one
		ADD{REG1, REG0},       // increment fontY by 1
		DATA{REG1, fontYAddr}, // fontY address
		STORE{REG1, REG0},     // store new value of fontY in memory

		// load font line from memory
		LOAD{REG3, REG0}, // load value from memory into reg0

		// write to display ram
		OUT{DATA_MODE, penPositionRegister}, // display RAM address
		OUT{DATA_MODE, REG0},                // display RAM value
		DATA{REG1, lineWidth},
		ADD{REG1, penPositionRegister}, // move pen down by 1 line

		// check if we have rendered all 8 lines
		DATA{REG0, fontYAddr}, // fontY addr
		LOAD{REG0, REG0},      //load fontY into reg0
		DATA{REG1, 0x0007},
		CMP{REG0, REG1}, // if fontY == 0x0007 then we have rendered the last line

		// if all 8 lines rendered, jump out of loop, we're done
		JMPF{[]string{"E"}, LABEL{labelPrefix + "-ENDLOOP"}},

		// otherwise jump back to start of loop and render next line of font
		JMP{LABEL{labelPrefix + "-STARTLOOP"}},
	)

	//increment pen position by 1 as we are moving to the next character
	instructions.Add(
		LABEL{labelPrefix + "-ENDLOOP"},
		DATA{REG0, PEN_POSITION_ADDR},
		LOAD{REG0, REG0},
		DATA{REG1, ONE},               // one
		ADD{REG1, REG0},               // increment pen position by 1
		DATA{REG1, PEN_POSITION_ADDR}, // pen position address
		STORE{REG1, REG0},             // store new value of pen position in memory
	)

	return instructions.Get()
}

func loadFontCharacterIntoFontRegion(char rune) []Instruction {
	fontDescription := CHARACTERS[char]

	instructions := []Instruction{}

	for i := uint16(0); i < 8; i++ {
		line := fontDescription[i]
		instructions = append(instructions, DATA{REG0, (uint16(char) << uint16(3)) + i})
		instructions = append(instructions, DATA{REG1, line})
		instructions = append(instructions, STORE{REG0, REG1})
	}

	return instructions
}
