package gsql

import "database/sql"

func Select[T any](db *sql.DB, table any, fields ...any) *SelectBuilder[T] {
	return NewSelectBuilder[T](db).Select(fields).From(table)
}

func Insert() {

}
