package spider

type Fetcher interface {
	Get(resq *Request) ([]byte, error)
}
