package cartridge

import "log"

func (c *Cartridge) Write(addr uint16, value uint8) {
	switch {
	case addr < 0x2000:
		c.MBC.enableRAM(value)
	case addr < 0x4000:
		c.MBC.SetROMBank(value)
	case addr < 0x6000:
		c.MBC.SetRAMBank(value)
	case addr < 0x8000:
		c.MBC.SetMode(value)

	case 0xA000 <= addr && addr < 0xC000:
		RAMAddress := c.MBC.computeRAMAddress(addr)
		c.MBC.RAM[RAMAddress] = value

	default:
		log.Printf("[WARN] Cartridge Write address is out of range: %04X\n", addr)
	}
}

func (c *Cartridge) Read(addr uint16) uint8 {
	switch {
	case addr < 0x8000:
		cartridgeAddress := c.MBC.computeROMAddress(addr)
		return c.Data[cartridgeAddress]

	case 0xA000 <= addr && addr < 0xC000:
		RAMAddress := c.MBC.computeRAMAddress(addr)
		return c.MBC.RAM[RAMAddress]

	default:
		log.Printf("[WARN] Cartridge Read address is out of range: %04X\n", addr)
		return 0xFF
	}
}
