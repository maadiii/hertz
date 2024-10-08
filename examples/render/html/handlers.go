package main

import (
	"context"

	"github.com/maadiii/hertz/server"
)

// [GET] /index/:title 200 index.html
func Index(_ context.Context, _ *server.Request, in *IndexRequest) (out *IndexResponse, err error) {
	out = &IndexResponse{
		Title:  in.Title,
		Name:   "Maadi",
		Family: "Azizi",
	}

	return
}

type IndexRequest struct {
	Title string `path:"title"`
}

func (i *IndexRequest) Validate(*server.Request) error {
	return nil
}

type IndexResponse struct {
	Title  string
	Name   string
	Family string
}
