package gnet

import (
	"errors"
	"fmt"
	"time"

	"gitee.com/liyp_admin/gnet/ghttp"

	_ "net/http/pprof"

	"gitee.com/liyp_admin/gnet/gfile"
)

var Common *commonST
var IsInit bool
var DB *dbST
var Format *formatST
var File *gfile.FileST
var Log *logST
var Config *configST
var Math *mathST
var Sys *sysST
var Date *dateST
var Web *ghttp.WebST
var Service *serviceST

var closeChan = make(chan string, 20)
var consoleChan = make(chan string, 20)

func init() {
	err := errors.New("gnet初始化失败...")
	IsInit = false
	Common = newCommon()
	if Common == nil {
		Close(err)
		return
	}

	DB = newDB()
	if DB == nil {
		Close(err)
		return
	}

	Format = newFormat()
	if Format == nil {
		Close(err)
		return
	}

	File = newFile()
	if File == nil {
		Close(err)
		return
	}

	Log = newLog()
	if Log == nil {
		Close(err)
		return
	}

	Config = newConfig()
	if Config == nil {
		Close(err)
		return
	}

	Math = newMath()
	if Math == nil {
		Close(err)
		return
	}

	Sys = newSys()
	if Sys == nil {
		Close(err)
		return
	}

	Date = newDate()
	if Date == nil {
		Close(err)
		return
	}

	Web = ghttp.NewWeb()
	if Web == nil {
		Close(err)
		return
	}

	Service = newService()
	if Service == nil {
		Close(err)
		return
	}
}

func inputMornitor() {
	for {
		var input string
		fmt.Scanln(&input)
		if input == "q" || input == "Q" {
			closeChan <- input
		} else {
			consoleChan <- input
		}
	}
}

func Start(handle IFIoservice, fps int) {
	Service.SetHandle(handle)
	Service.run(fps)
}

func Listen(port int, ssl bool, httpIf []string) {
	go func() {
		Service.GetHandle().Listen(port, ssl, httpIf)
	}()
}

func WaitClose() {
	go inputMornitor()

	Print("输入'q'退出程序")
	for {
		input := <-closeChan
		if input == "q" || input == "Q" {
			break
		} else {
			fmt.Println("无效命令:", input)
		}
	}

	Warning("正在退出...")
	Warning("************************************")
	time.Sleep(200 * time.Millisecond)
}

func Close(err error) {
	if len(closeChan) == 0 {
		if err != nil {
			Error("服务器关闭中:%s", err.Error())
		} else {
			Error("服务器关闭中")
		}
		closeChan <- "q"
	}
}

func Print(format string, a ...interface{}) {
	Log.Print(format, a...)
}

func Success(format string, a ...interface{}) {
	Log.Success(format, a...)
}

func Warning(format string, a ...interface{}) {
	Log.Warning(format, a...)
}

func Error(format string, a ...interface{}) {
	Log.Error(format, a...)
}
