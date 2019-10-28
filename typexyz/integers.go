// Package typexyz
// 方便将数据库中字符串字段值直接转换成int64数组，默认只支持半角逗号作为分隔符
// 开发者: 无双 (luciferlai@qq.com)

package typexyz

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
)

// Integers []int64匿名类型
type Integers []int64

// ArrayIntSep 分割符
var ArrayIntSep string = ","

// Scan 用于sql.Scan
func (integers *Integers) Scan(value interface{}) (err error) {
	if value == nil {
		return nil
	}
	var nullString sql.NullString
	err = nullString.Scan(value)
	if err != nil {
		return err
	}
	if nullString.Valid == true {
		if nullString.String == "" {
			return nil
		}
		var tmp []string
		str := strings.TrimSpace(nullString.String)
		tmp = strings.Split(str, ArrayIntSep)
		for _, v := range tmp {
			var newNumber int64
			newNumber, err := strconv.ParseInt(v, 10, 64)
			if err == nil {
				integers.Add(newNumber)
			}
		}
		return nil
	}
	return nil
}

// Add 增加一个值，如果存在将不添加
func (integers *Integers) Add(i int64) {
	if integers.Find(i) != true {
		*integers = append(*integers, i)
	}
}

/*
SetFromString 赋值函数
直接赋值
	var this ArrayInt64
	this = []int64{1,2,3}
	this.Set("1,2,3")
*/
func (integers *Integers) SetFromString(str string) error {
	var tmp []string = strings.Split(str, ArrayIntSep)
	for _, v := range tmp {
		newNumber, err := strconv.ParseInt(v, 10, 64)
		if err == nil {
			integers.Add(newNumber)
		} else {
			return err
		}
	}
	return nil
}

// ArrayInt64Init 将一个[]int64数据增加到当前对像里
func (integers *Integers) ArrayInt64Init(i []int64) error {
	var newI Integers
	for _, v := range i {
		newI.Add(v)
	}
	*integers = newI
	return nil
}

// Find 查找是否存在
func (integers Integers) Find(i int64) bool {
	for _, v := range integers {
		if v == i {
			return true
		}
	}
	return false
}

// Len 统计数量
func (integers Integers) Len() int {
	return len(integers)
}

// Value 用于sql写入时引用值
func (integers Integers) Value() (driver.Value, error) {
	var tmp string
	for _, v := range integers {
		if tmp != "" {
			tmp += ArrayIntSep
		}
		tmp += fmt.Sprintf("%d", v)
	}
	return tmp, nil
}

// ToInt64  输出[]int64
func (integers Integers) ToInt64() []int64 {
	if integers.Len() <= 0 {
		return []int64{}
	}
	var tmp []int64
	for _, v := range integers {
		tmp = append(tmp, v)
	}
	return tmp
}
