package pholcus_lib

// 基础包
import (
	"net/http"

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
	//	"strconv"
	"strings"
	// 其他包
	"fmt"
	// "math"
	// "time"
)

func init() {
	GirlHome.Register()
}

var GirlHome = &Spider{
	Name:         "pron",
	Description:  "pron [http://www.pronhub.com/]",
	Pausetime:    2000,
	Keyin:        KEYIN,
	Limit:        LIMIT,
	EnableCookie: true,
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {
			cookie := strings.Trim(ctx.GetKeyin(), " ")
			fmt.Println([]string{cookie})
			ctx.AddQueue(&request.Request{
				Url:    "https://www.pronhub.com/",
				Header: http.Header{"Cookie": []string{cookie}},
				Rule:   "登录后",
			})
		},
		Trunk: map[string]*Rule{
			"登录后": {
				ParseFunc: func(ctx *Context) {
					cookie := strings.Trim(ctx.GetKeyin(), " ")
					for i := 1; i < 2; i++ {
						ctx.AddQueue(&request.Request{
							Url:    "https://www.pronhub.com/view_video.php?viewkey=ph58fb80e3bbc91",
							Header: http.Header{"Cookie": []string{cookie}},
							Rule:   "获取明细",
						})
					}
				},
			},
			"获取明细": {
				ParseFunc: func(ctx *Context) {
					fmt.Printf("获取明细")
					query := ctx.GetDom()
					tabs := query.Find(".video-actions-container").Find(".video-actions-tabs")
					query.Find(".video-actions-container").Find(".video-actions-tabs").Each(func(i int, s *goquery.Selection) {
						fmt.Println("HTML:")
						fmt.Println(s.Html())
					})
					tabs.Find(".video-actions-tab").Eq(3).Find("a").Each(func(i int, s *goquery.Selection) {
						fmt.Println("HTMLTAB:")
						fmt.Println(s.Html())
					})
				},
			},
		},
	},
}
