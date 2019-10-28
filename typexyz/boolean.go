// Package typexyz
// 方便将数据库中tinyint字段值直接转换成bool数据 0为flase 其它数值则为true
// 开发者: 无双 (https://github.com/laixyz)

package typexyz

import (
	"database/sql"
	"database/sql/driver"
)

// Boolean bool匿名类型
type Boolean bool

// Scan 用于sql.Scan
func (boolean *Boolean) Scan(value interface{}) (err error) {
	if value == nil {
		*boolean = false
		return nil
	}
	var nullInt64 sql.NullInt64
	err = nullInt64.Scan(value)
	if err != nil {
		*boolean = false
		return nil
	}
	if nullInt64.Valid == true {
		if nullInt64.Int64 == 0 {
			*boolean = false
		} else {
			*boolean = true
		}
		return nil
	}
	*boolean = false

	return nil
}

// Value 用于sql写入时引用值
func (boolean Boolean) Value() (driver.Value, error) {
	var val int64
	if boolean == false {
		val = 0
	} else {
		val = 1
	}
	return val, nil
}

// IsTrue 适合某一些特定场景，方便使用
func (boolean Boolean) IsTrue() bool {
	return boolean == true
}

// IsFalse 适合某一些特定场景，方便使用
func (boolean Boolean) IsFalse() bool {
	return boolean == false
}
