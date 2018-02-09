package pholcus_lib

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
	//	"fmt"
	// "math"
	// "time"
)

func init() {
	GirlHome.Register()
}

var GirlHome = &Spider{
	Name:        "pronlogin",
	Description: "pronpronlogin [http://www.pronhub.com/]",
	//	Pausetime:    2000,
	Keyin:        KEYIN,
	Limit:        LIMIT,
	EnableCookie: true,
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {
			ctx.AddQueue(&request.Request{
				Url:    "https://www.pronhub.com",
				Method: "GET",
				Rule:   "首页",
			})
		},
		Trunk: map[string]*Rule{
			"首页": {
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()
					logintoken := query.Find(".js-loginFormModal")
					redirect, _ := logintoken.Find(".js-redirect").Attr("value")
					token, _ := logintoken.Find("[name='token']").Attr("value")
					remember_me, _ := logintoken.Find("[name='remember_me']").Attr("value")
					from, _ := logintoken.Find("[name='from']").Attr("value")

					//					fmt.Println("首页：redirect=" + redirect + "&token=" + token + "&remember_me=" + remember_me + "&from=" + from)

					ctx.AddQueue(&request.Request{
						Url:      "https://www.pronhub.com/login",
						Method:   "POST",
						PostData: "redirect=" + redirect + "&token=" + token + "&remember_me=" + remember_me + "&from=" + from + "&username=kenzhao&password=a123456",
						Rule:     "登录",
						Header:   http.Header{"Content-Type": []string{"application/x-www-form-urlencoded; charset=UTF-8"}},
					})
				},
			},
			"登录": {
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()
					logintoken := query.Find(".js-loginFormModal")
					redirect, _ := logintoken.Find(".js-redirect").Attr("value")
					token, _ := logintoken.Find("[name='token']").Attr("value")
					remember_me, _ := logintoken.Find("[name='remember_me']").Attr("value")
					from, _ := logintoken.Find("[name='from']").Attr("value")

					//					fmt.Println("登陆login：redirect=" + redirect + "&token=" + token + "&remember_me=" + remember_me + "&from=" + from)

					ctx.AddQueue(&request.Request{
						Url:      "https://www.pronhub.com/front/authenticate",
						Header:   http.Header{"Content-Type": []string{"application/x-www-form-urlencoded; charset=UTF-8"}},
						Method:   "POST",
						PostData: "redirect=" + redirect + "&token=" + token + "&remember_me=" + remember_me + "&from=" + from + "&username=kenzhao&password=a123456",
						Rule:     "登录后",
					})
				},
			},
			"登录后": {
				ParseFunc: func(ctx *Context) {
					var paramsStr = ctx.GetKeyin()
					//					params := strings.Split(paramsStr, "@")

					//					startIndex, error := strconv.Atoi(params[0])
					//					endIndex, error := strconv.Atoi(params[1])
					//					if error != nil {
					//						//							logs.Logs("字符串转换成整数失败")
					//					}
					//					searchStr := params[2]

					//							Url:  "https://www.pronhub.com/" + searchStr + "&page=" + strconv.Itoa(i),
					//					for i := startIndex; i <= endIndex; i++ {
					//						ctx.AddQueue(&request.Request{
					//							Url:  "https://www.pronhub.com/" + searchStr + strconv.Itoa(i),
					//							Rule: "home页面",
					//						})
					//					}

					//						jav:  "https://www.pronhub.com/video?c=111",
					//						koearn:  "https://www.pronhub.com/video?c=103",https://www.pronhub.com/video/search?search=asian
					ctx.AddQueue(&request.Request{
						Url:  "https://www.pronhub.com/" + paramsStr,
						Rule: "home页面",
					})

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
								Url:    "https://www.pronhub.com" + openUrl,
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
								Url:    "https://www.pronhub.com" + openUrl,
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
