package cartridge

import (
	"encoding/binary"
	"github.com/danielecanzoneri/gb-emulator/pkg/util"
	"log"
	"time"
)

const (
	// Range 0-59 (00-3B, bits 7-6 should always be 0)
	rtcSecondsMask = 0x3F
	rtcMinutesMask = 0x3F
	// Range 0-24 (00-18, bits 7-6-5 should always be 0)
	rtcHoursMask = 0x1F
	// Bit 7: carry bit, bit 6: RTC active, bit 0: bit 9 of day counter
	rtcControlMask = 0b11000001
)

type MBC3 struct {
	header  *Header
	battery bool // If battery is present RAM should be stored
	rtc     bool // If RTC is enabled store RTC registers in SAV file

	ROMBanks uint8
	RAMBanks uint8

	ROM []uint8
	RAM []uint8

	// Registers
	ramEnabled    bool  // Also enables and disables RTC
	romBankNumber uint8 // 7 bit register (if 0 is read as 1)
	// This register is also the RTC register if ranges between 08 and 0C,
	// otherwise it selects which RAM bank to use
	ramBankNumber uint8

	// When writing $00, and then $01 to this register, the current time becomes latched into the RTC registers.
	// The latched data will not change until it becomes latched again, by repeating the write $00->$01 procedure.
	// This provides a way to read the RTC registers while the clock keeps ticking.
	rtcS, lthRtcS   uint8
	rtcM, lthRtcM   uint8
	rtcH, lthRtcH   uint8
	rtcDL, lthRtcDL uint8
	rtcDH, lthRtcDH uint8
	lastWriteWas00  bool

	rtcClockCounter uint
}

func (mbc *MBC3) RAMDump() []uint8 {
	if mbc.battery {
		if mbc.rtc {
			dump := make([]uint8, len(mbc.RAM), len(mbc.RAM)+48)
			copy(dump, mbc.RAM)

			// Save RTC registers
			dump, _ = binary.Append(dump, binary.LittleEndian, int32(mbc.rtcS))
			dump, _ = binary.Append(dump, binary.LittleEndian, int32(mbc.rtcM))
			dump, _ = binary.Append(dump, binary.LittleEndian, int32(mbc.rtcH))
			dump, _ = binary.Append(dump, binary.LittleEndian, int32(mbc.rtcDL))
			dump, _ = binary.Append(dump, binary.LittleEndian, int32(mbc.rtcDH))
			dump, _ = binary.Append(dump, binary.LittleEndian, int32(mbc.lthRtcS))
			dump, _ = binary.Append(dump, binary.LittleEndian, int32(mbc.lthRtcM))
			dump, _ = binary.Append(dump, binary.LittleEndian, int32(mbc.lthRtcH))
			dump, _ = binary.Append(dump, binary.LittleEndian, int32(mbc.lthRtcDL))
			dump, _ = binary.Append(dump, binary.LittleEndian, int32(mbc.lthRtcDH))

			// Save current timestamp
			dump, _ = binary.Append(dump, binary.LittleEndian, time.Now().Unix())
			return dump
		}

		// Return only RAM
		return mbc.RAM
	}

	return nil
}

func (mbc *MBC3) Header() *Header {
	return mbc.header
}

func NewMBC3(rom []uint8, ram bool, savData []uint8, header *Header, battery bool, rtc bool) *MBC3 {
	mbc := &MBC3{
		header:        header,
		battery:       battery,
		rtc:           rtc,
		ROMBanks:      uint8(header.ROMBanks),
		RAMBanks:      uint8(header.RAMBanks),
		ROM:           rom,
		romBankNumber: 1,
	}
	if ram && header.RAMBanks == 0 {
		log.Println("[WARN] Cartridge header specifies RAM present, but RAM banks is set to 0")
		mbc.RAMBanks = 1
	}

	var rtcLen int
	if rtc {
		rtcLen = 48
	}

	if ram {
		ramLen := int(mbc.RAMBanks) * 0x2000
		switch {
		case battery && len(savData) != ramLen+rtcLen:
			log.Println("[WARN] sav file was of a different dimension than expected, resetting to zero")
			fallthrough

		case savData == nil:
			mbc.RAM = make([]uint8, ramLen)

		default: // SAV data is of correct format
			mbc.RAM = savData[:ramLen]
			if rtc {
				mbc.parseRTCData(savData[ramLen:])
			}
		}
	}

	return mbc
}
func (mbc *MBC3) Cycle() {
	// Check if RTC is enabled
	if util.ReadBit(mbc.rtcDH, 6) != 0 {
		return
	}

	// RTC clocking: Game Boy runs at 2^22 Hz and a cycle happens every 4 Hz
	mbc.rtcClockCounter++
	if mbc.rtcClockCounter == 1<<20 {
		mbc.rtcClockCounter = 0

		mbc.rtcS = (mbc.rtcS + 1) & rtcSecondsMask
		if mbc.rtcS == 60 {
			mbc.rtcS = 0
			mbc.rtcM = (mbc.rtcM + 1) & rtcMinutesMask
			if mbc.rtcM == 60 {
				mbc.rtcM = 0
				mbc.rtcH = (mbc.rtcH + 1) & rtcHoursMask
				if mbc.rtcH == 24 {
					mbc.rtcH = 0
					mbc.rtcDL++
					if mbc.rtcDL == 0 { // rtcDH bit 0 is bit 9 of day counter
						if util.ReadBit(mbc.rtcDH, 0) == 0 {
							util.SetBit(&mbc.rtcDH, 0, 1)
						} else { // Set carry bit
							util.SetBit(&mbc.rtcDH, 0, 0)
							util.SetBit(&mbc.rtcDH, 7, 1)
						}
					}
				}
			}
		}
	}
}

func (mbc *MBC3) Write(addr uint16, value uint8) {
	// Set MBC3 registers
	switch {
	case addr < 0x2000:
		// Low nibble = 0xA enables RAM
		mbc.ramEnabled = value&0x0F == 0xA

	case addr < 0x4000:
		// Only lower 7 bits are used
		value = value & 0x7F

		// Register 0 behaves as 1
		if value == 0 {
			value = 1
		}

		mbc.romBankNumber = value

	case addr < 0x6000:
		// Only lower 4 bits are used (bit 3 select RTC registers access)
		mbc.ramBankNumber = value & 0xF

	case addr < 0x8000:
		if value == 0x00 {
			mbc.lastWriteWas00 = true
		} else {
			mbc.lastWriteWas00 = false
			if value == 0x01 {
				// Latch registers
				mbc.lthRtcS = mbc.rtcS
				mbc.lthRtcM = mbc.rtcM
				mbc.lthRtcH = mbc.rtcH
				mbc.lthRtcDL = mbc.rtcDL
				mbc.lthRtcDH = mbc.rtcDH
			}
		}

	case 0xA000 <= addr && addr < 0xC000:
		if mbc.ramEnabled {
			// Access RTC registers
			switch mbc.ramBankNumber {
			case 0x8:
				mbc.rtcS = value & rtcSecondsMask
				mbc.rtcClockCounter = 0
			case 0x9:
				mbc.rtcM = value & rtcMinutesMask
			case 0xA:
				mbc.rtcH = value & rtcHoursMask
			case 0xB:
				mbc.rtcDL = value
			case 0xC:
				mbc.rtcDH = value & rtcControlMask
			default:
				// Access RAM
				RAMAddress := mbc.computeRamAddress(addr)
				mbc.RAM[RAMAddress] = value
			}
		}

	default:
		log.Printf("[WARN] MBC3 Write address is out of range: %04X\n", addr)
	}
}

func (mbc *MBC3) Read(addr uint16) uint8 {
	switch {
	case addr < 0x8000:
		cartridgeAddress := mbc.computeRomAddress(addr)
		return mbc.ROM[cartridgeAddress]

	case 0xA000 <= addr && addr < 0xC000:
		// Access RTC registers
		switch mbc.ramBankNumber {
		case 0x8:
			return mbc.lthRtcS
		case 0x9:
			return mbc.lthRtcM
		case 0xA:
			return mbc.lthRtcH
		case 0xB:
			return mbc.lthRtcDL
		case 0xC:
			return mbc.lthRtcDH
		default:
			// Access RAM
			RAMAddress := mbc.computeRamAddress(addr)
			return mbc.RAM[RAMAddress]
		}

	default:
		log.Printf("[WARN] MBC3 Read address is out of range: %04X\n", addr)
	}

	return 0xFF
}

func (mbc *MBC3) computeRomAddress(cpuAddress uint16) uint {
	// bank number: 7 bits, cpuAddress: 14 bits
	var bankNumber uint8 = 0

	switch {
	case cpuAddress < 0x4000: // bank number = 0
	case cpuAddress < 0x8000:
		bankNumber = mbc.romBankNumber

	default:
		panic("should never happen")
	}

	// Bank number is masked to the required number of bits
	bankNumber %= mbc.ROMBanks

	return uint(bankNumber)<<14 | uint(cpuAddress&0x3FFF)
}

func (mbc *MBC3) computeRamAddress(cpuAddress uint16) uint {
	// bank number: 2 bits, cpuAddress: 13 bits
	switch {
	case 0xA000 <= cpuAddress && cpuAddress < 0xC000:
		cpuAddress = cpuAddress & 0x1FFF

		// Bank number is masked to the required number of bits
		bank := mbc.ramBankNumber % mbc.RAMBanks
		return uint(bank)<<13 | uint(cpuAddress)

	default:
		panic("should never happen")
	}
}

func (mbc *MBC3) parseRTCData(data []uint8) {
	if len(data) != 48 {
		log.Println("[WARN] invalid RTC data length")
		return
	}

	// If RTC is present, additional data at end of SAV file should contain value of registers:
	// offset  size    desc
	// 0       4       time seconds
	// 4       4       time minutes
	// 8       4       time hours
	// 12      4       time days
	// 16      4       time days high
	// 20      4       latched time seconds
	// 24      4       latched time minutes
	// 28      4       latched time hours
	// 32      4       latched time days
	// 36      4       latched time days high
	// 40      8       unix timestamp when saving (64 bits little endian)
	// AFAIK a byte is saved as an int so extra bytes are 0
	mbc.rtcS = data[0]
	mbc.rtcM = data[4]
	mbc.rtcH = data[8]
	mbc.rtcDL = data[12]
	mbc.rtcDH = data[16]
	mbc.lthRtcS = data[20]
	mbc.lthRtcM = data[24]
	mbc.lthRtcH = data[28]
	mbc.lthRtcDL = data[32]
	mbc.lthRtcDH = data[36]

	// If RTC was enabled, advance registers for the time elapsed
	if util.ReadBit(mbc.rtcDH, 6) == 0 {
		timestamp := binary.LittleEndian.Uint64(data[40:])

		// Compute how much time has passed and update the registers accordingly
		saveTime := time.Unix(int64(timestamp), 0)
		elapsed := time.Since(saveTime)

		// Update seconds
		seconds := int(elapsed.Seconds()) % 60
		mbc.rtcS += uint8(seconds)
		extraMin := mbc.rtcS / 60
		mbc.rtcS %= 60

		// Update minutes
		minutes := int(elapsed.Minutes()) % 60
		mbc.rtcM += uint8(minutes) + extraMin
		extraHour := mbc.rtcM / 60
		mbc.rtcM %= 60

		// Update hours
		hours := int(elapsed.Hours()) % 24
		mbc.rtcH += uint8(hours) + extraHour
		extraDay := mbc.rtcH / 24
		mbc.rtcH %= 24

		// Update days
		days := int(elapsed.Hours()) / 24
		rtcDays := int(mbc.rtcDL) | (int(mbc.rtcDH&1) << 8)
		totalDays := days + rtcDays + int(extraDay)
		mbc.rtcDL = uint8(totalDays)                      // Set Day low
		util.SetBit(&mbc.rtcDH, 0, uint8(totalDays>>8)&1) // Set Day high
		if totalDays > 0x1FF {
			util.SetBit(&mbc.rtcDH, 7, 1) // Set Overflow
		}
	}
}
