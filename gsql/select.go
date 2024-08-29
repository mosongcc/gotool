package gsql

import (
	"database/sql"
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

type Select struct {
	db *sql.DB

	query string
	args  []any
}

func NewSelect(db *sql.DB) *Select {
	return &Select{db: db}
}

func (b *Select) Select(fields ...any) *Select {
	b.query = "SELECT"
	if fields == nil || len(fields) == 0 {
		b.query += " * "
		return b
	}
	for _, field := range fields {
		b.query += " " + getFieldName(field) + ","
	}
	b.query = strings.TrimRight(b.query, ",")
	return b
}

func (b *Select) From(table any) *Select {
	b.query += " FROM " + getTableName(table)
	return b
}

func (b *Select) Where(name any, opt string, v any) *Select {
	b.query += " WHERE " + getFieldName(name) + " " + opt + " ?"
	b.args = append(b.args, v)
	return b
}

func (b *Select) And(name any, opt string, v any) *Select {
	b.query += " AND " + getFieldName(name) + " " + opt + " ?"
	b.args = append(b.args, v)
	return b
}

func (b *Select) Or(name any, opt string, v any) *Select {
	b.query += " OR " + getFieldName(name) + " " + opt + " ?"
	b.args = append(b.args, v)
	return b
}

func (b *Select) GroupBy(names ...any) *Select {
	b.query += " GROUP BY "
	for _, name := range names {
		b.query += getFieldName(name) + ","
	}
	b.query = strings.TrimRight(b.query, ",")
	return b
}

func (b *Select) OrderBy(name any, sort string) *Select {
	b.query += " ORDER BY " + getFieldName(name) + " " + sort
	return b
}

func (b *Select) Find(result any) (err error) {

	//rows, err := b.db.Query(b.query, b.args...)

	return nil
}
