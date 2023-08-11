package application

import "errors"


var ErrConflictFound = errors.New("duplicate entity found")
var ErrEntityNotFound = errors.New("requested entity not found")
