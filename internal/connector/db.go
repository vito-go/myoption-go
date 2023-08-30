package connector

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/vito-go/mylog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"myoption/conf"
	"myoption/pkg/dblogger"
	"time"
)

func OpenDB(dbConf conf.DBConf) (*sql.DB, error) {
	db, err := sql.Open(dbConf.DriverName, dbConf.Dsn)
	if err != nil {
		return nil, err
	}
	const connectTimeOut = time.Second * 3
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeOut)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("dbConf: %+v err: %s", dbConf, err.Error())
	}
	mylog.Ctx(context.TODO()).WithField("dbConf", dbConf).Info("数据库已链接")
	return db, nil
}

func OpenGromDB(dbConf conf.DBConf) (*gorm.DB, error) {
	switch dbConf.DriverName {
	case "postgres":
	default:
		panic("暂不支持的数据库类型")
	}
	gdb, err := gorm.Open(postgres.Open(dbConf.Dsn),
		&gorm.Config{
			Logger: &dblogger.DBLogger{},
			// https://gorm.io/zh_CN/docs/create.html
			CreateBatchSize: 1000, // 批量插入
		})
	if err != nil {
		return nil, err
	}
	db, err := gdb.DB()
	if err != nil {
		return nil, err
	}
	const connectTimeOut = time.Second * 3
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeOut)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("dbConf: %+v err: %s", dbConf, err.Error())
	}
	mylog.Ctx(context.TODO()).WithField("dbConf", dbConf).Info("数据库已链接")
	return gdb, nil
}
