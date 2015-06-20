package error

const (
	SSLIsRequired             = 95
	InvalidSignature          = 96
	MissingSignature          = 97
	LoginFailed               = 98
	UserNotLoggedIn           = 99
	InvalidAPIKey             = 100
	ServiceCurrentUnavailable = 105
	WriteOperationFailed      = 106
	FormatNotFound            = 111
	MethodNotFound            = 112
	BadURLFound               = 116
)

var errors = map[int]string{
	SSLIsRequired:             "SSL is required to access the Flickr API.",
	InvalidSignature:          "The passed signature was invalid.",
	MissingSignature:          "The call required signing but no signature was sent.",
	LoginFailed:               "The login details or auth token passed were invalid.",
	UserNotLoggedIn:           "The login details or auth token passed were invalid.",
	InvalidAPIKey:             "The login details or auth token passed were invalid.",
	ServiceCurrentUnavailable: "The login details or auth token passed were invalid.",
	WriteOperationFailed:      "The login details or auth token passed were invalid.",
	FormatNotFound:            "The login details or auth token passed were invalid.",
	MethodNotFound:            "The login details or auth token passed were invalid.",
	BadURLFound:               "The login details or auth token passed were invalid.",
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
