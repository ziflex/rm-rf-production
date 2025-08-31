package database

import (
	"github.com/lib/pq"
)

func IsPgErr(err error) (*pq.Error, bool) {
	pqErr, ok := err.(*pq.Error)
	return pqErr, ok
}

func IsDbUniqueViolation(err *pq.Error) bool {
	return err.Code.Name() == "unique_violation"
}

func IsDbForeignKeyViolation(err *pq.Error) bool {
	return err.Code.Name() == "foreign_key_violation"
}
