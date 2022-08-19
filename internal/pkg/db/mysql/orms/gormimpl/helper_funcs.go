package gormimpl

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
)

type OrderColumn struct {
	Field string
	IsDesc bool
}

func OrderColumns(db *gorm.DB, orderColumns []*clause.OrderByColumn) {
	var sqlForOrderBy string
	for _, orderColumn := range orderColumns {
		sqlForOrderBy += orderColumn.Column.Table + "." + orderColumn.Column.Name
		if orderColumn.Desc {
			sqlForOrderBy += " DESC,"
			continue
		}
		sqlForOrderBy += "ASC,"
	}
	db.Order(strings.TrimRight(sqlForOrderBy, ","))
}
