package main

import (
	"context"
	"time"

	"github.com/maadiii/hertz/server"
)

// [GET] /index/:title 200 index.tmpl
func Index(_ context.Context, _ *server.Request, in *IndexRequest) (out *IndexResponse, err error) {
	out = &IndexResponse{
		Title: in.Title,
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
	Title string
}

// [GET] /raw 200 template1.html
func Raw(_ context.Context, _ *server.Request, _ any) (out *RawResponse, err error) {
	out = &RawResponse{
		Now: time.Date(2017, 0o7, 0, 0, 0, 0, 0, time.UTC),
	}

	return
}

type RawResponse struct {
	Now time.Time
}
