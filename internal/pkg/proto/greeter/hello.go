package greeter

import "context"

type Greeter struct{}

func (g *Greeter) Hello(ctx context.Context, req *Request, rsp *Response) error {
	rsp.Greeting = "Hello " + req.Name

	return nil
}
