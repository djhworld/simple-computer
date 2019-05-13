package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	goio "io"

	"github.com/djhworld/simple-computer/computer"
	"github.com/djhworld/simple-computer/io"
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

var binFile = flag.String("bin", "/dev/stdin", "the bin file to load into the computer")
var printState = flag.Bool("print-state", false, "print the computer state to stdout")
var printStateSampleSize = flag.Int("print-state-every", 512, "how often in steps to print the computer state. lower will decrease performance.")

func main() {
	flag.Parse()
	fmt.Println("\nDaniel's Simple Computer (based on the Scott CPU)")
	fmt.Println(strings.Repeat("-", 80))

	bin, err := read(*binFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error attempting to parse bin file", err)
		os.Exit(5)
	}

	run(bin)
}

func run(bin []uint16) {
	keyPressChannel := make(chan *io.KeyPress)
	screenChannel := make(chan *[160][240]byte)
	quitChannel := make(chan bool, 10)

	glfw := NewGlfwIO(screenChannel, keyPressChannel, quitChannel)
	if err := glfw.Init(fmt.Sprintf("%s", *binFile)); err != nil {
		fmt.Fprintln(os.Stderr, "error received initialising GLFW instnace", err)
		os.Exit(5)
	}

	comp := computer.NewComputer(screenChannel, quitChannel)
	keyboard := io.NewKeyboard(keyPressChannel, quitChannel)
	comp.ConnectKeyboard(keyboard)
	comp.LoadToRAM(0x0500, bin)

	go keyboard.Run()
	go comp.Run(time.Tick(1*time.Nanosecond), computer.PrintStateConfig{*printState, *printStateSampleSize})

	glfw.Run()
}

func read(filename string) ([]uint16, error) {
	var reader goio.ReadCloser
	if r, err := os.Open(filename); err != nil {
		return nil, err
	} else {
		reader = r
	}
	defer reader.Close()

	size, err := countUint16s(filename)
	if err != nil {
		return nil, err
	}

	instructions := make([]uint16, size)
	if err := binary.Read(reader, binary.LittleEndian, &instructions); err != nil {
		return nil, err
	}

	return instructions, nil
}

func countUint16s(filename string) (int64, error) {
	if stat, err := os.Stat(filename); err != nil {
		return -1, err
	} else {
		filesize := stat.Size()
		if filesize%2 != 0 {
			return -1, fmt.Errorf("size of file '%s' is not an even number (bytes = %d)", filename, filesize)
		}
		uint16s := filesize / 2
		log.Println("Size of", filename, "is", uint16s, "unsigned 16-bit integers")
		return uint16s, nil
	}
}
