package gsql

import (
	"database/sql"
)

type Delete struct {
	db *sql.DB

	sql  string
	args []any
}

func NewDelete(db *sql.DB) *Delete {
	return &Delete{db: db}
}

func (b *Delete) Delete(table any) *Delete {
	b.sql = "DELETE FROM " + GetTN(table)
	return b
}

func (b *Delete) Where(name any, opt string, v any) *Delete {
	b.sql += " WHERE " + GetFN(name) + " " + opt + " ?"
	b.args = append(b.args, v)
	return b
}

func (b *Delete) And(name any, opt string, v any) *Delete {
	b.sql += " AND " + GetFN(name) + " " + opt + " ?"
	b.args = append(b.args, v)
	return b
}

func (b *Delete) Or(name any, opt string, v any) *Delete {
	b.sql += " OR " + GetFN(name) + " " + opt + " ?"
	b.args = append(b.args, v)
	return b
}

func (b *Delete) Exec() (r sql.Result, err error) {
	return b.db.Exec(b.sql, b.args...)
}
