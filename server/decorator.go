package server

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

type decoratorFn func(context.Context, *Request)

var decorators = make(map[string]decoratorFn)

func AddDecorator(name string, f decoratorFn) {
	decorators[name] = f
}

func decorate(handlerName, decoratorName string) app.HandlerFunc {
	return func(c context.Context, rctx *app.RequestContext) {
		decorator, ok := decorators[decoratorName]
		if !ok {
			panic("decorator not exist for " + handlerName)
		}

		req := &Request{rctx}

		decorator(c, req)
	}
}
