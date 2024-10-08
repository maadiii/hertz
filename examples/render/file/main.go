package main

import "github.com/maadiii/hertz/server"

func main() {
	hertz := server.Hertz(true, server.WithHostPorts(":8080"))
	hertz.Spin()
}
