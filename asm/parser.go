package asm

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

var IS_DEFLABEL *regexp.Regexp = regexp.MustCompile("[A-Za-z0-9-]+:")
var IS_DEFSYMBOL *regexp.Regexp = regexp.MustCompile(`%([A-Za-z0-9-]+)\s*=\s*((0x)?[0-9a-fA-F]+)`)
var INSTRUCTION *regexp.Regexp = regexp.MustCompile(`(CALL)\s*([A-Za-z0-9-]+)|(DATA)\s*(R\d,\s*.+)|(CLF)|(JR)\s*(R\d)|(NOT)\s*(R\d)|(SHL)\s*(R\d)|(SHR)\s*(R\d)|(ADD)\s*(R\d,\s*R\d)|(CMP)\s*(R\d,\s*R\d)|(AND)\s*(R\d,\s*R\d)|(OR)\s*(R\d,\s*R\d)|(LD)\s*(R\d,\s*R\d)|(ST)\s*(R\d,\s*R\d)|(XOR)\s*(R\d,\s*R\d)|(OUT)\s*([A-Za-z]+,\s*R\d)|(IN)\s*([A-Za-z]+,\s*R\d)|(JMP[A-Z]+)\s*([A-Za-z0-9-]+)|(JMP)\s*([A-Za-z0-9-]+)`)
var TWO_REGISTER_EXTRACTOR *regexp.Regexp = regexp.MustCompile(`R(\d),\s*R(\d)\s*`)
var ONE_REGISTER_EXTRACTOR *regexp.Regexp = regexp.MustCompile(`R(\d)\s*`)
var DATA_EXTRACTOR *regexp.Regexp = regexp.MustCompile(`R(\d),\s*((0x)?[0-9a-fA-F]+|(%)([A-Za-z0-9-]+))`)
var IO_EXTRACTOR *regexp.Regexp = regexp.MustCompile(`(Addr|Data),\s*R(\d)`)
var LABEL_EXTRACTOR *regexp.Regexp = regexp.MustCompile(`([A-Za-z0-9-]+)`)
var FLAGS_EXTRACTOR *regexp.Regexp = regexp.MustCompile(`([CAEZ]+)`)

var REGISTERS map[string]REGISTER = map[string]REGISTER{
	"0": REG0,
	"1": REG1,
	"2": REG2,
	"3": REG3,
}

type Parser struct {
}

func (p *Parser) Parse(input io.Reader) ([]Instruction, error) {
	scanner := bufio.NewScanner(input)
	instructions := []Instruction{}

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		if IS_DEFLABEL.MatchString(line) {
			instructions = append(instructions, processLabel(line))
		} else if IS_DEFSYMBOL.MatchString(line) {
			if ins, err := parseDefSymbol(line); err == nil {
				instructions = append(instructions, ins)
			} else {
				return nil, err
			}
		} else if INSTRUCTION.MatchString(line) {
			if ins, err := parseInstruction(line); err == nil {
				instructions = append(instructions, ins)
			} else {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("unsupported/unparseable line: %s", line)
		}

	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return instructions, nil
}

func parseDefSymbol(line string) (Instruction, error) {
	tokens := IS_DEFSYMBOL.FindStringSubmatch(line)
	if len(tokens) != 4 {
		return nil, fmt.Errorf("could not parse the arguments correctly out of DEFSYMBOL %s", line)
	}

	var value uint64
	var err error
	// parse in base 16 or 10
	if tokens[3] == "0x" {
		value, err = strconv.ParseUint(strings.Replace(tokens[2], "0x", "", -1), 16, 16)
	} else {
		value, err = strconv.ParseUint(tokens[2], 10, 16)
	}

	if err != nil {
		return nil, err
	}

	return DEFSYMBOL{tokens[1], uint16(value)}, nil
}

func processLabel(line string) DEFLABEL {
	line = strings.Replace(line, ":", "", -1)

	return DEFLABEL{line}
}

func stripEmptyGroups(input []string) []string {
	result := []string{}
	for _, t := range input {
		if t != "" {
			result = append(result, t)
		}
	}
	return result
}

func parseInstruction(line string) (Instruction, error) {
	tokens := INSTRUCTION.FindStringSubmatch(line)
	tokens = stripEmptyGroups(tokens)
	instructionName := tokens[1]

	var operands string
	if len(tokens) > 2 {
		operands = tokens[2]
	}

	var instruction Instruction
	var err error
	switch instructionName {
	case "ADD", "AND", "XOR", "OR", "CMP", "LD", "ST":
		instruction, err = parseTwoRegisterInstruction(instructionName, operands)
	case "SHR", "SHL", "NOT", "JR":
		instruction, err = parseOneRegisterInstruction(instructionName, operands)
	case "DATA":
		instruction, err = parseDataInstruction(operands)
	case "CLF":
		instruction = CLF{}
	case "OUT", "IN":
		instruction, err = parseIOInstruction(instructionName, operands)
	case "CALL", "JMP", "JMPZ", "JMPE", "JMPEZ", "JMPA", "JMPAZ", "JMPAE", "JMPAEZ", "JMPC", "JMPCZ", "JMPCE", "JMPCEZ", "JMPCA", "JMPCAZ", "JMPCAE", "JMPCAEZ":
		instruction, err = parseLabelledJump(instructionName, operands)
	default:
		return nil, fmt.Errorf("unknown instruction name '%s'", instructionName)
	}
	return instruction, err
}

func parseLabelledJump(name string, operands string) (Instruction, error) {
	arguments := LABEL_EXTRACTOR.FindStringSubmatch(operands)
	if len(arguments) != 2 {
		return nil, fmt.Errorf("Instruction %s is not valid", name)
	}

	switch name {
	case "JMP":
		return JMP{LABEL{arguments[1]}}, nil
	case "CALL":
		return CALL{LABEL{arguments[1]}}, nil
	case "JMPZ", "JMPE", "JMPEZ", "JMPA", "JMPAZ", "JMPAE", "JMPAEZ", "JMPC", "JMPCZ", "JMPCE", "JMPCEZ", "JMPCA", "JMPCAZ", "JMPCAE", "JMPCAEZ":
		flags, err := extractFlagsFrom(name)
		if err != nil {
			return nil, err
		}
		return JMPF{flags, LABEL{arguments[1]}}, nil
	default:
		return nil, fmt.Errorf("Unsupported labelled jump instruction %s", name)
	}
}

func extractFlagsFrom(name string) ([]string, error) {
	tokens := FLAGS_EXTRACTOR.FindStringSubmatch(name)
	if len(tokens) != 2 {
		return nil, fmt.Errorf("Could not extract flags from instruction %s", name)
	}

	return strings.Split(tokens[1], ""), nil
}

func parseIOInstruction(name string, operands string) (Instruction, error) {
	arguments := IO_EXTRACTOR.FindStringSubmatch(operands)
	if len(arguments) != 3 {
		return nil, fmt.Errorf("Instruction %s is not valid", name)
	}

	var iomode IO_MODE
	switch arguments[1] {
	case "Data":
		iomode = DATA_MODE
	case "Addr":
		iomode = ADDRESS_MODE
	default:
		return nil, fmt.Errorf("Unsupported IO mode %s for instruction %s (must be Addr|Data)", arguments[1], name)
	}

	var register REGISTER
	if v, ok := REGISTERS[arguments[2]]; !ok {
		return nil, fmt.Errorf("Unknown register %s for instruction %s", arguments[2], name)
	} else {
		register = v
	}

	switch name {
	case "OUT":
		return OUT{iomode, register}, nil
	case "IN":
		return IN{iomode, register}, nil
	default:
		return nil, fmt.Errorf("Unknown IO instruction %s (only IN or OUT supported)", name)
	}
}

func parseTwoRegisterInstruction(name string, operands string) (Instruction, error) {
	registers := TWO_REGISTER_EXTRACTOR.FindStringSubmatch(operands)

	if len(registers) != 3 {
		return nil, fmt.Errorf("Instruction %s only supports two command separated registers", name)
	}

	var register1 REGISTER
	if v, ok := REGISTERS[registers[1]]; !ok {
		return nil, fmt.Errorf("Unknown register %s for instruction %s", registers[1], name)
	} else {
		register1 = v
	}

	var register2 REGISTER
	if v, ok := REGISTERS[registers[2]]; !ok {
		return nil, fmt.Errorf("Unknown register %s for instruction %s", registers[2], name)
	} else {
		register2 = v
	}

	switch name {
	case "ADD":
		return ADD{register1, register2}, nil
	case "AND":
		return AND{register1, register2}, nil
	case "XOR":
		return XOR{register1, register2}, nil
	case "OR":
		return OR{register1, register2}, nil
	case "LD":
		return LOAD{register1, register2}, nil
	case "ST":
		return STORE{register1, register2}, nil
	case "CMP":
		return CMP{register1, register2}, nil
	default:
		return nil, fmt.Errorf("unknown/unsupported instruction %s", name)
	}
}

func parseOneRegisterInstruction(name string, operands string) (Instruction, error) {
	registers := ONE_REGISTER_EXTRACTOR.FindStringSubmatch(operands)

	if len(registers) != 2 {
		return nil, fmt.Errorf("Instruction %s only supports one register", name)
	}

	var register1 REGISTER
	if v, ok := REGISTERS[registers[1]]; !ok {
		return nil, fmt.Errorf("Unknown register %s for instruction %s", registers[1], name)
	} else {
		register1 = v
	}

	switch name {
	case "SHL":
		return SHL{register1}, nil
	case "SHR":
		return SHR{register1}, nil
	case "NOT":
		return NOT{register1}, nil
	case "JR":
		return JR{register1}, nil
	default:
		return nil, fmt.Errorf("unknown/unsupported instruction %s", name)
	}
}

/*
func parseDataInstruction(operands string) (Instruction, error) {
	arguments := DATA_EXTRACTOR.FindStringSubmatch(operands)
	if len(arguments) != 4 {
		return nil, fmt.Errorf("could not parse the arguments correctly out of DATA %s", operands)
	}

	var register REGISTER
	if v, ok := REGISTERS[arguments[1]]; !ok {
		return nil, fmt.Errorf("Unknown register %s for DATA instruction", arguments[1])
	} else {
		register = v
	}

	var value uint64
	var err error
	// parse in base 16 or 10
	if arguments[3] == "0x" {
		value, err = strconv.ParseUint(strings.Replace(arguments[2], "0x", "", -1), 16, 16)
	} else {
		value, err = strconv.ParseUint(arguments[2], 10, 16)
	}

	if err != nil {
		return nil, err
	}

	return DATA{register, uint16(value)}, nil
}
*/

func parseDataInstruction(operands string) (Instruction, error) {
	arguments := DATA_EXTRACTOR.FindStringSubmatch(operands)
	if len(arguments) != 6 {
		return nil, fmt.Errorf("could not parse the arguments correctly out of DATA %s", operands)
	}

	var register REGISTER
	if v, ok := REGISTERS[arguments[1]]; !ok {
		return nil, fmt.Errorf("Unknown register %s for DATA instruction", arguments[1])
	} else {
		register = v
	}

	if arguments[4] == "%" {
		return DATA{register, SYMBOL{arguments[5]}}, nil
	} else {
		var value uint64
		var err error
		// parse in base 16 or 10
		if arguments[3] == "0x" {
			value, err = strconv.ParseUint(strings.Replace(arguments[2], "0x", "", -1), 16, 16)
		} else {
			value, err = strconv.ParseUint(arguments[2], 10, 16)
		}

		if err != nil {
			return nil, err
		}

		return DATA{register, NUMBER{uint16(value)}}, nil
	}

}
