package doubanbook

import (
	"crawler/internal/pkg/model"
	"crawler/internal/pkg/spider"
	"go.uber.org/zap"
	"regexp"
	"strconv"
)

const regexpStr = `<a href="([^"]+)" class="tag">([^<]+)</a>`

func ParseTag(ctx *spider.Context) (spider.ParseResult, error) {
	re := regexp.MustCompile(regexpStr)

	matches := re.FindAllSubmatch(ctx.Body, -1)
	res := spider.ParseResult{}

	for _, m := range matches {
		res.Requests = append(
			res.Requests, &spider.Request{
				Method:   "GET",
				Task:     ctx.Req.Task,
				Url:      "https://book.douban.com" + string(m[1]),
				Depth:    ctx.Req.Depth + 1,
				RuleName: "书籍列表",
			})
	}
	zap.S().Debugln("parse book tag,count:", len(res.Requests))
	// 在添加limit之前，临时减少抓取数量,防止被服务器封禁
	//res.Requests = res.Requests[:5]
	return res, nil
}

const BooklistRe = `<a.*?href="([^"]+)" title="([^"]+)"`

func ParseBookList(ctx *spider.Context) (spider.ParseResult, error) {
	re := regexp.MustCompile(BooklistRe)
	matches := re.FindAllSubmatch(ctx.Body, -1)
	res := spider.ParseResult{}
	for _, m := range matches {
		req := &spider.Request{
			Method:   "GET",
			Task:     ctx.Req.Task,
			Url:      string(m[1]),
			Depth:    ctx.Req.Depth + 1,
			RuleName: "书籍简介",
		}
		req.TmpData = &spider.Temp{}
		req.TmpData.Set("book_name", string(m[2]))
		res.Requests = append(res.Requests, req)
	}
	// 在添加limit之前，临时减少抓取数量,防止被服务器封禁
	//res.Requests = res.Requests[:5]

	return res, nil
}

var autoRe = regexp.MustCompile(`<span class="pl"> 作者</span>:[\d\D]*?<a.*?>([^<]+)</a>`)
var public = regexp.MustCompile(`<span class="pl">出版社:</span>([^<]+)<br/>`)
var pageRe = regexp.MustCompile(`<span class="pl">页数:</span> ([^<]+)<br/>`)
var priceRe = regexp.MustCompile(`<span class="pl">定价:</span>([^<]+)<br/>`)
var scoreRe = regexp.MustCompile(`<strong class="ll rating_num " property="v:average">([^<]+)</strong>`)
var intoRe = regexp.MustCompile(`<div class="intro">[\d\D]*?<p>([^<]+)</p></div>`)

func ParseBookDetail(ctx *spider.Context) (spider.ParseResult, error) {
	bookName := ctx.Req.TmpData.Get("book_name")
	page, _ := strconv.Atoi(ExtraString(ctx.Body, pageRe))

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

	res := spider.ParseResult{
		Items: []interface{}{data},
	}

	return res, nil
}

func ExtraString(contents []byte, re *regexp.Regexp) string {

	match := re.FindSubmatch(contents)

	if len(match) >= 2 {
		return string(match[1])
	} else {
		return ""
	}
}
