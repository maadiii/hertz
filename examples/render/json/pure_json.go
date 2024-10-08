package main

import (
	"context"

	"github.com/maadiii/hertz/server"
)

func init() {
	// server.Register(PureJSON)
}

// [GET] /api/v1/pureJSON 200 json_pure
func PureJSON(c context.Context, _ *server.Request, _ any) (out *PureJSONRespone, err error) {
	out = &PureJSONRespone{
		HTML: "<p> Hello World </p>",
	}

	return
}

type PureJSONRespone struct {
	HTML string `json:"html,omitempty"`
}
