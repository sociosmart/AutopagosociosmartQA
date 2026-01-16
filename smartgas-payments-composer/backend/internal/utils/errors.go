package utils

import (
	"errors"
	"fmt"
	"net/http"
	"smartgas-payment/internal/dto"
	"smartgas-payment/internal/lang"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

func CheckMySQLError(err error) error {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		return errors.New(mysqlErr.Message)
	}

	return nil
}

func CheckNotFoundRecord(c *gin.Context, err error) bool {

	ok := true
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.GeneralMessage{Detail: lang.NotFoundRecord})
			return false

		}
		c.JSON(http.StatusInternalServerError, dto.GeneralMessage{Detail: lang.InternalServerError})
		return false

	}
	return ok
}

func CheckDuplicatedEntry(err error) (duplicated bool) {
	var mysqlErr *mysql.MySQLError
	if err != nil {
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			duplicated = true
			return
		}
	}

	return
}

func CheckMysqlErrCode(err error, code uint16) (found bool) {
	var mysqlErr *mysql.MySQLError
	if err != nil {
		if errors.As(err, &mysqlErr) && mysqlErr.Number == code {
			found = true
			return
		}
	}

	return
}

func PrintExpectedValues(expected any, got any) string {
	return fmt.Sprintf("Expected %v, got %v", expected, got)
}
