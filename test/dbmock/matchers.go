package dbmock

import (
	"database/sql/driver"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

type (
	timeMatcher      struct{}
	timeOrNilMatcher struct{}
)

func MatchTime() sqlmock.Argument {
	return timeMatcher{}
}

func MatchTimeOrNil() sqlmock.Argument {
	return timeOrNilMatcher{}
}

func (timeOrNilMatcher) Match(val driver.Value) bool {
	if val == nil {
		return true
	}

	_, ok := val.(time.Time)

	return ok
}

func (timeMatcher) Match(val driver.Value) bool {
	_, ok := val.(time.Time)

	return ok
}
