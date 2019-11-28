package core

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// StoreImpl implementation for store interface
type DBStore struct {
	usr    string
	dbName string
	pass   string
	host   string
	port   string
	Db     *sqlx.DB
}

func (s *DBStore) connect() bool {
	connStr := fmt.Sprintf(
		"sslmode=disable user='%s' password='%s' host='%s' port='%s' dbname='%s'",
		s.usr, s.pass, s.host, s.port, s.dbName)
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
		return false
	}
	s.Db = db
	return true
}

// InitConnection to store
func (s *DBStore) InitConnection(user string, dbname string, password string, host string, port string) bool {
	s.usr = user
	s.dbName = dbname
	s.pass = password
	s.host = host
	s.port = port
	return s.connect()
}

// CloseConnection to store
func (s *DBStore) CloseConnection() bool {
	if s.Db == nil {
		return false
	}
	if err := s.Db.Close(); err != nil {
		return false
	}
	return true
}

// AssertConnectionOrReconnect assert connection to database or try to reconnect
func (s *DBStore) AssertConnectionOrReconnect() bool {
	db := s.Db
	if s.Db == nil || db.Ping() != nil {
		log.Info("No connection to DB")
		_ = s.CloseConnection()
		if ok := s.connect(); !ok {
			log.Error("Could not reconnect to DB")
			return false
		}
	}
	if db.Ping() != nil {
		log.Error("No connection to Database")
		return false
	}
	return true
}
