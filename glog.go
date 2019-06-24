package gnet

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/fatih/color"
)

//Print ...
const (
	_print = iota
	_sucess
	_warning
	_error
)

type queueData struct {
	color color.Attribute
	file  string
	line  int
	_str  string
}

var queueChan chan *queueData

func getLogFileName() string {
	tm := time.Now()
	t1 := tm.Year()
	t2 := tm.Month()
	t3 := tm.Day()
	t4 := tm.Hour()
	t5 := tm.Minute()
	t6 := tm.Second()
	t7 := tm.Nanosecond() / 1000000
	return fmt.Sprintf("%d-%02d-%02d %02d-%02d-%02d %dms.log", t1, t2, t3, t4, t5, t6, t7)
}

var logger *log.Logger

//Init ...
func InitLog() bool {
	queueChan = make(chan *queueData, 100)
	os.Mkdir("log", os.ModeDir)
	file, err := os.OpenFile("log/"+getLogFileName(), os.O_APPEND|os.O_CREATE, 666)
	if err != nil {
		println("glog 初始化失败:%s", err.Error())
		return false
	}

	logger = log.New(file, "", log.Ldate|log.Ltime)

	go func() {
		for {
			select {
			case data := <-queueChan:
				color.Set(data.color, color.Bold)
				log.Println(data._str, "			", data.file, data.line)
				logger.Println(data._str, "			", data.file, data.line)
				color.Set(color.FgWhite, color.Bold)
			}
		}
	}()
	return true
}

//Print ...
func Print(format string, a ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	str := fmt.Sprintf(format, a...)

	data := &queueData{}
	data.file = file
	data.line = line
	data._str = str
	data.color = color.FgWhite
	queueChan <- data

}

//Success ...
func Success(format string, a ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	str := fmt.Sprintf(format, a...)

	data := &queueData{}
	data.file = file
	data.line = line
	data._str = str
	data.color = color.FgGreen
	queueChan <- data
}

//Warning ...
func Warning(format string, a ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	str := fmt.Sprintf(format, a...)

	data := &queueData{}
	data.file = file
	data.line = line
	data._str = str
	data.color = color.FgYellow
	queueChan <- data
}

//Error ...
func Error(format string, a ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	str := fmt.Sprintf(format, a...)

	data := &queueData{}
	data.file = file
	data.line = line
	data._str = str
	data.color = color.FgRed
	queueChan <- data
}
