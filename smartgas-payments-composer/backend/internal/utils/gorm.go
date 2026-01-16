package utils

import (
	"math"
	"smartgas-payment/internal/schemas"
	"strings"

	"gorm.io/gorm"
)

func Paginate(pagination *schemas.Pagination, model any, db *gorm.DB, query string, filters any, tables ...string) func(db *gorm.DB) *gorm.DB {
	if query == "" {
		db.Model(model).Where(filters).Count(&pagination.TotalRows)
	} else {
		tx := db.Model(model)
		for _, table := range tables {
			if strings.Contains(table, ".") && !strings.Contains(table, " ") {
				tx.Preload(table)
			} else {
				tx.InnerJoins(table)
			}
		}

		tx.Where(query, filters).Count(&pagination.TotalRows)
	}

	pagination.TotalPages = int(math.Ceil(float64(pagination.TotalRows) / float64(pagination.GetLimit())))
	return func(db *gorm.DB) *gorm.DB {

		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit())
	}
}
