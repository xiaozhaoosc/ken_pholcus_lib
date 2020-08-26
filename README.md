# pholcus_lib

[Pholcus](https://github.com/henrylee2cn/pholcus) 用户公共维护的spider爬虫规则库。

## 维护规范

- 欢迎每位用户都来分享自己的爬虫规则
- 每个规则放在单一个独的子目录
- 新增规则最好提供README.md
- 新增规则时，须在根目录 `pholcus_lib.go` 文件的import组中添加类似`_ "github.com/henrylee2cn/pholcus_lib/jingdong"`的包引用声明
- 新增规则时，须在根目录README.md（本文档）的 `爬虫规则列表` 中按子目录名`a-z`的顺序插入一条相应的规则记录
- 维护旧规则时，应在规则文件或相应README.md中增加修改说明：如修改原因、修改时间、签名、联系方式等
- 凡爬虫规则的贡献者均可在其源码文件或相应README.md中留下在的签名、联系方式

## 使用方法
1. 下载并安装项目
```
go get -u -v github.com/henrylee2cn/pholcus
go get -u -v github.com/xiaozhaoosc/ken_pholcus_lib
```
2. 更改规则库
    1. example_main.go里面修改代码，如下 ：
```
package main

import (
	"github.com/henrylee2cn/pholcus/exec"
	//_ "github.com/henrylee2cn/pholcus_lib" // 此为公开维护的spider规则库
	_ "github.com/xiaozhaoosc/ken_pholcus_lib" // ken_pholcus_lib
)

func main() {
	// 设置运行时默认操作界面，并开始运行
	// 运行软件前，可设置 -a_ui 参数为"web"、"gui"或"cmd"，指定本次运行的操作界面
	// 其中"gui"仅支持Windows系统
	exec.DefaultRun("web")
}

```

3. 修改mysql
```
设置表名

	names := strings.Split(name, "_")
	self.tableName = wrapSqlKey(time.Now().Format("2006010215") + names[0] + names[4])
```
## 爬虫规则列表

|子目录|规则描述|
|---|---|