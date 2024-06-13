package server

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/maadiii/hertz/errors"
)

func Handle[IN any, OUT any](action func(context.Context, *Request, IN) (OUT, error)) {
	handler := &Handler[IN, OUT]{Action: action}
	handler.fixAPIDescriber()
	handler.fixIdentifierDesciber()

	key := fmt.Sprintf("%s::%s::%d::%s", handler.Method, handler.Path, handler.Status, handler.ActionType)

	if handler.identifierDescriber != nil {
		handlersMap[key] = append(handlersMap[key], identify(handler))
	}

	decorators := handler.getDecorators()
	for _, dec := range decorators {
		handlersMap[key] = append(handlersMap[key], decorate(handler.Method, dec))
	}

	handlersMap[key] = append(handlersMap[key], handle(handler))
}

func handle[IN any, OUT any](handler *Handler[IN, OUT]) app.HandlerFunc {
	return func(c context.Context, reqContext *app.RequestContext) {
		reqType, err := bind(handler, reqContext)
		if err != nil {
			reqContext.AbortWithStatusJSON(
				http.StatusUnprocessableEntity,
				errors.New(fmt.Sprintf( //nolint
					"%s\tAPI=%s\tMethod=%s\tHandler=%s",
					err.Error(),
					handler.Path,
					handler.Method,
					runtimeFunc(handler.Action).Name(),
				)),
			)

			return
		}

		req := &Request{reqContext}

		res, err := handler.Action(c, req, reqType)
		if err != nil {
			handleError(reqContext, err)

			return
		}

		handler.RespondFn(reqContext, res)
	}
}

type Handler[IN any, OUT any] struct {
	Action    func(context.Context, *Request, IN) (OUT, error)
	RespondFn func(rctx *app.RequestContext, response any)

	*apiDescriber
	*identifierDescriber
}

type apiDescriber struct {
	Path        string
	Method      string
	Status      int
	ContentType string
	ActionType  string
}

type identifierDescriber struct {
	Roles       []string
	Permissions []string
}

func (h *Handler[IN, OUT]) fixAPIDescriber() {
	comment := funcDescription(h.Action)
	comments := strings.Split(comment, "\n")
	name := runtimeFunc(h.Action).Name()
	apiDescriber := h.getFixedAPIDescriberFields(comments)

	if len(apiDescriber) == 0 {
		panic(name + " has not describer")
	}

	for i, d := range apiDescriber {
		if strings.HasPrefix(d, "[") && strings.HasSuffix(d, "]") {
			verb, ok := methods[strings.ToUpper(d)]
			if !ok {
				panic(name + " has invalid VERB")
			}

			h.Method = verb

			continue
		}

		if strings.HasPrefix(d, "/") {
			d = strings.TrimRight(d, "/")
			h.Path = d

			continue
		}

		if status, err := strconv.Atoi(d); err == nil {
			h.Status = status

			continue
		}

		if strings.Contains(d, "@") {
			typeAndContentType := strings.Split(d, "@")
			h.ActionType = typeAndContentType[0]
			h.ContentType = fmt.Sprintf("%s %s", typeAndContentType[1], apiDescriber[i+1])

			break
		}

		h.ActionType = d
	}

	h.setResponder(name)
}

func (h *Handler[IN, OUT]) getFixedAPIDescriberFields(comments []string) []string {
	h.apiDescriber = new(apiDescriber)

	for _, describer := range comments {
		if !strings.HasPrefix(describer, "[") {
			continue
		}

		desc := strings.Split(describer, " ")

		_, ok := methods[strings.ToUpper(desc[0])]
		if !ok {
			continue
		}

		return desc
	}

	return []string{}
}

func (h *Handler[IN, OUT]) fixIdentifierDesciber() {
	comment := funcDescription(h.Action)
	comments := strings.Split(comment, "\n")

	for _, describer := range comments {
		if !strings.HasPrefix(describer, "@authorize") {
			continue
		}

		h.identifierDescriber = new(identifierDescriber)

		describer, _ = strings.CutPrefix(describer, "@authorize")
		describer = strings.ReplaceAll(describer, " ", "")
		describer, _ = strings.CutPrefix(describer, "(")
		describer, _ = strings.CutSuffix(describer, ")")

		before, after, _ := strings.Cut(describer, "...")

		if len(before) > 0 {
			h.identifierDescriber.Roles = strings.Split(before, ",")
		}

		if len(after) > 0 {
			h.identifierDescriber.Permissions = strings.Split(after, ",")
		}
	}
}

func (h *Handler[IN, OUT]) getDecorators() (decorators []string) {
	comment := funcDescription(h.Action)
	comments := strings.Split(comment, "\n")

	for _, describer := range comments {
		if !strings.HasPrefix(describer, "@") ||
			strings.HasPrefix(describer, "@authorize") {
			continue
		}

		decorator, _ := strings.CutPrefix(describer, "@")

		decorators = append(decorators, decorator)
	}

	return
}

func handleError(rctx *app.RequestContext, err error) {
	if dev {
		devHandleError(rctx, err)
	} else {
		productHandleError(rctx, err)
	}
}

func productHandleError(rctx *app.RequestContext, err error) {
	switch t := err.(type) {
	case *errors.Error:
		status, ok := abortType[t]
		if !ok {
			rctx.AbortWithStatus(http.StatusInternalServerError)

			return
		}

		t.Stack = ""
		t.Message = strings.ToUpper(t.Message)
		t.Message = strings.ReplaceAll(t.Message, " ", "_")
		rctx.AbortWithStatusJSON(status, t)
	default:
		rctx.AbortWithStatus(http.StatusInternalServerError)
	}
}

func devHandleError(rctx *app.RequestContext, err error) {
	switch t := err.(type) {
	case *errors.Error:
		status, ok := abortType[t]
		if !ok {
			rctx.AbortWithStatusJSON(http.StatusInternalServerError, t)

			return
		}

		t.Message = strings.ToUpper(t.Message)
		t.Message = strings.ReplaceAll(t.Message, " ", "_")
		rctx.AbortWithStatusJSON(status, t)
	default:
		rctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
	}
}

func bind[IN any, OUT any](handler *Handler[IN, OUT], rctx *app.RequestContext) (req IN, err error) {
	p := reflect.TypeOf(handler.Action).In(2)
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
	errors.BadRequest:                    consts.StatusBadRequest,
	errors.Unauthorized:                  consts.StatusUnauthorized,
	errors.PaymentRequired:               consts.StatusPaymentRequired,
	errors.Forbidden:                     consts.StatusForbidden,
	errors.NotFound:                      consts.StatusNotFound,
	errors.MethodNotAllowed:              consts.StatusMethodNotAllowed,
	errors.NotAcceptable:                 consts.StatusNotAcceptable,
	errors.ProxyAuthRequired:             consts.StatusProxyAuthRequired,
	errors.RequestTimeout:                consts.StatusRequestTimeout,
	errors.Conflict:                      consts.StatusConflict,
	errors.Gone:                          consts.StatusGone,
	errors.LengthRequired:                consts.StatusLengthRequired,
	errors.PreconditionFailed:            consts.StatusPreconditionFailed,
	errors.RequestEntityTooLarge:         consts.StatusRequestEntityTooLarge,
	errors.RequestURITooLong:             consts.StatusRequestURITooLong,
	errors.UnsupportedMediaType:          consts.StatusUnsupportedMediaType,
	errors.RequestedRangeNotSatisfiable:  consts.StatusRequestedRangeNotSatisfiable,
	errors.ExpectationFailed:             consts.StatusExpectationFailed,
	errors.Teapot:                        consts.StatusTeapot,
	errors.UnprocessableEntity:           consts.StatusUnprocessableEntity,
	errors.Locked:                        consts.StatusLocked,
	errors.FailedDependency:              consts.StatusFailedDependency,
	errors.UpgradeRequired:               consts.StatusUpgradeRequired,
	errors.PreconditionRequired:          consts.StatusPreconditionFailed,
	errors.TooManyRequests:               consts.StatusTooManyRequests,
	errors.RequestHeaderFieldsTooLarge:   consts.StatusRequestHeaderFieldsTooLarge,
	errors.UnavailableForLegalReasons:    consts.StatusUnavailableForLegalReasons,
	errors.NotImplemented:                consts.StatusNotImplemented,
	errors.BadGateway:                    consts.StatusBadGateway,
	errors.ServiceUnavailable:            consts.StatusServiceUnavailable,
	errors.GatewayTimeout:                consts.StatusGatewayTimeout,
	errors.HTTPVersionNotSupported:       consts.StatusHTTPVersionNotSupported,
	errors.VariantAlsoNegotiates:         consts.StatusVariantAlsoNegotiates,
	errors.InsufficientStorage:           consts.StatusInsufficientStorage,
	errors.LoopDetected:                  consts.StatusLoopDetected,
	errors.NotExtended:                   consts.StatusNotExtended,
	errors.NetworkAuthenticationRequired: consts.StatusNetworkAuthenticationRequired,
	errors.InternalServerError:           consts.StatusInternalServerError,
	errors.AlreadyExist:                  consts.StatusConflict,
}
