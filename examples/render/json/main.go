package main

import (
	"context"
	"log"

	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/maadiii/hertz/server"
)

func main() {
	server.SetIdentifier(identify)
	server.AddDecorator("decorator", decorator)

	hertz := server.Hertz(server.WithHostPorts(":8080"))
	hertz.Spin()
}

func identify(_ context.Context, req *server.Request, roles []string, permissions ...string) {
	if len(req.GetHeader("authorize")) == 0 {
		req.AbortWithStatus(consts.StatusUnauthorized)

		return
	}

	if len(roles) != 0 {
		var matchRole bool

		incomeRole := string(req.GetHeader("role"))
		if len(incomeRole) == 0 {
			req.AbortWithStatus(consts.StatusUnauthorized)

			return
		}

		for _, r := range roles {
			if incomeRole == r {
				matchRole = true
			}
		}

		if !matchRole {
			req.AbortWithStatus(consts.StatusForbidden)

			return
		}
	}

	if len(permissions) != 0 {
		var matchPermission bool

		incomePerm := string(req.GetHeader("perm"))
		if len(incomePerm) == 0 {
			req.AbortWithStatus(consts.StatusForbidden)

			return
		}

		for _, p := range permissions {
			if incomePerm == p {
				matchPermission = true
			}
		}

		if !matchPermission {
			req.AbortWithStatus(consts.StatusForbidden)

			return
		}
	}

	identity := server.Identity{
		ID: "1",
	}

	req.SetIdentity(identity)
}

func decorator(c context.Context, req *server.Request) {
	log.Println("It's just decorator", "BEFORE")

	req.Next(c)

	log.Println("It's just decorator", "AFTER")
}
