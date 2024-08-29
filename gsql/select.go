package gsql

import (
	"database/sql"
	"github.com/mosongcc/gotool/gstring"
	"reflect"
	"strings"
)

const (
	Eq   = " = "
	Gt   = " > "
	Lt   = " < "
	Gte  = " >= "
	Lte  = " <= "
	In   = " in "
	Like = " like "
)

type SelectBuilder[T any] struct {
	db *sql.DB

	t     T //查询结果返回类型
	query string
	args  []any
}

func NewSelectBuilder[T any](db *sql.DB) *SelectBuilder[T] {
	return &SelectBuilder[T]{db: db}
}

func (b *SelectBuilder[T]) Select(fields ...any) *SelectBuilder[T] {
	b.query = "SELECT"
	if fields == nil || len(fields) == 0 {
		b.query += " * "
		return b
	}
	for _, field := range fields {
		b.query += " " + GetFN(field) + ","
	}
	b.query = strings.TrimRight(b.query, ",")
	return b
}

func (b *SelectBuilder[T]) From(table any) *SelectBuilder[T] {
	b.query += " FROM " + GetTN(table)
	return b
}

func (b *SelectBuilder[T]) Where(name any, opt string, v any) *SelectBuilder[T] {
	b.query += " WHERE " + GetFN(name) + " " + opt + " ?"
	b.args = append(b.args, v)
	return b
}

func (b *SelectBuilder[T]) And(name any, opt string, v any) *SelectBuilder[T] {
	b.query += " AND " + GetFN(name) + " " + opt + " ?"
	b.args = append(b.args, v)
	return b
}

func (b *SelectBuilder[T]) Or(name any, opt string, v any) *SelectBuilder[T] {
	b.query += " OR " + GetFN(name) + " " + opt + " ?"
	b.args = append(b.args, v)
	return b
}

func (b *SelectBuilder[T]) GroupBy(names ...any) *SelectBuilder[T] {
	b.query += " GROUP BY "
	for _, name := range names {
		b.query += GetFN(name) + ","
	}
	b.query = strings.TrimRight(b.query, ",")
	return b
}

func (b *SelectBuilder[T]) OrderBy(name any, sort string) *SelectBuilder[T] {
	b.query += " ORDER BY " + GetFN(name) + " " + sort
	return b
}

func (b *SelectBuilder[T]) Query(result any) (data []T, err error) {
	rows, err := b.db.Query(b.query, b.args...)
	if err != nil {
		return
	}
	return scan[T](rows)
}

func scan[T any](rows *sql.Rows) (data []T, err error) {
	var columns []string
	columns, err = rows.Columns()
	if err != nil {
		return
	}

	var t T
	rowType := reflect.TypeOf(t)
	for rows.Next() {
		rowValue := reflect.New(rowType).Elem()
		valuesPtr := make([]any, len(columns))
		for i := range columns {
			valuesPtr[i] = rowValue.FieldByName(gstring.Camel(columns[i])).Pointer()
		}
		err = rows.Scan(valuesPtr...)
		if err != nil {
			return
		}
		data = append(data, rowValue.Interface().(T))
	}
	return
}
