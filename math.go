package gnet

import (
	"math/rand"
	"time"
)

//BytesToInt 字节转换成整形
// func BytesToInt(b []byte) int {
// 	bytesBuffer := bytes.NewBuffer(b)
// 	var tmp int32
// 	binary.Read(bytesBuffer, binary.BigEndian, &tmp)
// 	return int(tmp)
// }
//Add ...
var randHanle *rand.Rand

func init() {
	randHanle = rand.New(rand.NewSource(time.Now().Unix()))
}

type mathST struct {
}

func newMath() *mathST {
	return &mathST{}
}

func (v *mathST) Add(x, y int) int {
	return x + y
}

//Random (1-3) 结果可能为 1,2,3
func (v *mathST) Random(min int, max int) int {
	if min > max {
		return 0
	} else if min == max {
		return min
	}

	dis := max - min
	return randHanle.Intn(dis) + min
}
