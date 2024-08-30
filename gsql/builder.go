package gsql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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
	In    Opt = " IN "
	NotIn Opt = " NOT IN "
	Like  Opt = " LIKE "
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
	typeOf := reflect.TypeOf(dest)
	if typeOf.Kind() != reflect.Struct {
		panic(fmt.Errorf("Builder.Insert 参数类型必须是 Struct"))
	}
	valueOf := reflect.ValueOf(dest)

	tableName := getTableName(typeOf, valueOf)
	keys, place, args := getStructFields(typeOf, valueOf)

	b.args = args
	b.sql = "INSERT INTO " + tableName + " (" + strings.Join(keys, ",") + ") VALUES (" + strings.Join(place, ",") + ")"
	return b
}

// Update 更新
// - table 表名，如果是字符串，字段名也要用字符串。如果是指针，字段名也可以用指针
// - set 类型:struct | map 。 如果是Map，当表是指针地址时key可以是指针地址
func (b *Builder) Update(table any, set any) *Builder {
	b.sql = "UPDATE " + TN(table) + " SET "

	typeOf := reflect.TypeOf(set)
	switch typeOf.Kind() {
	case reflect.Struct:
		valueOf := reflect.ValueOf(set)
		keys, _, args := getStructFields(typeOf, valueOf)
		b.sql += strings.Join(keys, " = ? ,") + " = ? "
		b.args = append(b.args, args...)
	case reflect.Map:
		for k, v := range set.(map[any]any) {
			b.sql += FN(k) + " = ? ,"
			b.args = append(b.args, v)
		}
		b.sql = strings.TrimRight(b.sql, ",")
	default:
		panic(fmt.Errorf("Builder.Update 参数 set 仅限 Struct 或者 Map 类型"))
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
	case Oracle:
		b.sql = "SELECT * FROM (SELECT t.*, ROWNUM rn FROM (" + b.sql + ") t WHERE ROWNUM <= ?) WHERE rn >= ?" //1结束行号  2开始行号
		b.args = append(b.args, offset+limit, offset)
	default:
		panic(errors.New("Unsupported driver " + string(b.db.DriverName)))
	}
	return b
}

func (b *Builder) safe() (err error) {
	if strings.HasPrefix(b.sql, "DELETE") && !strings.Contains(b.sql, "WHERE") {
		err = errors.New("Builder.Exec 禁止DELETE操作不加WHERE条件")
		return
	}
	if strings.HasPrefix(b.sql, "UPDATE") && !strings.Contains(b.sql, "WHERE") {
		err = errors.New("Builder.Exec 禁止UPDATE操作不加WHERE条件")
		return
	}
	return
}

// Exec 执行 INSERT UPDATE DELETE
func (b *Builder) Exec() (r sql.Result, err error) {
	if err = b.safe(); err != nil {
		return
	}
	return b.db.DB.ExecContext(b.ctx, b.sql, b.args...)
}

// ExecTx 事务执行
func (b *Builder) ExecTx(tx *sql.Tx) (r sql.Result, err error) {
	if err = b.safe(); err != nil {
		return
	}
	return tx.ExecContext(b.ctx, b.sql, b.args...)
}

// Query 执行 SELECT
func (b *Builder) Query() (*sql.Rows, error) {
	return b.db.DB.QueryContext(b.ctx, b.sql, b.args...)
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
