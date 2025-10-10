package server

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/go-playground/validator/v10"
)

func Register[IN any, OUT any](action func(context.Context, *Request, IN) (OUT, error)) {
	handler := &Handler[IN, OUT]{HandlerFn: action}
	handler.fixAPIDescriber()
	handler.fixIdentifierDesciber()

	key := fmt.Sprintf("%s::%s::%d::%s", handler.Verb, handler.Path, handler.Status, handler.ResponderType)

	if handler.identifierDescriber != nil {
		handlersMap[key] = append(handlersMap[key], identify(handler))
	}

	decorators := handler.getDecorators()
	for _, dec := range decorators {
		handlersMap[key] = append(handlersMap[key], decorate(handler.Path, handler.Verb, dec))
	}

	handlersMap[key] = append(handlersMap[key], register(handler))
}

func register[IN any, OUT any](handler *Handler[IN, OUT]) app.HandlerFunc {
	return func(c context.Context, r *app.RequestContext) {
		reqType, err := bind(handler, r)
		if err != nil {
			_ = r.Error(r.AbortWithError(http.StatusUnprocessableEntity, err))

			return
		}

		if err := validate.Struct(reqType); err != nil {
			_, ok := err.(validator.ValidationErrors)
			if ok || err.(*validator.InvalidValidationError).Type != nil {
				_ = r.Error(r.AbortWithError(http.StatusBadRequest, err))
				handleError(c, r, err)

				return
			}
		}

		req := &Request{r}

		res, err := handler.HandlerFn(c, req, reqType)
		if err != nil {
			if handleError != nil {
				handleError(c, r, err)

				return
			}

			r.AbortWithError(500, err) //nolint

			return
		}

		handler.RespondFn(r, res)
	}
}

type Handler[IN any, OUT any] struct {
	HandlerFn func(context.Context, *Request, IN) (OUT, error)
	RespondFn func(rctx *app.RequestContext, response any)

	*apiDescriber
	*identifierDescriber
}

func (h *Handler[IN, OUT]) fixAPIDescriber() {
	comment := funcDescription(h.HandlerFn)
	comments := strings.Split(comment, "\n")
	functionName := runtimeFunc(h.HandlerFn).Name()
	apiDescriber := h.getFixedAPIDescriberFields(functionName, comments)

	if len(apiDescriber) == 0 {
		panic(functionName + " has not describer")
	}

	for _, d := range apiDescriber {
		if strings.HasPrefix(d, "[") && strings.HasSuffix(d, "]") {
			verb := strings.Replace(d, "[", "", 1)
			verb = strings.Replace(verb, "]", "", 1)

			h.Verb = verb

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
			h.ResponderType = typeAndContentType[0]
			h.ContentType = typeAndContentType[1]

			break
		}

		h.ResponderType = d
	}

	h.setResponder(functionName)
}

func (h *Handler[IN, OUT]) getFixedAPIDescriberFields(functionName string, comments []string) []string {
	h.apiDescriber = new(apiDescriber)
	h.FunctionName = functionName

	for _, describer := range comments {
		if !strings.HasPrefix(describer, "[") {
			continue
		}

		desc := strings.Split(describer, " ")

		return desc
	}

	return []string{}
}

func (h *Handler[IN, OUT]) fixIdentifierDesciber() {
	comment := funcDescription(h.HandlerFn)
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

		before, after, _ := strings.Cut(describer, ":::")

		if len(before) > 0 {
			h.identifierDescriber.Roles = strings.Split(before, ",") //nolint
		}

		if len(after) > 0 {
			h.identifierDescriber.Permissions = strings.Split(after, ",") //nolint
		}
	}
}

func (h *Handler[IN, OUT]) getDecorators() (decorators []string) {
	comment := funcDescription(h.HandlerFn)
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

func bind[IN any, OUT any](handler *Handler[IN, OUT], rctx *app.RequestContext) (req IN, err error) {
	p := reflect.TypeOf(handler.HandlerFn).In(2)
	if p.Kind() == reflect.Interface {
		return
	}

	req = reflect.New(p.Elem()).Interface().(IN)

	err = rctx.Bind(req)

	return
}

var validate = validator.New()

type apiDescriber struct {
	FunctionName  string
	Path          string
	Verb          string
	Status        int
	ContentType   string
	ResponderType string
}

type identifierDescriber struct {
	Roles       []string
	Permissions []string
}

type ErrorHandler func(c context.Context, rctx *app.RequestContext, err error)
