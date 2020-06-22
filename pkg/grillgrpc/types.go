package grillgrpc

import "fmt"

type Request struct {
	Package string
	Service string
	Method  string
	MatchFn func(request interface{}) bool
}

func (r *Request) String() string {
	return fmt.Sprintf("/%s.%s/%s", r.Package, r.Service, r.Method)
}

type Response struct {
	// Data should be a pointer as ProtoMessage Method is implemented on Pointer object.
	Data interface{}

	Err error

	FixedDelayMilliseconds int

	// Modify the Response Data with values from the request
	TemplateFn func(request interface{}, response interface{})
}

type Stub struct {
	Request  Request
	Response Response
}
