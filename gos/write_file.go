package gos

import (
	"os"
	"path/filepath"
)

// FileExist 判断文件是否存在
func FileExist(name string) bool {
	_, err := os.Stat(name)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	panic(err)
}

// WriteFileCreate 新建文件，不存在创建，已经存在覆盖
func WriteFileCreate(name string, data []byte) (err error) {
	// 文件不存在则先创建目录
	err = os.MkdirAll(filepath.Dir(name), 0666)
	if err != nil {
		return
	}

	//新建文件，不存在创建，已经存在覆盖
	file, err := os.Create(name)
	defer func(f *os.File) {
		_ = f.Close()
	}(file)

	_, err = file.Write(data)
	if err != nil {
		return
	}
	return
}

// WriteFileAppend 写文件，路径不存在自动创建，已经存在追加内容
func WriteFileAppend(name string, data []byte) (err error) {

	// 文件不存在则先创建目录
	err = os.MkdirAll(filepath.Dir(name), 0666)
	if err != nil {
		return
	}

	// 以打开或创建文件的模式打开文件，并设置为追加模式
	file, err := os.OpenFile(name, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(file)

	_, err = file.Write(data)
	if err != nil {
		return
	}
	return
}
