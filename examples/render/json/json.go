package main

import (
	"context"

	"github.com/maadiii/hertz/errors"
	"github.com/maadiii/hertz/server"
)

func init() {
	server.Register(JSON)
}

// @authorize(role1, role2 ::: perm1, perm2)
// @decorator
// [GET] /api/v1/json/:id 200 json
func JSON(c context.Context, req *server.Request, in *JSONRequest) (out *JSONResponse, err error) {
	if err = in.Validate(req); err != nil {
		return
	}

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
		err = errors.BadRequest
	}

	return
}

type JSONResponse struct {
	ID       int    `json:"id,omitempty"`
	Company  string `json:"company,omitempty"`
	Location string `json:"location,omitempty"`
	Number   int    `json:"number,omitempty"`
}
