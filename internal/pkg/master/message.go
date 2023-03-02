package master

type Message struct {
	Cmd   Command
	Specs []*ResourceSpec
}
