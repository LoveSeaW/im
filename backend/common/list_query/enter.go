package list_query

import (
	"fim_server/common/models"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type Option struct {
	PageInfo models.PageInfo
	Where    *gorm.DB
	Debug    bool
	Joins    string
	Likes    []string
	Preload  []string
	Table    func() (string, any)
	Groups   []string
}

func ListQuery[T any](db *gorm.DB, model T, option Option) (list []T, count int64, err error) {
	if option.Debug {
		db = db.Debug()
	}

	query := db
	if option.Table != nil {
		table, data := option.Table()
		query = query.Table(table, data)
	} else {
		query = query.Model(model)
	}
	query = query.Where(model)

	if option.PageInfo.Key != "" && len(option.Likes) > 0 {
		conditions := make([]string, 0, len(option.Likes))
		args := make([]any, 0, len(option.Likes))
		for _, column := range option.Likes {
			conditions = append(conditions, fmt.Sprintf("%s ILIKE ?", column))
			args = append(args, fmt.Sprintf("%%%s%%", option.PageInfo.Key))
		}
		query = query.Where("("+strings.Join(conditions, " OR ")+")", args...)
	}

	if option.Joins != "" {
		query = query.Joins(option.Joins)
	}

	if option.Where != nil {
		if whereClause, ok := option.Where.Statement.Clauses["WHERE"]; ok && whereClause.Expression != nil {
			query = query.Where(whereClause.Expression)
		}
	}

	for _, group := range option.Groups {
		query = query.Group(group)
	}

	err = query.Count(&count).Error
	if err != nil {
		return
	}

	for _, preload := range option.Preload {
		query = query.Preload(preload)
	}

	if option.PageInfo.Page <= 0 {
		option.PageInfo.Page = 1
	}
	if option.PageInfo.Limit != -1 && option.PageInfo.Limit <= 0 {
		option.PageInfo.Limit = 10
	}

	offset := (option.PageInfo.Page - 1) * option.PageInfo.Limit

	if option.PageInfo.Sort != "" {
		query = query.Order(option.PageInfo.Sort)
	}

	err = query.Limit(option.PageInfo.Limit).Offset(offset).Find(&list).Error
	return
}
