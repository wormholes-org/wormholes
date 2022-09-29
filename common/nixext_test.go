package common

import (
	"testing"
)

func TestT(t *testing.T) {
	t.Log("hello world")
	res := LocalRandomBytes()
	t.Log(res)
	t.Log("123")
}
