package sqlxyz

import (
	"errors"
	"time"

	// 加载mysql官方库
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// SQLXConfig sqlx使用的全局参数变量
var SQLXConfig = make(map[string]MySQLConfig)

// SQLXConnects sqlx使用的全局连接变量数组
var SQLXConnects = make(map[string]*sqlx.DB)

// Close 关闭所有注册打开的连接
func Close() error {
	RWLock.RLock()
	for ConnectID, conn := range SQLXConnects {
		RWLock.RUnlock()
		if conn != nil {
			if err := conn.Close(); err != nil {
				return err
			}
			RWLock.Lock()
			delete(SQLXConnects, ConnectID)
			RWLock.Unlock()
			return nil
		}
	}
	RWLock.RUnlock()
	return nil
}

// Register 注册一个数据库配置并指定连接ID
func Register(ConnectID string, mc MySQLConfig) error {
	if ConnectID == "" {
		mc.ID = "default"
		ConnectID = "default"
	}
	var db *sqlx.DB
	var err error
	db, err = sqlx.Open("mysql", mc.String())
	if err != nil {
		return err
	}
	db.Close()
	RWLock.Lock()
	SQLXConfig[ConnectID] = mc
	RWLock.Unlock()
	return nil
}

// Using 使用一个连接
func Using(ConnectID string) (db *sqlx.DB, err error) {
	if ConnectID == "" {
		ConnectID = "default"
	}
	RWLock.RLock()
	mc, ok := SQLXConfig[ConnectID]
	RWLock.RUnlock()
	if !ok {
		RWLock.RLock()
		mc, ok = SQLXConfig[ConnectID]
		RWLock.RUnlock()
		if !ok {
			err = errors.New("the mysql server has not register [ connect id: " + ConnectID + " ]")
			return
		}
		err = Register(ConnectID, mc)
		if err != nil {
			return
		}
	}
	RWLock.RLock()
	db, ok = SQLXConnects[ConnectID]
	RWLock.RUnlock()
	if ok {
		dbStats := db.Stats()
		if dbStats.OpenConnections <= 0 {
			err = db.Ping()
			if err == nil {
				return db, nil
			}
		}
	}
	db, err = sqlx.Open("mysql", mc.String())
	if err == nil {
		db.SetMaxIdleConns(mc.MaxIdleConns)
		db.SetMaxOpenConns(mc.MaxOpenConns)
		d, err := time.ParseDuration(mc.ConnMaxLifetime)
		if err == nil {
			db.SetConnMaxLifetime(d)
		}
		RWLock.Lock()
		SQLXConnects[ConnectID] = db
		RWLock.Unlock()
	}
	return db, nil
}
