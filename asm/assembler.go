package asm

import (
	"fmt"
	"strings"

	"github.com/djhworld/simple-computer/utils"
)

type Assembler struct {
	labels map[string]uint16
}

func (a *Assembler) ResolveLabel(label LABEL) (uint16, error) {
	if v, ok := a.labels[label.Name]; !ok {
		return 0x0000, fmt.Errorf("Cannot find label: %s in label map", label.Name)
	} else {
		return v, nil
	}
}

func (a *Assembler) Process(codeStartOffset uint16, instructions []Instruction) ([]uint16, error) {
	a.labels = make(map[string]uint16)
	position := uint16(0)

	//calculate labels
	for _, ins := range instructions {
		position += uint16(ins.Size())

		if label, ok := ins.(LABEL); ok {
			if _, ok := a.labels[label.Name]; ok {
				return nil, fmt.Errorf("label '%s' already exists, all labels should be unique", label.Name)
			}

			a.labels[label.Name] = position + codeStartOffset
		}
	}

	emitted := []uint16{}

	for _, ins := range instructions {
		if _, ok := ins.(LABEL); !ok {
			emit, err := ins.Emit(a.ResolveLabel)
			if err != nil {
				return nil, err
			}

			for i := 0; i < ins.Size(); i++ {
				emitted = append(emitted, emit[i])
			}
		}
	}

	return emitted, nil
}

func (a *Assembler) ToString(codeStartOffset uint16, instructions []Instruction) (string, error) {
	a.labels = make(map[string]uint16)
	position := uint16(0)

	//calculate lengths
	for _, ins := range instructions {
		position += uint16(ins.Size())

		if label, ok := ins.(LABEL); ok {
			a.labels[label.Name] = position + codeStartOffset
		}
	}

	result := strings.Builder{}

	pos := uint16(0)
	for _, ins := range instructions {
		if _, ok := ins.(LABEL); !ok {
			result.WriteString("\t")
			result.WriteString(fmt.Sprintf("%s:\t", utils.ValueToString(pos+codeStartOffset)))

			emit, err := ins.Emit(a.ResolveLabel)
			if err != nil {
				return "", err
			}
			result.WriteString("{")
			switch ins.Size() {
			case 1:
				result.WriteString(fmt.Sprintf("%s\t      ", utils.ValueToString(emit[0])))
			case 2:
				result.WriteString(fmt.Sprintf("%s\t%s", utils.ValueToString(emit[0]), utils.ValueToString(emit[1])))
			}
			result.WriteString("}\t\t")
			result.WriteString(ins.String())
		} else {
			l := ins.(LABEL)
			result.WriteString("\n")
			result.WriteString(l.Name)
			result.WriteString(":\n")
		}
		result.WriteString("\n")
		pos += uint16(ins.Size())
	}

	return result.String(), nil
}
