package utils

import (
	"fmt"
)

func ValueToString(val uint16) string {
	if val <= 0x000F {
		return fmt.Sprintf("0x000%X", val)
	} else if val <= 0x00FF {
		return fmt.Sprintf("0x00%X", val)
	} else if val <= 0x0FFF {
		return fmt.Sprintf("0x0%X", val)
	}
	return fmt.Sprintf("0x%X", val)
}
