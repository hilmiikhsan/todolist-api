package utils

import (
	"errors"
	"fmt"
	"todolist-api/constants"
)

// ErrDataNotFound function to handle data not found
// Params:
// m: error string
// Returns *constants.ErrorResponse
func ErrDataNotFound(m interface{}) error {
	errs := fmt.Sprintf(constants.ErrDataNotFound, m)

	return errors.New(errs)
}
