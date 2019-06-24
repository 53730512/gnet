package el

import (
	"fmt"
	"net/url"
	"runtime"

	//"el/player"

	"time"

	"github.com/gorilla/websocket"
)

//FPS ...
var FPS int

//ChanUpdate ...
var chanUpdate chan float32

//lobbyTimer ...
var lobbyTimer <-chan time.Time

//tickerNumber...
var tickerNumber int64

var fpsTimeNano int64

var count int

//ChanConnected ...
var ChanConnected chan *websocket.Conn

var IsInit bool

//Init ...
func Init() bool {
	ChanConnected = make(chan *websocket.Conn, 10000)
	IsInit = false
	InitCommon()
	return InitLog()
}

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

	ChanConnected <- conn
	return conn
}

//Listen 监听接口
func (v *STIoservice) Listen(port int, ssl bool, httpIf []string) *Interval {

	StartHttp(port, ssl, httpIf)
	return nil
}

//onOnConnected ...
func (v *STIoservice) OnConnected(itv *Interval) {
}

//OnAccept ...
func (v *STIoservice) OnAccept(itv *Interval) {
}

//OnClose  ...
func (v *STIoservice) OnClose(itv *Interval) {
}

//OnReceive ...
func (v *STIoservice) OnReceive(itv *Interval, data []byte, isText bool) {

}

//OnPingpong ...
func (v *STIoservice) OnPingpong(itv *Interval, ping int64) {

}

//OnUpdate ...
func (v *STIoservice) OnUpdate(delta float32) {

}

//OnRequset ...
func (v *STIoservice) OnRequset(cmd string, Form *url.Values) []byte {
	return []byte{}
}

func (v *STIoservice) OnInit() {

}

//Close  ...
func (v *STIoservice) Close(socket int64) {
	itv := FindInerval(socket)
	if itv != nil {
		itv.valid = false
		itv.Close()
	}
}

//SendItv  ...
func (v *STIoservice) SendItv(itv *Interval, data []byte, isText bool) {

	if isText {
		itv.Send(websocket.TextMessage, data)
	} else {
		itv.Send(websocket.BinaryMessage, data)
	}
}

//Send  ...
func (v *STIoservice) Send(socket int64, data []byte, isText bool) {
	itv := FindInerval(socket)
	if itv != nil {
		if isText {
			itv.Send(websocket.TextMessage, data)
		} else {
			itv.Send(websocket.BinaryMessage, data)
		}
	}
}

var serviceHandle IFIoservice

//Create ...
func Create(handle IFIoservice) {
	serviceHandle = handle
}

//Run ...
func Run(fps int) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	lobbyTimer = time.Tick(time.Second / 30)
	tickerNumber = time.Now().UnixNano()
	fpsTimeNano = int64(1000000000 / fps)
	chanUpdate = make(chan float32)

	go mainLoop()
	go updateDriver()
}

func mainLoop() {
	for {
		select {
		case conn := <-ChanConnected:
			onWSConnected(conn)
		case conn := <-ChanAccept:
			onWSAccept(conn)
		case itv := <-ChanClose:
			onWSClose(itv)
		case context := <-ChanReceive:
			onWSReceive(context)
		case _time := <-chanUpdate:
			serviceHandle.OnUpdate(_time)
		case _httpData := <-ChanHTTP:
			onWSRequest(_httpData)
		default:
			if !IsInit {
				OnInit()
			}
			time.Sleep(1 * time.Millisecond)
		}
	}
}

func updateDriver() {
	for {
		if time.Now().UnixNano()-tickerNumber >= fpsTimeNano {
			tickerNumber += fpsTimeNano
			chanUpdate <- float32(fpsTimeNano) / 1000000000
		}

		time.Sleep(1 * time.Millisecond)
	}
}

func onWSConnected(conn *websocket.Conn) {
	itv := CreateInterval(conn)
	itv.forConnector = true
	itv.Run()
	serviceHandle.OnConnected(itv)
	//glog.Error("onWSConnected")
}

func onWSAccept(conn *websocket.Conn) {
	itv := CreateInterval(conn)
	itv.Run()
	serviceHandle.OnAccept(itv)

}

func OnInit() {
	IsInit = true
	serviceHandle.OnInit()
}

func onWSClose(itv *Interval) {
	itv.valid = false
	serviceHandle.OnClose(itv)
	RemoveInterval(itv)
}

func onWSPongMessage(itv *Interval, data []byte) {
	num := BytesToInt(data)
	//glog.Warning("%d", num)
	serviceHandle.OnPingpong(itv, time.Now().UnixNano()/1000000-int64(num))
}

func onWSReceive(context *Context) {
	//fmt.Println("recive:", context.data)
	if !context.itv.valid {
		return
	}
	switch context.messageType {
	case websocket.TextMessage:
		serviceHandle.OnReceive(context.itv, context.data, false)
	case websocket.BinaryMessage:
		serviceHandle.OnReceive(context.itv, context.data, true)
	case websocket.PongMessage:
		onWSPongMessage(context.itv, context.data)
		// case websocket.PingMessage:
		// 	println("###########")

	}
}

func onWSRequest(data *HTTPData) {
	// fmt.Println("onWSRequest")

	response := serviceHandle.OnRequset(data.Req, data.Form)
	go func() {
		data.ChanBack <- response
	}()
}
