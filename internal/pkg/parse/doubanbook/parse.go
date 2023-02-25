package doubanbook

import (
	"crawler/internal/pkg/collect"
	"crawler/internal/pkg/model"
	"go.uber.org/zap"
	"regexp"
	"strconv"
)

const regexpStr = `<a href="([^"]+)" class="tag">([^<]+)</a>`

func ParseTag(ctx *collect.Context) (collect.ParseResult, error) {
	re := regexp.MustCompile(regexpStr)

	matches := re.FindAllSubmatch(ctx.Body, -1)
	result := collect.ParseResult{}

	for _, m := range matches {
		result.Requests = append(
			result.Requests, &collect.Request{
				Method:   "GET",
				Task:     ctx.Req.Task,
				Url:      "https://book.douban.com" + string(m[1]),
				Depth:    ctx.Req.Depth + 1,
				RuleName: "书籍列表",
			})
	}
	zap.S().Debugln("parse book tag,count:", len(result.Requests))
	// 在添加limit之前，临时减少抓取数量,防止被服务器封禁
	result.Requests = result.Requests[:5]
	return result, nil
}

const BooklistRe = `<a.*?href="([^"]+)" title="([^"]+)"`

func ParseBookList(ctx *collect.Context) (collect.ParseResult, error) {
	re := regexp.MustCompile(BooklistRe)
	matches := re.FindAllSubmatch(ctx.Body, -1)
	result := collect.ParseResult{}
	for _, m := range matches {
		req := &collect.Request{
			Method:   "GET",
			Task:     ctx.Req.Task,
			Url:      string(m[1]),
			Depth:    ctx.Req.Depth + 1,
			RuleName: "书籍简介",
		}
		req.TmpData = &collect.Temp{}
		req.TmpData.Set("book_name", string(m[2]))
		result.Requests = append(result.Requests, req)
	}
	// 在添加limit之前，临时减少抓取数量,防止被服务器封禁
	result.Requests = result.Requests[:5]

	return result, nil
}

var autoRe = regexp.MustCompile(`<span class="pl"> 作者</span>:[\d\D]*?<a.*?>([^<]+)</a>`)
var public = regexp.MustCompile(`<span class="pl">出版社:</span>([^<]+)<br/>`)
var pageRe = regexp.MustCompile(`<span class="pl">页数:</span> ([^<]+)<br/>`)
var priceRe = regexp.MustCompile(`<span class="pl">定价:</span>([^<]+)<br/>`)
var scoreRe = regexp.MustCompile(`<strong class="ll rating_num " property="v:average">([^<]+)</strong>`)
var intoRe = regexp.MustCompile(`<div class="intro">[\d\D]*?<p>([^<]+)</p></div>`)

func ParseBookDetail(ctx *collect.Context) (collect.ParseResult, error) {
	bookName := ctx.Req.TmpData.Get("book_name")
	page, _ := strconv.Atoi(ExtraString(ctx.Body, pageRe))

	//book := map[string]interface{}{
	//	"Name":      bookName,
	//	"Author":    ExtraString(ctx.Body, autoRe),
	//	"Page":      page,
	//	"Publisher": ExtraString(ctx.Body, public),
	//	"Score":     ExtraString(ctx.Body, scoreRe),
	//	"Price":     ExtraString(ctx.Body, priceRe),
	//	"Intro":     ExtraString(ctx.Body, intoRe),
	//	"Url":       ctx.Req.Url,
	//}
	book := &model.Book{
		Name:      bookName.(string),
		Author:    ExtraString(ctx.Body, autoRe),
		Page:      page,
		Publisher: ExtraString(ctx.Body, public),
		Score:     ExtraString(ctx.Body, scoreRe),
		Price:     ExtraString(ctx.Body, priceRe),
		Intro:     ExtraString(ctx.Body, intoRe),
		Url:       ctx.Req.Url,
	}
	data := ctx.Output(book)

	result := collect.ParseResult{
		Items: []interface{}{data},
	}

	return result, nil
}

func ExtraString(contents []byte, re *regexp.Regexp) string {

	match := re.FindSubmatch(contents)

	if len(match) >= 2 {
		return string(match[1])
	} else {
		return ""
	}
}
