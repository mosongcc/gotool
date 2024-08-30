# 数据库操作

目标：

1. 扩展 `database/sql` 包
2. 实现零字符串操作数据库，与结构体强绑定
3. 单表简单增删改查操作（复杂查询通过多次查询与缓存的方案实现）


功能规划：

- model 生成 insert
- model 生成 update
- model 生成 delete
- model 生成 select | where | group by | order by  | limit  


要求：
1. 不写字符串形式的字段名，利用编译检查错误
2. 只做单表的增删改查操作，特殊场景可配合缓存提高性能






 