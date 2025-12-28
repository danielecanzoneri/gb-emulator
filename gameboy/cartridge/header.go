package cartridge

import (
	"fmt"
	"log"
)

type CGBMode int

const (
	DmgOnly CGBMode = iota
	CgbCompatible
	CgbOnly
)

const (
	cgbFlag       = 0x0143
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

type Header struct {
	// ROMSize = 16 KiB * ROMBanks
	ROMBanks uint

	// RAMSize = 8 KiB * RAMBanks
	RAMBanks uint

	// Bytes 0134-0143
	Title string
	// Licensee code
	Licensee string
	// Destination (00=Japan, 01=Overseas)
	Destination uint8
	// Game version (usually 00)
	GameVersion uint8

	// Byte 0143
	CgbMode CGBMode
}

func parseHeader(data []byte) *Header {
	// Parse title
	Title := parseTitle(data[title : title+titleLen])

	// Parse licensee code
	var LicenseeCode string
	if data[oldLicenseeCode] != 0x33 {
		LicenseeCode = fmt.Sprintf("%02X", data[oldLicenseeCode])
	} else {
		LicenseeCode = string(data[newLicenseeCode : newLicenseeCode+2])
	}

	RAMBanks := computeRAMSize(data[ramSize])

	var cgbMode CGBMode
	switch data[cgbFlag] {
	case 0x80:
		cgbMode = CgbCompatible
	case 0xC0:
		cgbMode = CgbOnly
	default:
		cgbMode = DmgOnly
	}

	return &Header{
		ROMBanks:    computeROMBanks(data[romSize]),
		RAMBanks:    RAMBanks,
		Title:       Title,
		Licensee:    LicenseeCode,
		Destination: data[destinationCode],
		GameVersion: data[gameVersion],
		CgbMode:     cgbMode,
	}
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

func computeRAMSize(v uint8) uint {
	switch v {
	case 0x00:
		return 0
	case 0x02:
		return 1
	case 0x03:
		return 4
	case 0x04:
		return 16
	case 0x05:
		return 8
	default:
		log.Fatalf("Unsupported RAM size: %02X\n", v)
		return 0
	}
}
