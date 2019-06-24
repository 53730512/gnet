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

type LogST struct {
	queueChan chan *queueData
	logger    *log.Logger
}

func NewLog() *LogST {
	ptr := &LogST{}
	if ptr.Init() {
		return ptr
	} else {
		return nil
	}
}

func (v *LogST) Init() bool {
	v.queueChan = make(chan *queueData, 100)

	os.Mkdir("log", os.ModeDir)
	file, err := os.OpenFile("log/"+v.getLogFileName(), os.O_APPEND|os.O_CREATE, 666)
	if err != nil {
		println("glog 初始化失败:%s", err.Error())
		return false
	}

	v.logger = log.New(file, "", log.Ldate|log.Ltime)

	go func() {
		for {
			select {
			case data := <-v.queueChan:
				color.Set(data.color, color.Bold)
				log.Println(data._str, "			", data.file, data.line)
				v.logger.Println(data._str, "			", data.file, data.line)
				color.Set(color.FgWhite, color.Bold)
			}
		}
	}()
	return true
}
func (v *LogST) getLogFileName() string {
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

//Print ...
func (v *LogST) Print(format string, a ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	str := fmt.Sprintf(format, a...)

	data := &queueData{}
	data.file = file
	data.line = line
	data._str = str
	data.color = color.FgWhite
	v.queueChan <- data

}

//Success ...
func (v *LogST) Success(format string, a ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	str := fmt.Sprintf(format, a...)

	data := &queueData{}
	data.file = file
	data.line = line
	data._str = str
	data.color = color.FgGreen
	v.queueChan <- data
}

//Warning ...
func (v *LogST) Warning(format string, a ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	str := fmt.Sprintf(format, a...)

	data := &queueData{}
	data.file = file
	data.line = line
	data._str = str
	data.color = color.FgYellow
	v.queueChan <- data
}

//Error ...
func (v *LogST) Error(format string, a ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	str := fmt.Sprintf(format, a...)

	data := &queueData{}
	data.file = file
	data.line = line
	data._str = str
	data.color = color.FgRed
	v.queueChan <- data
}
