package gsql

import (
	"database/sql"
	"log/slog"
	"reflect"
	"strings"
)

type InsertBuilder struct {
	db *sql.DB

	sql  string
	args []any
}

func NewInsertBuilder(db *sql.DB) *InsertBuilder {
	return &InsertBuilder{db: db}
}

func (b *InsertBuilder) Insert(dest any) {
	typeOf := reflect.TypeOf(dest)
	valueOf := reflect.ValueOf(dest)

	tableName := GetTN(dest)

	var keys []string
	var place []string
	for i := 0; i < typeOf.Elem().NumField(); i++ {
		isNotNull := valueOf.Elem().Field(i).Field(0).Field(1).Bool()
		if isNotNull {
			key := GetFN(valueOf.Elem().Field(i))
			keys = append(keys, key)
			place = append(place, "?")
			b.args = append(b.args, valueOf.Elem().Field(i).Field(0).Field(0).Interface())
		}
	}

	b.sql = "INSERT INTO " + tableName + " (" + strings.Join(keys, ",") + ") VALUES (" + strings.Join(place, ",") + ")"

	slog.Info(b.sql)
}

func (b *InsertBuilder) Exec() (r sql.Result, err error) {
	return b.db.Exec(b.sql, b.args...)
}
