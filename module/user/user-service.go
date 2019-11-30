package user

import (
	"context"
	"go_poc/core"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
)

const (
	// apiVersion is version of API is provided by server
	apiVersion = "v1"
)

// userServiceServer is implementation of v1.userServiceServer proto interface
type userServiceServer struct {
	dbStore *core.DBStore
}

func (m User) ToModel() UserModel {
	return UserModel{
		Id:    m.Id,
		Name:  m.Name,
		Email: m.Email,
	}
}

func ModelToUser(u *UserModel) User {
	return User{
		Id:    u.Id,
		Name:  u.Name,
		Email: u.Email,
	}
}

func ModelToUserList(ul *UserModelList) []*User {
	res := []*User{}
	for _, u := range *ul {
		proto := ModelToUser(u)
		res = append(res, &proto)
	}
	return res
}

func NewUserServiceServer(dbSettings core.DbConnection) UserServiceServer {
	dbStore := core.DBStore{}
	dbStore.InitConnection(dbSettings.User, dbSettings.Database, dbSettings.Password, dbSettings.Host, dbSettings.Port)
	return &userServiceServer{dbStore: &dbStore}
}

// checkAPI checks if the API version requested by client is supported by server
func (s *userServiceServer) checkAPI(api string) error {
	// API version is "" means use current version of the service
	if len(api) > 0 {
		if apiVersion != api {
			return status.Errorf(codes.Unimplemented,
				"unsupported API version: service implements API version '%s', but asked for '%s'", apiVersion, api)
		}
	}
	return nil
}

// Create new todo task
func (s *userServiceServer) Create(ctx context.Context, req *CreateRequest) (*CreateResponse, error) {
	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	usr := req.User.ToModel()

	// insert ToDo entity data
	err := Insert(&usr, s.dbStore.Db)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve id for created ToDo-> "+err.Message)
	}

	return &CreateResponse{
		Api: apiVersion,
		Id:  usr.Id,
	}, nil
}

// Update user
func (s *userServiceServer) Update(ctx context.Context, req *UpdateRequest) (*UpdateResponse, error) {
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	mUsr := req.User.ToModel()

	err := Insert(&mUsr, s.dbStore.Db)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve id for updated ToDo-> "+err.Message)
	}

	usr := ModelToUser(&mUsr)

	return &UpdateResponse{
		Api:  apiVersion,
		User: &usr,
	}, nil
}

// Create new todo task
func (s *userServiceServer) Read(ctx context.Context, req *ReadRequest) (*ReadResponse, error) {
	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	// insert ToDo entity data
	usr := GetOne(req.Id, s.dbStore.Db)
	if usr != nil {
		return nil, status.Error(codes.Unknown, "id not found-> "+strconv.FormatInt(req.Id, 10))
	}
	usrResp := ModelToUser(usr)
	return &ReadResponse{
		Api:  apiVersion,
		User: &usrResp,
	}, nil
}

// Create new todo task
func (s *userServiceServer) ReadAll(ctx context.Context, req *ReadAllRequest) (*ReadAllResponse, error) {
	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	// insert ToDo entity data
	usrList := GetAll(s.dbStore.Db)
	usrResp := ModelToUserList(usrList)
	return &ReadAllResponse{
		Api:   apiVersion,
		Users: usrResp,
	}, nil
}
