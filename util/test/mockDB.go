package test

import (
	"database/sql/driver"
	"log"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/bindu-server/config"
)

// AnyTime To match time for test sqlmock queries
type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

// SetupMockDB setups the mock sql db
func SetupMockDB() sqlmock.Sqlmock {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}

	config.SetupDB(db)

	return mock
}

//ExpectationsMet checks if all the expectations are fulfilled
func ExpectationsMet(t *testing.T, mock sqlmock.Sqlmock) {
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
