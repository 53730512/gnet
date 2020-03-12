package main

import "gitee.com/liyp_admin/gnet"

func init() {
	TestFun()
}

func main() {
	reqList := []string{"roleNum"}
	var _ioservice = &gnet.STIoservice{}
	gnet.Start(_ioservice, 1)
	gnet.Listen(10080, false, reqList)
	gnet.WaitClose()
}

func TestFun() {

}
