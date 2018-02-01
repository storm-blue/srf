package main

import (
	"github.com/zhangyueshan/srf/controller"
	"github.com/zhangyueshan/srf/srf"
)

func main() {
	server := srf.NewServer("127.0.0.1", 8080)
	server.Register("/book", controller.BookMapper)
	server.Register("/user", controller.UserMapper)
	e := server.Start()
	if e != nil {
		println(e.Error())
	}
}
