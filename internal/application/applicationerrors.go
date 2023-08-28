package application

import "errors"

var ErrConflictFound = errors.New("duplicate entity found")
var ErrEntityNotFound = errors.New("requested entity not found")
var ErrAuthorization = errors.New("user not allowed to perform use case")
