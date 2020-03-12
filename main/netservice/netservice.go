package netservice

import (
	"net/url"
	"time"

	"gitee.com/liyp/gnet"
)

type IOserver struct {
	gnet.STIoservice
}

func (v *IOserver) OnInit() {
}

//OnAccept ...
func (v *IOserver) OnAccept(itv *gnet.Interval) {
}

//OnClose  ...
func (v *IOserver) OnClose(itv *gnet.Interval) {

}

//OnReceive ...
func (v *IOserver) OnReceive(itv *gnet.Interval, data []byte, isText bool) {
}

//OnPingpong ...
func (v *IOserver) OnPingpong(itv *gnet.Interval, tm int64) {
	println("pingpong:", tm)
}

func (v *IOserver) OnConsoleCMD(cmd string) {
	switch cmd {
	default:
		gnet.Error("无效命令:%s", cmd)
	}
}

//OnUpdate ...
func (v *IOserver) OnUpdate(delta float32) {
}

func (v *IOserver) OnLoop() bool {

	return true
}

func (v *IOserver) OnRequset(cmd string, Form *url.Values) []byte {
	var tt []byte
	tt = make([]byte, 1)
	return tt
}

var (
	consoleUpdateTime float32
	curSecond         int64
)

//Start ...
func init() {
	consoleUpdateTime = 0
	curSecond = time.Now().Unix()
}

func GetTime() int64 {
	return curSecond
}
