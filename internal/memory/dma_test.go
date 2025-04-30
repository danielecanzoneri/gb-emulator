package memory

import (
	"github.com/danielecanzoneri/gb-emulator/internal/cartridge"
	"github.com/danielecanzoneri/gb-emulator/internal/ppu"
	"testing"
)

func TestDMA(t *testing.T) {
	mmu := &MMU{PPU: &ppu.PPU{}}
	mmu.SetMBC(&cartridge.Header{})

	// Write data to RAM
	startAddr := 0xA000
	for i := 0; i < dmaDuration; i++ {
		mmu.write(uint16(startAddr+i), uint8(0xA0-i))
	}

	testDMAAddresses := []uint16{0x0000, 0x3FFF, 0x7FFF, 0x8000, 0x9000, 0xA000, 0xC000, 0xD000, 0xF000, 0xFF7F}
	var validAddress uint16 = 0xFF90 // Address in hRAM that is accessible during DMA
	mmu.write(validAddress, 0x12)

	mmu.write(dmaAddress, uint8(startAddr>>8))
	for i := 0; i < dmaDuration; i++ {
		mmu.Cycle()

		// Check that at every cycle everything read from the MMU will return the value being currently written in DMA
		for _, addr := range testDMAAddresses {
			if mmu.Read(addr) != uint8(0xA0-i) {
				t.Errorf("Incorrect DMA value read on 0x%04X at cycle %d: %02X", addr, i, mmu.Read(addr))
			}

			if mmu.Read(validAddress) != 0x12 {
				t.Errorf("Incorrect value read in hRAM during DMA")
			}
		}
	}
}
