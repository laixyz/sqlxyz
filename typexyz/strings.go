// Package typexyz
// 方便将数据库中字符串字段值直接转换成字符串数组，默认只支持半角逗号作为分隔符
// 开发者: 无双 (luciferlai@qq.com)

package typexyz

import (
	"database/sql"
	"database/sql/driver"
	"strings"
)

/*
ArrayString 自定义类型 字符串数组, 支持sql查询

实例：
	type Member struct {
		MemberID int64       `db:"MemberID"`
		Email    string      `db:"Email"`
		Password string      `db:"Password"`
		Lives    ArrayString `db:"Lives"`
		State    int64       `db:"State"`
		Created  time.Time   `db:"Created"`
		Updated  time.Time   `db:"Updated"`
		Deleted  time.Time   `db:"Deleted"`
	}
	var members []Member
	err = db.Select(&members, "select * from members order by MemberID desc limit 1")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(members)
	}
输出
	[{20 luciferlai@qq.com wushuang {[111 222 333 444]} 0 2018-04-25 13:08:00 +0800 CST 1970-01-01 08:00:00 +0800 CST 1970-01-01 08:00:00 +0800 CST}]
*/
type ArrayString []string

// ArrayStringSep 数组分割符
var ArrayStringSep string = ","

// Scan 用于sql.Scan
func (arrayString *ArrayString) Scan(value interface{}) (err error) {
	if value == nil {
		return nil
	}
	var nullString sql.NullString
	err = nullString.Scan(value)
	if err != nil {
		return err
	}
	if nullString.Valid == false {
		return nil
	}
	if nullString.String == "" {
		return nil
	}
	var tmp []string
	tmp = strings.Split(nullString.String, ArrayStringSep)
	for _, v := range tmp {
		arrayString.Add(v)
	}
	return nil
}

// Add 增加一个值，如果存在将不添加
func (arrayString *ArrayString) Add(str string) {
	if arrayString.Find(str) != true {
		*arrayString = append(*arrayString, str)
	}
}

/*
Set 直接赋值
	var this ArrayString
	this.Set("lucifer,wushuang")
*/
func (arrayString *ArrayString) Set(str string) {
	*arrayString = strings.Split(str, ArrayStringSep)
}

//Find 查找是否存在
func (arrayString ArrayString) Find(str string) bool {
	for _, v := range arrayString {
		if v == str {
			return true
		}
	}
	return false
}

// Len 统计数量
func (arrayString ArrayString) Len() int {
	return len(arrayString)
}

// Value 用于sql写入时引用值
func (arrayString ArrayString) Value() (driver.Value, error) {
	return strings.Join(arrayString, ArrayStringSep), nil
}

// String 用于输出 fmt.Println(this)
func (arrayString ArrayString) String() []string {
	return arrayString
}

// Join 用于输出 strings.Join(this, ",")
func (arrayString ArrayString) Join(sep string) string {
	return strings.Join(arrayString, sep)
}
