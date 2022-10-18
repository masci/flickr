package error

import (
	"testing"
)

func TestError(t *testing.T) {
	var e *Error
	e = NewError(ApiError, "foo")
	if e.Error() != errors[ApiError]+"foo" {
		t.Errorf("Expected %+v found %+v", errors[ApiError], e.Error())

	}
}
