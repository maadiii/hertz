package server

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server/render"
)

func (h *Handler[IN, OUT]) setResponder(name string) {
	if h.setTemplateResponder() {
		return
	}

	switch h.ActionType {
	case "":
		h.RespondFn = func(ctx *app.RequestContext, _ any) { ctx.Status(h.Status) }
	case "json":
		h.RespondFn = func(ctx *app.RequestContext, res any) { ctx.JSON(h.Status, res) }
	case "json_pure":
		h.RespondFn = func(ctx *app.RequestContext, res any) { ctx.PureJSON(h.Status, res) }
	case "xml":
		h.RespondFn = func(ctx *app.RequestContext, res any) { ctx.XML(h.Status, res) }
	case "file":
		h.RespondFn = func(ctx *app.RequestContext, res any) { ctx.File(fmt.Sprintf("%v", res)) }
	case "text":
		h.setTextResponder()
	case "redirect":
		h.setRedirectResponder()
	case "attachment":
		h.setAttachmentResponder()
	case "stream":
		h.setStreamResponder()
	case "data":
		h.setDataResponder()
	case "render":
		h.setRenderResponder()
	default:
		panic(name + " action describer not acceptable")
	}
}

func (h *Handler[IN, OUT]) setTemplateResponder() bool {
	if strings.Contains(h.ActionType, "html") || strings.Contains(h.ActionType, "tmpl") {
		h.RespondFn = func(ctx *app.RequestContext, res any) {
			ctx.HTML(h.Status, h.ActionType, res)
		}

		return true
	}

	return false
}

func (h *Handler[IN, OUT]) setTextResponder() {
	h.RespondFn = func(ctx *app.RequestContext, res any) {
		_, err := ctx.WriteString(fmt.Sprintf("%s", res))
		if err != nil {
			panic(err)
		}
	}
}

func (h *Handler[IN, OUT]) setRedirectResponder() {
	h.RespondFn = func(ctx *app.RequestContext, res any) {
		ctx.Redirect(h.Status, []byte(fmt.Sprintf("%v", res)))
	}
}

func (h *Handler[IN, OUT]) setAttachmentResponder() {
	h.RespondFn = func(ctx *app.RequestContext, res any) {
		filepath := fmt.Sprintf("%v", res)
		filename := strings.Split(filepath, "/")

		ctx.SetContentType(h.ContentType)
		ctx.FileAttachment(filepath, filename[len(filename)-1])
	}
}

func (h *Handler[IN, OUT]) setStreamResponder() {
	h.RespondFn = func(ctx *app.RequestContext, res any) {
		ctx.SetContentType(h.ContentType)

		reader := bytes.NewReader(reflect.ValueOf(res).Bytes())
		if _, err := reader.WriteTo(ctx.Response.BodyWriter()); err != nil {
			panic(err)
		}
	}
}

func (h *Handler[IN, OUT]) setDataResponder() {
	h.RespondFn = func(ctx *app.RequestContext, res any) {
		ctx.SetContentType(h.ContentType)
		ctx.Data(h.Status, h.ContentType, reflect.ValueOf(res).Bytes())
	}
}

func (h *Handler[IN, OUT]) setRenderResponder() {
	h.RespondFn = func(ctx *app.RequestContext, res any) {
		ctx.Render(h.Status, res.(render.Render))
	}
}
