package gnet

import (
	"testing"
)

//TestAdd
func TestAdd_1(t *testing.T) {
	if Add(1, 2) != 3 {
		t.Error("测试失败,")
	}
}
