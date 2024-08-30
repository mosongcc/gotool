package gsql

import (
	"context"
	"database/sql"
	"errors"
	"github.com/mosongcc/gotool/gstring"
	"reflect"
	"strings"
)

type Opt string

const (
	Eq    Opt = " = "
	Ne    Opt = " <> "
	Gt    Opt = " > "
	Lt    Opt = " < "
	Gte   Opt = " >= "
	Lte   Opt = " <= "
	In    Opt = " in "
	NotIn Opt = " not in "
	Like  Opt = " like "
)

type Builder struct {
	ctx context.Context
	db  *DB

	sql  string
	args []any
}

func NewBuilder(ctx context.Context, db *DB) *Builder {
	if ctx == nil {
		ctx = context.Background()
	}
	return &Builder{ctx: ctx, db: db}
}

func (b *Builder) WithValue(k, v any) *Builder {
	b.ctx = context.WithValue(b.ctx, k, v)
	return b
}

func (b *Builder) Insert(dest any) *Builder {
	tableName := TN(dest)

	typeOf := reflect.TypeOf(dest)
	valueOf := reflect.ValueOf(dest)

	var keys []string
	var place []string
	for i := 0; i < typeOf.Elem().NumField(); i++ {
		isNotNull := valueOf.Elem().Field(i).Field(0).Field(1).Bool()
		if isNotNull {
			key := FN(valueOf.Elem().Field(i))
			keys = append(keys, key)
			place = append(place, "?")
			b.args = append(b.args, valueOf.Elem().Field(i).Field(0).Field(0).Interface())
		}
	}
	b.sql = "INSERT INTO " + tableName + " (" + strings.Join(keys, ",") + ") VALUES (" + strings.Join(place, ",") + ")"
	return b
}

func (b *Builder) Update(table any, set any) *Builder {
	b.sql = "UPDATE " + TN(table)

	//TODO SET
	if reflect.TypeOf(set).Kind() == reflect.Map {
	}
	b.sql = " SET "
	for k, v := range set.(map[string]any) {
		b.sql += FN(k) + " = ? "
		b.args = append(b.args, v)
	}
	return b
}

func (b *Builder) Delete(table any) *Builder {
	b.sql = "DELETE FROM " + TN(table)
	return b
}

func (b *Builder) Select(table any, fields ...any) *Builder {
	tableName := TN(table)

	b.sql = "SELECT"
	if fields == nil || len(fields) == 0 {
		b.sql += " * "
		return b
	}
	for _, field := range fields {
		b.sql += " " + FN(field) + ","
	}
	b.sql = strings.TrimRight(b.sql, ",")
	b.sql += " FROM " + tableName
	return b
}

func (b *Builder) Where(name any, opt Opt, v any) *Builder {
	b.sql += " WHERE " + FN(name) + string(opt) + "?"
	b.args = append(b.args, v)
	return b
}

func (b *Builder) And(name any, opt Opt, v any) *Builder {
	b.sql += " AND " + FN(name) + string(opt) + "?"
	b.args = append(b.args, v)
	return b
}

func (b *Builder) Or(name any, opt Opt, v any) *Builder {
	b.sql += " OR " + FN(name) + string(opt) + "?"
	b.args = append(b.args, v)
	return b
}

func (b *Builder) GroupBy(names ...any) *Builder {
	b.sql += " GROUP BY "
	for _, name := range names {
		b.sql += FN(name) + ","
	}
	b.sql = strings.TrimRight(b.sql, ",")
	return b
}

func (b *Builder) OrderBy(name any, sort string) *Builder {
	b.sql += " ORDER BY " + FN(name) + " " + sort
	return b
}

func (b *Builder) Limit(offset int64, limit int64) *Builder {
	switch b.db.DriverName {
	case Mysql:
		b.sql += " LIMIT ?,?"
		b.args = append(b.args, offset, limit)
	case Postgres, Sqlite3:
		b.sql += "OFFSET ? LIMIT ? "
		b.args = append(b.args, offset, limit)
	case Mssql:
		b.sql += " OFFSET ? ROWS FETCH NEXT ? ROWS ONLY "
		b.args = append(b.args, offset, limit)
	default:
		panic(errors.New("Unsupported driver " + string(b.db.DriverName)))
	}
	return b
}

// Exec 执行 INSERT UPDATE DELETE
func (b *Builder) Exec() (r sql.Result, err error) {
	return b.db.ExecContext(b.ctx, b.sql, b.args...)
}

// ExecTx 事务执行
func (b *Builder) ExecTx(tx *sql.Tx) (r sql.Result, err error) {
	return tx.ExecContext(b.ctx, b.sql, b.args...)
}

// Query 执行 SELECT
func (b *Builder) Query() (*sql.Rows, error) {
	return b.db.QueryContext(b.ctx, b.sql, b.args...)
}

// QueryTx 事务执行
func (b *Builder) QueryTx(tx *sql.Tx) (*sql.Rows, error) {
	return tx.QueryContext(b.ctx, b.sql, b.args...)
}

// Sql 输出sql
func (b *Builder) Sql() (sql string) {
	sql = b.sql
	for _, v := range b.args {
		sql = strings.Replace(sql, "?", "\""+v.(string)+"\"", 1)
	}
	return
}

// RowsScan 查询结果转为结构数据
func RowsScan[T any](rows *sql.Rows) (data []T, err error) {
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

// Find 查询结果列表解析  参数 rows 是 db.Query() 返回结果
func Find[T any](rows *sql.Rows, e ...error) (list []T, err error) {
	if len(e) > 0 {
		err = e[0]
		return
	}
	return RowsScan[T](rows)
}

func First[T any](rows *sql.Rows, e ...error) (entity T, err error) {
	var list []T
	list, err = Find[T](rows, e...)
	if err != nil {
		return
	}
	if list == nil || len(list) == 0 {
		err = errors.New("records not found")
		return
	}
	entity = list[0]
	return
}
