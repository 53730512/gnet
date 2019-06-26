package gnet

//DEBUG ...
type commonST struct {
	DEBUG bool
}

func newCommon() *commonST {
	ptr := &commonST{}
	if ptr.init() {
		return ptr
	} else {
		return nil
	}
}

func (v *commonST) init() bool {
	v.DEBUG = false
	return true
}

//IsLittleEndian ...
func (v *commonST) IsLittleEndian() bool {
	s := int16(0x1234)
	b := int8(s)
	return b == 0x34
}
