package net

import (
	"encoding/binary"
	"net"
	"time"

	"smig/common"
	"smig/event"
)

type NetworkManager struct {
	conns    *common.Vector
	sendlink chan interface{}
}

func MakeNetworkManager() *NetworkManager {
	nm := NetworkManager{
		common.MakeVector(),
		make(chan interface{}),
	}

	return &nm
}

func (nm *NetworkManager) Tick() {
	conns := nm.conns.Array()
	for i := range conns {
		for j := true; j; {
			select {
			case j := <-nm.sendlink:
				err := binary.Write(conns[i].(*net.UDPConn), binary.BigEndian, j)
				if err != nil {
					common.LogErr.Println(err)
				}
			default:
				j = false
			}
		}
	}
}

func (nm *NetworkManager) Send(evt event.Event) {
	go func() {
		nm.sendlink <- evt
	}()
}

func (nm *NetworkManager) Connect(address string) error {
	conn, err := net.Dial("udp", address)
	if err != nil {
		return err
	}
	nm.conns.Push_back(conn, 1, 1)
	return nil
}

func (nm *NetworkManager) SendBytes(b []byte) {
	go func() {
		nm.sendlink <- b
	}()
}

func (nm *NetworkManager) RecieveBytes(size int, timeout float32) ([]byte, error) {
	conns := nm.conns.Array()
	b := make([]byte, size)
	for i := range conns {
		conn := conns[i].(*net.UDPConn)
		conn.SetReadDeadline(time.Now().Add(time.Duration(timeout) * time.Millisecond))
		_, _, err := conn.ReadFromUDP(b)
		if err != nil {
			return []byte{}, err
		}
	}
	return b, nil
}
