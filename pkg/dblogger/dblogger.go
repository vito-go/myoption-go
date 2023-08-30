package dblogger

import (
	"context"
	"github.com/vito-go/mylog"
	"gorm.io/gorm/logger"
	"time"
)

type DBLogger struct {
}

func (D *DBLogger) LogMode(level logger.LogLevel) logger.Interface {
	//TODO implement me
	return D
}

func (D *DBLogger) Info(ctx context.Context, s string, i ...interface{}) {
	mylog.Ctx(ctx).Infof(s, i...)
}

func (D *DBLogger) Warn(ctx context.Context, s string, i ...interface{}) {
	mylog.Ctx(ctx).Warnf(s, i...)
}

func (D *DBLogger) Error(ctx context.Context, s string, i ...interface{}) {
	mylog.Ctx(ctx).Errorf(s, i...)

}

func (D *DBLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sql, rowsAffected := fc()
	timeElapsed := time.Since(begin)
	if err != nil {
		mylog.Ctx(ctx).WithFields("rowsAffected", rowsAffected, "timeElapsed", timeElapsed.String(), "sql", sql).Errorf("SQL ==>> Error: %s", err.Error())
	} else {
		//mylog.Ctx(ctx).WithFields("rowsAffected", rowsAffected, "timeElapsed", timeElapsed.String(), "sql", sql).Info("SQL ==>> Successfully")
	}
}
