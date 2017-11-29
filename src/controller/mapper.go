package controller

import "fmt"

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
	func(fucker Fucker) Response {
		fmt.Println("now in controller!")
		fmt.Println(fucker)
		return Response{Code: "000000", Message: "FUCKERS OK!"}
	},
}
