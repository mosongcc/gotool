package gsql

import (
	"context"
	"database/sql"
)

// Open 连接数据库
func Open(driverName Driver, dataSourceName string) (db *DB, err error) {
	database, err := sql.Open(string(driverName), dataSourceName)
	if err != nil {
		return
	}
	db = &DB{DriverName: driverName, DB: database}
	return
}

func Insert(db *DB, dest any) *Builder {
	return NewBuilder(context.TODO(), db).Insert(dest)
}

func Update(db *DB, table any, set any) *Builder {
	return NewBuilder(context.TODO(), db).Update(table, set)
}

func Delete(db *DB, table any) *Builder {
	return NewBuilder(context.TODO(), db).Delete(table)
}

func Select(db *DB, table any, fields ...any) *Builder {
	return NewBuilder(context.TODO(), db).Select(table, fields...)
}

type DB struct {
	DriverName Driver
	*sql.DB
}

// Tx 事务
func (db *DB) Tx(f func(tx *sql.Tx) (any, error)) (v any, err error) {
	var tx *sql.Tx
	tx, err = db.Begin()
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
