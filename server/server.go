package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/middlewares/server/recovery"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/maadiii/hertz/errors"
)

func Hertz(devMode bool, opts ...config.Option) *server.Hertz {
	dev = devMode
	s = server.New(opts...)
	s.Use(recovery.Recovery())

	for i := range uses {
		s.Use(uses[i])
	}

	for relativePath, root := range static {
		s.Static(relativePath, root)
	}

	for relativePath, filePath := range staticFile {
		s.StaticFile(relativePath, filePath)
	}

	s.NoMethod(noMethodHandlers...)
	s.NoRoute(noRouteHandlers...)

	for key, handlers := range handlersMap {
		action := strings.Split(key, "::")

		s.Handle(action[0], action[1], handlers...)
	}

	return s
}

var (
	s                *server.Hertz
	dev              bool
	uses             = make([]app.HandlerFunc, 0)
	static           = make(map[string]string)
	staticFile       = make(map[string]string)
	noMethodHandlers = make([]app.HandlerFunc, 0)
	noRouteHandlers  = make([]app.HandlerFunc, 0)
	handlersMap      = make(map[string][]app.HandlerFunc, 0)
)

func Handle[IN any, OUT any](handlers ...func(*Context, IN) (OUT, error)) {
	befores := make([]*Handler[IN, OUT], 0)
	afters := make([]*Handler[IN, OUT], 0)
	main := &Handler[IN, OUT]{}

	for _, h := range handlers {
		handler := &Handler[IN, OUT]{Action: h}
		handler.fix()

		if len(handler.Path) == 0 && len(handler.Method) == 0 {
			befores = append(befores, handler)

			continue
		}

		if len(handler.Path) != 0 && len(handler.Method) != 0 {
			main = handler

			continue
		}

		afters = append(afters, handler)
	}

	key := fmt.Sprintf("%s::%s::%d::%s", main.Method, main.Path, main.Status, main.ActionType)

	for _, h := range befores {
		h.describer = main.describer
		handlersMap[key] = append(handlersMap[key], handle(h))
	}

	handlersMap[key] = append(handlersMap[key], handle(main))

	for _, h := range afters {
		h.describer = main.describer
		handlersMap[key] = append(handlersMap[key], handle(h))
	}
}

func handle[IN any, OUT any](handler *Handler[IN, OUT]) app.HandlerFunc {
	return func(c context.Context, rctx *app.RequestContext) {
		req, err := bind(handler, rctx)
		if err != nil {
			rctx.AbortWithStatusJSON(
				http.StatusUnprocessableEntity,
				errors.New(fmt.Sprintf( //nolint
					"%s\n#Api=%s#Method=%s#Action=%s",
					err.Error(),
					handler.Path,
					handler.Method,
					funcPathAndName(handler.Action),
				)),
			)

			return
		}

		ctx := &Context{
			Context: c,
			request: rctx,
		}

		res, err := handler.Action(ctx, req)
		if err != nil {
			handleError(rctx, err)

			return
		}

		handler.RespondFn(rctx, res)
	}
}

func handleError(ctx *app.RequestContext, err error) {
	if dev {
		devHandleError(ctx, err)
	} else {
		productHandleError(ctx, err)
	}
}

func productHandleError(ctx *app.RequestContext, err error) {
	switch t := err.(type) {
	case *errors.Error:
		status, ok := abortType[t]
		if !ok {
			ctx.AbortWithStatus(http.StatusInternalServerError)

			return
		}

		t.Stack = ""
		ctx.JSON(status, t)
	default:
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}
}

func devHandleError(ctx *app.RequestContext, err error) {
	switch t := err.(type) {
	case *errors.Error:
		status, ok := abortType[t]
		if !ok {
			ctx.JSON(http.StatusInternalServerError, t)

			return
		}

		ctx.JSON(status, t)
	default:
		ctx.JSON(http.StatusInternalServerError, err)
	}
}

// NoMethod sets the handlers called when the HTTP method does not match.
func NoMethod(handlers ...app.HandlerFunc) {
	noMethodHandlers = append(noMethodHandlers, handlers...)
}

// NoRoute adds handlers for NoRoute. It returns a 404 code by default.
func NoRoute(handlers ...app.HandlerFunc) {
	noRouteHandlers = append(noRouteHandlers, handlers...)
}

// Static serves files from the given file system root.
// To use the operating system's file system implementation,
// use :
//
//	router.Static("/static", "/var/www")
func Static(relativePath, root string) {
	static[relativePath] = root
}

// StaticFile registers a single route in order to Serve a single file of the local filesystem.
// router.StaticFile("favicon.ico", "./resources/favicon.ico")
func StaticFile(relativePath, filepath string) {
	staticFile[relativePath] = filepath
}

// Use attaches a global middleware to the router. ie. the middleware attached though Use() will be
// included in the handlers chain for every single request. Even 404, 405, static files...
//
// For example, this is the right place for a logger or error management middleware.
func Use(handlers ...app.HandlerFunc) {
	uses = append(uses, handlers...)
}
