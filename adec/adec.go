package ken_pholcus_lib

// 基础包
import (
	"github.com/henrylee2cn/pholcus/app/downloader/request" //必需
	. "github.com/henrylee2cn/pholcus/app/spider"           //必需
	"github.com/henrylee2cn/pholcus/common/goquery"         //DOM解析

	//	"github.com/henrylee2cn/pholcus/logs"                   //信息输出
	// . "github.com/henrylee2cn/pholcus/app/spider/common"          //选用

	// net包
	"net/http" //设置http.Header
	// "net/url"

	// 编码包
	// "encoding/xml"
	// "encoding/json"

	// 字符串处理包
	"regexp"
	"strconv"
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
	Name:        "adec",
	Description: "adec [https://adcabi.com/index.html]",
	//	Pausetime:    2000,
	Keyin:        KEYIN,
	Limit:        LIMIT,
	EnableCookie: true,
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {
			ctx.AddQueue(&request.Request{
				Url:    "https://adcabi.com/index.html",
				Method: "GET",
				Rule:   "首页",
			})
		},
		Trunk: map[string]*Rule{
			"首页": {
				ParseFunc: func(ctx *Context) {
					// var paramsStr = ctx.GetKeyin()

					fmt.Printf("test%v", "test")
					query := ctx.GetDom()
					tabs := query.Find(".col-xs-6 col-sm-6 col-md-3")
					fmt.Printf("tabs%v", tabs)
					isOk := false
					tabs.Find(".list-item").Each(func(i int, s *goquery.Selection) {
						if url, ok := s.Attr("href"); ok {
							if isOk {
								return
							}

							isOk = true

							var title = s.Find(".name")
							// var pic = s.Find(".img-responsive").Attr("src")
							fmt.Printf("title%v", title)
							// fmt.Printf("pic%v", pic)
							// "pic":   pic,
							fmt.Printf("url%v", url)
							ctx.Output(map[string]interface{}{
								"title": title,
								"url":   url,
							})
						}

					})

					// ctx.AddQueue(&request.Request{
					// 	Url:  "https://adcabi.com" + url,
					// 	Rule: "home页面",
					// })

				},
			},
			"获取明细": {
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()
					var titleName = ctx.GetTemp("title", "").(string)
					var valid = ctx.GetTemp("valid", "").(string)
					var viewsNum = ctx.GetTemp("viewsNum", "").(string)
					tabs := query.Find(".video-actions-container").Find(".video-actions-tabs")
					isOk := false
					tabs.Find(".video-action-tab").Find("a.downloadBtn").Each(func(i int, s *goquery.Selection) {
						if url, ok := s.Attr("href"); ok {
							if isOk {
								return
							}

							isOk = true

							var title = url
							title = strings.Replace(title, "https://", " ", -1)
							a := strings.Split(title, "?")
							b := strings.Split(a[0], "/")
							fileName := b[2] + b[3] + b[4] + b[5]
							ctx.Output(map[string]interface{}{
								"titleName": titleName,
								"valid":     valid,
								"viewsNum":  viewsNum,
								"title":     fileName,
								"url":       url,
							})

							ctx.AddQueue(&request.Request{
								Url:          url,
								Rule:         "下载视频",
								Temp:         map[string]interface{}{"titleName": titleName + fileName, "valid": valid, "viewsNum": viewsNum},
								ConnTimeout:  -1,
								DownloaderID: 0, //图片等多媒体文件必须使用0（surfer surf go原生下载器）
							})
						}

					})

					relatedVideosCenter := query.Find(".underplayer-thumbs").Find(".videoblock")
					relatedVideosCenter.Each(func(i int, s *goquery.Selection) {
						//评分 73%!<(MISSING)
						valNum := s.Find(".thumbnail-info-wrapper").Find(".value").Text()
						var reg = regexp.MustCompile("[0-9]*")

						var valid = reg.FindString(valNum)
						b, error := strconv.Atoi(valid)
						if error != nil {
							//							logs.Logs("字符串转换成整数失败")
						}
						//观看次数
						viewsNum := s.Find(".thumbnail-info-wrapper").Find(".views").Find("var").Text()
						//						var viewsid = reg.FindString(viewsNum)

						vals := strings.Split(viewsNum, ",")

						viewsNum = strings.Join(vals, "")

						c, error := strconv.Atoi(viewsNum)
						if error != nil {
							//							logs.Logs("字符串转换成整数失败")
						}

						if b > 78 && c > 300000 {
							videoInfo := s.Find("a")
							//详细页面
							openUrl, _ := videoInfo.Attr("href")
							//名称
							title := videoInfo.Text()

							ctx.AddQueue(&request.Request{
								Url:    "https://www.pornhub.com" + openUrl,
								Header: http.Header{"Content-Type": []string{"application/x-www-form-urlencoded; charset=UTF-8"}},
								Temp:   map[string]interface{}{"title": title, "valid": valid, "viewsNum": viewsNum},
								Rule:   "获取明细",
							})
						}

					})
				},
			},
			"下载视频": {
				ParseFunc: func(ctx *Context) {
					var title = ctx.GetUrl()
					var titleName = ctx.GetTemp("title", "").(string)
					ctx.Output(map[string]interface{}{
						"title": titleName,
						"url":   title,
					})
					ctx.FileOutput(titleName)
				},
			},
			"home页面": {
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()

					//其它页面
					//					profileMenuDropdown := query.Find("#profileMenuDropdown")
					//					profile, _ := profileMenuDropdown.Find(".profile").Attr("href")
					//					videos, _ := profileMenuDropdown.Find(".videos").Attr("href")
					//					favorites, _ := profileMenuDropdown.Find(".favorites").Attr("href")
					//					playlists, _ := profileMenuDropdown.Find(".playlists").Attr("href")
					//					photos, _ := profileMenuDropdown.Find(".photos").Attr("href")
					//					recommended  video/search?search=
					//					ht := "/video?o=ht" hot
					//					mv := "/video?o=mv" Most Viewed
					//					tr := "/video?o=tr" Top Rate

					query.Find(".thumbnail-info-wrapper").Each(func(i int, s *goquery.Selection) {
						//评分 73%!<(MISSING)
						valNum := s.Find(".value").Text()
						var reg = regexp.MustCompile("[0-9]*")

						var valid = reg.FindString(valNum)
						b, error := strconv.Atoi(valid)
						if error != nil {
							//							logs.Logs("字符串转换成整数失败")
						}
						//观看次数
						viewsNum := s.Find(".views").Find("var").Text()
						//						var viewsid = reg.FindString(viewsNum)

						vals := strings.Split(viewsNum, ",")

						viewsNum = strings.Join(vals, "")

						c, error := strconv.Atoi(viewsNum)
						if error != nil {
							//							logs.Logs("字符串转换成整数失败")
						}

						if b > 78 || c > 300000 {
							videoInfo := s.Find("a")
							//详细页面
							openUrl, _ := videoInfo.Attr("href")
							//名称
							title := videoInfo.Text()

							ctx.AddQueue(&request.Request{
								Url:    "https://www.pornhub.com" + openUrl,
								Header: http.Header{"Content-Type": []string{"application/x-www-form-urlencoded; charset=UTF-8"}},
								Temp:   map[string]interface{}{"title": title, "valid": valid, "viewsNum": viewsNum},
								Rule:   "获取明细",
							})
						}

					})
				},
			},
		},
	},
}
