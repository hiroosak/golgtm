package golgtm

import "testing"

func TestString(t *testing.T) {
	s, _ := NewText("LGTM")
	if "LGTM" != s.String() {
		t.Errorf("expect: LGTM, but %v", s.String())
	}
}
