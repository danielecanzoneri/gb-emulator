package mmu

import (
	"github.com/danielecanzoneri/lucky-boy/util"
)

const (
	dmaDuration = 0xA0
	dmaAddress  = 0xFF46
)

// DMA transfer from XX00-XX9F to FE00-FE9F where XX = 00 to DF lasting 160 cycles
func (mmu *MMU) DMA(xx uint8) {
	if xx > 0xDF {
		xx &= 0xDF
	}

	// Wait two cycles before starting dma
	mmu.delayDmaTicks = 8
	mmu.dmaReg = xx
}

func (mmu *MMU) VDMA(length uint8) {
	if util.ReadBit(length, 7) > 0 {
		// Delay VDMA until HBlank
		mmu.vDMAHBlank = true
		mmu.ppu.HBlankCallback = func() {
			mmu.vDMAActive = true
		}
	} else {
		// Immediately start VDMA
		mmu.vDMAActive = true
	}

	mmu.vDMALength = length & 0x7F
}

// VDMAActive returns true if a vRAM DMA transfer is being performed
func (mmu *MMU) VDMAActive() bool {
	return mmu.vDMAActive
}

func (mmu *MMU) vDMATransfer() {
	if mmu.IsCPUHalted() {
		return
	}

	// Transfer 0x10 bytes
	for i := uint16(0); i < 0x10; i++ {
		src := mmu.read(mmu.vDMASrcAddress + i)
		mmu.ppu.VDMAWrite(mmu.vDMADestAddress+i, src)
	}
	mmu.vDMASrcAddress += 0x10
	mmu.vDMADestAddress += 0x10

	// Stop VDMA
	if mmu.vDMALength == 0 {
		mmu.vDMAActive = false
		mmu.vDMAHBlank = false
		mmu.ppu.HBlankCallback = nil
		mmu.vDMALength = 0x7F
		return
	}

	// Set up next Transfer
	if mmu.vDMAHBlank {
		mmu.vDMAActive = false
		mmu.vDMATicks = 0
	}

	mmu.vDMALength--
}
