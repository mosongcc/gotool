package gsql

import (
	"database/sql"
	"reflect"
)

type Update struct {
	db *sql.DB

	sql  string
	args []any
}

func NewUpdate(db *sql.DB) *Update {
	return &Update{db: db}
}

func (b *Update) Update(table any) *Update {
	b.sql = "UPDATE " + getTableName(table)
	return b
}

func (b *Update) SET(dest any) *Update {
	if reflect.TypeOf(dest).Kind() == reflect.Map {
	}
	b.sql = " SET "
	for k, v := range dest.(map[string]any) {
		b.sql += getFieldName(k) + " = ? "
		b.args = append(b.args, v)
	}
	return b
}

func (b *Update) Where(name any, opt string, v any) *Update {
	b.sql += " WHERE " + getFieldName(name) + " " + opt + " ?"
	b.args = append(b.args, v)
	return b
}

func (b *Update) And(name any, opt string, v any) *Update {
	b.sql += " AND " + getFieldName(name) + " " + opt + " ?"
	b.args = append(b.args, v)
	return b
}

func (b *Update) Or(name any, opt string, v any) *Update {
	b.sql += " OR " + getFieldName(name) + " " + opt + " ?"
	b.args = append(b.args, v)
	return b
}
