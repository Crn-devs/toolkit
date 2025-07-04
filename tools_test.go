package toolkit

import "testing"

func TestTools_RandomString(t *testing.T) {
	var tool Tools

	s := tool.RandomString(10)

	if len(s) != 10 {
		t.Error("fail: string is the wrong length")
	}
}
