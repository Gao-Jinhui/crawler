package mysql

import (
	"crawler/internal/pkg/collector"
	"crawler/internal/pkg/model"
)

func (sql *SqlClient) Save(datas ...*collector.DataCell) error {
	for _, data := range datas {
		d := data.Data["Data"].(*model.Book)
		//d := data.Data
		if err := sql.CreateBook(d); err != nil {
			return err
		}
	}
	return nil
}
