# srf
A simple framework for restful web app developers.
You can start a restful web application with the follow steps easily.

step1:
```
var BookMapper = map[string]interface{}{

    "GET /books":
    func(book Book) Response {
        return Response{Code: "000000", Message: "BOOKS OK!"}
    },

    "POST /books":
    func(book Book) Response {
        return Response{Code: "000000", Message: "BOOKS OK!"}
    },

    "/fuckers":
    func(fucker Fucker, session srf.Session) Response {
        if session.GetAttribute("fucker") != nil {
            fmt.Println("last fucker:" + session.GetAttribute("fucker").(string))
        }
        session.SetAttribute("fucker", fucker.Name)
        fmt.Println("now in controller!")
        fmt.Println(fucker)
        return Response{Code: "000000", Message: "FUCKERS OK!"}
    },
}
```
step2:
```
func main() {
    server := srf.NewServer("127.0.0.1", 8080)
    server.Register("/book", controller.BookMapper)
    e := server.Start()
    if e != nil {
        println(e.Error())
    }
}
```

