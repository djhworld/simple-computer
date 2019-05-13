package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/djhworld/simple-computer/asm"
)

// important RAM areas
// 0x0000 - 0x03FF ASCII table
// 0x0400 - 0x0400 pen position
// 0x0401 - 0x0401 keycode register
// 0x0500 - 0xFEFD user code + memory
// 0xFEFE - 0xFEFF used to jump back to user code
// 0xFF00 - 0xFFFF temporary variables
func main() {
	instructions := asm.Instructions{}

	instructions.AddBlocks(initialiseCommonCode())

	if len(os.Args) < 2 {
		log.Fatalf("Please provide a program name to generate")
	}

	program := os.Args[1]

	switch program {
	case "ascii":
		asciiTable(instructions)
		return
	case "brush":
		brush(instructions)
		return
	case "text-writer":
		textWriter(instructions)
		return
	case "me":
		me(instructions)
		return
	default:
		log.Fatalf("unknown program: %s", program)
	}

}

func me(instructions asm.Instructions) {
	// MAIN FUNCTION
	instructions.Add(
		asm.DEFLABEL{"main"},
	)
	instructions.AddBlocks(
		updatePenPosition(0x00F7),
		renderString("Daniel Harper"),
		resetLinex(),
		updatePenPosition(0x01E0),
		renderString(strings.Repeat("-", 30)),
		resetLinex(),
		updatePenPosition(0x03C1),
		renderString("djhworld.github.io"),
		resetLinex(),
		updatePenPosition(0x05A1),
		renderString("@djhworld"),
		resetLinex(),
		updatePenPosition(0x087C),
		renderString(":^)"),
		resetLinex(),
	)

	instructions.Add(
		asm.DEFLABEL{"noop"},
		asm.CLF{},
		asm.JMP{asm.LABEL{"noop"}},
	)

	fmt.Println(instructions.String())
}

func asciiTable(instructions asm.Instructions) {
	// MAIN FUNCTION
	instructions.Add(
		asm.DEFLABEL{"main"},
		asm.DATA{asm.REG0, asm.NUMBER{0x0020}},
		asm.DATA{asm.REG2, asm.NUMBER{0xFF23}},
		asm.STORE{asm.REG2, asm.REG0},
		asm.DATA{asm.REG2, asm.SYMBOL{"LINEX"}},
		asm.DATA{asm.REG1, asm.NUMBER{0x0000}},
		asm.STORE{asm.REG2, asm.REG1},
	)
	instructions.AddBlocks(
		updatePenPosition(0x00F0),
	)

	instructions.Add(asm.DEFLABEL{"main-loop"})

	instructions.Add(
		asm.DATA{asm.REG2, asm.NUMBER{0xFF23}},
		asm.LOAD{asm.REG2, asm.REG0},
		asm.DATA{asm.REG1, asm.SYMBOL{"ONE"}},
		asm.ADD{asm.REG1, asm.REG0},
		asm.DATA{asm.REG2, asm.SYMBOL{"KEYCODE-REGISTER"}},
		asm.STORE{asm.REG2, asm.REG0},
		asm.DATA{asm.REG2, asm.NUMBER{0xFF23}},
		asm.STORE{asm.REG2, asm.REG0},
	)
	instructions.AddBlocks(
		callRoutine("ROUTINE-io-drawFontCharacter"),
	)
	instructions.Add(
		asm.DATA{asm.REG2, asm.NUMBER{0xFF23}},
		asm.LOAD{asm.REG2, asm.REG0},
		asm.DATA{asm.REG2, asm.NUMBER{0x7E}},
		asm.CMP{asm.REG0, asm.REG2},
		asm.JMPF{[]string{"E"}, asm.LABEL{"main"}},
	)

	instructions.Add(asm.JMP{asm.LABEL{"main-loop"}})

	fmt.Println(instructions.String())

}

func brush(instructions asm.Instructions) {
	// MAIN FUNCTION
	instructions.Add(asm.DEFLABEL{"main"})
	instructions.AddBlocks(
		updatePenPosition(0x00F0),
	)

	instructions.Add(asm.DEFLABEL{"main-getInput"})
	instructions.AddBlocks(
		callRoutine("drawBrush"),
		callRoutine("ROUTINE-io-pollKeyboard"),
		callRoutine("drawBrush"),
	)
	instructions.Add(asm.JMP{asm.LABEL{"main-getInput"}})
	instructions.AddBlocks(routine_drawBrush("drawBrush"))

	fmt.Println(instructions.String())

}

func textWriter(instructions asm.Instructions) {
	// MAIN FUNCTION
	instructions.Add(asm.DEFLABEL{"main"})

	instructions.AddBlocks(
		updatePenPosition(0x00F0),
	)

	instructions.Add(asm.DEFLABEL{"main-getInput"})
	instructions.AddBlocks(
		callRoutine("ROUTINE-io-pollKeyboard"),
		callRoutine("ROUTINE-io-drawFontCharacter"),
	)
	instructions.Add(asm.JMP{asm.LABEL{"main-getInput"}})

	fmt.Println(instructions.String())

}

func routine_drawBrush(labelPrefix string) []asm.Instruction {
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

		asm.DATA{asm.REG1, asm.NUMBER{0x0107}}, // load keycode
		asm.CMP{asm.REG3, asm.REG1},            // load keycode
		asm.JMPF{[]string{"E"}, asm.LABEL{labelPrefix + "-left"}},
		asm.DATA{asm.REG1, asm.NUMBER{0x0106}}, // load keycode
		asm.CMP{asm.REG3, asm.REG1},            // load keycode
		asm.JMPF{[]string{"E"}, asm.LABEL{labelPrefix + "-right"}},
		asm.DATA{asm.REG1, asm.NUMBER{0x0108}}, // load keycode
		asm.CMP{asm.REG3, asm.REG1},            // load keycode
		asm.JMPF{[]string{"E"}, asm.LABEL{labelPrefix + "-down"}},
		asm.DATA{asm.REG1, asm.NUMBER{0x0109}}, // load keycode
		asm.CMP{asm.REG3, asm.REG1},            // load keycode
		asm.JMPF{[]string{"E"}, asm.LABEL{labelPrefix + "-up"}},
		asm.JMP{asm.LABEL{labelPrefix + "-start"}},
	)

	instructions.Add(asm.DEFLABEL{labelPrefix + "-right"},
		asm.DATA{asm.REG1, asm.SYMBOL{"ONE"}},               // load keycode
		asm.ADD{asm.REG1, asm.REG2},                         // load keycode
		asm.DATA{asm.REG3, asm.SYMBOL{"PEN-POSITION-ADDR"}}, // load keycode
		asm.STORE{asm.REG3, asm.REG2},                       // load keycode
		asm.JMP{asm.LABEL{labelPrefix + "-start"}},
	)

	instructions.Add(asm.DEFLABEL{labelPrefix + "-down"},
		asm.DATA{asm.REG1, asm.NUMBER{0x00F0}},              // load keycode
		asm.ADD{asm.REG1, asm.REG2},                         // load keycode
		asm.DATA{asm.REG3, asm.SYMBOL{"PEN-POSITION-ADDR"}}, // load keycode
		asm.STORE{asm.REG3, asm.REG2},                       // load keycode
		asm.JMP{asm.LABEL{labelPrefix + "-start"}},
	)

	instructions.Add(asm.DEFLABEL{labelPrefix + "-up"},
		asm.DATA{asm.REG0, asm.SYMBOL{"ONE"}},  // load keycode
		asm.DATA{asm.REG1, asm.NUMBER{0x00F0}}, // load keycode
		asm.NOT{asm.REG1},
		asm.ADD{asm.REG0, asm.REG1},
		asm.CLF{},
		asm.ADD{asm.REG1, asm.REG2},
		asm.DATA{asm.REG3, asm.SYMBOL{"PEN-POSITION-ADDR"}}, // load keycode
		asm.STORE{asm.REG3, asm.REG2},                       // load keycode
		asm.JMP{asm.LABEL{labelPrefix + "-start"}},
	)

	instructions.Add(asm.DEFLABEL{labelPrefix + "-left"},
		asm.DATA{asm.REG0, asm.SYMBOL{"ONE"}}, // load keycode
		asm.DATA{asm.REG1, asm.SYMBOL{"ONE"}}, // load keycode
		asm.NOT{asm.REG1},
		asm.ADD{asm.REG0, asm.REG1},
		asm.CLF{},
		asm.ADD{asm.REG1, asm.REG2},
		asm.DATA{asm.REG3, asm.SYMBOL{"PEN-POSITION-ADDR"}}, // load keycode
		asm.STORE{asm.REG3, asm.REG2},                       // load keycode
		asm.JMP{asm.LABEL{labelPrefix + "-start"}},
	)

	instructions.Add(asm.DEFLABEL{labelPrefix + "-start"})
	instructions.AddBlocks(
		selectDisplayAdapter(asm.REG3),
	)
	// calculate memory position of font line
	// start of loop:
	instructions.Add(
		asm.DEFLABEL{labelPrefix + "-STARTLOOP"},
		asm.DATA{asm.REG3, asm.NUMBER{0x0000}}, // load keycode
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
		asm.DATA{asm.REG1, asm.NUMBER{0x0008}},
		asm.CMP{asm.REG0, asm.REG1}, // if fontY == 0x0007 then we have rendered the last line

		// if all 8 lines rendered, jump out of loop, we're done
		asm.JMPF{[]string{"E"}, asm.LABEL{labelPrefix + "-ENDLOOP"}},

		// otherwise jump back to start of loop and render next line of font
		asm.JMP{asm.LABEL{labelPrefix + "-STARTLOOP"}},
	)

	//update pen position we are moving to the next character
	instructions.Add(
		asm.DEFLABEL{labelPrefix + "-ENDLOOP"},
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
