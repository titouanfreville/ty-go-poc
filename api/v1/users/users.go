package users

import (
	"context"
	"github.com/micro/go-micro"
	"go_poc/core"
	"go_poc/module/user"
)

const (
	DefaultClientName = "cli"
	DefaultServerName = "grpc"
)

type UserHandler struct {
	dbStore *core.DBStore
}

//  ---------- INITIALIZER ----------------
// NewHandler initialize the handler for TaskHandler processing
func NewHandler(connection core.DbConnection) UserHandler {
	taskServ := UserHandler{dbStore: &core.DBStore{}}
	taskServ.dbStore.InitConnection(
		connection.User, connection.Database,
		connection.Password, connection.Host,
		connection.Port,
	)
	return taskServ
}

// NewClient init a task client
func NewClient(name *string, server *string) UserService {
	if name == nil {
		cname := DefaultClientName
		name = &cname
	}
	if server == nil {
		sname := DefaultServerName
		server = &sname
	}
	service := micro.NewService(
		micro.Name("task."+*name),
		micro.WrapHandler(core.EnforceVersion),
	)
	service.Init()
	return NewUserService(*server, service.Client())
}

//  ---------- FEATURE ----------------
func (u UserHandler) Create(ctx context.Context, req *CreateRequest, rsp *CreateResponse) error {
	uM := req.User.toModel()

	err := user.Insert(&uM, u.dbStore.Db)
	if err != nil {
		return err
	}
	rsp.Id = uM.Id
	rsp.Api = "v1"
	return nil
}

func (u UserHandler) Update(ctx context.Context, req *UpdateRequest, rsp *UpdateResponse) error {
	uM := req.User.toModel()

	// err := user.Update(&uM, u.dbStore.Db)
	//	// if err != nil {
	//	// 	return err
	//	// }

	rsp.User = modelToProto(&uM)
	rsp.Api = "v1"
	return nil
}

func (u UserHandler) Read(ctx context.Context, req *ReadRequest, rsp *ReadResponse) error {
	rsp.User = modelToProto(user.GetOne(req.Id, u.dbStore.Db))
	rsp.Api = "v1"
	return nil
}

func (u UserHandler) ReadAll(ctx context.Context, req *ReadAllRequest, rsp *ReadAllResponse) error {
	rsp.Users = modelToProtoList(user.GetAll(u.dbStore.Db))
	rsp.Api = "v1"
	return nil
}

// ---------- CONVERTERS ----------------
func (m User) toModel() user.User {
	return user.User{
		Id:    m.Id,
		Name:  m.Name,
		Email: m.Email,
	}
}

func modelToProto(u *user.User) *User {
	return &User{
		Id:    u.Id,
		Name:  u.Name,
		Email: u.Email,
	}
}

func modelToProtoList(ul *user.UserList) []*User {
	var res []*User
	for _, u := range *ul {
		proto := modelToProto(u)
		res = append(res, proto)
	}
	return res
}
