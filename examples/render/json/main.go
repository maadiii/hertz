package main

import "github.com/maadiii/hertz/server"

func main() {
	server.Handle(Json)
	server.Handle(PureJson)
	server.Handle(SomeData)

	hertz := server.Hertz(true, server.WithHostPorts(":8080"))
	hertz.Spin()
}
