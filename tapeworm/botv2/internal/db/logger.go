package db

import (
	"context"
	"fmt"

	"github.com/jaxsax/projects/tapeworm/botv2/internal/logging"
	sqldblogger "github.com/simukti/sqldb-logger"
)

type QueryLogger struct {
}

func (q *QueryLogger) Log(ctx context.Context, level sqldblogger.Level, msg string, data map[string]interface{}) {
	if level == sqldblogger.LevelError {
		logging.FromContext(ctx).Error(nil, fmt.Sprintf("mysql_error:%s", msg), "data", data)
		return
	}

	if level == sqldblogger.LevelDebug {
		return
	}

	delete(data, "time")
	delete(data, "conn_id")

	logging.FromContext(ctx).V(0).Info(fmt.Sprintf("mysql:%s", msg), "data", data)
}
