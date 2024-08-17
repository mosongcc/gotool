package gos

import (
	"testing"
	"time"
)

// 文件追加内容测试
func TestWriteFileAppend(t *testing.T) {
	for i := 0; i < 10; i++ {
		err := WriteFileAppend(".data/t1.txt", []byte("\nContent "+time.Now().String()))
		if err != nil {
			t.Error(err)
		}
	}
}

// 文件覆盖内容测试
func TestWriteFileCreate(t *testing.T) {
	err := WriteFileCreate(".data/t2.txt", []byte("Content2"))
	if err != nil {
		t.Error(err)
	}
}
