package gsql

import (
	"database/sql"
	"log/slog"
	"reflect"
	"strings"
)

type Insert struct {
	db *sql.DB

	sql  string
	args []any
}

func NewInsert(db *sql.DB) *Insert {
	return &Insert{db: db}
}

func (b *Insert) Insert(dest any) {
	typeOf := reflect.TypeOf(dest)
	valueOf := reflect.ValueOf(dest)

	tableName := getTableName(dest)

	var keys []string
	var args []any
	var place []string
	for i := 0; i < typeOf.Elem().NumField(); i++ {
		key := getFieldName(valueOf.Elem().Field(i))

		isNotNull := valueOf.Elem().Field(i).Field(0).Field(1).Bool()
		if isNotNull {
			val := valueOf.Elem().Field(i).Field(0).Field(0).Interface()

			keys = append(keys, key)
			args = append(args, val)
			place = append(place, "?")
		}
	}

	query := "INSERT INTO " + tableName + " (" + strings.Join(keys, ",") + ") VALUES (" + strings.Join(place, ",") + ")"

	slog.Info(query)
}
