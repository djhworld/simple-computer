package main

import (
	"fmt"

	"github.com/djhworld/simple-computer/components"
)

func main() {
	clock := components.Clock{}
	clock.Start()

	b := false
	for {
		select {
		case <-clock.BaseClock.C:
			if b {
				b = false
			} else {
				b = true
			}
		}
		fmt.Printf("b: %v\n", b)
	}
}
