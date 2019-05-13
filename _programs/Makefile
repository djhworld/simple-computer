all: ascii brush text-writer me

ascii:
	../bin/generator ascii > ascii.asm
	../bin/assembler -i ascii.asm -o ascii.bin

text-writer:
	../bin/generator text-writer > text-writer.asm
	../bin/assembler -i text-writer.asm -o text-writer.bin

brush:
	../bin/generator brush > brush.asm
	../bin/assembler -i brush.asm -o brush.bin

me:
	../bin/generator me > me.asm
	../bin/assembler -i me.asm -o me.bin