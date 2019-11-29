package tasks

import (
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"go_poc/core"
)

var (
	log = logrus.New()
)

func GetAll(db *sqlx.DB) *TaskList {
	res := TaskList{}
	err := db.Select(&res, "SELECT * from tasks")
	if err != nil {
		log.Error(err)
		return &res
	}
	for _, tsk := range res {
		tsk.Populate(db)
	}
	return &res
}

func Insert(tsk *Task, db *sqlx.DB) *core.TYPoc {
	if err := tsk.IsValid(db); err != nil {
		return err
	}
	tsk.BeforeDB()

	if _, exErr := db.NamedExec("INSERT INTO tasks (resume, content, reporter_id, worker_id) VALUES (:resume, :content, :reporter_id, :worker_id)", *tsk); exErr != nil {
		log.Info(exErr)
		return core.NewDatastoreError("Task.INSERT", "query", exErr.Error())
	}
	var id int64
	if err := db.Get(&id, "SELECT id FROM tasks order by id desc limit 1"); err != nil {
		return core.NewDatastoreError("Task.INSERT", "get_id", err.Error())
	}

	tsk.Id = id

	return nil
}

func GetOne(id int64, db *sqlx.DB) *Task {
	tsk := Task{}
	if err := db.Get(&tsk, "SELECT * FROM tasks where id = $1", id); err != nil {
		return nil
	}
	tsk.Populate(db)
	return &tsk
}

func Update(tsk *Task, db *sqlx.DB) *core.TYPoc {
	if err := tsk.IsValid(db); err != nil {
		return err
	}

	tsk.BeforeDB()
	if _, exErr := db.NamedExec("UPDATE tasks SET resume = :resume,  content = :content, reporter_id = :reporter_id, worker_id = :worker_id WHERE id = :id", *tsk); exErr != nil {
		log.Info(exErr)
		return core.NewDatastoreError("Task.USER", "query", exErr.Error())
	}
	tsk.Populate(db)

	return nil
}

func CheckId(id int64, db *sqlx.DB) bool {
	tsk := Task{}
	if err := db.Get(&tsk, "SELECT * FROM tasks where id = $1", id); err != nil {
		log.Error(err.Error())
		return false
	}
	return tsk != Task{}
}
