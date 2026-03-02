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
