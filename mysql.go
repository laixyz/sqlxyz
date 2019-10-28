package sqlxyz

import (
	"database/sql"
	"errors"
	"fmt"

	//加载mysql官方库
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// Prepare 实现fmt.Sprintf("%d", 1000) 用于防止SQL注入
func Prepare(format string, a ...interface{}) string {
	return fmt.Sprintf(format, a...)
}

// DataBaseExist 查询数据库是否存在
func DataBaseExist(db *sqlx.DB, dbname string) (exists bool, err error) {
	var getdbname string
	err = db.QueryRow("show DATABASES like \"?\"", dbname).Scan(&getdbname)
	switch {
	case err == sql.ErrNoRows:
		return false, nil
	case err != nil:
		return false, err
	default:
		if getdbname == dbname {
			return true, nil
		}
		return false, nil
	}
}

// TableExist 查询数据库的一个表是否存在
func TableExist(db *sqlx.DB, dbname string, tablename string) (exists bool, err error) {
	var result string
	err = db.QueryRow("show tables from ? like \"?\"", dbname, tablename).Scan(&result)
	switch {
	case err == sql.ErrNoRows:
		return false, nil
	case err != nil:
		return false, err
	default:
		if result == dbname {
			return true, nil
		}
		return false, nil
	}
}

// DataBaseCreate 创建数据库
func DataBaseCreate(db *sqlx.DB, dbname string, charset string) (err error) {
	//utf8_bin 区分大小写 而 utf8_general_ci 不区分大小写
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + dbname + " CHARACTER SET " + charset + " COLLATE " + charset + "_general_ci")
	if err != nil {
		return err
	}
	_, err = db.Exec("flush privileges")
	if err != nil {
		return err
	}
	return nil
}

// DataBaseDrop 删除数据库
func DataBaseDrop(db *sqlx.DB, dbname string) (err error) {
	_, err = db.Exec("drop database if exists " + dbname)
	if err != nil {
		return err
	}
	return nil
}

// TableDrop 删除表
func TableDrop(db *sqlx.DB, dbname, tablename string) (err error) {
	_, err = db.Exec("drop table IF EXISTS " + dbname + ".`" + tablename + "`")
	if err != nil {
		return err
	}
	return nil
}

// DataBaseUserDrop 删除用户
func DataBaseUserDrop(db *sqlx.DB, user string, ip string) (err error) {
	sql := ""
	if ip != "localhost" {
		sql = "drop user " + user + "@`" + ip + "`"
	} else {
		sql = "drop user " + user + "@" + ip
	}

	_, err = db.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}

// DataBaseUserCreate 创建数据库用户，并指定拥有数据库所有权限
func DataBaseUserCreate(db *sqlx.DB, dbname string, user string, password string, ip string) (err error) {
	_, err = db.Exec("GRANT Select ON " + dbname + ".* to '" + user + "'@'" + ip + "' IDENTIFIED BY '" + password + "';")
	if err != nil {
		return err
	}
	return nil
}

// Affected 获取sql.exec的结果，得到影响的行数和错误
func Affected(Result sql.Result, ResultError error) (RowsAffected int64, err error) {
	if ResultError != nil {
		return 0, ResultError
	}
	RowsAffected, err = Result.RowsAffected()
	if err != nil {
		return 0, err
	}
	if RowsAffected == 0 {
		return 0, errors.New("没有可操作的数据")
	}
	return RowsAffected, nil
}

// LastInsertId 获取新建记录的自增长主键id值和错误
func LastInsertId(Result sql.Result, ResultError error) (LastInsertId int64, err error) {
	if ResultError != nil {
		return 0, ResultError
	}
	LastInsertId, err = Result.LastInsertId()
	if err != nil {
		return 0, err
	}
	if LastInsertId == 0 {
		return 0, errors.New("没有获取到自增长ID数据")
	}
	return LastInsertId, nil
}
