package gnet

import (
	"fmt"
	"net/http"
	"time"
)

type serverHandler struct {
	time string
	req  string
}

func NewServerHandle(req string) *serverHandler {
	handle := new(serverHandler)
	handle.time = time.Stamp
	handle.req = req

	return handle
}

func (th *serverHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("回应http:失败:", err)
		}
	}()

	//glog.Print("http request:", th.time, th.req)
	r.ParseForm()

	//	fmt.Println(r)
	waitChan := make(chan []byte)
	Web.ChanHTTP <- &HTTPData{Req: th.req, Form: &r.Form, ChanBack: waitChan}

	data := <-waitChan

	if !r.Close {
		w.Write(data)
	}

}
