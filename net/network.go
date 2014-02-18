package net

import (
	"encoding/json"
	"fmt"
	"io"
	"net"

	"github.com/stnma7e/betuol/common"
	"github.com/stnma7e/betuol/event"
)

type NetworkManager struct {
	em *event.EventManager

	listener *net.TCPListener
	conns    *common.Vector
}

func MakeNetworkManager(em *event.EventManager, hostToListenTo string, eventlink chan event.Event) *NetworkManager {
	hostAddr, err := net.ResolveTCPAddr("tcp", hostToListenTo)
	if err != nil {
		common.LogErr.Fatal(err)
	}

	listener, err := net.ListenTCP("tcp", hostAddr)
	if err != nil {
		common.LogErr.Fatal(err)
	}

	nm := NetworkManager{
		em,
		listener,
		common.MakeVector(),
	}

	go func() {
		for evt := range eventlink {
			allConnChans := nm.conns.Array()
			for i := range allConnChans {
				allConnChans[i].(chan event.Event) <- evt
			}
		}
	}()

	return &nm
}

func (nm *NetworkManager) Tick() {
	conn, err := nm.listener.AcceptTCP()
	if err != nil {
		common.LogErr.Println(err)
		return
	}

	eventlink := make(chan event.Event)
	nm.conns.Insert(eventlink)
	go nm.TCPTick(conn, eventlink)
}

func (nm *NetworkManager) TCPTick(conn *net.TCPConn, eventlink chan event.Event) {
	go func() {
		for evt := range eventlink {
			common.LogInfo.Println("got a new event")
			SendEvent(evt, conn)
		}
	}()

	defer func() {
		conn.Close()
		conn = nil
		eventlink = nil
	}()

	for {
		data := make([]byte, 256)
		num, err := conn.Read(data)
		if err != nil {
			if err != io.EOF {
				common.LogErr.Println(err)
			}
			continue
		}
		if num == 0 {
			continue
		}

		nm.Dispatch(conn, data[:num])
	}

}

func (nm *NetworkManager) Dispatch(conn net.Conn, data []byte) {
	evt := event.NetworkEvent{}
	err := json.Unmarshal(data, &evt)
	if err != nil {
		common.LogErr.Println(err)
		common.LogErr.Println(string(data)+"\n", data)
		return
	}
	common.LogInfo.Println("received network event: ", evt)
	nm.em.Send(evt)
}

func SendEvent(evt event.Event, conn net.Conn) {
	json, err := json.Marshal(event.NetworkEvent{
		Type:  evt.GetType(),
		Event: fmt.Sprint(evt),
	})
	if err != nil {
		common.LogErr.Println(err)
	}

	n, err := conn.Write(json)
	if err != nil {
		common.LogErr.Println(err)
		if n == 0 {
			common.LogErr.Println("no data sent; attempted payload of", json)
		}
		return
	}
	if n == 0 {
		common.LogErr.Println("no data sent; attempted payload of", json)
		return
	}
}
