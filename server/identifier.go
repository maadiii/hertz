package server

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/maadiii/hertz/errors"
)

// SetIdentifier set authentication and authorization method
func SetIdentifier(identifierFn identifierFn) {
	identifier = identifierFn
}

func (ctx *Context) SetIdentity(identity Identity) {
	ctx.rc.Set(identityKey, identity)
}

func (ctx *Context) Identity() Identity {
	return ctx.rc.Value(identityKey).(Identity)
}

type (
	Identity     map[string]any
	identifierFn func(ctx *Context, roles []string, permissions ...string) error
)

const identityKey = "identity"

var identifier identifierFn

func identify[IN Request, OUT any](handler *Handler[IN, OUT]) app.HandlerFunc {
	return func(c context.Context, reqCtx *app.RequestContext) {
		ctx := &Context{Context: c, rc: reqCtx}

		err := identifier(ctx, handler.identifierDescriber.Roles, handler.identifierDescriber.Permissions...)
		if err != nil {
			switch err {
			case errors.Unauthorized:
				reqCtx.AbortWithStatus(consts.StatusUnauthorized)
			case errors.Forbidden:
				reqCtx.AbortWithStatus(consts.StatusForbidden)
			default:
				reqCtx.AbortWithStatus(consts.StatusInternalServerError)
			}
		}
	}
}
