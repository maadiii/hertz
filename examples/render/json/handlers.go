package main

import (
	"context"

	"github.com/maadiii/hertz/errors"
	"github.com/maadiii/hertz/server"
)

func init() {
	server.Register(JSON)
	server.Register(PureJSON)
	server.Register(SomeData)
}

// @authorize(role1 ... perm1)
// @decorator
// [GET] /api/v1/json/:id 200 json
func JSON(c context.Context, _ *server.Request, in *JSONRequest) (out *JSONResponse, err error) {
	out = &JSONResponse{
		ID:       in.ID,
		Company:  "company",
		Location: "location",
		Number:   123,
	}

	return
}

type JSONRequest struct {
	ID int `path:"id"`
}

func (r *JSONRequest) Validate(*server.Request) (err error) {
	if r.ID < 1 {
		err = errors.BadRequest.Format("invalid id")
	}

	return
}

type JSONResponse struct {
	ID       int    `json:"id,omitempty"`
	Company  string `json:"company,omitempty"`
	Location string `json:"location,omitempty"`
	Number   int    `json:"number,omitempty"`
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

// [POST] /api/v1/someJSON 200 data@application/yaml; charset=utf-8
func SomeData(c context.Context, _ *server.Request, _ any) (out []byte, err error) {
	out = []byte(`{"library": "hertzwrapper", "author": "Maadi Azizi"}`)

	return
}
