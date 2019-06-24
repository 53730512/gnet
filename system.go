package el

import (
	"fmt"
	"runtime"
	"syscall"
	"unsafe"
)

//os
const (
	OSWindows = iota
	OSLinux
)

var (
	_SetConsoleTitle uintptr
)

//获得系统
func OS() uint8 {
	switch runtime.GOOS {
	case "windows":
		return OSWindows
	case "linux":
		return OSLinux
	default:
		return 0
	}

}

func SystemInit() {
	if OS() == OSWindows {
		kernel32, loadErr := syscall.LoadLibrary("kernel32.dll")
		if loadErr != nil {
			fmt.Println("loadErr", loadErr)
		}
		defer syscall.FreeLibrary(kernel32)
		_SetConsoleTitle, _ = syscall.GetProcAddress(kernel32, "SetConsoleTitleW")
	}
}
func SetConsoleTitle(title string) int {
	if OS() != OSWindows {
		return 0
	}

	ret, _, callErr := syscall.Syscall(_SetConsoleTitle, 1, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))), 0, 0)
	if callErr != 0 {
		fmt.Println("callErr", callErr)
	}
	return int(ret)
}
