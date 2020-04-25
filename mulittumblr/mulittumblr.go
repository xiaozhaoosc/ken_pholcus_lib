package ken_pholcus_lib

// 基础包
import (
	"encoding/json"
	"fmt"

	//	"math"

	"github.com/henrylee2cn/pholcus/app/downloader/request" //必需
	. "github.com/henrylee2cn/pholcus/app/spider"           //必需

	//	"github.com/henrylee2cn/pholcus/common/goquery"         //DOM解析
	//	"github.com/henrylee2cn/pholcus/logs"                   //信息输出
	// . "github.com/henrylee2cn/pholcus/app/spider/common"          //选用

	// net包
	//	"net/http" //设置http.Header
	// "net/url"

	// 编码包
	// "encoding/xml"
	// "encoding/json"

	// 字符串处理包
	//	"regexp"
	//	"strconv"
	"strings"
	// 其他包
	//	"fmt"
	// "math"
	// "time"
)

func init() {
	GirlHome.Register()
}

type IndexPage struct {
	meta     string
	response string
}

type ResponseData struct {
	blog        string
	posts       string
	total_posts int64
}

type original struct {
	Url    string
	Width  int64
	Height int64
}
type photo struct {
	Original_size original
}
type postsObj struct {
	Blog_name string
	Post_url  string
	Photos    []photo
}

type metaObj struct {
	Status float64
	Msg    string
}
type responseObj struct {
	Posts       []postsObj
	Total_posts int
}
type dataObj struct {
	Meta     metaObj
	Response responseObj
}

var GirlHome = &Spider{
	Name:        "mulit tumblr API",
	Description: "mulit tumblr API [https://www.tumblr.com/]",
	//	Pausetime:    2000,
	Keyin:        KEYIN,
	Limit:        LIMIT,
	EnableCookie: true,
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {
			var paramsStr = ctx.GetKeyin()
			params := strings.Split(paramsStr, "@")
			for i := 0; i < len(params)-1; i++ {
				searchStr := params[i+2]
				var openUrl = "https://api.tumblr.com/v2/blog/" + searchStr + "/posts/photo?api_key=nXcMfImiJuDIhaO7qNT1VF234UhRID8yab3f5tvUoOhCMDUk3y&offset=" + params[0] + "&limit=" + params[1]
				ctx.AddQueue(&request.Request{
					Url:    openUrl,
					Method: "GET",
					Rule:   "首页",
				})
			}
		},
		Trunk: map[string]*Rule{
			"首页": {
				ParseFunc: func(ctx *Context) {
					query := ctx.GetText()
					var msg map[string]interface{}
					err := json.Unmarshal([]byte(query), &msg)
					if err != nil {
						fmt.Println("fuck", err)
					}
					var do dataObj
					if err := json.Unmarshal([]byte(query), &do); err != nil {
						fmt.Println("data error", err)
					}
					for _, post := range do.Response.Posts {
						for _, photo := range post.Photos {
							var picUrl = photo.Original_size.Url

							ctx.AddQueue(&request.Request{
								Url:          picUrl,
								Rule:         "下载图片",
								ConnTimeout:  -1,
								DownloaderID: 0, //图片等多媒体文件必须使用0（surfer surf go原生下载器）
							})
						}
					}
				},
			},
			"获取图片": {
				ParseFunc: func(ctx *Context) {
					query := ctx.GetText()
					var msg map[string]interface{}
					err := json.Unmarshal([]byte(query), &msg)
					if err != nil {
						fmt.Println("fuck", err)
					}

					var do dataObj

					if err := json.Unmarshal([]byte(query), &do); err != nil {
						fmt.Println("data error", err)
					}

					for _, post := range do.Response.Posts {
						for _, photo := range post.Photos {
							var picUrl = photo.Original_size.Url

							ctx.AddQueue(&request.Request{
								Url:          picUrl,
								Rule:         "下载图片",
								ConnTimeout:  -1,
								DownloaderID: 0, //图片等多媒体文件必须使用0（surfer surf go原生下载器）
							})
						}
					}

				},
			},
			"下载图片": {
				ParseFunc: func(ctx *Context) {
					paths := strings.Split(ctx.GetUrl(), "/")
					len := len(paths)
					fileName := paths[len-1]
					ctx.FileOutput(fileName) // 等价于ctx.AddFile("baidu")
				},
			},
		},
	},
}
