package utils

import (
	"os"
	"path"
	"runtime"
)

// 获取当前执行文件绝对路径（go run）
func GetCurrentAbPathByCaller() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
	}
	return path.Join(abPath, "../../../")
}

// 获取当前工作目录
func GetCurrentWorkDirectory() string {
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}

	return wd
}

// 反转数组

func ReverseArray(array []string) []string {
	newarray := array
	for i, j := 0, len(newarray)-1; i < j; i, j = i+1, j-1 {
		newarray[i], newarray[j] = newarray[j], newarray[i]
	}

	return newarray
}
