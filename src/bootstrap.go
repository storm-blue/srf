package main

import (
    "./restframework"
    "./controller"
)

func main() {
    server := restframework.NewRestServer("127.0.0.1", 8080, controller.MAPPER)
    e := server.Start()
    if e != nil {
        println(e.Error())
    }
}
