package gnet

import (
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

type serviceST struct {
	socket       int64
	rwmutex      sync.RWMutex
	ChanReceive  chan *Context
	ChanPing     chan *Interval
	ChanClose    chan *Interval
	intervalList map[int64]*Interval

	serviceHandle IFIoservice

	FPS           int
	chanUpdate    chan float32
	lobbyTimer    <-chan time.Time
	tickerNumber  int64
	fpsTimeNano   int64
	count         int
	ChanConnected chan *websocket.Conn
	IsInit        bool
}

func NewService() *serviceST {
	ptr := &serviceST{}
	if ptr.Init() {
		return ptr
	} else {
		return nil
	}
}

func (v *serviceST) Init() bool {
	v.intervalList = make(map[int64]*Interval)
	v.ChanReceive = make(chan *Context, 20)
	v.ChanClose = make(chan *Interval, 20)

	v.ChanConnected = make(chan *websocket.Conn, 10000)
	v.IsInit = false
	return true
}

func (v *serviceST) SetHandle(handle IFIoservice) {
	v.serviceHandle = handle
	v.serviceHandle.OnInit()
}

//CreateInterval ..
func (v *serviceST) CreateInterval(conn *websocket.Conn) *Interval {
	atomic.AddInt64(&v.socket, 1)
	pinterval := &Interval{}
	pinterval.init(v.socket, conn)
	v.rwmutex.Lock()
	v.intervalList[v.socket] = pinterval
	v.rwmutex.Unlock()

	return pinterval
}

//FindInerval ...
func (v *serviceST) FindInerval(socket int64) *Interval {
	v.rwmutex.RLock()
	defer v.rwmutex.RUnlock()

	pitv, ok := v.intervalList[socket]
	if !ok {
		return nil
	}

	return pitv
}

//RemoveIntervalByID ...
func (v *serviceST) RemoveIntervalByID(socket int64) {
	v.rwmutex.Lock()
	delete(v.intervalList, socket)
	v.rwmutex.Unlock()
}

//RemoveInterval ...
func (v *serviceST) RemoveInterval(itv *Interval) {
	v.rwmutex.Lock()
	delete(v.intervalList, itv.ID)
	v.rwmutex.Unlock()
}

//GetIntervalSize ...
func (v *serviceST) GetIntervalSize() int {
	v.rwmutex.RLock()
	defer v.rwmutex.RUnlock()
	return len(v.intervalList)
}

func (v *serviceST) Run(fps int) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	v.lobbyTimer = time.Tick(time.Second / 30)
	v.tickerNumber = time.Now().UnixNano()
	v.fpsTimeNano = int64(1000000000 / fps)
	v.chanUpdate = make(chan float32)

	go v.mainLoop()
	go v.updateDriver()
}

func (v *serviceST) mainLoop() {
	for {
		select {
		case conn := <-v.ChanConnected:
			v.onWSConnected(conn)
		case conn := <-Web.ChanAccept:
			v.onWSAccept(conn)
		case itv := <-Service.ChanClose:
			v.onWSClose(itv)
		case context := <-Service.ChanReceive:
			v.onWSReceive(context)
		case _time := <-v.chanUpdate:
			v.serviceHandle.OnUpdate(_time)
		case _httpData := <-Web.ChanHTTP:
			v.onWSRequest(_httpData)
		default:
			if !v.IsInit {
				v.serviceHandle.OnInit()
				v.IsInit = true
			}
			time.Sleep(1 * time.Millisecond)
		}
	}
}

func (v *serviceST) updateDriver() {
	for {
		if time.Now().UnixNano()-v.tickerNumber >= v.fpsTimeNano {
			v.tickerNumber += v.fpsTimeNano
			v.chanUpdate <- float32(v.fpsTimeNano) / 1000000000
		}

		time.Sleep(1 * time.Millisecond)
	}
}

func (v *serviceST) onWSConnected(conn *websocket.Conn) {
	itv := v.CreateInterval(conn)
	itv.forConnector = true
	itv.Run()
	v.serviceHandle.OnConnected(itv)
	//glog.Error("onWSConnected")
}

func (v *serviceST) onWSAccept(conn *websocket.Conn) {
	itv := v.CreateInterval(conn)
	itv.Run()
	v.serviceHandle.OnAccept(itv)

}

func (v *serviceST) onWSClose(itv *Interval) {
	itv.valid = false
	v.serviceHandle.OnClose(itv)
	v.RemoveInterval(itv)
}

func (v *serviceST) onWSPongMessage(itv *Interval, data []byte) {
	num := Format.BytesToInt(data)
	//glog.Warning("%d", num)
	v.serviceHandle.OnPingpong(itv, time.Now().UnixNano()/1000000-int64(num))
}

func (v *serviceST) onWSReceive(context *Context) {
	//fmt.Println("recive:", context.data)
	if !context.itv.valid {
		return
	}
	switch context.messageType {
	case websocket.TextMessage:
		v.serviceHandle.OnReceive(context.itv, context.data, false)
	case websocket.BinaryMessage:
		v.serviceHandle.OnReceive(context.itv, context.data, true)
	case websocket.PongMessage:
		v.onWSPongMessage(context.itv, context.data)
		// case websocket.PingMessage:
		// 	println("###########")

	}
}

func (v *serviceST) onWSRequest(data *HTTPData) {
	// fmt.Println("onWSRequest")

	response := v.serviceHandle.OnRequset(data.Req, data.Form)
	go func() {
		data.ChanBack <- response
	}()
}
