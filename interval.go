package el

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

//Context ...
type Context struct {
	itv         *Interval
	messageType int
	data        []byte
}

var socket int64
var rwmutex sync.RWMutex

//ChanReceive ...
var ChanReceive chan *Context

//ChanPing ...
var ChanPing chan *Interval

//ChanClose ...
var ChanClose chan *Interval

//Interval ...
type Interval struct {
	ID              int64
	wsocket         *websocket.Conn
	forConnector    bool
	valid           bool
	recivePackCount int64
	userdata        interface{}
	localTiker      *time.Ticker
	chanSend        chan *Context
	chanClose       chan bool
	closed          bool
}

func _init() {
	intervalList = make(map[int64]*Interval)
	ChanReceive = make(chan *Context, 20)
	ChanClose = make(chan *Interval, 20)
}

func (v *Interval) init(ID int64, conn *websocket.Conn) {
	v.ID = ID
	v.forConnector = false
	v.recivePackCount = 0
	v.localTiker = time.NewTicker(10 * time.Second)
	v.chanSend = make(chan *Context, 5)
	v.chanClose = make(chan bool)
	v.wsocket = conn
	v.wsocket.SetPingHandler(v.pingCallback)
	v.wsocket.SetPongHandler(v.pongCallback)
	v.closed = false
}

func (v *Interval) GetConn() *websocket.Conn {
	return v.wsocket
}

func (v *Interval) pingCallback(appData string) error {
	v.wsocket.WriteMessage(websocket.PongMessage, []byte(appData))
	return nil
}

func (v *Interval) pongCallback(appData string) error {

	context := new(Context)
	context.itv = v
	context.messageType = websocket.PongMessage
	context.data = []byte(appData)
	// glog.Success("%d", len(context.data))
	ChanReceive <- context
	return nil
}

//Send ...
func (v *Interval) Send(_type int, data []byte) {
	context := new(Context)
	context.itv = v
	context.messageType = _type
	context.data = data
	v.chanSend <- context
}

//Run ...
func (v *Interval) Run() {
	v.valid = true
	go v.update()
	go v.reciveRuntime()
	// if !v.forConnector {
	// }
}

//Close ...
func (v *Interval) Close() bool {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("关闭interval失败:", err)
		}
	}()
	if v.closed {
		return false
	}

	v.chanClose <- true
	v.closed = true
	return true
}

func (v *Interval) reciveRuntime() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("reciveRuntime:", err)
			v.Close()
		}

	}()

	for {
		messageType, message, err := v.wsocket.ReadMessage()

		// if messageType == -1 {
		// 	continue
		// }
		//fmt.Println("recive:", message, "type:", messageType)
		if err != nil || len(message) == 0 {
			fmt.Println("close:", v.ID, err)
			v.Close()
			break
		}

		if messageType == websocket.CloseMessage {
			v.Close()
			break
		}

		v.recivePackCount++
		context := new(Context)
		context.itv = v
		context.messageType = messageType
		context.data = message
		ChanReceive <- context

		//fmt.Println("recived")
	}
}

//SetUserData ...
func (v *Interval) SetUserData(data interface{}) {
	v.userdata = data
}

//GetUserData ...
func (v *Interval) GetUserData() interface{} {
	return v.userdata
}

func (v *Interval) update() {
	if v.closed {
		return
	}

	defer func() {
		//fmt.Println("defer on Read pack")
		if err := recover(); err != nil {
			fmt.Println("update:", err)
			v.Close()
		}
	}()

	for {

		select {
		case context, ok := <-v.chanSend:
			if !ok {
				v.Close()
				return
			}
			v.wsocket.WriteMessage(context.messageType, context.data)
		case <-v.localTiker.C:
			if v.forConnector {
				tm := time.Now().UnixNano() / 1000000

				//strTm := strconv.FormatInt(tm, 10)
				bytes := IntToBytes(int(tm))
				// glog.Error("%d", len(bytes))
				v.wsocket.WriteMessage(websocket.PingMessage, bytes)
				//fmt.Println("send ping")
			}
		case _, ok := <-v.chanClose:
			if ok {
				v.wsocket.Close()
				v.wsocket = nil
				close(v.chanClose)
				ChanClose <- v
			}
			return
		default:
			time.Sleep(1 * time.Millisecond)
			break
			// default:
			// 	{
			// 		v.reciveRuntime()
			// 	}
			// 	break
		}

		//fmt.Println("update")

	}
}

var intervalList map[int64]*Interval

//CreateInterval ..
func CreateInterval(conn *websocket.Conn) *Interval {
	if intervalList == nil {
		_init()
	}
	atomic.AddInt64(&socket, 1)

	pinterval := &Interval{}
	pinterval.init(socket, conn)
	rwmutex.Lock()
	intervalList[socket] = pinterval
	rwmutex.Unlock()
	return pinterval
}

//FindInerval ...
func FindInerval(socket int64) *Interval {
	rwmutex.RLock()
	defer rwmutex.RUnlock()

	pitv, ok := intervalList[socket]
	if !ok {
		return nil
	}

	return pitv
}

//RemoveIntervalByID ...
func RemoveIntervalByID(socket int64) {
	rwmutex.Lock()
	delete(intervalList, socket)
	rwmutex.Unlock()
}

//RemoveInterval ...
func RemoveInterval(itv *Interval) {
	rwmutex.Lock()
	delete(intervalList, itv.ID)
	rwmutex.Unlock()
}

//GetIntervalSize ...
func GetIntervalSize() int {
	rwmutex.RLock()
	defer rwmutex.RUnlock()
	return len(intervalList)
}
