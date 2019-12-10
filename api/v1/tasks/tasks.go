package tasks

import (
	"context"
	"github.com/micro/go-micro"
	"go_poc/core"
	"go_poc/module/tasks"
	"net/http"
)

const (
	DefaultClientName = "cli"
	DefaultServerName = "grpc"
)

type TaskHandler struct {
	dbStore *core.DBStore
}

//  ---------- INITIALIZER ----------------
// NewHandler initialize the handler for TaskHandler processing
func NewHandler(connection core.DbConnection) TaskHandler {
	taskServ := TaskHandler{dbStore: &core.DBStore{}}
	taskServ.dbStore.InitConnection(
		connection.User, connection.Database,
		connection.Password, connection.Host,
		connection.Port,
	)
	return taskServ
}

// NewClient init a task client
func NewClient(name *string, server *string) ToDoService {
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
	return NewToDoService(*server, service.Client())
}

//  ---------- FEATURE ----------------
func (t TaskHandler) Create(ctx context.Context, req *CreateRequest, rsp *CreateResponse) error {
	task := req.ToDo.toModel()

	err := tasks.Insert(&task, t.dbStore.Db)
	if err != nil {
		return err
	}
	rsp.Id = task.Id
	rsp.Api = "v1"
	return nil
}

func (t TaskHandler) Read(ctx context.Context, req *ReadRequest, rsp *ReadResponse) error {
	rsp.Api = "v1"
	rsp.ToDo = modelToProto(tasks.GetOne(req.Id, t.dbStore.Db))
	return nil
}

func (t TaskHandler) ReadAll(ctx context.Context, req *ReadAllRequest, rsp *ReadAllResponse) error {
	rsp.Api = "v1"
	rsp.ToDos = modelToProtoList(tasks.GetAll(t.dbStore.Db))
	return nil
}

func (t TaskHandler) Update(ctx context.Context, req *UpdateRequest, rsp *UpdateResponse) error {
	task := req.ToDo.toModel()
	if !tasks.CheckId(task.Id, t.dbStore.Db) {
		return core.NewAPIError(
			"service", "tasks",
			"Could not retrieve task to update", http.StatusNotFound,
			"update",
		)
	}

	// update ToDo
	err := tasks.Update(&task, t.dbStore.Db)

	if err != nil {
		return err
	}
	rsp.Api = "v1"
	rsp.Updated = task.Id
	return nil
}

// ---------- CONVERTERS ----------------
func modelToProto(t *tasks.Task) *ToDo {
	return &ToDo{
		Id:         t.Id,
		Resume:     t.Resume,
		Content:    t.Content,
		ReporterId: t.ReporterId,
		WorkerId:   t.WorkerId,
		Status:     t.StatusStr,
	}
}
func modelToProtoList(tl *tasks.TaskList) []*ToDo {
	var res []*ToDo
	for _, t := range *tl {
		proto := modelToProto(t)
		res = append(res, proto)
	}
	return res
}

func (m ToDo) toModel() tasks.Task {
	return tasks.Task{
		Id:         m.Id,
		Resume:     m.Resume,
		Content:    m.Content,
		ReporterId: m.ReporterId,
		WorkerId:   m.WorkerId,
		StatusStr:  m.Status,
	}
}
