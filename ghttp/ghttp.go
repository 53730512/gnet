package ghttp

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

//HttpData ...
type HttpData struct {
	Req      string
	Form     *url.Values
	ChanBack chan []byte
}

type WebST struct {
	ChanAccept chan *websocket.Conn
	ChanHTTP   chan *HttpData
	httpserver *http.Server
}

func NewWeb() *webST {
	ptr := &webST{}
	if ptr.init() {
		return ptr
	} else {
		return nil
	}
}

func (v *WebST) init() bool {
	v.ChanAccept = make(chan *websocket.Conn, 100)
	v.ChanHTTP = make(chan *HttpData, 100)
	v.httpserver = &http.Server{
		ReadTimeout: 3 * time.Second,
		//WriteTimeout: 5 * time.Second,
		//IdleTimeout: 3 * time.Second,
	}

	return true
}

//Close ...
func (v *WebST) Close() {

}

//Start ...
func (v *WebST) Start(port int, ssl bool, httpIf []string) bool {
	v.httpserver.Addr = fmt.Sprintf(":%d", port)

	v.RegisterHandles(httpIf)
	go func() {
		if ssl {
			//Log.Print("start https service")
			err := v.httpserver.ListenAndServeTLS(File.GetFilePath("assets/keys/ssl.crt"), File.GetFilePath("assets/keys/ssl.key"))
			fmt.Println(err)
			if err != nil {
				panic(err)
			}

		} else {
			//Log.Print("start http service")
			err := v.httpserver.ListenAndServe()
			fmt.Println(err)
			if err != nil {
				panic(err)
			}
		}

	}()

	return true
}

//RegisterHandles ...
func (v *WebST) RegisterHandles(httpIf []string) {
	fs := http.FileServer(http.Dir(File.GetFilePath("assets/static")))
	http.Handle("/", fs)
	http.HandleFunc("/ws", v.WsPage)
	for i := 0; i < len(httpIf); i++ {
		req := httpIf[i]
		http.HandleFunc("/"+req, v.Requst)
	}
}

func (v *WebST) Requst(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("回应http:失败:", err)
		}
	}()

	//glog.Print("http request:", th.time, th.req)
	r.ParseForm()

	//fmt.Println(r)
	waitChan := make(chan []byte)
	v.ChanHTTP <- &HttpData{Req: r.RequestURI, Form: &r.Form, ChanBack: waitChan}

	data := <-waitChan

	if !r.Close {
		w.Write(data)
	}
}

//WsPage ...
func (v *WebST) WsPage(res http.ResponseWriter, req *http.Request) {
	//	fmt.Println("wsPage:", req.Header)
	conn, _error := (&websocket.Upgrader{HandshakeTimeout: 3 * time.Second, CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(res, req, nil)
	if _error != nil {
		println(_error)
		http.NotFound(res, req)
		return
	}

	//\\conn.WriteMessage(websocket.TextMessage, []byte("hello world"))
	//ChanConnected <- conn
	go func() {
		v.ChanAccept <- conn
	}()

}

//Get ...
func (v *WebST) Get(url string) map[string]string {
	mp := make(map[string]string)
	mp["result"] = "ok"

	var resp *http.Response
	var err error
	if strings.Index(url, "https") == -1 {
		resp, err = http.Get(url)
	} else {
		tr := &http.Transport{

			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		client := &http.Client{Transport: tr}
		resp, err = client.Get(url)
	}

	if err != nil {
		mp["result"] = "failed"
		return mp
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		mp["result"] = "failed"
		return mp
	}
	mp["data"] = string(body)
	return mp
}