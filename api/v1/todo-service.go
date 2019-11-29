package v1

import (
	"context"
	"github.com/sirupsen/logrus"
	"go_poc/core"
	"go_poc/module/tasks"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	// apiVersion is version of API is provided by server
	apiVersion = "v1"
)

var (
	log = logrus.New()
)

// toDoServiceServer is implementation of v1.ToDoServiceServer proto interface
type toDoServiceServer struct {
	dbStore *core.DBStore
}

func (m ToDo) ToTask() tasks.Task {
	return tasks.Task{
		Id:         m.Id,
		Resume:     m.Resume,
		Content:    m.Content,
		ReporterId: m.ReporterId,
		WorkerId:   m.WorkerId,
		StatusStr:  m.Status,
	}
}

func TaskToToDo(t *tasks.Task) ToDo {
	return ToDo{
		Id:         t.Id,
		Resume:     t.Resume,
		Content:    t.Content,
		ReporterId: t.ReporterId,
		WorkerId:   t.WorkerId,
		Status:     t.StatusStr,
	}
}
func TaskToToDoList(tl *tasks.TaskList) []*ToDo {
	res := []*ToDo{}
	for _, t := range *tl {
		proto := TaskToToDo(t)
		res = append(res, &proto)
	}
	return res
}

// NewToDoServiceServer creates ToDo service
func NewToDoServiceServer(dbSettings core.DbConnection) ToDoServiceServer {
	dbStore := core.DBStore{}
	dbStore.InitConnection(dbSettings.User, dbSettings.Database, dbSettings.Password, dbSettings.Host, dbSettings.Port)
	return &toDoServiceServer{dbStore: &dbStore}
}

// checkAPI checks if the API version requested by client is supported by server
func (s *toDoServiceServer) checkAPI(api string) error {
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
func (s *toDoServiceServer) Create(ctx context.Context, req *CreateRequest) (*CreateResponse, error) {
	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	task := req.ToDo.ToTask()

	// insert ToDo entity data
	err := tasks.Insert(&task, s.dbStore.Db)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve id for created ToDo-> "+err.Message)
	}

	return &CreateResponse{
		Api: apiVersion,
		Id:  task.Id,
	}, nil
}

// Read todo task
func (s *toDoServiceServer) Read(ctx context.Context, req *ReadRequest) (*ReadResponse, error) {
	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	// query ToDo by ID
	task := tasks.GetOne(req.Id, s.dbStore.Db)
	toDo := TaskToToDo(task)
	return &ReadResponse{
		Api:  apiVersion,
		ToDo: &toDo,
	}, nil

}

// Update todo task
func (s *toDoServiceServer) Update(ctx context.Context, req *UpdateRequest) (*UpdateResponse, error) {
	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	task := req.ToDo.ToTask()
	log.Info(task.Id)
	if !tasks.CheckId(task.Id, s.dbStore.Db) {
		return nil, status.Error(codes.NotFound, "failed to update ToDo-> provided id does not exists")
	}

	// update ToDo
	err := tasks.Update(&task, s.dbStore.Db)

	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve rows affected value-> "+err.Message)
	}

	return &UpdateResponse{
		Api:     apiVersion,
		Updated: task.Id,
	}, nil
}

//
// // Delete todo task
// func (s *toDoServiceServer) Delete(ctx context.Context, req *v1.DeleteRequest) (*v1.DeleteResponse, error) {
// 	// check if the API version requested by client is supported by server
// 	if err := s.checkAPI(req.Api); err != nil {
// 		return nil, err
// 	}
//
// 	// get SQL connection from pool
// 	c, err := s.connect(ctx)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer c.Close()
//
// 	// delete ToDo
// 	res, err := c.ExecContext(ctx, "DELETE FROM ToDo WHERE `ID`=?", req.Id)
// 	if err != nil {
// 		return nil, status.Error(codes.Unknown, "failed to delete ToDo-> "+err.Error())
// 	}
//
// 	rows, err := res.RowsAffected()
// 	if err != nil {
// 		return nil, status.Error(codes.Unknown, "failed to retrieve rows affected value-> "+err.Error())
// 	}
//
// 	if rows == 0 {
// 		return nil, status.Error(codes.NotFound, fmt.Sprintf("ToDo with ID='%d' is not found",
// 			req.Id))
// 	}
//
// 	return &v1.DeleteResponse{
// 		Api:     apiVersion,
// 		Deleted: rows,
// 	}, nil
// }
//
// Read all todo tasks
func (s *toDoServiceServer) ReadAll(ctx context.Context, req *ReadAllRequest) (*ReadAllResponse, error) {
	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	// get ToDo list
	taskList := tasks.GetAll(s.dbStore.Db)

	return &ReadAllResponse{
		Api:   apiVersion,
		ToDos: TaskToToDoList(taskList),
	}, nil
}
