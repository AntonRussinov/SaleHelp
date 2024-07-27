package dbmock

import (
	"database/sql/driver"

	"github.com/DATA-DOG/go-sqlmock"
)

func MockAnyArgs(num int) []driver.Value {
	args := make([]driver.Value, num)
	for i := range args {
		args[i] = sqlmock.AnyArg()
	}

	return args
}
