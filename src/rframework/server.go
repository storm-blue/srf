package rframework

import (
    "net/http"
    "fmt"
    "strconv"
    "reflect"
    "strings"
    "io/ioutil"
    "encoding/json"
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
            "POST /books" : func(),
            "GET /books" : func()
        }
    }
    */
    innerMapper map[string]map[string]interface{}
}

func NewServer(bind string, port int, mapper map[string]interface{}) Server {
    restServer := &restServer{Bind: bind, Port: port, Mapper: mapper}

    /*construct the url mapper*/
    restServer.innerMapper = make(map[string]map[string]interface{})
    for k, v := range mapper {

        typeV := reflect.TypeOf(v)

        if typeV.Kind() != reflect.Func {
            panic("Mapper value must be func !")
        }

        if typeV.NumIn() != 1 || typeV.NumOut() != 1 ||
            (typeV.In(0).Kind() != reflect.Ptr && typeV.In(0).Kind() != reflect.Struct) ||
            (typeV.Out(0).Kind() != reflect.Ptr && typeV.Out(0).Kind() != reflect.Struct) {
            panic("Wrong func definition: " + reflect.TypeOf(v).String())
        }

        s := strings.Split(k, " ")

        /* process key like "/books" */ if len(s) == 1 {
            if restServer.innerMapper[s[0]] == nil {
                restServer.innerMapper[s[0]] = make(map[string]interface{})
            }
            for _, m := range METHODS {
                if restServer.innerMapper[s[0]][m+" "+s[0]] == nil {
                    restServer.innerMapper[s[0]][m+" "+s[0]] = v
                }
            }
        } /* process key like "POST /books" */ else if len(s) == 2 {
            if !contains(METHODS, s[0]) {
                panic("Unknown method type: " + s[0])
            }

            if restServer.innerMapper[s[1]] == nil {
                restServer.innerMapper[s[1]] = make(map[string]interface{})
            }
            restServer.innerMapper[s[1]][k] = v
        } else {
            panic("Wrong mapper key: \"" + k + "\"")
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
    for k, v := range srv.innerMapper {
        http.HandleFunc(k, srv.getHandler(k, v))
    }
    return http.ListenAndServe(srv.Bind+":"+strconv.Itoa(srv.Port), nil)
}

func (srv *restServer) getHandler(pattern string, m map[string]interface{}) func(http.ResponseWriter, *http.Request) {

    return func(writer http.ResponseWriter, request *http.Request) {
        request.ParseForm()

        key := request.Method + " " + pattern

        fmt.Println(m)

        if m[key] == nil {
            fmt.Fprintf(writer, "Method %v not supported!", request.Method)
            return
        }

        inIsPtr := false
        if reflect.TypeOf(m[key]).In(0).Kind() == reflect.Ptr {
            inIsPtr = true
        }
        var tIn reflect.Type

        if inIsPtr {
            tIn = reflect.TypeOf(m[key]).In(0).Elem()
        } else {
            tIn = reflect.TypeOf(m[key]).In(0)
        }

        objectIn := reflect.New(tIn)

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
        if !inIsPtr {
            objectIn = objectIn.Elem()
        }
        f := reflect.ValueOf(m[key])
        params := make([]reflect.Value, 1)
        params[0] = objectIn
        objectOut := f.Call(params)[0].Interface()
        bs, _ := json.Marshal(objectOut)
        writer.Write(bs)
    }
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
