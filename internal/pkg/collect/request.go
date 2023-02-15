package collect

type Request struct {
	Url       string
	Cookie    string
	ParseFunc func([]byte, *Request) ParseResult
}

type ParseResult struct {
	Requests []*Request
	Items    []interface{}
}

func NewCollectRequest(url, cookie string, parseFunc func([]byte, *Request) ParseResult) *Request {
	return &Request{
		Url:       url,
		Cookie:    cookie,
		ParseFunc: parseFunc,
	}
}
