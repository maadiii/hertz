package main

import (
	"context"

	"github.com/maadiii/hertz/server"
)

func init() {
	server.Register(Stream)
}

// [GET] /stream 200 stream@application/octet-stream
func Stream(_ context.Context, _ *server.Request, _ any) (out []byte, err error) {
	data := `address: localhost
port: 5432
`

	out = []byte(data)

	return
}
