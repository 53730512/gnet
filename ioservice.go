package gnet

import (
	"fmt"
	"net/url"

	//"el/player"

	"time"

	ghttp "gitee.com/liyp_admin/gnet/ghttp"
	"github.com/gorilla/websocket"
)

//IFIoservice 接口
type IFIoservice interface {
	OnConnected(itv *Interval)
	OnConsoleCMD(cmd string)
	OnAccept(itv *Interval)
	OnReceive(itv *Interval, data []byte, isText bool)
	OnClose(itv *Interval)
	OnPingpong(itv *Interval, ping int64)
	Listen(port int, ssl bool, httpIf []string) *Interval
	Connect(address string) *websocket.Conn
	OnUpdate(delta float32)
	OnLoop() bool
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

	select {
	case Service.ChanConnected <- conn:
	default:
		return nil
	}

	return conn
}

func (v *STIoservice) Listen(port int, ssl bool, httpIf []string) *Interval {

	ghttp.Start(port, ssl, httpIf)
	Log.Success("listen at:%d", port)
	listenOk <- true
	return nil
}

func (v *STIoservice) OnConsoleCMD(cmd string) {
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

func (v *STIoservice) OnLoop() bool {

	return false
}

func (v *STIoservice) OnRequset(cmd string, Form *url.Values) []byte {
	return []byte{}
}

func (v *STIoservice) OnInit() {

}

func (v *STIoservice) Close(socket int64) bool {
	itv := Service.FindInerval(socket)
	if itv != nil {
		itv.valid = false
		return itv.Close()
	}

	return false
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
