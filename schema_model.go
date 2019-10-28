package sqlxyz

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

// SchemaModel mysql数据库表格模板基础结构
type SchemaModel struct {
	ConnectID string   // 数据库连接 ID
	DB        *sqlx.DB //sqlx.DB对像
}

// Prepare 可以实现基础版的防注入
func (sm SchemaModel) Prepare(format string, a ...interface{}) string {
	return fmt.Sprintf(format, a...)
}

// Ping 用于检查数据库连接是否正常
func (sm *SchemaModel) Ping() (err error) {
	sm.DB, err = Using(sm.ConnectID)
	if err != nil {
		return err
	}
	return nil
}

/*
func (this *SQLXYZ_MODEL) MapDataX(db *sqlx.DB, table, keyField, valueField string, vars ...string) (data map[interface{}]interface{}, err error) {
	err = db.Ping()
	if err != nil {
		return data, err
	}
	var where, order string
	for index, v := range vars {
		if index == 0 {
			where = v
		} else if index == 1 {
			order = v
		}
	}
	if where != "" {
		where = " WHERE " + where
	}
	if order != "" {
		order = " ORDER BY " + order
	}
	rows, err := db.Query("select " + keyField + ", " + valueField + " from " + table + where + order)
	if err != nil {
		return data, err
	}
	defer rows.Close()
	data = make(map[interface{}]interface{})
	for rows.Next() {
		var key, value interface{}
		err = rows.Scan(&key, &value)
		if err != nil {
			return data, err
		}
		//if _, ok := data[key]; !ok {
		//data[key] = value
		fmt.Println(key, value)
		//}
	}
	return data, nil
}
*/
