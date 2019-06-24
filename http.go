package gnet

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

//ChanHTTP ...
var ChanHTTP chan *HTTPData

//ChanAccept ...
var ChanAccept chan *websocket.Conn

//HTTPData ...
type HTTPData struct {
	Req      string
	Form     *url.Values
	ChanBack chan []byte
}

var mux *http.ServeMux

//Init ...
func InitHttp() {
	mux = http.NewServeMux()
	ChanAccept = make(chan *websocket.Conn, 100)
	ChanHTTP = make(chan *HTTPData, 100)
}

//Close ...
func Close() {

}

//Start ...
func StartHttp(port int, ssl bool, httpIf []string) bool {
	InitHttp()
	RegisterHandles(httpIf)
	go func() {
		if ssl {
			Print("start https service")
			err := http.ListenAndServeTLS(":8081", GetFilePath("assets/keys/ssl.crt"), GetFilePath("assets/keys/ssl.key"), mux)
			fmt.Println(err)
			if err != nil {
				panic(err)
			}

		} else {
			Print("start http service")
			err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
			fmt.Println(err)
			if err != nil {
				panic(err)
			}
		}

	}()

	return true
}

//RegisterHandles ...
func RegisterHandles(httpIf []string) {
	fs := http.FileServer(http.Dir(GetFilePath("assets/static")))
	//fmt.Println(fs)
	mux.Handle("/", fs)

	//websocket
	mux.HandleFunc("/ws", WsPage)

	addHandle := func(req string) {
		mux.Handle("/"+req, NewServerHandle(req))
	}
	for i := 0; i < len(httpIf); i++ {
		// println(httpIf[i])
		addHandle(httpIf[i])
	}

}

//WsPage ...
func WsPage(res http.ResponseWriter, req *http.Request) {
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
		ChanAccept <- conn
	}()

}

//Get ...
func Get(url string) map[string]string {
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
