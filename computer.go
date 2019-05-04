package main

import (
	"fmt"
	"log"
	"time"

	"github.com/djhworld/simple-computer/components"
	"github.com/djhworld/simple-computer/cpu"
	"github.com/djhworld/simple-computer/io"
	"github.com/djhworld/simple-computer/memory"
)

type PrintStateConfig struct {
	PrintState      bool
	PrintStateEvery int
}

type SimpleComputer struct {
	memory  *memory.Memory64K
	cpu     *cpu.CPU
	mainBus *components.Bus

	displayAdapter  *io.DisplayAdapter
	screenControl   *io.ScreenControl
	keyboardAdapter *io.KeyboardAdapter

	screenChannel chan *[160][240]byte
	quitChannel   chan bool
}

func NewComputer(screenChannel chan *[160][240]byte, quitChannel chan bool) *SimpleComputer {
	c := new(SimpleComputer)

	c.screenChannel = screenChannel
	c.quitChannel = quitChannel

	c.mainBus = components.NewBus(16)
	c.memory = memory.NewMemory64K(c.mainBus)
	c.cpu = cpu.NewCPU(c.mainBus, c.memory)

	c.keyboardAdapter = io.NewKeyboardAdapter()
	c.cpu.ConnectPeripheral(c.keyboardAdapter)

	c.displayAdapter = io.NewDisplaydAdapter()
	c.screenControl = io.NewScreenControl(c.displayAdapter, c.screenChannel, c.quitChannel)
	c.cpu.ConnectPeripheral(c.displayAdapter)

	return c
}

func (c *SimpleComputer) ConnectKeyboard(keyboard *io.Keyboard) {
	keyboard.ConnectTo(c.keyboardAdapter.KeyboardInBus)
}

func (c *SimpleComputer) LoadToRAM(offset uint16, values []uint16) {
	log.Printf("Loading %d words to RAM at offset 0x%X", len(values), offset)
	for i := 0; i < len(values); i++ {
		c.putValueInRAM(offset+uint16(i), values[i])
	}
}

func (c *SimpleComputer) putValueInRAM(address, value uint16) {
	c.memory.AddressRegister.Set()
	c.mainBus.SetValue(address)
	c.memory.Update()

	c.memory.AddressRegister.Unset()
	c.memory.Update()

	c.mainBus.SetValue(value)
	c.memory.Set()
	c.memory.Update()

	c.memory.Unset()
	c.memory.Update()
}

func (c *SimpleComputer) Run(tickInterval <-chan time.Time, printStateConfig PrintStateConfig) {
	log.Println("Starting computer....")
	go c.screenControl.Run()

	steps := 0
	for {
		<-tickInterval
		c.cpu.Step()

		if printStateConfig.PrintState {
			if steps%printStateConfig.PrintStateEvery == 0 {
				fmt.Println("COMPUTER\n-----------------------------------------------------------")
				fmt.Printf("Step count = %d, printing state every %d steps\n\n", steps, printStateConfig.PrintStateEvery)
				fmt.Println("CPU\n----------------------------------------")
				fmt.Println(c.cpu.String())
				fmt.Println()
			}
		}
		steps++
	}
}
