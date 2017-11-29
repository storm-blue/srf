package main

import (
    "./rframework"
    "./controller"
)

func main() {
    server := rframework.NewServer("127.0.0.1", 8080, controller.MAPPER)
    e := server.Start()
    if e != nil {
        println(e.Error())
    }
}
