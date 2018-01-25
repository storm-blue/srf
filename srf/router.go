package srf

import "strings"

const RootPath = "/"

var (
	METHODS = []string{"GET", "POST", "PUT", "DELETE"}

	//mark => handler
	Routers = make(map[string]interface{})
)

func (_ *restServer) Register(nameSpace string, mapper map[string]interface{}) {
	nameSpace = fmtNameSpace(nameSpace)

	for mark, handler := range mapper {
		uri, method := getUriAndMethod(mark)
		if method == "" {
			for _, v := range METHODS {
				register(v, nameSpace, uri, handler)
			}
		} else {
			register(method, nameSpace, uri, handler)
		}
	}
}

func register(method, nameSpace, uri string, handler interface{}) {
	mark := buildMark(method, nameSpace, uri)
	if _, ok := Routers[mark]; ok {
		panic("Repeated mark: \"" + mark + "\"")
	}
	Routers[mark] = handler
}

func buildMark(method, nameSpace, uri string) string {
	return method + " " + nameSpace + uri
}

func getUriAndMethod(mark string) (string, string) {
	var uri string
	var method string
	s := strings.Split(mark, " ")
	if len(s) == 1 {
		uri = s[0]
	} else if len(s) == 2 {
		method = s[0]
		uri = s[1]
	} else {
		panic("Error mark: \"" + mark + "\"")
	}
	return fmtUri(uri), method
}

func fmtNameSpace(nameSpace string) string {
	if nameSpace == "" {
		nameSpace = RootPath
	}

	if !strings.HasPrefix(nameSpace, RootPath) {
		nameSpace = RootPath + nameSpace
	}

	if !strings.HasSuffix(nameSpace, RootPath) {
		nameSpace = nameSpace + RootPath
	}
	return nameSpace
}

func fmtUri(uri string) string {
	if strings.HasPrefix(uri, RootPath) {
		uri = uri[1:]
	}
	return uri
}
