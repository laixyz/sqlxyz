package sqlxyz

import (
	"sync"
)

// RWLock 读写锁
var RWLock sync.RWMutex

// MySQLConfigs 全局mysql配置参数变量
var MySQLConfigs = make(map[string]MySQLConfig)
