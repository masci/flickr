package flickr

import (
	"testing"
)

func TestExpect(t *testing.T) {
	t2 := testing.T{}
	Expect(&t2, 1, 2)
	if !t2.Failed() {
		t.Errorf("Expect should fail")
	}
}
