package cpu

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var Debug = false
var steps = 1

func (cpu *CPU) logState(opcode uint8) {
	if Debug {
		steps--
		if steps > 0 {
			return
		}
		steps = 1

		if opcode == PREFIX_OPCODE {
			fmt.Printf(
				"PC: %04X | OP: %02X %02X | ",
				cpu.PC-1, opcode, cpu.MMU.Read(cpu.PC),
			)
		} else {
			fmt.Printf(
				"PC: %04X | OP: %02X | ",
				cpu.PC-1, opcode,
			)
		}
		fmt.Printf(
			"A:%02X B:%02X C:%02X D:%02X E:%02X H:%02X L:%02X | F:%02X | SP:%04X | Steps:%d\n",
			cpu.A, cpu.B, cpu.C, cpu.D, cpu.E, cpu.H, cpu.L,
			cpu.F, cpu.SP, cpu.steps,
		)
		fmt.Println("Timer:", cpu.Timer)

		// Get user input
		scanner := bufio.NewScanner(os.Stdin)
		for {
			fmt.Print("> ")
			if !scanner.Scan() {
				break
			}
			input := scanner.Text()

			if len(input) == 0 {
				break
			} else if input[0] == 'x' {
				// read address
				addr, err := parseAddr(input[2:])

				if err != nil {
					fmt.Println("Invalid address:", err)
					continue
				} else {
					fmt.Printf("[%04X] = %02X\n", addr, cpu.MMU.Read(addr))
				}
			} else {
				s, err := strconv.Atoi(input)
				if err != nil {
					fmt.Println("Unexpected command:", input)
					continue
				}
				steps = s
				break
			}
		}
	}
}

func parseAddr(addr string) (uint16, error) {
	addr = strings.TrimPrefix(addr, "0x")
	addr = strings.TrimSpace(addr)
	_addr, err := strconv.ParseUint(addr, 16, 16)
	return uint16(_addr), err
}
