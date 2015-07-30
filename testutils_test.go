package flickr

import (
	"testing"
)

func TestExpect(t *testing.T) {
	t2 := *t
	Expect(&t2, 1, 2)
	if !t2.Failed() {
		t.Errorf("Expect should fail")
	}
}
