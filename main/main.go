package main

import (
	"gitee.com/liyp/gnet"
)

func init() {
	TestFun()
}

func main() {
	gnet.Success("abc")
	gnet.WaitClose()
}

func TestFun() {

}
