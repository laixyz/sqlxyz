package sqlxyz

import (
	"errors"
	"fmt"
	"net/url"

	// 加载mysql官方库
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// MySQLConfig mysql连接参数
type MySQLConfig struct {
	ID              string `ini:"-"`
	User            string `ini:"user"`
	Password        string `ini:"password"`
	Host            string `ini:"host"`
	Port            string `ini:"port"`
	DB              string `ini:"db"`
	Charset         string `ini:"charset"`
	Loc             string `ini:"loc"`
	ReadTimeout     string `ini:"read_timeout"`  //readTimeout=30s
	WriteTimeout    string `ini:"write_timeout"` //writeTimeout=30s
	Timeout         string `ini:"timeout"`       //timeout=30s
	MaxIdleConns    int    `ini:"max_idle"`
	MaxOpenConns    int    `ini:"max_open"`
	ConnMaxLifetime string `ini:"max_life_time"`
	ParseTime       bool   `ini:"-"`
}

// String 参数输出dsn字符串
func (mc *MySQLConfig) String() string {
	var dsn string
	if len(mc.Port) == 0 {
		mc.Port = "3306"
	}
	if mc.Timeout == "" {
		mc.Timeout = "5s"
	}
	if mc.ReadTimeout == "" {
		mc.ReadTimeout = "2s"
	}
	if mc.WriteTimeout == "" {
		mc.WriteTimeout = "2s"
	}
	if mc.ConnMaxLifetime == "" {
		mc.ConnMaxLifetime = "60s"
	}

	if len(mc.Loc) == 0 {
		mc.Loc = "Asia/Shanghai"
	}
	if len(mc.Charset) == 0 {
		mc.Charset = "utf8"
	}

	params := url.Values{}

	params.Add("charset", mc.Charset)
	params.Add("loc", mc.Loc)
	params.Add("timeout", mc.Timeout)
	params.Add("readTimeout", mc.ReadTimeout)
	params.Add("writeTimeout", mc.WriteTimeout)

	params.Add("parseTime", "true") //parseTime=true 会将date 或 datetime 字段自动专换成 time.Time

	var strParam string
	if len(params) > 0 {
		strParam = ("?" + params.Encode())
	}

	if len(mc.DB) > 0 {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s%s", mc.User, mc.Password, mc.Host, mc.Port, mc.DB, strParam)
	} else {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", mc.User, mc.Password, mc.Host, mc.Port, strParam)
	}
	return dsn
}

// ConnectCheck 检查连接
func (mc *MySQLConfig) ConnectCheck() error {
	db, err := sqlx.Open("mysql", mc.String())
	if err != nil {
		return errors.New(mc.ID + "connect test:" + err.Error() + ", end")
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		return errors.New(mc.ID + " ping test:" + err.Error() + ", end")
	}
	return nil
}
