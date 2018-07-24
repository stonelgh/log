package log

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	. "github.com/stonelgh/log/error"
)

const (
	ReqTypeSetProbe = iota
	ReqTypeGetProbes
	ReqTypeGetPackages
	ReqTypeGetModules
	ReqTypeGetTags
)

type ReqHeader struct {
	Type int
}

type Request struct {
	Head ReqHeader
	Body interface{}
	body []byte
}

type ReqGetProbes struct {
	Tags     []string
	Packages []string
	Modules  []string
	Names    []string
	Level    int
}

type ReqSetProbe struct {
	Tags      []string
	Packages  []string
	Modules   []string
	Names     []string
	Level     int
	Mode      int
	Condition Enabler
}

type RespHeader struct {
	Code int
	Msg  string
}

type Response struct {
	Head RespHeader
	Body interface{}
}

func NewRespFromError(err *Error) *Response {
	r := &Response{}
	r.FromError(err)
	return r
}

func (r *Response) FromError(err *Error) {
	if err != nil {
		r.Head.Code = err.Code()
		r.Head.Msg = err.Error()
	}
}

type Handler struct {
	Type    int
	ReqBody reflect.Type
	//RespBody reflect.Type
	Handler func(req *Request) (*Response, *Error)
}

var handlers = make(map[int]*Handler)

func init() {
	handlers[ReqTypeSetProbe] = &Handler{ReqTypeSetProbe, reflect.TypeOf(ReqSetProbe{}), handleSetProbe}
	handlers[ReqTypeGetProbes] = &Handler{ReqTypeGetProbes, reflect.TypeOf(ReqGetProbes{}), handleGetProbes}
}

func handleRequest(head, body []byte) (*Response, *Error) {
	req := &Request{body: body}
	if e := json.Unmarshal(head, &req.Head); e != nil {
		e1 := NewError(ErrFormat, "Failed to parse head: "+e.Error())
		e1.SetDetail(string(head))
		fmt.Println(e1.Error())
		return nil, e1
	}

	handler := handlers[req.Head.Type]
	if handler == nil {
		e := NewError(ErrUndefined, "Undefined request type "+strconv.Itoa(req.Head.Type))
		fmt.Println(e.Error())
		return nil, e
	}
	if handler.ReqBody != nil {
		req.Body = reflect.New(handler.ReqBody).Interface()
		if e := json.Unmarshal(body, req.Body); e != nil {
			e1 := NewError(ErrFormat, "Failed to parse body: "+e.Error())
			e1.SetDetail(string(req.body))
			fmt.Println(e1.Error())
			return nil, e1
		}
	}
	return handler.Handler(req)
}

func handleGetProbes(req *Request) (*Response, *Error) {
	resp := &Response{}
	body := req.Body.(*ReqGetProbes)

	// TODO: handle tags/packages/modules
	_ = body

	names := []string{}
	for k := range allProbes {
		names = append(names, k)
	}
	resp.Body = names
	return resp, nil
}

func handleSetProbe(req *Request) (*Response, *Error) {
	global := true
	resp := &Response{}
	body := req.Body.(*ReqSetProbe)

	// TODO: handle tags/packages/modules

	// handle Names
	for _, name := range body.Names {
		if body.Level > 0 {
			SetProbeLvlByName(name, body.Level)
			fmt.Println("Set probe level", name, body.Level)
		}
		global = false
	}
	if global {
		globalLevel = body.Level
	}
	return resp, nil
}
