package main

import "github.com/maadiii/hertz/server"

func main() {
	hertz := server.Hertz(server.WithHostPorts(":8080"))
	hertz.Spin()
}
