// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: tasks/tasks.proto

package tasks

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	client "github.com/micro/go-micro/client"
	server "github.com/micro/go-micro/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for ToDoService service

type ToDoService interface {
	// Create new todo task
	Create(ctx context.Context, in *CreateRequest, opts ...client.CallOption) (*CreateResponse, error)
	// Read todo task
	Read(ctx context.Context, in *ReadRequest, opts ...client.CallOption) (*ReadResponse, error)
	// Update todo task
	Update(ctx context.Context, in *UpdateRequest, opts ...client.CallOption) (*UpdateResponse, error)
	// Read all todo tasks
	ReadAll(ctx context.Context, in *ReadAllRequest, opts ...client.CallOption) (*ReadAllResponse, error)
}

type toDoService struct {
	c    client.Client
	name string
}

func NewToDoService(name string, c client.Client) ToDoService {
	if c == nil {
		c = client.NewClient()
	}
	if len(name) == 0 {
		name = "tasks"
	}
	return &toDoService{
		c:    c,
		name: name,
	}
}

func (c *toDoService) Create(ctx context.Context, in *CreateRequest, opts ...client.CallOption) (*CreateResponse, error) {
	req := c.c.NewRequest(c.name, "ToDoService.Create", in)
	out := new(CreateResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *toDoService) Read(ctx context.Context, in *ReadRequest, opts ...client.CallOption) (*ReadResponse, error) {
	req := c.c.NewRequest(c.name, "ToDoService.Read", in)
	out := new(ReadResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *toDoService) Update(ctx context.Context, in *UpdateRequest, opts ...client.CallOption) (*UpdateResponse, error) {
	req := c.c.NewRequest(c.name, "ToDoService.Update", in)
	out := new(UpdateResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *toDoService) ReadAll(ctx context.Context, in *ReadAllRequest, opts ...client.CallOption) (*ReadAllResponse, error) {
	req := c.c.NewRequest(c.name, "ToDoService.ReadAll", in)
	out := new(ReadAllResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for ToDoService service

type ToDoServiceHandler interface {
	// Create new todo task
	Create(context.Context, *CreateRequest, *CreateResponse) error
	// Read todo task
	Read(context.Context, *ReadRequest, *ReadResponse) error
	// Update todo task
	Update(context.Context, *UpdateRequest, *UpdateResponse) error
	// Read all todo tasks
	ReadAll(context.Context, *ReadAllRequest, *ReadAllResponse) error
}

func RegisterToDoServiceHandler(s server.Server, hdlr ToDoServiceHandler, opts ...server.HandlerOption) error {
	type toDoService interface {
		Create(ctx context.Context, in *CreateRequest, out *CreateResponse) error
		Read(ctx context.Context, in *ReadRequest, out *ReadResponse) error
		Update(ctx context.Context, in *UpdateRequest, out *UpdateResponse) error
		ReadAll(ctx context.Context, in *ReadAllRequest, out *ReadAllResponse) error
	}
	type ToDoService struct {
		toDoService
	}
	h := &toDoServiceHandler{hdlr}
	return s.Handle(s.NewHandler(&ToDoService{h}, opts...))
}

type toDoServiceHandler struct {
	ToDoServiceHandler
}

func (h *toDoServiceHandler) Create(ctx context.Context, in *CreateRequest, out *CreateResponse) error {
	return h.ToDoServiceHandler.Create(ctx, in, out)
}

func (h *toDoServiceHandler) Read(ctx context.Context, in *ReadRequest, out *ReadResponse) error {
	return h.ToDoServiceHandler.Read(ctx, in, out)
}

func (h *toDoServiceHandler) Update(ctx context.Context, in *UpdateRequest, out *UpdateResponse) error {
	return h.ToDoServiceHandler.Update(ctx, in, out)
}

func (h *toDoServiceHandler) ReadAll(ctx context.Context, in *ReadAllRequest, out *ReadAllResponse) error {
	return h.ToDoServiceHandler.ReadAll(ctx, in, out)
}