package common

import (
	"testing"
	"time"
)

func TestT(t *testing.T) {
	for i := 0; i < 10; i++ {
		res := LocalRandomBytes()
		t.Log(res)
		time.Sleep(1 * time.Nanosecond)
	}
}
