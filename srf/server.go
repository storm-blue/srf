package srf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

var METHODS = []string{"GET", "POST", "PUT", "DELETE"}

type Server interface {
	Start() error
}

type restServer struct {
	Bind string

	Port int

	Mapper map[string]interface{}

	/*
	   {
	       "/books" :
	       {
	           "POST" : func metadata,
	           "GET" : func metadata
	       }
	   }
	*/
	metadata map[string]map[string]*restFuncMeta
}

func getUriAndMethod(k string) (string, string) {
	var uri string
	var method string
	s := strings.Split(k, " ")
	if len(s) == 1 {
		uri = s[0]
	} else if len(s) == 2 {
		method = s[0]
		uri = s[1]
	} else {
		panic("Wrong mapper key: \"" + k + "\"")
	}
	return uri, method
}

func getMappers() map[string]interface{} {
	var mappers = make(map[string]interface{})
	for nameSpace, mapper := range Routers {
		for k, handler := range mapper {
			uri, method := getUriAndMethod(k)
			var url string
			if strings.HasPrefix(uri, RootPath) {
				url = nameSpace + uri[1:]
			} else {
				url = nameSpace + uri
			}

			if method == "" {
				mappers[url] = handler
			} else {
				mappers[method+" "+url] = handler
			}

		}
	}
	return mappers
}

func NewServer(bind string, port int) Server {

	mapper := getMappers()

	restServer := &restServer{Bind: bind, Port: port, Mapper: mapper}

	/*construct the url mapper*/
	restServer.metadata = make(map[string]map[string]*restFuncMeta)
	for k, v := range mapper {

		s := strings.Split(k, " ")

		/* process key like "/books" */
		if len(s) == 1 {
			if restServer.metadata[s[0]] == nil {
				restServer.metadata[s[0]] = make(map[string]*restFuncMeta)
			}
			for _, m := range METHODS {
				if restServer.metadata[s[0]][m] == nil {
					restServer.metadata[s[0]][m] = getFuncMeta(v)
				}
			}
		} else /* process key like "POST /books" */ if len(s) == 2 {
			if !contains(METHODS, s[0]) {
				panic("Unknown method type: " + s[0])
			}

			if restServer.metadata[s[1]] == nil {
				restServer.metadata[s[1]] = make(map[string]*restFuncMeta)
			}
			restServer.metadata[s[1]][s[0]] = getFuncMeta(v)
		}
	}
	return restServer
}

func contains(slice []string, s string) bool {
	for _, a := range slice {
		if a == s {
			return true
		}
	}
	return false
}

func (srv *restServer) Start() error {
	for k, v := range srv.metadata {
		http.HandleFunc(k, srv.getHandler(v))
	}
	return http.ListenAndServe(srv.Bind+":"+strconv.Itoa(srv.Port), nil)
}

func (srv *restServer) getHandler(m map[string]*restFuncMeta) func(http.ResponseWriter, *http.Request) {

	return func(writer http.ResponseWriter, request *http.Request) {
		request.ParseForm()
		session := obtainSession(writer, request)

		metadataf := m[request.Method]

		if metadataf == nil {
			fmt.Fprintf(writer, "Method %v not supported!", request.Method)
			return
		}

		f := metadataf.funcValue
		params := make([]reflect.Value, len(metadataf.inMeta))
		for i, v := range metadataf.inMeta {
			if v.argType == IN_DATA {
				dm := v.meta.(*dataMeta)
				params[i] = buildDataParameter(dm, request)
			} else if v.argType == IN_SESSION {
				params[i] = reflect.ValueOf(session)
			}
		}
		objectOut := f.Call(params)[0].Interface()
		bs, _ := json.Marshal(objectOut)
		writer.Write(bs)
	}
}

var goSessionKey = "gsessionId"

func obtainSession(writer http.ResponseWriter, request *http.Request) (session Session) {
	sessionId, err := request.Cookie(goSessionKey)
	if err == nil {
		session = GetSession(sessionId.Value)
		if session == nil {
			session = CreateSession()
			cookie := http.Cookie{Name: goSessionKey, Value: session.GetId(), Path: "/", HttpOnly: true}
			http.SetCookie(writer, &cookie)
		}
	} else {
		session = CreateSession()
		cookie := http.Cookie{Name: goSessionKey, Value: session.GetId(), Path: "/", HttpOnly: true}
		http.SetCookie(writer, &cookie)
	}
	return
}

func upCaseFirstLetter(f string) string {
	r := ""
	for i, c := range f {
		if i == 0 && c >= 97 && c <= 122 {
			c = c - 32
		}
		r += string(c)
	}
	return r
}

func buildDataParameter(dm *dataMeta, request *http.Request) reflect.Value {
	objectIn := reflect.New(dm.dataType)

	switch request.Method {
	case "GET":
		for k, v := range request.Form {
			field := objectIn.Elem().FieldByName(upCaseFirstLetter(k))
			fmt.Println(field.Kind())
			switch field.Kind() {
			case reflect.String:
				field.SetString(v[0])
			case reflect.Bool:
				boolValue, err := strconv.ParseBool(v[0])
				if err == nil {
					field.SetBool(boolValue)
				}
			case reflect.Int:
				fallthrough
			case reflect.Int8:
				fallthrough
			case reflect.Int16:
				fallthrough
			case reflect.Int32:
				fallthrough
			case reflect.Int64:
				intValue, err := strconv.ParseInt(v[0], 10, 64)
				if err == nil {
					field.SetInt(intValue)
				}
			case reflect.Float32:
				fallthrough
			case reflect.Float64:
				floatValue, err := strconv.ParseFloat(v[0], 64)
				if err == nil {
					field.SetFloat(floatValue)
				}
			case reflect.Uint:
				fallthrough
			case reflect.Uint8:
				fallthrough
			case reflect.Uint16:
				fallthrough
			case reflect.Uint32:
				fallthrough
			case reflect.Uint64:
				uintValue, err := strconv.ParseUint(v[0], 10, 64)
				if err == nil {
					field.SetUint(uintValue)
				}
			}
		}
	case "POST":
		fallthrough
	case "PUT":
		fallthrough
	case "DELETE":
		bs, _ := ioutil.ReadAll(request.Body)
		json.Unmarshal(bs, objectIn.Interface())
	}

	//invoke
	if !dm.isPtr {
		objectIn = objectIn.Elem()
	}
	return objectIn
}
