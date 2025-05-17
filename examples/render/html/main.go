package main

import (
	"fmt"
	"time"

	"github.com/maadiii/hertz/server"
)

func main() {
	server.Register(Index)

	hertz := server.Hertz(server.WithHostPorts(":8080"))

	hertz.LoadHTMLGlob("./*")

	hertz.Spin()
}

func formatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d/%02d/%02d", year, month, day)
}
