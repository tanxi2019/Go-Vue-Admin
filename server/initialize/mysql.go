package initialize

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"server/config"
	"server/global"
	"time"
)

// InitMysql 初始化mysql数据库
func InitMysql() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&collation=%s&%s",
		config.Conf.Mysql.Username,
		config.Conf.Mysql.Password,
		config.Conf.Mysql.Host,
		config.Conf.Mysql.Port,
		config.Conf.Mysql.Database,
		config.Conf.Mysql.Charset,
		config.Conf.Mysql.Collation,
		config.Conf.Mysql.Query,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// 禁用外键
		DisableForeignKeyConstraintWhenMigrating: true,
		// 禁用默认事务（提高运行速度）
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			// 使用单数表名
			SingularTable: true,
		},
	})
	if err != nil {
		global.Log.Panicf("初始化mysql数据库异常: %v", err)
		panic(fmt.Errorf("初始化mysql数据库异常: %v", err))
	}

	// 开启mysql日志
	if config.Conf.Mysql.LogMode {
		db.Debug()
	}
	// 全局DB赋值
	global.DB = db

	// 自动迁移表结构
	dbAutoMigrate()

	sqldb, _ := db.DB()

	// SetMaxIdleCons 设置连接池中的最大闲置连接数。
	sqldb.SetMaxIdleConns(10)

	// SetMaxOpenCons 设置数据库的最大连接数量。
	sqldb.SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置连接的最大可复用时间。
	sqldb.SetConnMaxLifetime(10 * time.Second)

	global.Log.Infof("初始化mysql数据库完成!")
}

// 自动迁移表结构
func dbAutoMigrate() {
	err := global.DB.AutoMigrate(
	//&system.User{},
	// &system.Role{},
	//&system.Menu{},
	// &system.Api{},
	// &system.OperationLog{},
	//&system.DictDetail{},
	//&system.Dict{},

	// project
	//&example.Example{},
	)
	if err != nil {
		return
	}
}
