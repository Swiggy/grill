package grillgrpc

import (
	"context"
	"fmt"
	"testing"

	"github.com/Swiggy/grill/pkg/grillgrpc/hello"

	"github.com/Swiggy/grill"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	request = Request{
		Package: "hello",
		Service: "HelloAPI",
		Method:  "Hello",
	}

	response = Response{
		Data: &hello.HelloResponse{Message: "Hi!"},
		Err:  nil,
	}

	errResponse = Response{
		Data: nil,
		Err:  status.Errorf(codes.InvalidArgument, "invalid argument"),
	}

	templateResponse = Response{
		Data: &hello.HelloResponse{Message: "Hi!"},
		Err:  nil,
		TemplateFn: func(request interface{}, response interface{}) {
			req, _ := request.(*hello.HelloRequest)
			res, _ := response.(*hello.HelloResponse)
			if req.Message == "hello" {
				res.Message = req.Message
			}
		},
	}

	requestMatchFn = Request{
		Package: "hello",
		Service: "HelloAPI",
		Method:  "Hello",
		MatchFn: func(request interface{}) bool {
			req := request.(*hello.HelloRequest)
			return req.Message == "namastey"
		},
	}
)

type codeAssertion struct {
	got      error
	expected codes.Code
}

func (c *codeAssertion) Assert() error {
	got := status.Code(c.got)
	if got != c.expected {
		return fmt.Errorf("invalid grpc code, got=%v, want=%v", got, c.expected)
	}
	return nil
}

func (c *codeAssertion) SetOutput(output interface{}) {
	if err, ok := output.(error); ok {
		c.got = err
	}
}

func Test_GrillGRPC(t *testing.T) {
	helper := &GRPC{}
	helper.RegisterServices(func(server *grpc.Server) {
		hello.RegisterHelloAPIServer(server, &hello.UnimplementedHelloAPIServer{})
	})

	if err := helper.Start(context.TODO()); err != nil {
		t.Errorf("error starting grpc grill, error=%v", err)
		return
	}

	conn, err := grpc.Dial(fmt.Sprintf("localhost:%s", helper.Port()), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		t.Errorf("error connecting to grpc grill, error=%v", err)
		return
	}
	fmt.Println("connected")
	client := hello.NewHelloAPIClient(conn)

	tests := []grill.TestCase{
		{
			Name:  "Test_NoStubs",
			Stubs: []grill.Stub{},
			Action: func() interface{} {
				_, err := client.Hello(context.Background(), &hello.HelloRequest{Message: "Hi"})
				return err
			},
			Assertions: []grill.Assertion{
				&codeAssertion{expected: codes.Unimplemented},
			},
			Cleaners: []grill.Cleaner{
				helper.ResetAllStubs(),
			},
		},
		{
			Name: "Test_ReturnsStubResponse",
			Stubs: []grill.Stub{
				helper.Stub(request, response),
			},
			Action: func() interface{} {
				res, err := client.Hello(context.Background(), &hello.HelloRequest{Message: "Hi"})
				return grill.ActionOutput(res.Message, err)
			},
			Assertions: []grill.Assertion{
				grill.AssertOutput("Hi!", nil),
				&codeAssertion{expected: codes.OK},
			},
			Cleaners: []grill.Cleaner{
				helper.ResetAllStubs(),
			},
		},
		{
			Name: "Test_ResponseTemplateTest",
			Stubs: []grill.Stub{
				helper.Stub(request, templateResponse),
			},
			Action: func() interface{} {
				res1, err := client.Hello(context.Background(), &hello.HelloRequest{Message: "hello"})
				res2, err := client.Hello(context.Background(), &hello.HelloRequest{Message: "noHello"})
				return grill.ActionOutput(res1.Message, res2.Message, err)
			},
			Assertions: []grill.Assertion{
				grill.AssertOutput("hello", "Hi!", nil),
				&codeAssertion{expected: codes.OK},
			},
			Cleaners: []grill.Cleaner{
				helper.ResetAllStubs(),
			},
		},
		{
			Name: "Test_RequestMatchTest-Failure",
			Stubs: []grill.Stub{
				helper.Stub(requestMatchFn, response),
			},
			Action: func() interface{} {
				_, err := client.Hello(context.Background(), &hello.HelloRequest{Message: "hello"})
				return err
			},
			Assertions: []grill.Assertion{
				&codeAssertion{expected: codes.Unimplemented},
			},
			Cleaners: []grill.Cleaner{
				helper.ResetAllStubs(),
			},
		},
		{
			Name: "Test_RequestMatchTest-Success",
			Stubs: []grill.Stub{
				helper.Stub(requestMatchFn, response),
			},
			Action: func() interface{} {
				res1, err := client.Hello(context.Background(), &hello.HelloRequest{Message: "namastey"})
				res2, err := client.Hello(context.Background(), &hello.HelloRequest{Message: "namastey"})
				_, _ = client.Hello(context.Background(), &hello.HelloRequest{Message: "hello"})
				return grill.ActionOutput(res1.Message, res2.Message, err)
			},
			Assertions: []grill.Assertion{
				grill.AssertOutput("Hi!", "Hi!", nil),
				&codeAssertion{expected: codes.OK},
				helper.AssertCount(requestMatchFn, 2),
			},
			Cleaners: []grill.Cleaner{
				helper.ResetAllStubs(),
			},
		},
	}

	grill.Run(t, tests)
}
