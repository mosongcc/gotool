package gos

import (
	"log"
	"os"
	"path/filepath"
)

// GetExecutablePath 获取当前执行文件绝对路径
func GetExecutablePath() string {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	res, _ := filepath.EvalSymlinks(filepath.Dir(exePath))
	return res
}
