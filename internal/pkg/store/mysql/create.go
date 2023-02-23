package mysql

import (
	"crawler/internal/pkg/model"
	"github.com/pkg/errors"
)

func (sql *SqlClient) CreateBook(books ...*model.Book) error {
	if err := sql.db.Model(&model.Book{}).Create(&books).Error; err != nil {
		return errors.New("create error")
	}
	return nil
}
