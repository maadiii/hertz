package main

import (
	"github.com/maadiii/hertz/errors"
	"github.com/maadiii/hertz/server"
)

func main() {
	server.SetIdentifier(identifier)

	hertz := server.Hertz(true, server.WithHostPorts(":8080"))
	hertz.Spin()
}

func identifier(c *server.Context, roles []string, permissions ...string) error {
	incomeRole := string(c.GetHeader("role"))
	if len(incomeRole) == 0 {
		return errors.Unauthorized
	}

	var matchRole bool

	for _, r := range roles {
		if incomeRole == r {
			matchRole = true
		}
	}

	if !matchRole {
		return errors.Forbidden
	}

	incomePerm := string(c.GetHeader("perm"))
	if len(incomePerm) == 0 {
		return errors.Forbidden
	}

	var matchPermission bool

	for _, p := range permissions {
		if incomePerm == p {
			matchPermission = true
		}
	}

	if !matchPermission {
		return errors.Forbidden
	}

	identity := server.Identity{"id": 1, "username": "maadi"}

	c.SetIdentity(identity)

	return nil
}
