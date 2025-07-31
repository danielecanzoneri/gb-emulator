package ui

import (
	"log"
	"net"
)

func (ui *UI) CreateSocket(addr string) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	// Wait for an incoming connection
	go func() {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("[ERROR] Creating master socket: ", err)
			return
		}

		// Important for low latency
		tcpConn := conn.(*net.TCPConn)
		err = tcpConn.SetNoDelay(true)
		if err != nil {
			log.Println("[ERROR] Setting socket no delay: ", err)
		}

		ui.gameBoy.SerialPort.Conn = conn
		ui.gameBoy.SerialPort.Listen()
	}()

	return nil
}

func (ui *UI) ConnectToSocket(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	// Important for low latency
	tcpConn := conn.(*net.TCPConn)
	err = tcpConn.SetNoDelay(true)
	if err != nil {
		log.Println("[ERROR] Setting socket no delay: ", err)
	}

	ui.gameBoy.SerialPort.Conn = conn
	ui.gameBoy.SerialPort.Listen()

	return nil
}
