/*
* the rest function metadata cache.
 */
package srf

import "reflect"

//argMeta type
const (
	IN_DATA    = iota //-> dataMeta
	IN_SESSION        //-> sessionMeta
	OUT_DATA          //-> resultMeta
	OUT_ERROR         //-> errorMeta
)

var session_type_name = "srf.Session"
var error_type_name = "error"

//function metadata main struct
type restFuncMeta struct {
	funcValue reflect.Value //the reflect value of a rest function
	inMeta    []argMeta     //in arguments metadata
	outMeta   []argMeta     //out arguments metadata
}

type argMeta struct {
	argType int
	meta    interface{}
}

type dataMeta struct {
	index    int
	isPtr    bool
	dataType reflect.Type
}

type sessionMeta struct {
	index int
}

type resultMeta struct {
	index int
	isPtr bool
}

type errorMeta struct {
	index int
}

//Get meta from func "f".
//It panics if "f" is not a valid func.
func getFuncMeta(f interface{}) *restFuncMeta {
	typeF := reflect.TypeOf(f)

	if typeF.Kind() != reflect.Func {
		panic("Mapper value must be func !")
	}

	inMeta := getInMeta(typeF)
	outMeta := getOutMeta(typeF)

	return &restFuncMeta{funcValue: reflect.ValueOf(f), inMeta: inMeta, outMeta: outMeta}
}

func getInMeta(typeF reflect.Type) (inMeta []argMeta) {
	if typeF.NumIn() == 0 {
		inMeta = make([]argMeta, 0)
	} else if typeF.NumIn() == 1 {
		inMeta = make([]argMeta, 1)
		in0 := typeF.In(0)
		switch in0.Kind() {
		case reflect.Struct:
			dm := new(dataMeta)
			dm.isPtr = false
			dm.index = 0
			dm.dataType = in0
			inMeta[0] = argMeta{argType: IN_DATA, meta: dm}
		case reflect.Ptr:
			dm := new(dataMeta)
			dm.isPtr = true
			dm.index = 0
			dm.dataType = in0.Elem()
			inMeta[0] = argMeta{argType: IN_DATA, meta: dm}
		case reflect.Interface:
			if in0.String() != session_type_name {
				wrongFunc(typeF.String())
			}
			sm := new(sessionMeta)
			sm.index = 0
			inMeta[0] = argMeta{argType: IN_SESSION, meta: sm}
		default:
			wrongFunc(typeF.String())
		}
	} else if typeF.NumIn() == 2 {
		inMeta = make([]argMeta, 2)
		in0 := typeF.In(0)
		in1 := typeF.In(1)
		if (in0.Kind() == reflect.Struct || in0.Kind() != reflect.Ptr) && (in1.Kind() == reflect.Interface && in1.String() == session_type_name) {
			if in0.Kind() == reflect.Struct {
				dm := new(dataMeta)
				dm.isPtr = false
				dm.index = 0
				dm.dataType = in0
				inMeta[0] = argMeta{argType: IN_DATA, meta: dm}
			} else if in0.Kind() == reflect.Ptr {
				dm := new(dataMeta)
				dm.isPtr = true
				dm.index = 0
				dm.dataType = in0.Elem()
				inMeta[0] = argMeta{argType: IN_DATA, meta: dm}
			}
			sm := new(sessionMeta)
			sm.index = 1
			inMeta[1] = argMeta{argType: IN_SESSION, meta: sm}
		} else if (in1.Kind() == reflect.Struct || in1.Kind() != reflect.Ptr) && (in0.Kind() == reflect.Interface && in0.String() == session_type_name) {
			if in1.Kind() == reflect.Struct {
				dm := new(dataMeta)
				dm.isPtr = false
				dm.index = 1
				dm.dataType = in1
				inMeta[1] = argMeta{argType: IN_DATA, meta: dm}
			} else if in1.Kind() == reflect.Ptr {
				dm := new(dataMeta)
				dm.isPtr = true
				dm.index = 1
				dm.dataType = in1.Elem()
				inMeta[1] = argMeta{argType: IN_DATA, meta: dm}
			}
			sm := new(sessionMeta)
			sm.index = 0
			inMeta[0] = argMeta{argType: IN_SESSION, meta: sm}
		} else {
			wrongFunc(typeF.String())
		}
	} else {
		wrongFunc(typeF.String())
	}
	return
}

func getOutMeta(typeF reflect.Type) (outMeta []argMeta) {
	if typeF.NumOut() == 0 {
		outMeta = make([]argMeta, 0)
	} else if typeF.NumOut() == 1 {
		outMeta = make([]argMeta, 1)
		out0 := typeF.Out(0)
		switch out0.Kind() {
		case reflect.Struct:
			rm := new(resultMeta)
			rm.index = 0
			rm.isPtr = false
			outMeta[0] = argMeta{argType: OUT_DATA, meta: rm}
		case reflect.Ptr:
			rm := new(resultMeta)
			rm.index = 0
			rm.isPtr = true
			outMeta[0] = argMeta{argType: OUT_DATA, meta: rm}
		case reflect.Interface:
			if out0.String() != error_type_name {
				wrongFunc(typeF.String())
			}
			em := new(errorMeta)
			em.index = 0
			outMeta[0] = argMeta{argType: OUT_ERROR, meta: em}
		default:
			wrongFunc(typeF.String())
		}
	} else if typeF.NumOut() == 2 {
		out0 := typeF.Out(0)
		out1 := typeF.Out(1)
		if (out0.Kind() == reflect.Struct || out0.Kind() != reflect.Ptr) && (out1.Kind() == reflect.Interface && out1.String() == error_type_name) {
			if out0.Kind() == reflect.Struct {
				rm := new(resultMeta)
				rm.index = 0
				rm.isPtr = false
				outMeta[0] = argMeta{argType: OUT_DATA, meta: rm}
			} else if out0.Kind() == reflect.Ptr {
				rm := new(resultMeta)
				rm.index = 0
				rm.isPtr = true
				outMeta[0] = argMeta{argType: OUT_DATA, meta: rm}
			}
			em := new(errorMeta)
			em.index = 1
			outMeta[1] = argMeta{argType: OUT_ERROR, meta: em}
		} else if (out1.Kind() == reflect.Struct || out1.Kind() != reflect.Ptr) && (out0.Kind() == reflect.Interface && out0.String() == error_type_name) {
			if out1.Kind() == reflect.Struct {
				rm := new(resultMeta)
				rm.index = 1
				rm.isPtr = false
				outMeta[1] = argMeta{argType: OUT_DATA, meta: rm}
			} else if out1.Kind() == reflect.Ptr {
				rm := new(resultMeta)
				rm.index = 1
				rm.isPtr = true
				outMeta[1] = argMeta{argType: OUT_DATA, meta: rm}
			}
			em := new(errorMeta)
			em.index = 0
			outMeta[0] = argMeta{argType: OUT_ERROR, meta: em}
		} else {
			wrongFunc(typeF.String())
		}
	} else {
		wrongFunc(typeF.String())
	}
	return
}

func wrongFunc(funcDef string) {
	panic("Wrong func definition: " + funcDef)
}
