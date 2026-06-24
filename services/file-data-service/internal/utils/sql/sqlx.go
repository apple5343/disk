package sqlutil

import (
	"errors"

	"github.com/lib/pq"
)

func IsUniqueViolationSQL(err error) bool {
	pqErr := &pq.Error{}
	if errors.As(err, &pqErr) {
		return pqErr.Code == "23505"
	}
	return false
}
