package gnet

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type fileST struct {
}

func newFile() *fileST {
	ptr := &fileST{}
	if ptr.init() {
		return ptr
	} else {
		return nil
	}
}

func (v *fileST) init() bool {
	return true
}

//PathExists ...
func (v *fileST) PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//GetFilePath ...
func (v *fileST) GetFilePath(localPath string) string {
	//命令启动时的路径
	exists, _ := v.PathExists(localPath)
	if !exists {
		AppPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			panic(err)
		}
		//二进制启动时路径
		localPath = filepath.Join(AppPath, localPath)
		exists, err := v.PathExists(localPath)
		if !exists {
			fmt.Println(localPath + "not found")
			return ""
		}
	}

	return localPath
}

//GetExePath ...
func (v *fileST) GetExePath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	index := strings.LastIndex(path, string(os.PathSeparator))

	return path[:index]
}
