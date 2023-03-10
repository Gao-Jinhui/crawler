package spider

import (
	"crawler/internal/pkg/collector"
	"time"
)

type Context struct {
	Body []byte
	Req  *Request
}

func (c *Context) Output(data interface{}) *collector.DataCell {
	res := &collector.DataCell{
		Storage: c.Req.Task.Storage,
	}
	res.Data = make(map[string]interface{})
	res.Data["Task"] = c.Req.Task.Name
	res.Data["Rule"] = c.Req.RuleName
	res.Data["Data"] = data
	res.Data["Url"] = c.Req.Url
	res.Data["Time"] = time.Now().Format("2006-01-02 15:04:05")
	return res
}
