package pholcus_lib

// 基础包
import (
	"github.com/henrylee2cn/pholcus/app/downloader/request" //必需
	"github.com/henrylee2cn/pholcus/common/goquery"         //DOM解析
	// "github.com/henrylee2cn/pholcus/logs"               //信息输出
	. "github.com/henrylee2cn/pholcus/app/spider" //必需
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
	Name:        "妹子图",
	Description: "妹子图 [http://www.meizitu.com/]",
	// Pausetime: 300,
	// Keyin:   KEYIN,
	// Limit:        LIMIT,
	EnableCookie: false,
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {
			for i := 1; i < 4; i++ {
				ctx.AddQueue(&request.Request{
					Url:  "http://www.meizitu.com/a/xinggan_2_" + strconv.Itoa(i) + ".html",
					Rule: "获取列表",
				})
			}
			for i := 1; i < 4; i++ {
				ctx.AddQueue(&request.Request{
					Url:  "http://www.meizitu.com/a/sifang_5_" + strconv.Itoa(i) + ".html",
					Rule: "获取列表",
				})
			}
			//			for i := 1; i < 4; i++ {
			//				ctx.AddQueue(&request.Request{
			//					Url:  "http://www.meizitu.com/a/xiaoqingxin_6_" + strconv.Itoa(i) + ".html",
			//					Rule: "获取列表",
			//				})
			//			}
			//			for i := 1; i < 2000; i++ {
			//				ctx.AddQueue(&request.Request{
			//					Url:  "http://www.meizitu.com/a/" + strconv.Itoa(i) + ".html",
			//					Rule: "详细内容",
			//				})
			//			}
		},

		Trunk: map[string]*Rule{
			"获取列表": {
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()
					query.Find(".pic").
						Each(func(i int, s *goquery.Selection) {
							openurl, _ := s.Find("a").Attr("href")
							imgUrl, _ := s.Find("img").Attr("src")

							ctx.Output(map[string]interface{}{
								"imgURL":  imgUrl,
								"openurl": openurl,
							})
							ctx.AddQueue(&request.Request{
								Url:  openurl,
								Rule: "详细内容",
							})
						})
				},
			},
			"详细内容": {
				ParseFunc: func(ctx *Context) {
					detail := ctx.GetDom().Find("div#picture p")
					detail.Find("img").
						Each(func(i int, s *goquery.Selection) {
							alt, _ := s.Attr("alt")
							imgUrl, _ := s.Attr("src")
							paths := strings.Split(imgUrl, "/")
							len := len(paths)
							fileName := paths[len-4] + paths[len-3] + paths[len-2] + paths[len-1]
							ctx.Output(map[string]interface{}{
								"alt":      alt,
								"fileName": fileName,
								"imgURL":   imgUrl,
							})

							ctx.AddQueue(&request.Request{
								Url:          imgUrl,
								Rule:         "下载图片",
								ConnTimeout:  -1,
								DownloaderID: 0, //图片等多媒体文件必须使用0（surfer surf go原生下载器）
							})
						})
				},
			},
			"下载图片": {
				ParseFunc: func(ctx *Context) {
					var title = ctx.GetUrl()
					paths := strings.Split(title, "/")
					len := len(paths)
					fileName := paths[len-4] + paths[len-3] + paths[len-2] + paths[len-1]
					ctx.Output(map[string]interface{}{
						"title": fileName,
					})
					ctx.FileOutput(fileName)
				},
			},
		},
	},
}
