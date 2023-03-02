package master

import "encoding/json"

type ResourceSpec struct {
	ID           string
	Name         string
	AssignedNode string
	CreationTime int64
}

func encode(s *ResourceSpec) string {
	b, _ := json.Marshal(s)
	return string(b)
}

func decode(ds []byte) (*ResourceSpec, error) {
	var s *ResourceSpec
	err := json.Unmarshal(ds, &s)
	return s, err
}
