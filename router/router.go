package router

import (
	"srf/controller"
	"srf/srf"
)

func init() {
	srf.Register("/", controller.BookMapper)
	srf.Register("/user/", controller.UserMapper)
}
