package master

import "go-micro.dev/v4/registry"

type NodeSpec struct {
	Node    *registry.Node
	Payload int
}
