package memory

import (
	"github.com/danielecanzoneri/gb-emulator/emulator/internal/gameboy/cartridge"
	"github.com/danielecanzoneri/gb-emulator/emulator/internal/gameboy/ppu"
	"testing"
)

func TestDMA(t *testing.T) {
	mmu := &MMU{Cartridge: cartridge.NewMBC1(make([]uint8, 0x8000), true, nil, &cartridge.Header{ROMBanks: 1, RAMBanks: 1}, false), PPU: &ppu.PPU{}}

	// Write data to RAM
	startAddr := 0xA000
	for i := 0; i < dmaDuration; i++ {
		mmu.write(uint16(startAddr+i), uint8(0xA0-i))
	}

	var oamAddress uint16 = 0xFE12
	var validAddress uint16 = 0xFF90 // Address in hRAM that is accessible during DMA
	mmu.write(validAddress, 0x12)

	mmu.write(dmaAddress, uint8(startAddr>>8))
	mmu.Tick(4)
	for i := 0; i < dmaDuration; i++ {
		mmu.Tick(4)

		// Check that at every cycle everything read from the OAM will return 0xFF
		if mmu.Read(oamAddress) != 0xFF {
			t.Errorf("Incorrect OAM value read at cycle %d: %02X", i, mmu.Read(oamAddress))
		}

		if mmu.Read(validAddress) != 0x12 {
			t.Errorf("Incorrect value read in hRAM during DMA")
		}
	}
}
