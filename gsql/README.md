# 数据库操作

Go操作数据库工具包。

## 简介

`gsql` 只是扩展 `database/sql` 包中`sql.DB`功能。

实现不写字符串形式的字段名进行数据库操作。

主要目的是解决当字段有变化时，可在编译时发现问题。

Go的几个热门ORM，都是通过写SQL字符串或者生成代码的方式，总感觉都不是很好。

`gsql`模块对比其它ORM的一些思考：

1. 字段名编译检查，`gorm`字段名直接写字符串，没有编译检查，写错一个字符不能及时发现
2. 不生成多余的代码，`ent`为结构体生成代码，增加了很多没用的代码，增加了体积，不喜欢这种用法

所以手搓一个简化版的ORM，实现不写SQL的单表增删改查操作。

对于多表关联查询的建议：

1. 简单场景的多表查询，可拆分多次查询实现，提高性能可配合使用缓存；
2. 超级复杂的多表查询，使用任何ORM可能都是一团乱的代码，不如原生sql看着清晰，这里保留了`database/sql`完整功能，直接写SQL吧。

## 示例

`gsql` 使用示例。

```sql

var u = &model.User{ }

-- 打开连接
db,err = gsql.Open(gsql.Postgres, "")
  
-- 新增数据
_,err = db.Insert(&model.User{Nickname: "喵喵",Ctime: "20240101"}).Exec()
  
-- 修改数据
_,err = db.Update(u,map[any]any{ &u.Nickname:"旺旺",&u.Ctime:"20240201" }).Where(&u.Nickname,gsql.Eq,"喵喵").And(&u.Ctime,gsql.Gt,"20240101").Exec()
  
-- 删除数据
_,err = db.Delete(u).Where(&u.Nickname,gsql.Eq,"旺旺").Exec()
  
-- 查询数据
userList,err := gsql.Find[model.User](db.Select(u,&u.Nickname,&u.Ctime).Where(&u.Ctime,gsql.Gt,"20240101").Limit(0,10).Query())

```

实现基于结构体的单表操作，当结构体有变动时，编译能做全局检查。


 