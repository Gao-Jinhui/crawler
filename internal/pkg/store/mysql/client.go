package mysql

import (
	"crawler/internal/pkg/model"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type SqlClient struct {
	db *gorm.DB
	options
}

func NewSqlClient(dsn string, opts ...Option) *SqlClient {
	options := defaultOptions
	for _, opt := range opts {
		opt(&options)
	}
	s := &SqlClient{options: options}
	var err error
	if s.db, err = newMysqlClient(dsn); err != nil {
		s.logger.Error("new mysql client error",
			zap.String("error", err.Error()))
		return nil
		//todo: replace errors
	}
	return s
}

func newMysqlClient(dsn string) (*gorm.DB, error) {
	DB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	DB.AutoMigrate(model.Book{})
	return DB, err
}
