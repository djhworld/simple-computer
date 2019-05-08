package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/djhworld/simple-computer/asm"
)

const USER_CODE_START = uint16(0x0500)

var inputFile = flag.String("i", "", "input file (default: stdin)")
var outputFile = flag.String("o", "", "output file (default: stdout)")

func exitWithError(message string, err error, exitCode int) {
	fmt.Fprintln(os.Stderr, message, err)
	fmt.Fprint(os.Stderr, "\n")
	flag.Usage()
	os.Exit(exitCode)
}

func main() {
	flag.Parse()

	reader, err := getReaderFor(*inputFile)
	if err != nil {
		exitWithError("error reading input: ", err, 5)
	}
	defer reader.Close()

	parser := asm.Parser{}
	instructions, err := parser.Parse(reader)
	if err != nil {
		exitWithError("error parsing input: ", err, 104)
	}

	asm := asm.Assembler{}
	rawIns, err := asm.Process(USER_CODE_START, instructions)
	if err != nil {
		exitWithError("error assembling input: ", err, 104)
	}

	writer, err := getWriterFor(*outputFile)
	if err != nil {
		exitWithError("error getting output handle: ", err, 104)
	}
	defer writer.Close()

	if err := binary.Write(writer, binary.LittleEndian, rawIns); err != nil {
		exitWithError("error writing output handle: ", err, 5)
	}
}

func getReaderFor(file string) (io.ReadCloser, error) {
	if file == "" {
		return os.Stdin, nil
	}

	return os.Open(file)
}

func getWriterFor(file string) (io.WriteCloser, error) {
	if file == "" {
		return os.Stdout, nil
	}

	return os.Create(file)
}
