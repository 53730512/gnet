package gnet

import (
	"bytes"
	"encoding/binary"
	"strings"
	"unsafe"
)

type FormatST struct {
}

func NewFormat() *FormatST {
	ptr := &FormatST{}
	if ptr.Init() {
		return ptr
	} else {
		return nil
	}
}

func (v *FormatST) Init() bool {
	return true
}

//ResetByte ...
func (v *FormatST) ResetByte(data []byte) {
	for i := 0; i < len(data); i++ {
		data[i] = 0
	}
}

//BytesToInt 字节转换成整形s
func (v *FormatST) BytesToInt(b []byte) int {
	bBuf := bytes.NewBuffer(b)
	var x int64
	if Common.IsLittleEndian() {
		binary.Read(bBuf, binary.LittleEndian, &x)
	} else {
		binary.Read(bBuf, binary.BigEndian, &x)
	}
	return int(x)
}

//IntToBytes 整形转换成字节
func (v *FormatST) IntToBytes(i int) []byte {
	size := unsafe.Sizeof(i)
	var buf = make([]byte, size)
	if Common.IsLittleEndian() {
		if size == 4 {
			binary.LittleEndian.PutUint32(buf, uint32(i))
		} else {
			binary.LittleEndian.PutUint64(buf, uint64(i))
		}
	} else {
		if size == 4 {
			binary.BigEndian.PutUint32(buf, uint32(i))
		} else {
			binary.BigEndian.PutUint64(buf, uint64(i))
		}
	}
	return buf
}

//StringClean ...
func (v *FormatST) StringClean(str *string) {
	*str = strings.Replace(*str, " ", "", -1)
	*str = strings.Replace(*str, "\n", "", -1)
}
