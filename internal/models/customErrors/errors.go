package customErrors

import "errors"

const MY_SQL_DUPLICATE_CODE = 1062
const DUPLICATE_PET_TYPE_KEY = "petType_uc_name"

const MY_SQL_CONSTRAINT_CODE = 3819

var (
	ErrDuplicatePetType  = errors.New("models: duplicate pet type")
	CheckConstraintError = errors.New("models: check constraint fail")
)
