package collector

type Storage interface {
	Save(datas ...*DataCell) error
}
