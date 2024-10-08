package main

import (
	"context"

	"github.com/maadiii/hertz/server"
)

func init() {
	server.Register(File)
	server.Register(Attachment)
}

// [GET] /file 200 file
func File(_ context.Context, _ *server.Request, _ any) (out string, err error) {
	out = "../json/json.go"

	return
}

// [GET] /attach 200 attachment
func Attachment(_ context.Context, _ *server.Request, _ any) (out string, err error) {
	out = "../json/json.go"

	return
}
