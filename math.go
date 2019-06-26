package gnet

//BytesToInt 字节转换成整形
// func BytesToInt(b []byte) int {
// 	bytesBuffer := bytes.NewBuffer(b)
// 	var tmp int32
// 	binary.Read(bytesBuffer, binary.BigEndian, &tmp)
// 	return int(tmp)
// }
//Add ...
type mathST struct {
}

func newMath() *mathST {
	return &mathST{}
}

func (v *mathST) Add(x, y int) int {
	return x + y
}
