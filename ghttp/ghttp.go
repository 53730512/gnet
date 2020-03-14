package ghttp

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"gitee.com/liyp_admin/gnet/ghttp/user"

	"gitee.com/liyp_admin/gnet/gfile"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var Handle *webST

//HttpData ...
type HttpData struct {
	Req      string
	Form     *url.Values
	ChanBack chan []byte
}

type webST struct {
	ChanAccept chan *websocket.Conn
	ChanHTTP   chan *HttpData
	httpserver *http.Server
}

func init() {
	Handle = &webST{}
	Handle.init()
}

func (v *webST) init() bool {
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
func (v *webST) Close() {

}

//Start ...
func Start(port int, ssl bool, httpIf []string) bool {
	Handle.httpserver.Addr = fmt.Sprintf(":%d", port)

	Handle.RegisterHandles(httpIf)
	go func() {
		if ssl {
			//Log.Print("start https service")
			err := Handle.httpserver.ListenAndServeTLS(gfile.GetFilePath("assets/keys/ssl.crt"), gfile.GetFilePath("assets/keys/ssl.key"))
			fmt.Println(err)
			if err != nil {
				panic(err)
			}

		} else {
			//Log.Print("start http service")
			err := Handle.httpserver.ListenAndServe()
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
	Route := mux.NewRouter()

	session := Route.PathPrefix("/user").Subrouter()
	session.HandleFunc("/login", user.Login)
	session.HandleFunc("/logout", user.Logout)
	session.HandleFunc("/Secret", user.Secret)

	fs := http.FileServer(http.Dir(gfile.GetFilePath("assets/static")))
	Route.Handle("/", fs)
	Route.HandleFunc("/ws", v.WsPage)

	for i := 0; i < len(httpIf); i++ {
		req := httpIf[i]
		Route.HandleFunc("/"+req, v.Requst)
	}

	http.Handle("/", Route)
}
func (v *webST) Requst(w http.ResponseWriter, r *http.Request) {
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
