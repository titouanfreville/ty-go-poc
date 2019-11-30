package tasks

import (
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"go_poc/core"
	"go_poc/module/user"
	"io"
)

const (
	TASK_BACKLOG     = 0
	TASK_IN_PROGRESS = 1
	TASK_DONE        = 2
)

// Task table model
type Task struct {
	Id         int64          `form:"-" json:"-" db:"id"`
	Resume     string         `form:"resume" json:"resume" db:"resume"`
	Content    string         `form:"content" json:"content" db:"content"`
	ReporterId int64          `form:"reporter_id" json:"reporter_id" db:"reporter_id"`
	WorkerId   int64          `form:"worker_id" json:"worker_id" db:"worker_id"`
	Reporter   user.UserModel `form:"reporter" json:"reporter" db:"-"`
	Worker     user.UserModel `form:"worker" json:"worker" db:"-"`
	Status     int64          `form:"-" json:"-" db:"status"`
	StatusStr  string         `form:"status" json:"status" db:"status"`
}

// TaskList is a shortcut to a list of Task
type TaskList []*Task

// IsValid check if Task object is valid
func (t *Task) IsValid(db *sqlx.DB) *core.TYPoc {
	if t.Resume == "" {
		return core.NewModelError("Task.IsValid", "resume", "resume required")
	}

	if t.ReporterId != 0 && !user.CheckId(t.ReporterId, db) {
		return core.NewModelError("Task.IsValid", "reporter_id", "invalid reporter")
	}
	if t.WorkerId != 0 && !user.CheckId(t.WorkerId, db) {
		return core.NewModelError("Task.IsValid", "reporter_id", "invalid reporter")
	}

	return nil
}

// BeforeDB
func (t *Task) BeforeDB() {
	t.Status = convertStatus(t.StatusStr).(int64)
}

func statusFromInt(s int64) string {
	switch s {
	case TASK_IN_PROGRESS:
		return "in_progress"
	case TASK_DONE:
		return "done"
	default:
		return "backlog"
	}
}

func statusFromString(s string) int64 {
	switch s {
	case "in_progress":
		return TASK_IN_PROGRESS
	case "done":
		return TASK_DONE
	default:
		return TASK_BACKLOG
	}
}

func convertStatus(status interface{}) interface{} {
	switch status.(type) {
	case string:
		return statusFromString(status.(string))
	case int64:
		return statusFromInt(status.(int64))
	}
	return nil
}

func (t *Task) Populate(db *sqlx.DB) {
	if t.ReporterId != 0 {
		t.Reporter = *user.GetOne(t.ReporterId, db)
	} else {
		t.Worker = user.UserModel{}
	}
	if t.WorkerId != 0 {
		t.Worker = *user.GetOne(t.WorkerId, db)
	} else {
		t.Worker = user.UserModel{}
	}
	t.StatusStr = convertStatus(t.Status).(string)
}

// ToJson serializes the bot patch to json.
func (t *Task) ToJson() []byte {
	data, err := json.Marshal(t)
	if err != nil {
		return nil
	}

	return data
}

// BotPatchFromJson deserializes a bot patch from json.
func TaskFromJson(data io.Reader) *Task {
	decoder := json.NewDecoder(data)
	var taskData Task
	err := decoder.Decode(&taskData)
	if err != nil {
		return nil
	}

	return &taskData
}

// ToJson serializes the bot patch to json.
func (tl *TaskList) ToJson() []byte {
	data, err := json.Marshal(tl)
	if err != nil {
		return nil
	}

	return data
}

// BotPatchFromJson deserializes a bot patch from json.
func TaskListFromJson(data io.Reader) *TaskList {
	decoder := json.NewDecoder(data)
	var taskList TaskList
	err := decoder.Decode(&taskList)
	if err != nil {
		return nil
	}

	return &taskList
}
