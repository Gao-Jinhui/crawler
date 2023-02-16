package doubangroup

import (
	"crawler/internal/pkg/collect"
	"regexp"
)

const urlListRe = `(https://www.douban.com/group/topic/[0-9a-z]+/)"[^>]*>([^<]+)</a>`

//const ContentRe = `<div class="rich-content topic-richtext"><p>[\s\S]*?初学[\s\S]*?</p>`
const ContentRe = `<div class="topic-content">[\s\S]*?阳台[\s\S]*?<div`

func ParseURL(contents []byte, req *collect.Request) collect.ParseResult {
	re := regexp.MustCompile(urlListRe)

	matches := re.FindAllSubmatch(contents, -1)
	result := collect.ParseResult{}

	for _, m := range matches {
		u := string(m[1])
		result.Requests = append(result.Requests,
			collect.NewCollectRequest(u, req.Cookie, req.Depth+1, req.MaxDepth, req.WaitTime, GetContent))
	}
	return result
}

func GetContent(contents []byte, req *collect.Request) collect.ParseResult {
	//fmt.Println(url)
	re := regexp.MustCompile(ContentRe)

	ok := re.Match(contents)
	if !ok {
		return collect.ParseResult{
			Items: []interface{}{},
		}
	}

	result := collect.ParseResult{
		Items: []interface{}{req.Url},
	}

	return result
}
