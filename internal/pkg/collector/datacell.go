package collector

type DataCell struct {
	Data map[string]interface{}
}

func (d *DataCell) GetTableName() string {
	return d.Data["Task"].(string)
}

func (d *DataCell) GetTaskName() string {
	return d.Data["Task"].(string)
}
