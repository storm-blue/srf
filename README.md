# simple-rest-framework
A simple framework of restful api.
With this, you can start a web application easily.

example:
```
import "fmt"
import "../srf"

var MAPPER = map[string]interface{}{

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
```

```
