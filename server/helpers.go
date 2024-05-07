package server

import (
	"net/http"
	"reflect"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/maadiii/hertz/errors"
)

func bind[IN any, OUT any](handler *Handler[IN, OUT], rctx *app.RequestContext) (req IN, err error) {
	p := reflect.TypeOf(handler.Action).In(1)
	if p.Kind() == reflect.Interface {
		return
	}

	req = reflect.New(p.Elem()).Interface().(IN)

	err = rctx.Bind(req)

	return
}

var methods = map[string]string{
	"[GET]":     http.MethodGet,
	"[HEAD]":    http.MethodHead,
	"[POST]":    http.MethodPost,
	"[PUT]":     http.MethodPut,
	"[PATCH]":   http.MethodPatch,
	"[DELETE]":  http.MethodDelete,
	"[CONNECT]": http.MethodConnect,
	"[OPTIONS]": http.MethodOptions,
	"[TRACE]":   http.MethodTrace,
}

var abortType = map[*errors.Error]int{
	errors.BadRequest:                    http.StatusBadRequest,
	errors.Unauthorized:                  http.StatusUnauthorized,
	errors.PaymentRequired:               http.StatusPaymentRequired,
	errors.Forbidden:                     http.StatusForbidden,
	errors.NotFound:                      http.StatusNotFound,
	errors.MethodNotAllowed:              http.StatusMethodNotAllowed,
	errors.NotAcceptable:                 http.StatusNotAcceptable,
	errors.ProxyAuthRequired:             http.StatusProxyAuthRequired,
	errors.RequestTimeout:                http.StatusRequestTimeout,
	errors.Conflict:                      http.StatusConflict,
	errors.Gone:                          http.StatusGone,
	errors.LengthRequired:                http.StatusLengthRequired,
	errors.PreconditionFailed:            http.StatusPreconditionFailed,
	errors.RequestEntityTooLarge:         http.StatusRequestEntityTooLarge,
	errors.RequestURITooLong:             http.StatusRequestURITooLong,
	errors.UnsupportedMediaType:          http.StatusUnsupportedMediaType,
	errors.RequestedRangeNotSatisfiable:  http.StatusRequestedRangeNotSatisfiable,
	errors.ExpectationFailed:             http.StatusExpectationFailed,
	errors.Teapot:                        http.StatusTeapot,
	errors.MisdirectedRequest:            http.StatusMisdirectedRequest,
	errors.UnprocessableEntity:           http.StatusUnprocessableEntity,
	errors.Locked:                        http.StatusLocked,
	errors.FailedDependency:              http.StatusFailedDependency,
	errors.TooEarly:                      http.StatusTooEarly,
	errors.UpgradeRequired:               http.StatusUpgradeRequired,
	errors.PreconditionRequired:          http.StatusPreconditionFailed,
	errors.TooManyRequests:               http.StatusTooManyRequests,
	errors.RequestHeaderFieldsTooLarge:   http.StatusRequestHeaderFieldsTooLarge,
	errors.UnavailableForLegalReasons:    http.StatusUnavailableForLegalReasons,
	errors.NotImplemented:                http.StatusNotImplemented,
	errors.BadGateway:                    http.StatusBadGateway,
	errors.ServiceUnavailable:            http.StatusServiceUnavailable,
	errors.GatewayTimeout:                http.StatusGatewayTimeout,
	errors.HTTPVersionNotSupported:       http.StatusHTTPVersionNotSupported,
	errors.VariantAlsoNegotiates:         http.StatusVariantAlsoNegotiates,
	errors.InsufficientStorage:           http.StatusInsufficientStorage,
	errors.LoopDetected:                  http.StatusLoopDetected,
	errors.NotExtended:                   http.StatusNotExtended,
	errors.NetworkAuthenticationRequired: http.StatusNetworkAuthenticationRequired,
	errors.InternalServerError:           http.StatusInternalServerError,
	errors.AlreadyExist:                  http.StatusConflict,
}
