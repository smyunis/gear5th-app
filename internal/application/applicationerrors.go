package application

import "errors"

var ErrConflictFound = errors.New("duplicate entity found")
var ErrEntityNotFound = errors.New("requested entity not found")
var ErrAuthorization = errors.New("user not allowed to perform use case")
var ErrInvalidData = errors.New("data is invalid or is in unexpected format")
var ErrSiteVerification = errors.New("site is not verified")
var ErrRequirementFailed = errors.New("requirement faild to take this action")