package server

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server/render"
)

type Handler[IN any, OUT any] struct {
	Action    func(*Context, IN) (OUT, error)
	RespondFn func(ctx *app.RequestContext, response any)

	*describer
}

type describer struct {
	Path        string
	Method      string
	Status      int
	ContentType string
	ActionType  string
}

func (h *Handler[IN, OUT]) fix() {
	name := funcPathAndName(h.Action)
	desc := h.getFixedDescriberArray()

	for i, d := range desc {
		if strings.HasPrefix(d, "/") {
			d = strings.TrimRight(d, "/")
			h.Path = d

			continue
		}

		if strings.HasPrefix(d, "[") && strings.HasSuffix(d, "]") {
			verb, ok := methods[strings.ToUpper(d)]
			if !ok {
				panic(fmt.Sprintf("%s has wrong @action describer, Invalid VERB", name))
			}

			h.Method = verb

			continue
		}

		if status, err := strconv.Atoi(d); err == nil {
			h.Status = status

			continue
		}

		if strings.Contains(d, "@") {
			typeAndContentType := strings.Split(d, "@")
			h.ActionType = typeAndContentType[0]
			h.ContentType = fmt.Sprintf("%s %s", typeAndContentType[1], desc[i+1])
			break
		} else {
			h.ActionType = d
		}
	}

	h.setResponder()
}

func (h *Handler[IN, OUT]) getFixedDescriberArray() []string {
	h.describer = new(describer)

	comment := funcDescription(h.Action)
	comments := strings.Split(comment, "\n")

	for _, actionDescription := range comments {
		if !strings.HasPrefix(actionDescription, "@action") {
			continue
		}

		desc := strings.Split(actionDescription, " ")
		index := 0

		for i := 0; i < len(desc); i++ {
			if len(desc[i]) != 0 {
				desc[index] = desc[i]
				index++
			}
		}

		desc = desc[1:index]

		return desc
	}

	return []string{}
}

func (h *Handler[IN, OUT]) setResponder() {
	if strings.Contains(h.ActionType, "html") || strings.Contains(h.ActionType, "tmpl") {
		h.RespondFn = func(ctx *app.RequestContext, res any) {
			ctx.HTML(h.Status, h.ActionType, res)
		}

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
		h.RespondFn = func(ctx *app.RequestContext, res any) {
			_, err := ctx.WriteString(fmt.Sprintf("%s", res))
			if err != nil {
				panic(err)
			}
		}
	case "redirect":
		h.RespondFn = func(ctx *app.RequestContext, res any) {
			ctx.Redirect(h.Status, []byte(fmt.Sprintf("%v", res)))
		}
	case "attachment":
		h.RespondFn = func(ctx *app.RequestContext, res any) {
			filepath := fmt.Sprintf("%v", res)
			filename := strings.Split(filepath, "/")
			ctx.FileAttachment(filepath, filename[len(filename)-1])
		}
	case "stream":
		h.RespondFn = func(ctx *app.RequestContext, res any) {
			ctx.SetContentType(h.ContentType)

			reader := bytes.NewReader(reflect.ValueOf(res).Bytes())
			if _, err := reader.WriteTo(ctx.Response.BodyWriter()); err != nil {
				panic(err)
			}
		}
	case "data":
		h.RespondFn = func(ctx *app.RequestContext, res any) {
			ctx.SetContentType(h.ContentType)

			ctx.Data(h.Status, h.ContentType, reflect.ValueOf(res).Bytes())
		}

	case "render":
		h.RespondFn = func(ctx *app.RequestContext, res any) {
			ctx.Render(h.Status, res.(render.Render))
		}
	default:
		panic("actionType in action describer not acceptable")
	}
}
