package error

// TODO
// here define ONLY errors from the library NOT from flickr
// error from flickr have already a code and a message that are returned
// along with the HTTP Response

const (
	FooError = 1
)

var errors = map[int]string{
	FooError: "This error is Foo.",
}

type Error struct {
	ErrorCode int
	Message   string
}

func NewError(errorCode int) *Error {
	return &Error{
		ErrorCode: errorCode,
		Message:   errors[errorCode],
	}
}
