package server

import (
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/middlewares/server/recovery"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/config"
)

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
		verbAndPath := strings.Split(key, "::")

		s.Handle(verbAndPath[0], verbAndPath[1], handlers...)
	}

	return s
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
