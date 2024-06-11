package errors

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

func New(format string, args ...any) error {
	msg := fmt.Sprintf(format, args...)

	return &Error{Message: msg, Stack: stack()}
}

func Wrap(err error) error {
	return &Error{Message: err.Error(), Stack: stack()}
}

type Error struct {
	Message string `json:"message,omitempty"`
	Params  []any  `json:"params,omitempty"`
	Stack   string `json:"stack,omitempty"`
}

func (e Error) Error() string {
	msg := strings.Split(e.Message, " ")
	index := 0

	for i, v := range msg {
		if strings.Contains(v, "[") {
			msg[i] = fmt.Sprintf("%v", e.Params[index])
			index++
		}
	}

	return strings.Join(msg, " ") + "\n" + e.Stack
}

func (e *Error) Format(format string, args ...any) *Error {
	for i := range args {
		index := strings.Index(format, "%")
		format = format[:index] + fmt.Sprintf("[%d]", i) + format[index+2:]
	}

	e.Message = format
	e.Params = args
	e.Stack = stack()

	return e
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
	buf := make([]byte, 512)
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
)
