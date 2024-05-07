package main

import (
	"fmt"
	"text/template"
	"time"

	"github.com/maadiii/hertz/server"
)

func main() {
	server.Handle(Index)
	server.Handle(Raw)

	hertz := server.Hertz(true, server.WithHostPorts(":8080"))

	hertz.Delims("{[{", "}]}")
	hertz.SetFuncMap(template.FuncMap{
		"formatAsDate": formatAsDate,
	})
	hertz.LoadHTMLGlob("examples/render/html/*")

	hertz.Spin()
}

func formatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d/%02d/%02d", year, month, day)
}
