package main

import (
	"context"

	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/maadiii/hertz/server"
	"gopkg.in/yaml.v3"
)

func main() {
	server.Register(Yaml)

	hertz := server.Hertz(server.WithHostPorts(":8080"))
	hertz.Spin()
}

// response type must implement github.com/cloudwego/hertz/pkg/app/server/render/render.Render interface
//
// [GET] /api/v1/yaml 200 render
func Yaml(_ context.Context, _ *server.Request, _ any) (out YAML, err error) {
	out = YAML{Data: "some yaml data"}

	return
}

type YAML struct {
	Data interface{}
}

func (r YAML) Render(resp *protocol.Response) error {
	r.WriteContentType(resp)

	yamlBytes, err := yaml.Marshal(r.Data)
	if err != nil {
		return err
	}

	resp.AppendBody(yamlBytes)

	return nil
}

func (r YAML) WriteContentType(resp *protocol.Response) {
	resp.Header.SetContentType("application/yaml; charset=utf-8")
}
