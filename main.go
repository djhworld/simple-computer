package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/djhworld/simple-computer/io"
)

/*
	instructions := []uint16{
		0x20,
		0x0100,

		0x21,
		0x01,
		0x22, // DATA R2
		0x000F,
		0x7E, // OUT Addr, R2
		0x73, // IN Data, R3
		0xEA, // XOR, R2, R2
		0x7E, // OUT Addr, R2
		0xCF, // AND R3, R3
		0x51, // JMPZ
		0x04,
		0x13, //ST R0, R3
		0x84, //ADD R1, R0
		0x60, //ADD R1, R0
		0x40,
		0x04,
	}
*/

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func main() {
	instructions := []uint16{
		0x20, //DATA R0
		0x0228,
		0x21, //DATA R1
		0x7E,
		0x11, // ST R0, R1

		0x20, //DATA R0
		0x0229,
		0x21, //DATA R1
		0x40,
		0x11, // ST R0, R1

		0x20, //DATA R0
		0x022A,
		0x21, //DATA R1
		0x40,
		0x11, // ST R0, R1

		0x20, //DATA R0
		0x022B,
		0x21, //DATA R1
		0x7E,
		0x11, // ST R0, R1

		0x20, //DATA R0
		0x022C,
		0x21, //DATA R1
		0x40,
		0x11, // ST R0, R1

		0x20, //DATA R0
		0x022D,
		0x21, //DATA R1
		0x40,
		0x11, // ST R0, R1

		0x20, //DATA R0
		0x022E,
		0x21, //DATA R1
		0x7E,
		0x11, // ST R0, R1

		0x20, //DATA R0
		0x022F,
		0x21, //DATA R1
		0x00,
		0x11, // ST R0, R1

		0x22, // DATA R2
		0x012F,

		0x21, // DATA R1
		0x001E,

		0x23, // DATA R3
		0x0007,
		0x7F, // OUT Addr, R3

		0x7A, // OUT Data, R2
		0x23, // DATA R3
		0x0228,
		0x000C, // LD R3, R0
		0x78,   // OUT Data, R0
		0x86,   // ADD R1, R2

		0x7A, // OUT Data, R2
		0x23, // DATA R3
		0x0229,
		0x000C, // LD R3, R0
		0x78,   // OUT Data, R0
		0x86,   // ADD R1, R2

		0x7A, // OUT Data, R2
		0x23, // DATA R3
		0x022A,
		0x000C, // LD R3, R0
		0x78,   // OUT Data, R0
		0x86,   // ADD R1, R2

		0x7A, // OUT Data, R2
		0x23, // DATA R3
		0x022B,
		0x000C, // LD R3, R0
		0x78,   // OUT Data, R0
		0x86,   // ADD R1, R2

		0x7A, // OUT Data, R2
		0x23, // DATA R3
		0x022C,
		0x000C, // LD R3, R0
		0x78,   // OUT Data, R0
		0x86,   // ADD R1, R2

		0x7A, // OUT Data, R2
		0x23, // DATA R3
		0x022D,
		0x000C, // LD R3, R0
		0x78,   // OUT Data, R0
		0x86,   // ADD R1, R2

		0x7A, // OUT Data, R2
		0x23, // DATA R3
		0x022E,
		0x000C, // LD R3, R0
		0x78,   // OUT Data, R0
		0x86,   // ADD R1, R2

		0x7A, // OUT Data, R2
		0x23, // DATA R3
		0x022F,
		0x000C, // LD R3, R0
		0x78,   // OUT Data, R0
		0x86,   // ADD R1, R2

		// unselect and exit
		0xEF, // XOR, R3, R3
		0x7F, // OUT Addr, R3
		//0xFF00,
	}

	keyPressChannel := make(chan *io.KeyPress)
	screenChannel := make(chan *[160][240]byte)
	quitChannel := make(chan bool, 10)

	glfw := NewGlfwIO(screenChannel, keyPressChannel, quitChannel)
	err := glfw.Init("simple-computer")
	if err != nil {
		fmt.Println("problem!", err)
		os.Exit(1)
	}

	computer := NewComputer(screenChannel, quitChannel)
	keyboard := io.NewKeyboard(keyPressChannel, quitChannel)
	computer.ConnectKeyboard(keyboard)
	computer.LoadToRAM(0x0000, instructions)

	go keyboard.Run()
	go computer.Run(time.Tick(1*time.Nanosecond), PrintStateConfig{true, 512})

	glfw.Run()
}
