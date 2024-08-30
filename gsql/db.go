package gsql

import (
	"context"
	"database/sql"
)

// Driver 驱动标识，具体值根据使用的驱动获取，这里给出几个最常用的驱动值
type Driver string

const (
	Mysql    Driver = "mysql"
	Mssql    Driver = "mssql"
	Postgres Driver = "postgres"
	Sqlite3  Driver = "sqlite3"
	Oracle   Driver = "ora"
)

type DB struct {
	DriverName Driver
	*sql.DB
}

// Open 连接数据库
func Open(driverName Driver, dataSourceName string) (db *DB, err error) {
	database, err := sql.Open(string(driverName), dataSourceName)
	if err != nil {
		return
	}
	db = &DB{DriverName: driverName, DB: database}
	return
}

func (db *DB) Insert(dest any) *Builder {
	return NewBuilder(context.TODO(), db).Insert(dest)
}

func (db *DB) Update(table any, set any) *Builder {
	return NewBuilder(context.TODO(), db).Update(table, set)
}

func (db *DB) Delete(table any) *Builder {
	return NewBuilder(context.TODO(), db).Delete(table)
}

func (db *DB) Select(table any, fields ...any) *Builder {
	return NewBuilder(context.TODO(), db).Select(table, fields...)
}

// Tx 事务
func (db *DB) Tx(f func(tx *sql.Tx) (any, error)) (v any, err error) {
	var tx *sql.Tx
	tx, err = db.DB.Begin()
	if err != nil {
		return
	}
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
		if err != nil {
			if e := tx.Rollback(); e != nil {
				panic(e)
			}
		}
	}()
	v, err = f(tx)
	if err != nil {
		return
	}
	err = tx.Commit()
	if err != nil {
		return
	}
	return
}
