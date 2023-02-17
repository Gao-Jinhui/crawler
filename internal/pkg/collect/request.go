package collect

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
)

type Request struct {
	unique    string
	Task      *Task
	Url       string
	Depth     int
	Method    string
	Priority  int
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

func (r *Request) Unique() string {
	block := md5.Sum([]byte(r.Url + r.Method))
	return hex.EncodeToString(block[:])
}
