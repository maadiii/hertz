package server

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

// SetIdentifier set authentication and authorization method
func SetIdentifier(identifierFn identifierFn) {
	identifier = identifierFn
}

func (req *Request) SetIdentity(identity Identity) {
	req.rc.Set(identityKey, identity)
}

func (req *Request) Identity() Identity {
	return req.rc.Value(identityKey).(Identity)
}

type (
	Identity     map[string]any
	identifierFn func(c context.Context, req *Request, roles []string, permissions ...string)
)

const identityKey = "identity"

var identifier identifierFn

func identify[IN any, OUT any](handler *Handler[IN, OUT]) app.HandlerFunc {
	return func(c context.Context, rctx *app.RequestContext) {
		req := &Request{rctx}

		identifier(c, req, handler.identifierDescriber.Roles, handler.identifierDescriber.Permissions...)
	}
}
