package constants

import "errors"

const (
	ErrDataNotFound = "Activity with ID %v Not Found"
)

var (
	ErrTitleCannotBeNull = errors.New("title cannot be null")
	ErrBeginTransaction  = errors.New("Failed To Begin Transaction")
)
