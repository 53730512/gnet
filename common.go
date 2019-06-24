package el

//DEBUG ...
var DEBUG bool

func InitCommon() {
	DEBUG = true
	SystemInit()
}

//IsLittleEndian ...
func IsLittleEndian() bool {
	s := int16(0x1234)
	b := int8(s)
	return b == 0x34
}
