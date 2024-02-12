package customErrors

import (
	"errors"
	"strings"

	"github.com/go-sql-driver/mysql"
)

const mySqlDuplicateKeyCode = 1062
const duplicatePetTypeKey = "petType_uc_name"

const mySqlConstraintCode = 3819

var (
	ErrDuplicatePetType = errors.New("models: duplicate pet type")
	ErrConstraintFail   = errors.New("models: check constraint fail")
)

func HandleMySqlError(err error) error {
	var mySqlError *mysql.MySQLError
	if errors.As(err, &mySqlError) {
		if mySqlError.Number == mySqlDuplicateKeyCode && strings.Contains(mySqlError.Message, duplicatePetTypeKey) {
			return ErrDuplicatePetType
		}

		if mySqlError.Number == mySqlConstraintCode {
			return ErrConstraintFail
		}
	}

	return err
}
