package asm

import (
	"fmt"
	"strings"

	"github.com/djhworld/simple-computer/utils"
)

const (
	CURRENTINSTRUCTION = "CURRENTINSTRUCTION"
	NEXTINSTRUCTION    = "NEXTINSTRUCTION"
)

var reservedSymbols = map[string]interface{}{
	NEXTINSTRUCTION:    new(interface{}),
	CURRENTINSTRUCTION: new(interface{}),
}

type Assembler struct {
	labels  map[string]uint16
	symbols map[string]uint16
}

func (a *Assembler) ResolveLabel(label LABEL) (uint16, error) {
	if v, ok := a.labels[label.Name]; !ok {
		return 0x0000, fmt.Errorf("Cannot find label: %s in label map", label.Name)
	} else {
		return v, nil
	}
}

func (a *Assembler) ResolveSymbol(symbol SYMBOL) (uint16, error) {
	if v, ok := a.symbols[symbol.Name]; !ok {
		fmt.Println(a.symbols)
		return 0x0000, fmt.Errorf("Cannot find symbol: %s in symbol map", symbol.Name)
	} else {
		return v, nil
	}
}

func (a *Assembler) Process(codeStartOffset uint16, instructions []Instruction) ([]uint16, error) {
	a.labels = make(map[string]uint16)
	a.symbols = make(map[string]uint16)
	position := uint16(0)

	//calculate labels and symbols
	for _, ins := range instructions {
		position += uint16(ins.Size())

		if label, ok := ins.(DEFLABEL); ok {
			if _, ok := a.labels[label.Name]; ok {
				return nil, fmt.Errorf("label '%s' already exists, all labels should be unique", label.Name)
			}

			a.labels[label.Name] = position + codeStartOffset
		}

		if symbol, ok := ins.(DEFSYMBOL); ok {
			if _, ok := a.symbols[symbol.Name]; ok {
				return nil, fmt.Errorf("symbol '%s' already exists, all symbols should be unique", symbol.Name)
			}

			if isReservedSymbol(symbol.Name) {
				return nil, fmt.Errorf("symbol '%s' is reserved for internal use, please use another symbol name", symbol.Name)
			}

			a.symbols[symbol.Name] = symbol.Value
		}
	}

	emitted := []uint16{}

	position = 0
	for index, ins := range instructions {
		if _, ok := ins.(DEFLABEL); ok {
			continue
		}
		if _, ok := ins.(DEFSYMBOL); ok {
			continue
		}

		a.symbols[CURRENTINSTRUCTION] = position + codeStartOffset
		a.symbols[NEXTINSTRUCTION] = getNextExecutableInstructionLoc(a.symbols[CURRENTINSTRUCTION], index, instructions)
		emit, err := ins.Emit(a.ResolveLabel, a.ResolveSymbol)
		if err != nil {
			return nil, err
		}

		emitted = append(emitted, emit...)
		position += uint16(ins.Size())
	}

	return emitted, nil
}

func (a *Assembler) ToString(codeStartOffset uint16, instructions []Instruction) (string, error) {
	a.labels = make(map[string]uint16)
	a.symbols = make(map[string]uint16)
	position := uint16(0)

	//calculate lengths
	for _, ins := range instructions {
		position += uint16(ins.Size())

		if label, ok := ins.(DEFLABEL); ok {
			a.labels[label.Name] = position + codeStartOffset
		}

		if symbol, ok := ins.(DEFSYMBOL); ok {
			a.symbols[symbol.Name] = symbol.Value
		}
	}

	result := strings.Builder{}

	position = 0
	for index, ins := range instructions {
		if _, ok := ins.(DEFLABEL); ok {
			l := ins.(DEFLABEL)
			result.WriteString("\n")
			result.WriteString(l.Name)
			result.WriteString(":")
		} else if _, ok := ins.(DEFSYMBOL); ok {
			s := ins.(DEFSYMBOL)
			result.WriteString(s.String())
		} else {
			a.symbols[CURRENTINSTRUCTION] = position + codeStartOffset
			a.symbols[NEXTINSTRUCTION] = getNextExecutableInstructionLoc(a.symbols[CURRENTINSTRUCTION], index, instructions)
			result.WriteString("\t")
			result.WriteString(fmt.Sprintf("%s:\t", utils.ValueToString(position+codeStartOffset)))

			emit, err := ins.Emit(a.ResolveLabel, a.ResolveSymbol)
			if err != nil {
				return "", err
			}
			result.WriteString("{")
			for i := 0; i < ins.Size(); i++ {
				result.WriteString(fmt.Sprintf("%s", utils.ValueToString(emit[i])))
				if i < ins.Size()-1 {
					result.WriteString(" ")
				}
			}
			result.WriteString("}")
			switch ins.Size() {
			case 4:
				result.WriteString("\t")
			default:
				result.WriteString(strings.Repeat("\t", 3))
			}
			result.WriteString(ins.String())
			position += uint16(ins.Size())
		}
		result.WriteString("\n")
	}

	return result.String(), nil
}

func isReservedSymbol(name string) bool {
	if _, ok := reservedSymbols[name]; ok {
		return true
	}
	return false
}

func getNextExecutableInstructionLoc(currentOffset uint16, currentInstrIndex int, instructions []Instruction) uint16 {
	// if at the end then just return the location outside the loop
	if currentInstrIndex == (len(instructions) - 1) {
		return uint16(currentOffset + uint16(instructions[currentInstrIndex].Size()))
	}

	nextInstructionPos := 0

	for i := currentInstrIndex; i < len(instructions); i++ {
		instruction := instructions[i]
		if _, ok := instruction.(DEFLABEL); ok {
			continue
		} else if _, ok := instruction.(DEFSYMBOL); ok {
			continue
		} else {
			if currentInstrIndex == i {
				nextInstructionPos += instruction.Size()
			} else {
				break
			}
		}

	}
	return uint16(currentOffset) + uint16(nextInstructionPos)
}
