package mysql

import (
	"crawler/internal/pkg/errno"
	"crawler/internal/pkg/model"
)

func (sql *SqlClient) CreateBook(books ...*model.Book) error {
	if err := sql.db.Model(&model.Book{}).Create(&books).Error; err != nil {
		return errno.ErrCreateDocument
	}
	return nil
}
