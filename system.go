package gnet

import (
	//"fmt"
	"runtime"
	// "syscall"
	// "unsafe"
)

//os
const (
	OSWindows = iota
	OSLinux
)


type sysST struct {
	_SetConsoleTitle uintptr
}

func newSys() *sysST {
	ptr := &sysST{}
	if ptr.Init() {
		return ptr
	} else {
		return nil
	}
}

func (v *sysST) Init() bool {
	if v.OS() == OSWindows {
		// kernel32, loadErr := syscall.LoadLibrary("kernel32.dll")
		// if loadErr != nil {
		// 	fmt.Println("loadErr", loadErr)
		// }
		// defer syscall.FreeLibrary(kernel32)
		// v._SetConsoleTitle, _ = syscall.GetProcAddress(kernel32, "SetConsoleTitleW")
	}

	return true
}

//获得系统
func (v *sysST) OS() uint8 {
	switch runtime.GOOS {
	case "windows":
		return OSWindows
	case "linux":
		return OSLinux
	default:
		return 0
	}

}

func (v *sysST) SetConsoleTitle(title string) int {
	// if v.OS() != OSWindows {
	// 	return 0
	// }

	// ret, _, callErr := syscall.Syscall(v._SetConsoleTitle, 1, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))), 0, 0)
	// if callErr != 0 {
	// 	fmt.Println("callErr", callErr)
	// }
	// return int(ret)

	return 0
}
