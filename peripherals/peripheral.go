package peripherals

import (
	"github.com/djhworld/simple-computer/components"
)

type Peripheral interface {
	Connect(*components.IOBus, *components.Bus)
	Update()
}

