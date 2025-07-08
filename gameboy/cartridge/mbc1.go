package cartridge

import "log"

type MBC1 struct {
	header  *Header
	battery bool // If battery is present RAM should be stored

	ROMBanks uint8
	RAMBanks uint8

	ROM []uint8
	RAM []uint8

	// Registers
	ramEnabled    bool
	romBankNumber uint8 // 5 bit register (if 0 is read as 1)
	ramBankNumber uint8 // 2 bit register (can be used as high bits of romBankNumber if number of banks is > 32)
	// 0: ROM and RAM banks 0 fixed, 1: ROM bank 0 and RAM banks can be switched
	bankingMode uint8

	useRamBankNumberAsHighRomBankNumber bool
}

func (mbc *MBC1) RAMDump() []uint8 {
	if mbc.battery {
		return mbc.RAM
	}

	return nil
}

func (mbc *MBC1) Header() *Header {
	return mbc.header
}

func NewMBC1(rom []uint8, ram bool, savData []uint8, header *Header, battery bool) *MBC1 {
	// TODO - Detect and handle MBC1M (Multi-Game carts)
	mbc := &MBC1{
		header:                              header,
		battery:                             battery,
		ROMBanks:                            uint8(header.ROMBanks),
		RAMBanks:                            uint8(header.RAMBanks),
		ROM:                                 rom,
		romBankNumber:                       1,
		useRamBankNumberAsHighRomBankNumber: header.ROMBanks > 32,
	}
	if ram && header.RAMBanks == 0 {
		log.Println("[WARN] Cartridge header specifies RAM present, but RAM banks is set to 0")
		mbc.RAMBanks = 1
	}

	if ram {
		switch {
		case battery && len(savData) != int(mbc.RAMBanks)*0x2000:
			log.Println("[WARN] sav file was of a different dimension than expected, resetting to zero")
			fallthrough
		case savData == nil:
			savData = make([]uint8, int(mbc.RAMBanks)*0x2000)
		}
		mbc.RAM = savData
	}

	return mbc
}

func (mbc *MBC1) Write(addr uint16, value uint8) {
	// Set MBC1 registers
	switch {
	case addr < 0x2000:
		// Low nibble = 0xA enables RAM
		mbc.ramEnabled = value&0x0F == 0xA

	case addr < 0x4000:
		// Only lower 5 bits are used
		value = value & 0x1F

		// Register 0 behaves as 1
		if value == 0 {
			value = 1
		}

		mbc.romBankNumber = value

	case addr < 0x6000:
		// Only lower 2 bits are used
		mbc.ramBankNumber = value & 0x3

	case addr < 0x8000:
		// 1 bit register
		mbc.bankingMode = value & 1

	case 0xA000 <= addr && addr < 0xC000:
		if mbc.ramEnabled {
			RAMAddress := mbc.computeRamAddress(addr)
			mbc.RAM[RAMAddress] = value
		}
	}
}

func (mbc *MBC1) Read(addr uint16) uint8 {
	switch {
	case addr < 0x8000:
		cartridgeAddress := mbc.computeRomAddress(addr)
		return mbc.ROM[cartridgeAddress]

	case 0xA000 <= addr && addr < 0xC000:
		if mbc.ramEnabled {
			RAMAddress := mbc.computeRamAddress(addr)
			return mbc.RAM[RAMAddress]
		}
	}

	return 0xFF
}

func (mbc *MBC1) computeRomAddress(cpuAddress uint16) uint {
	// bank number: 2 bits - 5 bits, cpuAddress: 14 bits
	var bankNumber uint8 = 0

	switch {
	case cpuAddress < 0x4000:
		if mbc.bankingMode == 1 && mbc.useRamBankNumberAsHighRomBankNumber {
			bankNumber = mbc.ramBankNumber << 5
		}

	case cpuAddress < 0x8000:
		bankNumber = mbc.romBankNumber
		if mbc.useRamBankNumberAsHighRomBankNumber {
			bankNumber |= mbc.ramBankNumber << 5
		}

	default:
		panic("should never happen")
	}

	// Bank number is masked to the required number of bits
	bankNumber %= mbc.ROMBanks

	return uint(bankNumber)<<14 | uint(cpuAddress&0x3FFF)
}

func (mbc *MBC1) computeRamAddress(cpuAddress uint16) uint {
	// bank number: 2 bits, cpuAddress: 13 bits
	switch {
	case 0xA000 <= cpuAddress && cpuAddress < 0xC000:
		cpuAddress = cpuAddress & 0x1FFF

		if mbc.bankingMode == 1 && !mbc.useRamBankNumberAsHighRomBankNumber {
			// Bank number is masked to the required number of bits
			bank := mbc.ramBankNumber % mbc.RAMBanks
			return uint(bank)<<13 | uint(cpuAddress)
		}
		return uint(cpuAddress)

	default:
		panic("should never happen")
	}
}
