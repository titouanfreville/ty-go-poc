package user

import (
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"go_poc/core"
)

var (
	log = logrus.New()
)

func GetAll(db *sqlx.DB) *UserModelList {
	res := UserModelList{}
	err := db.Select(&res, "SELECT * from users")
	if err != nil {
		log.Error(err)
		return &res
	}
	return &res
}

func Insert(usr *UserModel, db *sqlx.DB) *core.TYPoc {
	if err := usr.IsValid(); err != nil {
		return err
	}

	if _, exErr := db.NamedExec("INSERT INTO users (name, email) VALUES (:name, :email)", *usr); exErr != nil {
		log.Info(exErr)
		return core.NewDatastoreError("Test.INSERT", "query", exErr.Error())
	}

	return nil
}

func GetOne(id int64, db *sqlx.DB) *UserModel {
	user := UserModel{}
	if err := db.Get(&user, "SELECT * FROM users where id = $1", id); err != nil {
		log.Error(err.Error())
		return nil
	}
	return &user
}

func CheckId(id int64, db *sqlx.DB) bool {
	user := UserModel{}
	if err := db.Get(&user, "SELECT * FROM users where id = $1", id); err != nil {
		log.Error(err.Error())
		return false
	}
	return user != UserModel{}
}
