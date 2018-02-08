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


## 爬虫规则列表

|子目录|规则描述|
|---|---|