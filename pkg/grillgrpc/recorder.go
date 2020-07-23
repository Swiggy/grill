package grillgrpc

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type record struct {
	stub     *Stub
	requests []interface{}
}

type recorder struct {
	mu         sync.Mutex
	recordings map[string]record
}

func newGRPCRecorder() *recorder {
	return &recorder{
		mu:         sync.Mutex{},
		recordings: map[string]record{},
	}
}

func (r *recorder) add(stub *Stub) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.recordings[stub.Request.String()]; ok {
		return fmt.Errorf("stub already exists")
	}

	r.recordings[stub.Request.String()] = record{
		stub:     stub,
		requests: []interface{}{},
	}

	return nil
}

func (r *recorder) find(str string) (record, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	rec, ok := r.recordings[str]
	return rec, ok
}

func (r *recorder) addRequest(str string, req interface{}) {
	r.mu.Lock()
	defer r.mu.Unlock()

	rec, _ := r.recordings[str]
	rec.requests = append(rec.requests, req)
	r.recordings[str] = rec
}

func (r *recorder) resetAll() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.recordings = map[string]record{}
}

func (r *recorder) count(request *Request) int {
	r.mu.Lock()
	defer r.mu.Unlock()

	rec, ok := r.recordings[request.String()]
	if !ok || len(rec.requests) == 0 {
		return 0
	}
	matched := 0
	for _, req := range rec.requests {
		copyReq := proto.Clone(req.(proto.Message))
		if request.MatchFn == nil {
			matched++
		} else {
			if request.MatchFn(copyReq) {
				matched++
			}
		}
	}

	return matched
}

func (r *recorder) unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	record, ok := r.find(info.FullMethod)
	if !ok {
		return nil, status.Errorf(codes.Unimplemented, "no stubs found for request method=%v", info.FullMethod)
	}
	stub := record.stub
	if stub.Request.MatchFn != nil {
		copyReq := proto.Clone(req.(proto.Message))
		if !stub.Request.MatchFn(copyReq) {
			return nil, status.Errorf(codes.Unimplemented, "request match failed for method=%v, request=%v", info.FullMethod, req)
		}
	}

	r.addRequest(info.FullMethod, req)

	time.Sleep(time.Millisecond * time.Duration(record.stub.Response.FixedDelayMilliseconds))

	if stub.Response.TemplateFn != nil {
		response := proto.Clone(stub.Response.Data.(proto.Message))
		stub.Response.TemplateFn(req, response)
		return response, stub.Response.Err
	}

	return stub.Response.Data, stub.Response.Err
}
