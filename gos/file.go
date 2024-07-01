package gos

import "os"

// WriteFileCreate TODO 写文件，路径不存在自动创建，已经存在覆盖内容
func WriteFileCreate(name string, data []byte, perm os.FileMode) error {

	return os.WriteFile(name, data, perm)
}

// WriteFileAppend TODO 写文件，路径不存在自动创建，已经存在追加内容
func WriteFileAppend(name string, data []byte, perm os.FileMode) error {

	return os.WriteFile(name, data, perm)
}
