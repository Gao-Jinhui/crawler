package collect

import (
	"errors"
	"time"
)

type Request struct {
	Url       string
	Cookie    string
	Depth     int
	MaxDepth  int
	WaitTime  time.Duration
	ParseFunc func([]byte, *Request) ParseResult
}

func NewCollectRequest(url, cookie string, depth, maxDepth int, waitTime time.Duration, parseFunc func([]byte, *Request) ParseResult) *Request {
	return &Request{
		Url:       url,
		Cookie:    cookie,
		Depth:     depth,
		MaxDepth:  maxDepth,
		WaitTime:  waitTime,
		ParseFunc: parseFunc,
	}
}

func (req *Request) CheckDepth() error {
	if req.Depth > req.MaxDepth {
		return errors.New("Max Depth limit reached ")
	}
	return nil
}
