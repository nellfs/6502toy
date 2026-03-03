package cpu

import "fmt"

type CPU struct {
	A   byte
	PC  uint16
	Mem [65536]byte
}

func (c *CPU) Step() {
	opcode := c.Mem[c.PC]
	c.PC++

	switch opcode {
	case 0xA9:
		value := c.Mem[c.PC]
		c.PC++
		c.A = value
	default:
		fmt.Printf("unknown opcode: %02X\n", opcode)
	}
}

// Run executes instructions until a BRK (0x00) is encountered or
// the given maximum number of steps has been executed.
func (c *CPU) Run(maxSteps int) {
	for i := 0; i < maxSteps; i++ {
		opcode := c.Mem[c.PC]
		if opcode == 0x00 {
			fmt.Printf("halt (BRK) at %04X\n", c.PC)
			return
		}

		c.Step()
	}

	fmt.Printf("warning: max steps (%d) reached at PC=%04X\n", maxSteps, c.PC)
}
