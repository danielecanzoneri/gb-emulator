package cartridge

import (
	"errors"
	"fmt"
)

const (
	cartridgeType = 0x0147
	romSize       = 0x0148
	ramSize       = 0x0149

	title    = 0x134
	titleLen = 16

	destinationCode = 0x014A
	oldLicenseeCode = 0x014B
	newLicenseeCode = 0x0144
	gameVersion     = 0x014C
)

const (
	NoMBC = iota
	MBC1
)

type Header struct {
	// Memory Bank Controller (MBC)
	MBC int

	// ROMSize = 16 KiB * ROMBanks
	ROMBanks uint

	// RAMSize = 8 KiB * RAMBanks
	RAMBanks uint

	// Bits 0134-0143
	Title string
	// Licensee code
	Licensee string
	// Destination (00=Japan, 01=Overseas)
	Destination uint8
	// Game version (usually 00)
	GameVersion uint8
}

func (h *Header) String() string {
	return fmt.Sprintf(
		"MBC: %x, ROM banks = %d (%d KiB), RAM banks = %d (%d KiB) \nGame title: %v, version %x, licensee: %v, destination: %x",
		h.MBC, h.ROMBanks, h.ROMBanks*16, h.RAMBanks, h.RAMBanks*8,
		h.Title, h.GameVersion, h.Licensee, h.Destination,
	)
}

func parseHeader(data []byte) (*Header, error) {
	// Parse title
	Title := parseTitle(data[title : title+titleLen])

	// Parse licensee code
	var LicenseeCode string
	if data[oldLicenseeCode] != 0x33 {
		LicenseeCode = fmt.Sprintf("%02X", data[oldLicenseeCode])
	} else {
		LicenseeCode = string(data[newLicenseeCode : newLicenseeCode+2])
	}

	RAMBanks, err := computeRAMSize(data[ramSize])
	if err != nil {
		return nil, err
	}

	var mbc int
	switch data[cartridgeType] {
	case 0:
		mbc = NoMBC
	case 1:
		fallthrough
	case 2:
		fallthrough
	case 3:
		mbc = MBC1
	default:
		return nil, fmt.Errorf("cartridge type %02X not supported", data[cartridgeType])
	}

	return &Header{
		MBC:         mbc,
		ROMBanks:    computeROMBanks(data[romSize]),
		RAMBanks:    RAMBanks,
		Title:       Title,
		Licensee:    LicenseeCode,
		Destination: data[destinationCode],
		GameVersion: data[gameVersion],
	}, nil
}

func parseTitle(titleData []byte) string {
	first0 := len(titleData)
	for i, b := range titleData {
		if b == 0 {
			first0 = i
			break
		}
	}
	return string(titleData[:first0])
}

func computeROMBanks(v uint8) uint {
	return 1 << (v + 1)
}

func computeRAMSize(v uint8) (uint, error) {
	switch v {
	case 0x00:
		return 0, nil
	case 0x02:
		return 1, nil
	case 0x03:
		return 4, nil
	case 0x04:
		return 16, nil
	case 0x05:
		return 8, nil
	default:
		return 0, errors.New("unsupported RAM size")
	}
}
