package gnet

import (
	"fmt"
	"net/url"

	//"el/player"

	"time"

	"github.com/gorilla/websocket"
)

//IFIoservice 接口
type IFIoservice interface {
	OnConnected(itv *Interval)
	OnAccept(itv *Interval)
	OnReceive(itv *Interval, data []byte, isText bool)
	OnClose(itv *Interval)
	OnPingpong(itv *Interval, ping int64)
	Listen(port int, ssl bool, httpIf []string) *Interval
	Connect(address string) *websocket.Conn
	OnUpdate(delta float32)
	OnRequset(cmd string, Form *url.Values) []byte
	OnInit()
	//Close(socket int64)
}

//STIoservice 结构
type STIoservice struct {
}

//Connect 监听接口
func (v *STIoservice) Connect(address string) *websocket.Conn {
	dialer := &websocket.Dialer{HandshakeTimeout: 3 * time.Second}
	conn, _, err := dialer.Dial(address, nil)
	if err != nil {
		fmt.Println("Connect error:", err)
		return nil
	}
	//u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}

	Service.ChanConnected <- conn
	return conn
}

func (v *STIoservice) Listen(port int, ssl bool, httpIf []string) *Interval {

	Web.Start(port, ssl, httpIf)
	return nil
}

func (v *STIoservice) OnConnected(itv *Interval) {
}

func (v *STIoservice) OnAccept(itv *Interval) {
}

func (v *STIoservice) OnClose(itv *Interval) {
}

func (v *STIoservice) OnReceive(itv *Interval, data []byte, isText bool) {

}

func (v *STIoservice) OnPingpong(itv *Interval, ping int64) {

}

func (v *STIoservice) OnUpdate(delta float32) {

}

func (v *STIoservice) OnRequset(cmd string, Form *url.Values) []byte {
	return []byte{}
}

func (v *STIoservice) OnInit() {

}

func (v *STIoservice) Close(socket int64) {
	itv := Service.FindInerval(socket)
	if itv != nil {
		itv.valid = false
		itv.Close()
	}
}

func (v *STIoservice) SendItv(itv *Interval, data []byte, isText bool) {

	if isText {
		itv.Send(websocket.TextMessage, data)
	} else {
		itv.Send(websocket.BinaryMessage, data)
	}
}

func (v *STIoservice) Send(socket int64, data []byte, isText bool) {
	itv := Service.FindInerval(socket)
	if itv != nil {
		if isText {
			itv.Send(websocket.TextMessage, data)
		} else {
			itv.Send(websocket.BinaryMessage, data)
		}
	}
}
