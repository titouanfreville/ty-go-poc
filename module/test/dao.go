package test

import (
	"database/sql"
	"github.com/sirupsen/logrus"
	"github.com/titouanfreville/go-namedParameterQuery"
	"go_poc/core"
)

var (
	log = logrus.New()
)

func GetAll(db *sql.DB) *TestList {
	res := TestList{}
	rows, err := db.Query("SELECT * from test")
	if err != nil {
		log.Error(err)
		return &res
	}
	defer rows.Close()
	for rows.Next() {
		var scan Test
		err = rows.Scan(&scan.Id, &scan.Name)
		if err != nil {
			log.Error(err)
			return &TestList{}
		}
		res = append(res, &scan)
	}
	return &res
}

func Insert(t *Test, db *sql.DB) *core.TYPoc {
	if err := t.IsValid(); err != nil {
		return err
	}

	query := namedParameterQuery.NewNamedParameterQuery("INSERT INTO test (name) VALUES ( :name )", "$")
	if err := query.SetValuesFromStruct(*t); err != nil {
		log.Info(err)
		return core.NewDatastoreError("Test.INSERT", "query", err.Error())
	}

	if _, exErr := db.Query(query.GetParsedQuery(), query.GetParsedParameters()...); exErr != nil {
		log.Info(exErr)
		return core.NewDatastoreError("Test.INSERT", "query", exErr.Error())
	}

	return nil
}
