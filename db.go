package gnet

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql" //mysql driver
)

//None ...
const (
	None = iota
	MYSQL
)

type dbST struct {
}

func newDB() *dbST {
	ptr := &dbST{}
	if ptr.init() {
		return ptr
	} else {
		return nil
	}
}

func (v *dbST) init() bool {
	return true
}

//DBConnect connect to db
func (v *dbST) DBConnect(dbtype int8, addr string, account string, pwd string, dbname string) (*sql.DB, error) {
	var _type string
	switch dbtype {
	case MYSQL:
		_type = "mysql"
	default:
		return nil, errors.New("invalid type")
	}

	formatStr := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", account, pwd, addr, dbname)
	var db *sql.DB
	var err error
	if db, err = sql.Open(_type, formatStr); err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
