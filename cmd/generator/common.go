package main

import (
	"github.com/djhworld/simple-computer/asm"
)

const (
	USER_CODE_AREA = uint16(0x0500)
)

var CHARACTERS map[rune][8]uint16 = map[rune][8]uint16{
	' ':  [8]uint16{0x0000, 0x0000, 0x0000, 0x0000, 0x0000, 0x0000, 0x0000, 0x0000},
	'!':  [8]uint16{0x0010, 0x0010, 0x0010, 0x0010, 0x0010, 0x000, 0x0010, 0x000},
	'"':  [8]uint16{0x0028, 0x0028, 0x000, 0x000, 0x000, 0x000, 0x000, 0x000},
	'\'': [8]uint16{0x0020, 0x0020, 0x0020, 0x0000, 0x0000, 0x0000, 0x0000, 0x0000},
	'#':  [8]uint16{0x0028, 0x0028, 0x007C, 0x0028, 0x007C, 0x0028, 0x0028, 0x000},
	'%':  [8]uint16{0x00C2, 0x00C4, 0x008, 0x0010, 0x0020, 0x004C, 0x008C, 0x000},
	'$':  [8]uint16{0x0010, 0x007E, 0x0090, 0x007C, 0x0012, 0x00FC, 0x0010, 0x000},
	'&':  [8]uint16{0x0038, 0x0028, 0x0038, 0x00E0, 0x0094, 0x0088, 0x00F4, 0x000},
	'(':  [8]uint16{0x008, 0x0010, 0x0020, 0x0020, 0x0020, 0x0010, 0x008, 0x000},
	')':  [8]uint16{0x0020, 0x0010, 0x008, 0x008, 0x008, 0x0010, 0x0020, 0x000},
	'*':  [8]uint16{0x000, 0x0092, 0x0054, 0x0038, 0x0038, 0x0054, 0x0092, 0x000},
	'+':  [8]uint16{0x000, 0x0010, 0x0010, 0x007C, 0x0030, 0x0010, 0x000, 0x000},
	'/':  [8]uint16{0x002, 0x004, 0x008, 0x0010, 0x0020, 0x0040, 0x0080, 0x000},
	'.':  [8]uint16{0x000, 0x000, 0x000, 0x000, 0x000, 0x000, 0x0010, 0x000},
	',':  [8]uint16{0x000, 0x000, 0x000, 0x000, 0x008, 0x008, 0x0010, 0x000},
	'-':  [8]uint16{0x000, 0x000, 0x000, 0x007C, 0x000, 0x000, 0x000, 0x000},
	'=':  [8]uint16{0x000, 0x000, 0x00FE, 0x000, 0x00FE, 0x000, 0x000, 0x000},
	'>':  [8]uint16{0x0040, 0x0020, 0x0010, 0x008, 0x0010, 0x0020, 0x0040, 0x000},
	'<':  [8]uint16{0x002, 0x004, 0x008, 0x0010, 0x008, 0x004, 0x002, 0x000},
	'|':  [8]uint16{0x0010, 0x0010, 0x0010, 0x0010, 0x0010, 0x0010, 0x0010, 0x000},
	']':  [8]uint16{0x0030, 0x0010, 0x0010, 0x0010, 0x0010, 0x0010, 0x0030, 0x000},
	'[':  [8]uint16{0x0030, 0x0020, 0x0020, 0x0020, 0x0020, 0x0020, 0x0030, 0x000},
	'\\': [8]uint16{0x0080, 0x0040, 0x0020, 0x0010, 0x008, 0x004, 0x002, 0x000},
	'~':  [8]uint16{0x000, 0x000, 0x000, 0x0032, 0x004C, 0x000, 0x000, 0x000},
	'}':  [8]uint16{0x0030, 0x008, 0x00C, 0x002, 0x00C, 0x008, 0x0030, 0x000},
	'{':  [8]uint16{0x0010, 0x0020, 0x0060, 0x0080, 0x0060, 0x0020, 0x0010, 0x000},
	'_':  [8]uint16{0x000, 0x000, 0x000, 0x000, 0x000, 0x000, 0x007E, 0x000},
	'`':  [8]uint16{0x000, 0x0020, 0x0010, 0x008, 0x000, 0x000, 0x000, 0x000},
	'^':  [8]uint16{0x0010, 0x0028, 0x0044, 0x000, 0x000, 0x000, 0x000, 0x000},
	':':  [8]uint16{0x000, 0x0010, 0x000, 0x000, 0x0010, 0x000, 0x000, 0x000},
	';':  [8]uint16{0x000, 0x0010, 0x000, 0x000, 0x0010, 0x0020, 0x000, 0x000},
	'?':  [8]uint16{0x007C, 0x0042, 0x002, 0x004, 0x008, 0x000, 0x008, 0x000},
	'@':  [8]uint16{0x007C, 0x008A, 0x009C, 0x00A8, 0x0098, 0x0084, 0x0078, 0x000},
	'A':  [8]uint16{0x007C, 0x00C6, 0x0082, 0x00FE, 0x0082, 0x0082, 0x0082, 0x0000},
	'B':  [8]uint16{0x00FC, 0x0086, 0x0082, 0x00FE, 0x0082, 0x0086, 0x00FC, 0x0000},
	'C':  [8]uint16{0x007E, 0x00C0, 0x0080, 0x0080, 0x0080, 0x00C0, 0x007E, 0x0000},
	'D':  [8]uint16{0x00F8, 0x0086, 0x0082, 0x0082, 0x0082, 0x0086, 0x00F8, 0x0000},
	'E':  [8]uint16{0x007E, 0x00C0, 0x0080, 0x00FE, 0x0080, 0x00C0, 0x007E, 0x0000},
	'F':  [8]uint16{0x007E, 0x0080, 0x0080, 0x00FC, 0x0080, 0x0080, 0x0080, 0x0000},
	'G':  [8]uint16{0x007E, 0x0080, 0x0080, 0x009C, 0x0082, 0x0082, 0x00FE, 0x0000},
	'H':  [8]uint16{0x0082, 0x0082, 0x0082, 0x00FE, 0x0082, 0x0082, 0x0082, 0x0000},
	'I':  [8]uint16{0x00FE, 0x0010, 0x0010, 0x0010, 0x0010, 0x0010, 0x00FE, 0x0000},
	'J':  [8]uint16{0x0002, 0x0002, 0x0002, 0x0002, 0x0002, 0x0002, 0x00FC, 0x0000},
	'K':  [8]uint16{0x00C4, 0x00C8, 0x00F0, 0x00E0, 0x00D8, 0x00C4, 0x00C6, 0x000},
	'L':  [8]uint16{0x0080, 0x0080, 0x0080, 0x0080, 0x0080, 0x0080, 0x007E, 0x0000},
	'M':  [8]uint16{0x0066, 0x00aa, 0x0092, 0x0092, 0x0082, 0x0082, 0x0082, 0x0000},
	'N':  [8]uint16{0x00C2, 0x00a2, 0x0092, 0x0092, 0x008A, 0x008A, 0x0086, 0x0000},
	'O':  [8]uint16{0x007C, 0x0082, 0x0082, 0x0082, 0x0082, 0x0082, 0x007C, 0x000},
	'P':  [8]uint16{0x00FC, 0x0082, 0x0082, 0x001FC, 0x0080, 0x0080, 0x0080, 0x000},
	'Q':  [8]uint16{0x0078, 0x0084, 0x0084, 0x0084, 0x0094, 0x008C, 0x0076, 0x007},
	'R':  [8]uint16{0x00FC, 0x0082, 0x0082, 0x00FC, 0x00A0, 0x0090, 0x008E, 0x000},
	'S':  [8]uint16{0x007C, 0x0080, 0x0080, 0x007C, 0x004, 0x004, 0x00F8, 0x000},
	'T':  [8]uint16{0x00FE, 0x0010, 0x0010, 0x0010, 0x0010, 0x0010, 0x0010, 0x000},
	'U':  [8]uint16{0x00C6, 0x0042, 0x0042, 0x0042, 0x0042, 0x0042, 0x003C, 0x000},
	'V':  [8]uint16{0x0082, 0x0082, 0x0082, 0x0082, 0x0044, 0x006C, 0x0010, 0x000},
	'W':  [8]uint16{0x0082, 0x0082, 0x0082, 0x0092, 0x00BA, 0x00AA, 0x0044, 0x000},
	'X':  [8]uint16{0x00C6, 0x0044, 0x0028, 0x0010, 0x0028, 0x0044, 0x00C6, 0x000},
	'Y':  [8]uint16{0x00C6, 0x0044, 0x0028, 0x0010, 0x0010, 0x0010, 0x0038, 0x000},
	'Z':  [8]uint16{0x00FE, 0x0082, 0x00C, 0x0038, 0x0060, 0x0082, 0x007E, 0x000},
	'a':  [8]uint16{0x007c, 0x00c6, 0x0082, 0x00fe, 0x0082, 0x0082, 0x0082, 0x0000},
	'b':  [8]uint16{0x00fc, 0x0086, 0x0082, 0x00fe, 0x0082, 0x0086, 0x00fc, 0x0000},
	'c':  [8]uint16{0x007e, 0x00c0, 0x0080, 0x0080, 0x0080, 0x00c0, 0x007e, 0x0000},
	'd':  [8]uint16{0x00f8, 0x0086, 0x0082, 0x0082, 0x0082, 0x0086, 0x00f8, 0x0000},
	'e':  [8]uint16{0x007e, 0x00c0, 0x0080, 0x00fe, 0x0080, 0x00c0, 0x007e, 0x0000},
	'f':  [8]uint16{0x007e, 0x0080, 0x0080, 0x00fc, 0x0080, 0x0080, 0x0080, 0x0000},
	'g':  [8]uint16{0x007e, 0x0080, 0x0080, 0x009c, 0x0082, 0x0082, 0x00fe, 0x0000},
	'h':  [8]uint16{0x0082, 0x0082, 0x0082, 0x00fe, 0x0082, 0x0082, 0x0082, 0x0000},
	'i':  [8]uint16{0x00fe, 0x0010, 0x0010, 0x0010, 0x0010, 0x0010, 0x00fe, 0x0000},
	'j':  [8]uint16{0x0002, 0x0002, 0x0002, 0x0002, 0x0002, 0x0002, 0x00fc, 0x0000},
	'k':  [8]uint16{0x00c4, 0x00c8, 0x00f0, 0x00e0, 0x00d8, 0x00c4, 0x00c6, 0x000},
	'l':  [8]uint16{0x0080, 0x0080, 0x0080, 0x0080, 0x0080, 0x0080, 0x007e, 0x0000},
	'm':  [8]uint16{0x0066, 0x00aa, 0x0092, 0x0092, 0x0082, 0x0082, 0x0082, 0x0000},
	'n':  [8]uint16{0x00c2, 0x00a2, 0x0092, 0x0092, 0x008a, 0x008a, 0x0086, 0x0000},
	'o':  [8]uint16{0x007c, 0x0082, 0x0082, 0x0082, 0x0082, 0x0082, 0x007c, 0x000},
	'p':  [8]uint16{0x00fc, 0x0082, 0x0082, 0x001fc, 0x0080, 0x0080, 0x0080, 0x000},
	'q':  [8]uint16{0x0078, 0x0084, 0x0084, 0x0084, 0x0094, 0x008c, 0x0076, 0x007},
	'r':  [8]uint16{0x00fc, 0x0082, 0x0082, 0x00fc, 0x00a0, 0x0090, 0x008e, 0x000},
	's':  [8]uint16{0x007c, 0x0080, 0x0080, 0x007c, 0x004, 0x004, 0x00f8, 0x000},
	't':  [8]uint16{0x00fe, 0x0010, 0x0010, 0x0010, 0x0010, 0x0010, 0x0010, 0x000},
	'u':  [8]uint16{0x00c6, 0x0042, 0x0042, 0x0042, 0x0042, 0x0042, 0x003c, 0x000},
	'v':  [8]uint16{0x0082, 0x0082, 0x0082, 0x0082, 0x0044, 0x006c, 0x0010, 0x000},
	'w':  [8]uint16{0x0082, 0x0082, 0x0082, 0x0092, 0x00ba, 0x00aa, 0x0044, 0x000},
	'x':  [8]uint16{0x00c6, 0x0044, 0x0028, 0x0010, 0x0028, 0x0044, 0x00c6, 0x000},
	'y':  [8]uint16{0x00c6, 0x0044, 0x0028, 0x0010, 0x0010, 0x0010, 0x0038, 0x000},
	'z':  [8]uint16{0x00fe, 0x0082, 0x00c, 0x0038, 0x0060, 0x0082, 0x007e, 0x000},
	'0':  [8]uint16{0x007C, 0x00E2, 0x00A2, 0x0092, 0x008A, 0x008E, 0x007C, 0x000},
	'1':  [8]uint16{0x0038, 0x0058, 0x0018, 0x0018, 0x0018, 0x0018, 0x007E, 0x000},
	'2':  [8]uint16{0x007C, 0x0082, 0x001C, 0x0020, 0x0040, 0x0080, 0x00FE, 0x000},
	'3':  [8]uint16{0x007C, 0x002, 0x002, 0x001E, 0x002, 0x002, 0x00FC, 0x000},
	'4':  [8]uint16{0x001C, 0x0024, 0x0044, 0x0084, 0x00FE, 0x004, 0x004, 0x000},
	'5':  [8]uint16{0x00FE, 0x0080, 0x00F8, 0x004, 0x002, 0x006, 0x00FC, 0x000},
	'6':  [8]uint16{0x003E, 0x0040, 0x00F8, 0x0084, 0x0082, 0x0086, 0x00FC, 0x000},
	'7':  [8]uint16{0x00FE, 0x002, 0x004, 0x008, 0x0010, 0x0020, 0x0040, 0x000},
	'8':  [8]uint16{0x007C, 0x0082, 0x0082, 0x007C, 0x0082, 0x0082, 0x007C, 0x000},
	'9':  [8]uint16{0x007C, 0x0082, 0x0082, 0x007E, 0x002, 0x0082, 0x007C, 0x000},
	0:    [8]uint16{0xFFFF, 0xFFFF, 0xFFFF, 0xFFFF, 0xFFFF, 0xFFFF, 0xFFFF, 0xFFFF},
}

func initialiseCommonCode() []asm.Instruction {
	instructions := asm.Instructions{}

	instructions.Add(
		asm.DEFSYMBOL{"LINE-WIDTH", 0x001E},
		asm.DEFSYMBOL{"CALL-RETURN-ADDRESS", 0xFF33},
		asm.DEFSYMBOL{"ONE", 0x0001},
		asm.DEFSYMBOL{"LINEX", 0xFF01},
		asm.DEFSYMBOL{"PEN-POSITION-ADDR", 0x0400},
		asm.DEFSYMBOL{"KEYCODE-REGISTER", 0x0401},

		asm.DEFSYMBOL{"DISPLAY-ADAPTER-ADDR", 0x0007},
		asm.DEFSYMBOL{"KEY-ADAPTER-ADDR", 0x000F},
	)

	instructions.Add(
		asm.DATA{asm.REG0, asm.SYMBOL{"LINEX"}},
		asm.DATA{asm.REG1, asm.NUMBER{0x0000}},
		asm.STORE{asm.REG0, asm.REG1},
	)

	// jump to main
	instructions.Add(asm.DEFLABEL{"start"})
	instructions.AddBlocks(callRoutine("ROUTINE-init-fontDescriptions"))
	instructions.Add(asm.JMP{asm.LABEL{"main"}})

	instructions.Add(asm.DEFLABEL{"ROUTINES"})
	instructions.AddBlocks(routine_loadFontDescriptions("ROUTINE-init-fontDescriptions"))
	instructions.AddBlocks(routine_drawFontCharacter("ROUTINE-io-drawFontCharacter"))
	instructions.AddBlocks(routine_pollKeyboard("ROUTINE-io-pollKeyboard"))
	return instructions.Get()
}

func routine_loadFontDescriptions(label string) []asm.Instruction {
	instructions := asm.Instructions{}
	instructions.Add(asm.DEFLABEL{label})

	for char, _ := range CHARACTERS {
		instructions.AddBlocks(
			loadFontCharacterIntoFontRegion(char),
		)
	}

	instructions.Add(asm.JR{asm.REG3})
	return instructions.Get()
}

func routine_pollKeyboard(labelPrefix string) []asm.Instruction {
	instructions := asm.Instructions{}
	instructions.Add(asm.DEFLABEL{labelPrefix})

	// push retun address
	instructions.Add(
		asm.DATA{asm.REG2, asm.SYMBOL{"CALL-RETURN-ADDRESS"}},
		asm.STORE{asm.REG2, asm.REG3},
	)

	instructions.Add(
		asm.DATA{asm.REG2, asm.NUMBER{0x000F}}, //select keyboard keyboard
		asm.OUT{asm.ADDRESS_MODE, asm.REG2},
	)

	instructions.Add(
		asm.DEFLABEL{labelPrefix + "-STARTLOOP"},
		asm.IN{asm.DATA_MODE, asm.REG3},                                // request key from keyboard adapter
		asm.AND{asm.REG3, asm.REG3},                                    // check if value is zero
		asm.JMPF{[]string{"Z"}, asm.LABEL{labelPrefix + "-STARTLOOP"}}, // if it is - keep polling

		asm.DEFLABEL{labelPrefix + "-ENDLOOP"}, //otherwise
		asm.DATA{asm.REG0, asm.SYMBOL{"KEYCODE-REGISTER"}},
		asm.STORE{asm.REG0, asm.REG3}, //store key in here

		// deselect keyboard
		asm.XOR{asm.REG2, asm.REG2},
		asm.OUT{asm.ADDRESS_MODE, asm.REG2},
	)

	// return to callee
	instructions.Add(
		asm.CLF{},
		asm.DATA{asm.REG3, asm.SYMBOL{"CALL-RETURN-ADDRESS"}},
		asm.LOAD{asm.REG3, asm.REG3},
		asm.JR{asm.REG3},
	)

	return instructions.Get()
}

func resetLinex() []asm.Instruction {
	instructions := asm.Instructions{}
	instructions.Add(
		asm.DATA{asm.REG2, asm.SYMBOL{"LINEX"}}, //reset linex
		asm.DATA{asm.REG3, asm.NUMBER{0x000}},
		asm.STORE{asm.REG2, asm.REG3},
	)
	return instructions.Get()
}

func routine_drawFontCharacter(labelPrefix string) []asm.Instruction {
	fontYAddr := uint16(0xFF00)

	instructions := asm.Instructions{}
	instructions.Add(asm.DEFLABEL{labelPrefix})

	instructions.Add(
		asm.DATA{asm.REG2, asm.SYMBOL{"CALL-RETURN-ADDRESS"}},
		asm.STORE{asm.REG2, asm.REG3},
		asm.DATA{asm.REG2, asm.SYMBOL{"PEN-POSITION-ADDR"}},
		asm.LOAD{asm.REG2, asm.REG2},
	)

	// we can keep this value in reg2 to track where in display RAM we are writing
	penPositionRegister := asm.REG2

	// counter for what line of the font are we rendering
	instructions.Add(
		asm.DATA{asm.REG0, asm.NUMBER{fontYAddr}},
		asm.DATA{asm.REG1, asm.NUMBER{0x0000}},
		asm.STORE{asm.REG0, asm.REG1},
	)

	instructions.Add(
		asm.DATA{asm.REG3, asm.SYMBOL{"KEYCODE-REGISTER"}}, // load keycode
		asm.LOAD{asm.REG3, asm.REG3},
		asm.DATA{asm.REG1, asm.NUMBER{0x0101}}, // load keycode
		asm.CMP{asm.REG3, asm.REG1},            // load keycode

		asm.JMPF{[]string{"E"}, asm.LABEL{labelPrefix + "-carriage-return"}},
	)

	instructions.AddBlocks(
		selectDisplayAdapter(asm.REG3),
	)
	// calculate memory position of font line
	// start of loop:
	instructions.Add(
		asm.DEFLABEL{labelPrefix + "-STARTLOOP"},
		asm.DATA{asm.REG3, asm.SYMBOL{"KEYCODE-REGISTER"}}, // load keycode
		asm.LOAD{asm.REG3, asm.REG3},
		asm.SHL{asm.REG3},
		asm.SHL{asm.REG3},
		asm.SHL{asm.REG3},                         // memory address in RAM for start of font
		asm.DATA{asm.REG0, asm.NUMBER{fontYAddr}}, // fontY address
		asm.LOAD{asm.REG0, asm.REG0},              // load fontY
		asm.ADD{asm.REG0, asm.REG3},               // calculate memory position of fontstart+fontYinstructions = append(instructions, ADD{asm.REG0, asm.REG3})       // calculate memory position of fontstart+fontY

		//increment fontY by 1
		asm.DATA{asm.REG1, asm.SYMBOL{"ONE"}},     // one
		asm.ADD{asm.REG1, asm.REG0},               // increment fontY by 1
		asm.DATA{asm.REG1, asm.NUMBER{fontYAddr}}, // fontY address
		asm.STORE{asm.REG1, asm.REG0},             // store new value of fontY in memory

		// load font line from memory
		asm.LOAD{asm.REG3, asm.REG0}, // load value from memory into reg0

		// write to display ram
		asm.OUT{asm.DATA_MODE, penPositionRegister}, // display RAM address
		asm.OUT{asm.DATA_MODE, asm.REG0},            // display RAM value
		asm.DATA{asm.REG1, asm.SYMBOL{"LINE-WIDTH"}},
		asm.ADD{asm.REG1, penPositionRegister}, // move pen down by 1 line

		// check if we have rendered all 8 lines
		asm.DATA{asm.REG0, asm.NUMBER{fontYAddr}}, // fontY addr
		asm.LOAD{asm.REG0, asm.REG0},              //load fontY into reg0
		asm.DATA{asm.REG1, asm.NUMBER{0x0007}},
		asm.CMP{asm.REG0, asm.REG1}, // if fontY == 0x0007 then we have rendered the last line

		// if all 8 lines rendered, jump out of loop, we're done
		asm.JMPF{[]string{"E"}, asm.LABEL{labelPrefix + "-ENDLOOP"}},

		// otherwise jump back to start of loop and render next line of font
		asm.JMP{asm.LABEL{labelPrefix + "-STARTLOOP"}},
	)

	//update pen position we are moving to the next character
	instructions.Add(
		asm.DEFLABEL{labelPrefix + "-ENDLOOP"},

		// increment line x
		asm.DATA{asm.REG1, asm.SYMBOL{"LINEX"}},
		asm.LOAD{asm.REG1, asm.REG1},
		asm.DATA{asm.REG2, asm.SYMBOL{"ONE"}},
		asm.ADD{asm.REG2, asm.REG1}, //increment line X
		asm.DATA{asm.REG2, asm.SYMBOL{"LINEX"}},
		asm.STORE{asm.REG2, asm.REG1},

		asm.DATA{asm.REG0, asm.SYMBOL{"PEN-POSITION-ADDR"}},
		asm.LOAD{asm.REG0, asm.REG0},

		// test if pen needs to be moved down
		asm.DATA{asm.REG3, asm.NUMBER{0x1E}}, // have we reached the end of the line?
		asm.CMP{asm.REG1, asm.REG3},
		asm.JMPF{[]string{"E"}, asm.LABEL{labelPrefix + "-carriage-return"}},
		asm.JMP{asm.LABEL{labelPrefix + "-increment-cursor"}},

		asm.DEFLABEL{labelPrefix + "-increment-cursor"},
		asm.DATA{asm.REG1, asm.SYMBOL{"ONE"}}, // one
		asm.ADD{asm.REG1, asm.REG0},           // increment pen position by 1
		asm.DATA{asm.REG1, asm.SYMBOL{"PEN-POSITION-ADDR"}},
		asm.STORE{asm.REG1, asm.REG0}, // store new value of pen position in memory
		asm.JMP{asm.LABEL{labelPrefix + "-deselectIO"}},

		// used for when the enter key is hit and you need to return
		// this is absolutely monstrous
		asm.DEFLABEL{labelPrefix + "-carriage-return"},
		asm.DATA{asm.REG1, asm.SYMBOL{"LINEX"}},
		asm.LOAD{asm.REG1, asm.REG1}, // retrieve linex

		// if linex == 0
		asm.DATA{asm.REG2, asm.NUMBER{0x0000}},
		asm.DATA{asm.REG3, asm.NUMBER{0x00F0}},
		asm.CMP{asm.REG1, asm.REG2},
		asm.JMPF{[]string{"E"}, asm.LABEL{labelPrefix + "-reposition-pen"}},

		// if linex == 1...
		asm.DATA{asm.REG2, asm.SYMBOL{"ONE"}},
		asm.DATA{asm.REG3, asm.NUMBER{0x00EF}},
		asm.CMP{asm.REG1, asm.REG2},
		asm.JMPF{[]string{"E"}, asm.LABEL{labelPrefix + "-reposition-pen"}},

		// if linex == end of line
		asm.DATA{asm.REG2, asm.SYMBOL{"LINE-WIDTH"}},
		asm.DATA{asm.REG3, asm.NUMBER{0x00F1}},
		asm.CMP{asm.REG1, asm.REG2},
		asm.JMPF{[]string{"E"}, asm.LABEL{labelPrefix + "-reposition-pen"}},

		// otherwise calculate difference and move pen down
		asm.DEFLABEL{labelPrefix + "-reposition-pen-when-midline"},
		asm.DATA{asm.REG2, asm.SYMBOL{"ONE"}},
		asm.DATA{asm.REG0, asm.NUMBER{0x00ef}}, // 239 pixels (next line)
		// subtract linex - 239
		asm.NOT{asm.REG1},
		asm.ADD{asm.REG2, asm.REG1},
		asm.CLF{},
		asm.ADD{asm.REG0, asm.REG1},

		// add subtraction to pen position
		asm.DATA{asm.REG0, asm.SYMBOL{"PEN-POSITION-ADDR"}},
		asm.LOAD{asm.REG0, asm.REG0},
		asm.ADD{asm.REG1, asm.REG0},
		asm.DATA{asm.REG1, asm.SYMBOL{"PEN-POSITION-ADDR"}},
		asm.STORE{asm.REG1, asm.REG0}, // store new value of pen position in memory
		asm.JMP{asm.LABEL{labelPrefix + "-resetlinex"}},

		// needs value in reg3
		asm.DEFLABEL{labelPrefix + "-reposition-pen"},
		// add subtraction to pen position
		asm.DATA{asm.REG0, asm.SYMBOL{"PEN-POSITION-ADDR"}},
		asm.LOAD{asm.REG0, asm.REG0},
		asm.ADD{asm.REG3, asm.REG0},
		asm.DATA{asm.REG1, asm.SYMBOL{"PEN-POSITION-ADDR"}},
		asm.STORE{asm.REG1, asm.REG0}, // store new value of pen position in memory
		asm.JMP{asm.LABEL{labelPrefix + "-resetlinex"}},

		// reset linex
		asm.DEFLABEL{labelPrefix + "-resetlinex"},
		asm.DATA{asm.REG2, asm.SYMBOL{"LINEX"}}, //reset linex
		asm.DATA{asm.REG3, asm.NUMBER{0x000}},
		asm.STORE{asm.REG2, asm.REG3},

		asm.JMP{asm.LABEL{labelPrefix + "-deselectIO"}},
	)

	instructions.Add(asm.DEFLABEL{labelPrefix + "-deselectIO"})

	// deselect IO adapter
	instructions.AddBlocks(
		deselectIO(asm.REG3),
	)

	// return to callee
	instructions.Add(
		asm.CLF{},
		asm.DATA{asm.REG3, asm.SYMBOL{"CALL-RETURN-ADDRESS"}},
		asm.LOAD{asm.REG3, asm.REG3},
		asm.JR{asm.REG3},
	)

	return instructions.Get()
}

func callRoutine(routine string) []asm.Instruction {
	return []asm.Instruction{
		asm.CALL{asm.LABEL{routine}},
	}
}

func updatePenPosition(position uint16) []asm.Instruction {
	return []asm.Instruction{
		asm.DATA{asm.REG0, asm.SYMBOL{"PEN-POSITION-ADDR"}},
		asm.DATA{asm.REG1, asm.NUMBER{position}},
		asm.STORE{asm.REG0, asm.REG1},
	}
}

func loadCharIntoKeycodeRegister(char rune) []asm.Instruction {
	return []asm.Instruction{
		asm.DATA{asm.REG0, asm.SYMBOL{"KEYCODE-REGISTER"}},
		asm.DATA{asm.REG1, asm.NUMBER{uint16(char)}},
		asm.STORE{asm.REG0, asm.REG1},
	}
}

func renderString(str string) []asm.Instruction {
	instructions := asm.Instructions{}

	for _, r := range str {
		instructions.AddBlocks(
			loadCharIntoKeycodeRegister(r),
			callRoutine("ROUTINE-io-drawFontCharacter"),
		)

	}

	return instructions.Get()
}

func selectDisplayAdapter(useRegister asm.REGISTER) []asm.Instruction {
	return []asm.Instruction{
		asm.DATA{useRegister, asm.SYMBOL{"DISPLAY-ADAPTER-ADDR"}},
		asm.OUT{asm.ADDRESS_MODE, useRegister},
	}
}

func deselectIO(useRegister asm.REGISTER) []asm.Instruction {
	return []asm.Instruction{
		asm.XOR{useRegister, useRegister},
		asm.OUT{asm.ADDRESS_MODE, useRegister},
	}
}

func loadFontCharacterIntoFontRegion(char rune) []asm.Instruction {
	fontDescription := CHARACTERS[char]

	instructions := []asm.Instruction{}

	for i := uint16(0); i < 8; i++ {
		line := fontDescription[i]
		instructions = append(instructions, asm.DATA{asm.REG0, asm.NUMBER{(uint16(char) << uint16(3)) + i}})
		instructions = append(instructions, asm.DATA{asm.REG1, asm.NUMBER{line}})
		instructions = append(instructions, asm.STORE{asm.REG0, asm.REG1})
	}

	return instructions
}
