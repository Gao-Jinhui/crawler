package master

import "fmt"

const (
	RESOURCEPATH = "/resources"
)

func getResourcePath(name string) string {
	return fmt.Sprintf("%s/%s", RESOURCEPATH, name)
}
