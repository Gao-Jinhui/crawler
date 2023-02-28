package spider

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"math/rand"
	"time"
)

type Request struct {
	Task     *Task
	Url      string
	Depth    int64
	Method   string
	Priority int
	RuleName string
	TmpData  *Temp
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

func (r *Request) Fetch() ([]byte, error) {
	if err := r.Task.Limit.Wait(context.Background()); err != nil {
		return nil, err
	}
	// 随机休眠，模拟人类行为
	sleeptime := rand.Int63n(r.Task.WaitTime * 1000)
	time.Sleep(time.Duration(sleeptime) * time.Millisecond)
	return r.Task.Fetcher.Get(r)
}
