package ui

import (
	"github.com/danielecanzoneri/lucky-boy/gameboy/serial"
	"log"
	"net"
)

// Listen on the specified port for incoming connections
func (ui *UI) Listen(socketPort string) {
	ln, err := net.Listen("tcp", "localhost:"+socketPort)
	if err != nil {
		log.Println("[ERROR] Listening: ", err)
		return
	}

	go func() {
		ui.GameBoy.SerialPort.State = serial.Connecting

		// Wait for an incoming connection
		conn, err := ln.Accept()
		if err != nil {
			log.Println("[ERROR] Accepting incoming connection", err)
			ui.GameBoy.SerialPort.State = serial.Disconnected
			return
		}

		// Important for low latency
		tcpConn := conn.(*net.TCPConn)
		err = tcpConn.SetNoDelay(true)
		if err != nil {
			log.Println("[ERROR] Setting socket no delay: ", err)
		}

		ui.GameBoy.SerialPort.Conn = conn
		go ui.GameBoy.SerialPort.Listen()
		ui.GameBoy.SerialPort.State = serial.Connected
	}()
}

// Connect on the specified port to another socket
func (ui *UI) Connect(socketPort string) {
	conn, err := net.Dial("tcp", "localhost:"+socketPort)
	if err != nil {
		log.Println("[ERROR] Connecting to socket: ", err)
	}

	// Important for low latency
	tcpConn := conn.(*net.TCPConn)
	err = tcpConn.SetNoDelay(true)
	if err != nil {
		log.Println("[ERROR] Setting socket no delay: ", err)
	}

	ui.GameBoy.SerialPort.Conn = conn
	go ui.GameBoy.SerialPort.Listen()
	ui.GameBoy.SerialPort.State = serial.Connected
}
