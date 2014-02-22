package net

import (
	"encoding/json"
	"io"
	"net"

	"github.com/stnma7e/betuol/common"
	"github.com/stnma7e/betuol/event"
	"github.com/stnma7e/betuol/math"
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
	for {
		conn, err := nm.listener.AcceptTCP()
		if err != nil {
			common.LogErr.Println(err)
			return
		}

		eventlink := make(chan event.Event)
		nm.conns.Insert(eventlink)
		go nm.TCPTick(conn, eventlink)
	}
}

func (nm *NetworkManager) TCPTick(conn *net.TCPConn, eventlink chan event.Event) {
	go func() {
		for evt := range eventlink {
			common.LogInfo.Println("got a new event to send over network")
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
	common.LogInfo.Println("received event from network")
	evt := event.NetworkEvent{}
	err := json.Unmarshal(data, &evt)
	if err != nil {
		common.LogErr.Println(err.Error()+
			"\n\tevent (string): "+string(data)+
			"\n\tevent (bytes): ", data)
		return
	}

	common.LogInfo.Println("received event data:", evt)
	nm.em.Send(ConvertJSONToEvent(evt.Type, evt.Event, conn))
}

func ConvertJSONToEvent(eventType string, eventData map[string]interface{}, conn net.Conn) event.Event {
	switch eventType {
	case "requestCharacterCreation":
		loc := math.Vec3{}
		interfaceLoc := eventData["Location"].([]interface{})
		for i := range loc {
			loc[i] = float32(interfaceLoc[i].(float64))
		}
		return &event.RequestCharacterCreationEvent{
			Type:          eventData["Type"].(string),
			Location:      loc,
			RequestOrigin: conn,
		}
	}

	common.LogWarn.Println("no eventType")
	return nil
}

func SendEvent(evt event.Event, conn net.Conn) {
	evtJSONbs, err := json.Marshal(evt)
	if err != nil {
		common.LogErr.Print(err)
		return
	}

	evtJSONis := make([]int, len(evtJSONbs))
	for i := range evtJSONbs {
		evtJSONis[i] = int(evtJSONbs[i])
	}

	obj := struct {
		Type  string
		Event []int
	}{
		evt.GetType(),
		evtJSONis,
	}

	json, err := json.Marshal(obj)
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

	common.LogInfo.Printf("sent %s event: %v", evt.GetType(), evt)
}
