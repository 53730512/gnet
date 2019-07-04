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

//httpData ...
type httpData struct {
	Req      string
	Form     *url.Values
	ChanBack chan []byte
}

type webST struct {
	mux        *http.ServeMux
	ChanAccept chan *websocket.Conn
	ChanHTTP   chan *httpData
}

func newWeb() *webST {
	ptr := &webST{}
	if ptr.init() {
		return ptr
	} else {
		return nil
	}
}

func (v *webST) init() bool {
	v.mux = http.NewServeMux()
	v.ChanAccept = make(chan *websocket.Conn, 100)
	v.ChanHTTP = make(chan *httpData, 100)
	return true
}

//Close ...
func (v *webST) Close() {

}

//Start ...
func (v *webST) Start(port int, ssl bool, httpIf []string) bool {
	v.RegisterHandles(httpIf)
	go func() {
		if ssl {
			//Log.Print("start https service")
			err := http.ListenAndServeTLS(":8081", File.GetFilePath("assets/keys/ssl.crt"), File.GetFilePath("assets/keys/ssl.key"), v.mux)
			fmt.Println(err)
			if err != nil {
				panic(err)
			}

		} else {
			//Log.Print("start http service")
			err := http.ListenAndServe(fmt.Sprintf(":%d", port), v.mux)
			fmt.Println(err)
			if err != nil {
				panic(err)
			}
		}

	}()

	return true
}

//RegisterHandles ...
func (v *webST) RegisterHandles(httpIf []string) {
	fs := http.FileServer(http.Dir(File.GetFilePath("assets/static")))
	//fmt.Println(fs)
	v.mux.Handle("/", fs)

	//websocket
	v.mux.HandleFunc("/ws", v.WsPage)

	addHandle := func(req string) {
		v.mux.Handle("/"+req, NewServerHandle(req))
	}
	for i := 0; i < len(httpIf); i++ {
		// println(httpIf[i])
		addHandle(httpIf[i])
	}

}

//WsPage ...
func (v *webST) WsPage(res http.ResponseWriter, req *http.Request) {
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
func (v *webST) Get(url string) map[string]string {
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
