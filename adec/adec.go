package ken_pholcus_lib

// 基础包
import (
	"github.com/henrylee2cn/pholcus/app/downloader/request" //必需
	. "github.com/henrylee2cn/pholcus/app/spider"           //必需
	"github.com/henrylee2cn/pholcus/common/goquery"         //DOM解析

	"github.com/henrylee2cn/pholcus/logs" //信息输出
	// . "github.com/henrylee2cn/pholcus/app/spider/common"          //选用

	// net包
	// "net/http" //设置http.Header
	// "net/url"

	// 编码包
	// "encoding/xml"
	// "encoding/json"

	// 字符串处理包
	// "regexp"
	"strconv"
	"strings"
	// 其他包
	// "fmt"
	// "math"
	// "time"
)

func init() {
	GirlHome.Register()
}

var GirlHome = &Spider{
	Name:        "adec",
	Description: "adec [https://adcabi.com/index.html]",
	//	Pausetime:    2000,
	Keyin:        KEYIN,
	Limit:        LIMIT,
	EnableCookie: true,
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {
			ctx.Aid(map[string]interface{}{"loop": [2]int{0, 1}, "Rule": "生成请求"}, "生成请求")
		},
		Trunk: map[string]*Rule{
			"生成请求": {
				AidFunc: func(ctx *Context, aid map[string]interface{}) interface{} {
					var paramsStr string = ctx.GetKeyin()
					if len(paramsStr) > 5 {
						params := strings.Split(paramsStr, "@")
						if len(params) < 2 {
							ctx.AddQueue(&request.Request{
								Url:    paramsStr,
								Method: "GET",
								Rule:   "列表页面",
							})
						} else {
							startIndex, error := strconv.Atoi(params[1])
							endIndex, error := strconv.Atoi(params[2])
							if error != nil {
							}
							searchStr := params[0]
							for i := startIndex; i < endIndex; i++ {
								ctx.AddQueue(&request.Request{
									Url:  searchStr + strconv.Itoa(i) + ".html",
									Rule: "列表页面",
								})
							}
						}
					} else {

						ctx.AddQueue(&request.Request{
							Url:    "https://adcabi.com/index.html",
							Method: "GET",
							Rule:   "列表页面",
						})
					}

					// logs.Log.Debug("详细内容startIndex%v,%v,%v", startIndex, endIndex, searchStr)
					return nil
				},
				ParseFunc: func(ctx *Context) {
				},
			},
			"列表页面": {
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()
					list := query.Find(".list-item")
					list.Each(func(i int, s *goquery.Selection) {
						logs.Log.Debug("s:%v", s.Text())
						url, _ := s.Find("a").Attr("href")
						var title = s.Find(".name").Text()
						pic, _ := s.Find(".img-responsive").Attr("src")
						var author = s.Find(".author").Text()
						var date = s.Find(".date").Text()

						if strings.Index(url, "http") < 0 {
							ctx.AddQueue(&request.Request{
								Url:  "https://adcabi.com" + url,
								Rule: "获取明细",
								Temp: map[string]interface{}{"title": title, "url": url, "pic": pic, "author": author, "date": date},
							})
						}
					})

				},
			},
			"获取明细": {
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()
					var titleName = ctx.GetTemp("title", "").(string)
					var pic = ctx.GetTemp("pic", "").(string)
					var date = ctx.GetTemp("date", "").(string)
					var author = ctx.GetTemp("author", "").(string)
					var url = ctx.GetTemp("url", "").(string)
					downUrl, _ := query.Find(".share_div").Find("a").Attr("href")

					logs.Log.Debug("title%v", titleName)
					// logs.Log.Debug("pic%v", pic)
					// logs.Log.Debug("date%v", date)
					// logs.Log.Debug("author%v", author)
					// logs.Log.Debug("url:%v", url)
					ctx.Output(map[string]interface{}{
						"title":   titleName,
						"downUrl": downUrl,
						"url":     url,
						"pic":     pic,
						"author":  author,
						"date":    date,
					})

					ctx.AddQueue(&request.Request{
						Url:          downUrl,
						Rule:         "下载视频",
						Temp:         map[string]interface{}{"titleName": titleName},
						ConnTimeout:  -1,
						DownloaderID: 0, //图片等多媒体文件必须使用0（surfer surf go原生下载器）
					})
				},
			},
			"下载视频": {
				ParseFunc: func(ctx *Context) {
					var titleName = ctx.GetTemp("title", "").(string)
					ctx.FileOutput(titleName)
				},
			},
		},
	},
}
