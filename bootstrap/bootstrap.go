package main

import (
	_ "srf/router"
	"srf/srf"
)

func main() {
	server := srf.NewServer("127.0.0.1", 8080)
	e := server.Start()
	if e != nil {
		println(e.Error())
	}
}
