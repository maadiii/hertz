package errors

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

type Error struct {
	Message string `json:"message,omitempty"`
	Key     string `json:"key,omitempty"`
	Stack   string `json:"stack,omitempty"`
}

func (e *Error) WithKey(key string) {
	key = strings.ReplaceAll(key, " ", "_")
	e.Key = key
}

func (e Error) Error() string {
	return e.Message
}

func (e *Error) Wrap(err error) error {
	e.Message = err.Error()
	e.Stack = stack()

	return e
}

func New(format string, args ...any) error {
	msg := fmt.Sprintf(format, args...)

	return &Error{Message: msg, Stack: stack()}
}

func Wrap(err error) error {
	return &Error{Message: err.Error(), Stack: stack()}
}

func Join(errs ...error) error {
	return errors.Join(errs...)
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target any) bool {
	return errors.As(err, target)
}

func stack() string {
	buf := make([]byte, 1024)
	runtime.Stack(buf, false)
	lines := strings.Split(string(buf), "\n")

	var stack string

	for i := 5; i < len(lines)-2; i += 2 {
		if i != 5 {
			stack += "\n"
		}

		stack += strings.Join(lines[i:i+2], "")
		stack = stack[:strings.LastIndex(stack, " ")]
	}

	return stack
}

var (
	BadRequest                   = &Error{} //nolint
	Unauthorized                 = &Error{} //nolint
	PaymentRequired              = &Error{} //nolint
	Forbidden                    = &Error{} //nolint
	NotFound                     = &Error{} //nolint
	MethodNotAllowed             = &Error{} //nolint
	NotAcceptable                = &Error{} //nolint
	ProxyAuthRequired            = &Error{} //nolint
	RequestTimeout               = &Error{} //nolint
	Conflict                     = &Error{} //nolint
	Gone                         = &Error{} //nolint
	LengthRequired               = &Error{} //nolint
	PreconditionFailed           = &Error{} //nolint
	RequestEntityTooLarge        = &Error{} //nolint
	RequestURITooLong            = &Error{} //nolint
	UnsupportedMediaType         = &Error{} //nolint
	RequestedRangeNotSatisfiable = &Error{} //nolint
	ExpectationFailed            = &Error{} //nolint
	Teapot                       = &Error{} //nolint
	MisdirectedRequest           = &Error{} //nolint
	UnprocessableEntity          = &Error{} //nolint
	Locked                       = &Error{} //nolint
	FailedDependency             = &Error{} //nolint
	TooEarly                     = &Error{} //nolint
	UpgradeRequired              = &Error{} //nolint
	PreconditionRequired         = &Error{} //nolint
	TooManyRequests              = &Error{} //nolint
	RequestHeaderFieldsTooLarge  = &Error{} //nolint
	UnavailableForLegalReasons   = &Error{} //nolint
	AlreadyExist                 = &Error{} //nolint

	InternalServerError           = &Error{} //nolint
	NotImplemented                = &Error{} //nolint
	BadGateway                    = &Error{} //nolint
	ServiceUnavailable            = &Error{} //nolint
	GatewayTimeout                = &Error{} //nolint
	HTTPVersionNotSupported       = &Error{} //nolint
	VariantAlsoNegotiates         = &Error{} //nolint
	InsufficientStorage           = &Error{} //nolint
	LoopDetected                  = &Error{} //nolint
	NotExtended                   = &Error{} //nolint
	NetworkAuthenticationRequired = &Error{} //nolint
	Retry                         = &Error{} //nolint
)
