package gnet

//DEBUG ...
type CommonST struct {
	DEBUG bool
}

func NewCommon() *CommonST {
	ptr := &CommonST{}
	if ptr.Init() {
		return ptr
	} else {
		return nil
	}
}

func (v *CommonST) Init() bool {
	v.DEBUG = false
	return true
}

//IsLittleEndian ...
func (v *CommonST) IsLittleEndian() bool {
	s := int16(0x1234)
	b := int8(s)
	return b == 0x34
}
