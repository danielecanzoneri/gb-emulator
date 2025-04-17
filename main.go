package main

import (
	"github.com/danielecanzoneri/gb-emulator/cpu"
	"github.com/danielecanzoneri/gb-emulator/memory"
)

func main() {
	mem := &memory.Memory{}
	gb := cpu.CPU{
		Mem: mem,
	}

	// Inizializza PC, carica ROM (a breve)
	gb.PC = 0x0100

	// Loop semplificato
	for i := 0; i < 10; i++ {
		opcode := gb.Mem.Read(gb.PC)
		gb.PC++
		// (decodifica ed esegui istruzione)
		println("Fetched opcode:", opcode)
	}
}
