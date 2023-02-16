package collect

import (
	"errors"
)

type Request struct {
	Task      *Task
	Url       string
	Depth     int
	ParseFunc func([]byte, *Request) ParseResult
}

func NewCollectRequest(url string, depth int, parseFunc func([]byte, *Request) ParseResult, task *Task) *Request {
	return &Request{
		Url:       url,
		Depth:     depth,
		ParseFunc: parseFunc,
		Task:      task,
	}
}

func (req *Request) CheckDepth() error {
	if req.Depth > req.Task.MaxDepth {
		return errors.New("Max Depth limit reached ")
	}
	return nil
}
