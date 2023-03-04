package master

import "fmt"

const (
	RESOURCEPATH = "/resources"

	ADDRESOURCE = iota
	DELETERESOURCE

	//ServiceName = "go.micro.server.worker"
)

func getResourcePath(name string) string {
	return fmt.Sprintf("%s/%s", RESOURCEPATH, name)
}
