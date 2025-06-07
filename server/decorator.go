package server

import (
	"context"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"
)

type decoratorFn func(context.Context, *Request)

var decorators = make(map[string]decoratorFn)

func AddDecorator(name string, f decoratorFn) {
	decorators[name] = f
}

func decorate(path, verb, decorator string) app.HandlerFunc {
	return func(c context.Context, rctx *app.RequestContext) {
		decorate, ok := decorators[decorator]
		if !ok {
			msg := fmt.Sprintf("%s decorator does not exist for [%s] %s", decorator, verb, path)
			panic(msg)
		}

		req := &Request{rctx}

		decorate(c, req)
	}
}
