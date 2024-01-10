package customErrors

import "errors"

const MY_SQL_DUPLICATE_CODE = 1062
const DUPLICATE_PET_TYPE_KEY = "petType_uc_name"

var (
	ErrDuplicatePetType = errors.New("models: duplicate pet type")
)
