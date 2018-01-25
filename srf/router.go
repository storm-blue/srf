package srf

var Routers = make(map[string]map[string]interface{})

const RootPath string = "/"

func Register(nameSpace string, mapper map[string]interface{}) {
	if nameSpace == "" {
		nameSpace = RootPath
	}

	end := string(nameSpace[len(nameSpace)-1])

	if end != RootPath {
		nameSpace = nameSpace + RootPath
	}

	if v, ok := Routers[nameSpace]; ok {
		merge(v, mapper)
	} else {
		Routers[nameSpace] = mapper
	}

}

func merge(target, src map[string]interface{}) {
	for k, v := range src {
		if _, ok := target[k]; ok {
			print("The same URL occurs", k)
		} else {
			target[k] = v
		}
	}
}
