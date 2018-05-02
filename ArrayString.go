/*

ArrayString类型

方便将数据库中字符串字段值直接转换成字符串数组，默认只支持半角逗号作为分隔符

开发者: 无双 (luciferlai@qq.com)

*/
package sqlxyz

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

/*
自定义类型ArrayString 字符串数组, 支持sql查询

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

// 数组分割符
var ArrayStringSep string = ","

//用于sql.Scan
func (this *ArrayString) Scan(value interface{}) (err error) {
	if value == nil {
		return nil
	}
	switch value.(type) {
	case []uint8:
		var tmp []string
		tmp = strings.Split(this.b2s(value.([]uint8)), ArrayStringSep)
		for _, v := range tmp {
			this.Add(v)
		}
		return nil
	case string:
		var tmp []string
		tmp = strings.Split(value.(string), ArrayStringSep)
		for _, v := range tmp {
			this.Add(v)
		}
		return nil
	}
	return fmt.Errorf("Can't convert %T to ArrayString", value)
}

//转换[]uint8 成 string类型
func (this *ArrayString) b2s(bs []uint8) string {
	b := make([]byte, len(bs))
	for i, v := range bs {
		b[i] = byte(v)
	}
	return string(b)
}

//增加一个值，如果存在将不添加
func (this *ArrayString) AddOnly(str string) {
	if this.Find(str) != true {
		*this = append(*this, str)
	}
}

//增加一个值，允许重复
func (this *ArrayString) Add(str string) {
	*this = append(*this, str)
}

/*
直接赋值
	var this ArrayString
	this.Set("lucifer,wushuang")
*/
func (this *ArrayString) Set(str string) {
	*this = strings.Split(str, ArrayStringSep)
}

//查找是否存在
func (this ArrayString) Find(str string) bool {
	for _, v := range this {
		if v == str {
			return true
		}
	}
	return false
}

//统计数量
func (this ArrayString) Len() int {
	return len(this)
}

//用于sql写入时引用值
func (this ArrayString) Value() (driver.Value, error) {
	return strings.Join(this, ArrayStringSep), nil
}

//用于输出 fmt.Println(this)
func (this ArrayString) String() []string {
	return this
}

//用于输出 strings.Join(this, ",")
func (this ArrayString) Join(sep string) string {
	return strings.Join(this, sep)
}

// 输出字符串形式, 如 "1","2","3","4"
func (this ArrayString) ToString() string {
	return "\"" + strings.Join(this, "\",\"") + "\""
}
