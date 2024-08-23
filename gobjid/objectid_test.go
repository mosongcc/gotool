package gobjid

import (
	"testing"
	"time"
)

func TestNewObjectID(t *testing.T) {
	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)
		go func() {
			t.Log(NewObjectID().Hex())
		}()
	}
	time.Sleep(3 * time.Second)
}
