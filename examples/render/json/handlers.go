package main

import (
	"github.com/maadiii/hertz/errors"
	"github.com/maadiii/hertz/server"
)

// @action /api/v1/json/:id [GET] 200 json
func Json(ctx *server.Context, in *JsonRequest) (out *JsonResponse, err error) {
	out = &JsonResponse{
		ID:       in.ID,
		Company:  "company",
		Location: "location",
		Number:   123,
	}

	return nil, errors.NotFound.Format("user sldkjf")
}

type JsonRequest struct {
	ID string `path:"id"`
}

type JsonResponse struct {
	ID       string `json:"id,omitempty"`
	Company  string `json:"company,omitempty"`
	Location string `json:"location,omitempty"`
	Number   int    `json:"number,omitempty"`
}

// @action /api/v1/pureJson [GET] 200 json_pure
func PureJson(ctx *server.Context, _ any) (out *PureJsonRespone, err error) {
	out = &PureJsonRespone{
		Html: "<p> Hello World </p>",
	}

	return
}

type PureJsonRespone struct {
	Html string `json:"html,omitempty"`
}

// @action /api/v1/someJson [GET] 200 data@application/yaml; charset=utf-8
func SomeData(ctx *server.Context, _ any) (out []byte, err error) {
	out = []byte(`{"library": "hertzwrapper", "author": "Maadi Azizi"}`)

	return
}
