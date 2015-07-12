package error

import (
	"testing"
)

func TestError(t *testing.T) {
	var e *Error
	e = NewError(ApiError)
	if e.Error() != errors[ApiError] {
		t.Errorf("Expected", errors[ApiError], "found", e.Error())
	}
}
