package ken_pholcus_lib

// 基础包
import (
	"github.com/henrylee2cn/pholcus/app/downloader/request" //必需
	"github.com/henrylee2cn/pholcus/common/goquery"         //DOM解析

	. "github.com/henrylee2cn/pholcus/app/spider" //必需
	// "github.com/henrylee2cn/pholcus/logs"         //信息输出

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
	Name:        "mzitu",
	Description: "mzitu[https://www.mzitu.com/]",
	// Pausetime: 300,https://www.mzitu.com/jiepai/comment-page-1/#comments
	Keyin: KEYIN,
	// Limit:        LIMIT,
	EnableCookie: false,
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {
			ctx.Aid(map[string]interface{}{"loop": [2]int{0, 1}, "Rule": "生成请求"}, "生成请求")

			// logs.Log.Debug("详细内容startIndex%v,%v,%v", startIndex, endIndex, searchStr)
		},

		Trunk: map[string]*Rule{
			"生成请求": {
				AidFunc: func(ctx *Context, aid map[string]interface{}) interface{} {
					var paramsStr = ctx.GetKeyin()
					params := strings.Split(paramsStr, "@")
					startIndex, error := strconv.Atoi(params[0])
					endIndex, error := strconv.Atoi(params[1])
					if error != nil {
					}
					searchStr := params[2]
					for i := startIndex; i < endIndex; i++ {
						ctx.AddQueue(&request.Request{
							Url:  "https://www.mzitu.com/" + searchStr + "/comment-page-" + strconv.Itoa(i) + "/#comments",
							Rule: "详细内容",
						})
					}
					// for i := 1; i < 89; i++ {
					// 	ctx.AddQueue(&request.Request{
					// 		Url:  "https://www.mzitu.com/jiepai/comment-page-" + strconv.Itoa(i) + "/#comments",
					// 		Rule: "详细内容",
					// 	})
					// }
					return nil
				},
				ParseFunc: func(ctx *Context) {
				},
			},
			"详细内容": {
				ParseFunc: func(ctx *Context) {
					detail := ctx.GetDom().Find("img.lazy")
					detail.Each(func(i int, s *goquery.Selection) {
						imgUrl, _ := s.Attr("data-original")
						paths := strings.Split(imgUrl, "/")
						len := len(paths)
						fileName := paths[len-2] + paths[len-1]
						ctx.Output(map[string]interface{}{
							"fileName": fileName,
							"imgURL":   imgUrl,
						})

						// ctx.AddQueue(&request.Request{
						// 	Url:          imgUrl,
						// 	Rule:         "下载图片",
						// 	ConnTimeout:  -1,
						// 	DownloaderID: 0, //图片等多媒体文件必须使用0（surfer surf go原生下载器）
						// })
					})
				},
			},
			"下载图片": {
				ParseFunc: func(ctx *Context) {
					var title = ctx.GetUrl()
					paths := strings.Split(title, "/")
					len := len(paths)
					fileName := paths[len-2] + paths[len-1]
					ctx.Output(map[string]interface{}{
						"title": fileName,
					})
					// ctx.FileOutput(fileName)
				},
			},
		},
	},
}
