# Game Boy Emulator

A feature-rich, cross-platform Game Boy emulator written in Go, with a modern graphical interface and an integrated graphical debugger.

## Features

- **Accurate CPU Emulation**: Implements the full Game Boy Z80-like CPU, with extensive opcode tests.
- **PPU (Graphics) Emulation**: Renders original Game Boy graphics with accurate timing and palette.
- **APU (Audio) Emulation**: Emulates all four Game Boy sound channels.
- **Memory and MBC Support**: Supports various Memory Bank Controllers (MBC0, MBC1, MBC2, MBC3 with RTC, MBC5).
- **Joypad Input**: Keyboard-mapped controls for all Game Boy buttons and D-Pad.
- **Save States**: Automatically loads and saves `.sav` files alongside ROMs.
- **Debugger**: Integrated graphical debugger with disassembly, memory viewer, register viewer, breakpoints, and step/continue/reset controls.
- **Cross-platform GUI**: Built with [Ebiten](https://ebiten.org/) for the emulator and [Fyne](https://fyne.io/) for the debugger.

## Project Structure

```
.
├── emulator/         # Main emulator code (CPU, PPU, APU, memory, UI)
│   ├── cmd/          # Entry point for running the emulator
│   ├── internal/     # Emulator internals (gameboy, server)
│   └── ui/           # Ebiten-based GUI and input handling
├── debugger/         # Standalone graphical debugger (Fyne-based)
│   ├── cmd/          # Entry point for running the debugger
│   ├── internal/     # Debugger client logic
│   └── ui/           # Debugger GUI components
├── pkg/              # Shared protocol, debug, and utility packages
└── README.md         # This file
```

## Getting Started

### Prerequisites

- **Go 1.24+** (see `go.mod` for version)
- [Ebiten](https://ebiten.org/) and [Fyne](https://fyne.io/) dependencies are managed via Go modules.

### Building

Clone the repository and build the emulator and debugger:

```sh
git clone https://github.com/danielecanzoneri/gb-emulator.git
cd gb-emulator
cd emulator
go build -o gbemu ./cmd
cd ../debugger
go build -o gbdebugger ./cmd
```

### Running

#### Emulator

From the `emulator` directory:

```sh
go run ./cmd
```
or, if built:
```sh
./gbemu
```

- On launch, you will be prompted to select a Game Boy ROM file (`.gb` or `.gbc`).
- Controls:
  - **D-Pad**: Arrow keys
  - **A**: S
  - **B**: A
  - **Start**: X
  - **Select**: Z
  - **Ctrl+L**: Load a new game
  - **1-4**: Toggle audio channels
  - **Esc**: Launch the debugger

#### Debugger

The debugger can be launched automatically from the emulator (press `Esc`), or manually:

From the `debugger` directory:

```sh
go run ./cmd
```
or, if built:
```sh
./gbdebugger
```

- Connects to the emulator via WebSocket on port 8080.
- Features:
  - Disassembly view with breakpoints (click to toggle)
  - Memory viewer
  - Register and interrupt viewer
  - Step (`F3`), Continue (`F8`), and Reset controls

## Resources

- [Pandocs](https://gbdev.io/pandocs/OAM.html)
- [Gameboy Development Wiki](https://gbdev.gg8.se/wiki/articles/Main_Page) ([sound hardware](https://gbdev.gg8.se/wiki/articles/Gameboy_sound_hardware))
- [Data Crystal](https://datacrystal.tcrf.net/wiki/Data_Crystal) (MBC testing)
- [Opcode table](https://gbdev.io/gb-opcodes/optables/) and [opcode reference](https://rgbds.gbdev.io/docs/v0.9.2/gbz80.7)
- [GBops](https://izik1.github.io/gbops/) (opcode timing)
- [Blargg](https://github.com/retrio/gb-test-roms), [Gekkio](https://github.com/Gekkio/mooneye-test-suite), [DMG acid](https://github.com/mattcurrie/dmg-acid2), [MBC3 RTC test](https://github.com/aaaaaa123456789/rtc3test) (test ROMs)
- [This reddit post](https://www.reddit.com/r/EmuDev/comments/59pawp/gb_mode3_sprite_timing/) for fixing PPU timing with sprites

## TODO

- Redesign the debugger using ebitengine/debugui instead of Fyne, to allow for a fully integrated debugger without the need for WebSocket communication.
- Implement a feature to speed up emulation (fast-forward).
- Add real-time save states, allowing users to save and load game state instantly during gameplay.
- Expand support for additional cartridge types and MBC variants.