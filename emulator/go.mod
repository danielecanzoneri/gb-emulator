module github.com/danielecanzoneri/gb-emulator/emulator

go 1.24.2

require (
	github.com/danielecanzoneri/gb-emulator/pkg v0.0.0
	github.com/ebitengine/oto/v3 v3.3.3
	github.com/gorilla/websocket v1.5.3
	github.com/hajimehoshi/ebiten/v2 v2.8.8
	github.com/sqweek/dialog v0.0.0-20240226140203-065105509627
)

require (
	github.com/TheTitanrain/w32 v0.0.0-20180517000239-4f5cfb03fabf // indirect
	github.com/ebitengine/gomobile v0.0.0-20240911145611-4856209ac325 // indirect
	github.com/ebitengine/hideconsole v1.0.0 // indirect
	github.com/ebitengine/purego v0.8.0 // indirect
	github.com/jezek/xgb v1.1.1 // indirect
	golang.org/x/sync v0.12.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
)

replace github.com/danielecanzoneri/gb-emulator/pkg => ../pkg
