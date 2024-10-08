package main

import (
	"context"

	"github.com/maadiii/hertz/server"
)

func init() {
	server.Register(XML)
}

// [GET] /api/v1/xml/:id 200 xml
func XML(ctx context.Context, req *server.Request, in *XMLRequest) (out *XMLResponse, err error) {
	out = &XMLResponse{
		ID:     in.ID,
		Name:   "Maadi",
		Family: "Azizi",
	}

	return
}

type XMLRequest struct {
	ID string `path:"id"`
}

type XMLResponse struct {
	ID     string
	Name   string
	Family string
}
