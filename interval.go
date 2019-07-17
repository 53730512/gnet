package gnet

import (
	"fmt"
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
	chanPing        chan string
	iclosed         int32
}

func (v *Interval) init(ID int64, conn *websocket.Conn) {
	v.ID = ID
	v.forConnector = false
	v.recivePackCount = 0
	v.localTiker = time.NewTicker(10 * time.Second)
	v.chanSend = make(chan *Context, 500)
	v.chanClose = make(chan bool, 2)
	v.chanPing = make(chan string, 2)
	v.wsocket = conn
	v.wsocket.SetPingHandler(v.pingCallback)
	v.wsocket.SetPongHandler(v.pongCallback)
	v.iclosed = 0
}

func (v *Interval) GetConn() *websocket.Conn {
	return v.wsocket
}

func (v *Interval) pingCallback(appData string) error {
	//v.wsocket.WriteMessage(websocket.PongMessage, []byte(appData))
	v.chanPing <- appData
	return nil
}

func (v *Interval) pongCallback(appData string) error {
	context := new(Context)
	context.itv = v
	context.messageType = websocket.PongMessage
	context.data = []byte(appData)
	// glog.Success("%d", len(context.data))
	Service.ChanReceive <- context
	return nil
}

//Send ...
func (v *Interval) Send(_type int, data []byte) {
	if v.iclosed > 0 {
		return
	}
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
			Log.Print("关闭interval失败:%s, id:%d", err, v.ID)
		}
	}()
	// fmt.Println("try to close", v.ID)

	if v.iclosed > 0 {
		fmt.Println("尝试关闭一个已经关闭的对象")
		return false
	}

	v.chanClose <- true
	atomic.AddInt32(&v.iclosed, 1)
	return true
}

func (v *Interval) reciveRuntime() {
	defer func() {
		if err := recover(); err != nil {
			Log.Error("reciveRuntime:%s", err)
			v.Close()
		}

	}()

	for {
		messageType, message, err := v.wsocket.ReadMessage()

		if v.iclosed == 1 {
			break
		}
		// if messageType == -1 {
		// 	continue
		// }
		//fmt.Println("recive:", message, "type:", messageType)
		if err != nil || len(message) == 0 {
			// fmt.Println("close:", v.ID, err)
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
		Service.ChanReceive <- context

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

func (v *Interval) doChanClose() {
	// Log.Print("%s:%d", "dochanclose", v.ID)
	v.wsocket.Close()
	v.wsocket = nil
	close(v.chanClose)
	Service.ChanClose <- v
}
func (v *Interval) update() {
	defer func() {
		//fmt.Println("defer on Read pack")
		if err := recover(); err != nil {
			Log.Error("update:%s", err)
			v.Close()
		}
	}()

	for {

		select {
		case context, ok := <-v.chanSend:
			if !ok {
				if atomic.LoadInt32(&v.iclosed) == 1 {
					// Log.Success("chanClose33:%d", v.ID)
					v.doChanClose()
				} else {
					// Log.Success("chanClose333:%d", v.ID)
					v.Close()
				}
				break
			}
			v.wsocket.WriteMessage(context.messageType, context.data)
			//Log.Success("send len:%d", len(v.chanSend))
			if atomic.LoadInt32(&v.iclosed) > 0 && len(v.chanSend) == 0 {
				// Log.Success("chanClose3:%d", v.ID)
				v.doChanClose()
				return
			}
		case <-v.localTiker.C:
			if v.forConnector {
				tm := time.Now().UnixNano() / 1000000

				//strTm := strconv.FormatInt(tm, 10)
				bytes := Format.IntToBytes(int(tm))
				// glog.Error("%d", len(bytes))
				v.wsocket.WriteMessage(websocket.PingMessage, bytes)
				//fmt.Println("send ping")
			}
		case _, ok := <-v.chanClose:
			if ok {
				// Log.Success("chanClose1:%d", v.ID)
				if len(v.chanSend) > 0 {
					//v.closed = true
					// Log.Success("chanClose111:%d", v.ID)
					atomic.AddInt32(&v.iclosed, 1)
					break
				}
				// Log.Success("chanClose2:%d", v.ID)
				v.doChanClose()
				return
			}
		case data, ok := <-v.chanPing:
			if ok {
				v.wsocket.WriteMessage(websocket.PongMessage, []byte(data))
			}
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
