package net

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/stnma7e/betuol/common"
	"github.com/stnma7e/betuol/event"
)

type NetworkEvent struct {
	Type  string `json:"eventType"`
	Event string `json:"event"`
}

type NetworkManager struct {
	listener *net.TCPListener
	conns    *common.Vector
}

func MakeNetworkManager(hostToListenTo string, eventlink chan event.Event) *NetworkManager {
	hostAddr, err := net.ResolveTCPAddr("tcp", hostToListenTo)
	if err != nil {
		common.LogErr.Fatal(err)
	}

	listener, err := net.ListenTCP("tcp", hostAddr)
	if err != nil {
		common.LogErr.Fatal(err)
	}

	nm := NetworkManager{
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
			SendEvent(evt, conn)
		}
	}()

	defer func() {
		conn.Close()
		conn = nil
		eventlink = nil
	}()

	for {
		data := []byte{}
		num, err := conn.Read(data)
		if err != nil {
			//common.LogErr.Println(err)
			continue
		}
		if num == 0 {
			continue
		}

		nm.Dispatch(conn, data)
	}

}

func (nm *NetworkManager) Dispatch(addr net.Conn, data []byte) {
	var obj struct {
		Type  string
		Event event.Event
	}
	err := json.Unmarshal(data, &obj)
	if err != nil {
		common.LogErr.Println(err)
	}
	common.LogInfo.Println(obj)
}

func SendEvent(evt event.Event, conn net.Conn) {
	json, err := json.Marshal(NetworkEvent{
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
