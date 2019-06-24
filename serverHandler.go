package gnet

import (
	"fmt"
	"net/http"
	"time"
)

//NewServerHandle ...
func NewServerHandle(req string) *ServerHandler {
	handle := new(ServerHandler)
	handle.time = time.Stamp
	handle.req = req

	return handle
}

//ServerHandler ...
type ServerHandler struct {
	time string
	req  string
}

//ServeHTTP ...
func (th *ServerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("回应http:失败:", err)
		}
	}()

	//glog.Print("http request:", th.time, th.req)
	r.ParseForm()

	//	fmt.Println(r)
	waitChan := make(chan []byte)
	ChanHTTP <- &HTTPData{Req: th.req, Form: &r.Form, ChanBack: waitChan}

	data := <-waitChan

	if !r.Close {
		w.Write(data)
	}

}
