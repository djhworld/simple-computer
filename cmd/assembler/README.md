Really crude assembler for the simple computer.

See main README for a description of all the available instructions and pseudo instructions provided by the assembler.

Features:

* Parses a text based description of a program, and assembles into a little-endian binary file
* Label address calculation
* Maintains two reserved `SYMBOL`s that can be referenced
  * `CURRENTINSTRUCTION`: contains the memory address of the current instruction
  * `NEXTINSTRUCTION`: contains the memory address of the next instruction relative to the current instruction
* Writes output as little-endian binary data

# Usage

```
  -i string
        input file (default: stdin)
  -o string
        output file (default: stdout)
  -s    output assembly as string
```

Example: 

```
go run github.com/djhworld/simple-computer/cmd/assembler -i myprogram.asm -o myprogram.bin
```

Note you can see the output of the assembler as a string for debugging purposes by passing the `-s` flag

Example: 

```
go run github.com/djhworld/simple-computer/cmd/assembler -i myprogram.asm -s
```

# Assembler directives

## Labels

Labels can be defined by an alpha-numeric sequence of characters followed by a colon, e.g.

```
MAIN:
   <instructions>
```


## Symbols

Symbols can be defined by a percentage sign `%` followed by an alpha-numeric sequence of characters, an equals sign `=` and a numeric value

```
%LINE-WIDTH = 0x001E
%ONE = 1
%DISPLAY-ADAPTER-ADDR = 0x7
```

