封装了一些对mysql的常用方法：

1. 快捷的使用sqlx操作mysql
2. 增加了一些mysql字段类型，Boolean（int转换成布尔型）, Integers（varchar 自动转换成 整型数组）,ArrayString(varchar 自动转换成 字符串数组),Timestamp(int转换成unix时间戳)
3. 分页计算
4. 一些操作数据库，用户，表的操作
5. SchemaModel表模型
6. 支持go mod

Lucifer Lai 
2019-10-28