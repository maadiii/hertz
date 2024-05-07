package main

import (
	"time"

	"github.com/maadiii/hertz/server"
)

// @action /index/:title [GET] 200 index.tmpl
func Index(ctx *server.Context, in *IndexRequest) (out *IndexResponse, err error) {
	out = &IndexResponse{
		Title: in.Title,
	}

	return
}

type IndexRequest struct {
	Title string `path:"title"`
}

type IndexResponse struct {
	Title string
}

// @action /raw [GET] 200 template1.html
func Raw(ctx *server.Context, _ any) (out *RawResponse, err error) {
	out = &RawResponse{
		Now: time.Date(2017, 0o7, 0, 0, 0, 0, 0, time.UTC),
	}

	return
}

type RawResponse struct {
	Now time.Time
}
