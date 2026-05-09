package asm

// Base opcodes for two-register instructions (indexed by A-register: R0=0, R1=1, R2=2, R3=3).
// Add the B-register index to get the full opcode.
var ldBases  = [4]uint16{0x0000, 0x0004, 0x0008, 0x000C}
var stBases  = [4]uint16{0x0010, 0x0014, 0x0018, 0x001C}
var addBases = [4]uint16{0x0080, 0x0084, 0x0088, 0x008C}
var andBases = [4]uint16{0x00C0, 0x00C4, 0x00C8, 0x00CC}
var orBases  = [4]uint16{0x00D0, 0x00D4, 0x00D8, 0x00DC}
var xorBases = [4]uint16{0x00E0, 0x00E4, 0x00E8, 0x00EC}
var cmpBases = [4]uint16{0x00F0, 0x00F4, 0x00F8, 0x00FC}

// Full opcodes for single-register instructions (indexed by register: R0=0 … R3=3).
var dataOpcodes = [4]uint16{0x0020, 0x0021, 0x0022, 0x0023}
var jrOpcodes   = [4]uint16{0x0030, 0x0031, 0x0032, 0x0033}
var shrOpcodes  = [4]uint16{0x0090, 0x0095, 0x009A, 0x009F}
var shlOpcodes  = [4]uint16{0x00A0, 0x00A5, 0x00AA, 0x00AF}
var notOpcodes  = [4]uint16{0x00B0, 0x00B5, 0x00BA, 0x00BF}

// Fixed opcodes.
const (
	opJMP     = uint16(0x0040)
	opCLF     = uint16(0x0060)
	opINData  = uint16(0x0070)
	opINAddr  = uint16(0x0074)
	opOUTData = uint16(0x0078)
	opOUTAddr = uint16(0x007C)
)

// Conditional jump opcodes keyed by flag string (e.g. "CAEZ").
var jmpfOpcodes = map[string]uint16{
	"Z":    0x0051,
	"E":    0x0052,
	"EZ":   0x0053,
	"A":    0x0054,
	"AZ":   0x0055,
	"AE":   0x0056,
	"AEZ":  0x0057,
	"C":    0x0058,
	"CZ":   0x0059,
	"CE":   0x005A,
	"CEZ":  0x005B,
	"CA":   0x005C,
	"CAZ":  0x005D,
	"CAE":  0x005E,
	"CAEZ": 0x005F,
}
