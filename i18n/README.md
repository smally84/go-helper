# go-i18n
Golang国际化实现

# 特色
- 简单易用
- 支持嵌套翻译
- 支持模板翻译

# 使用步骤
go get -u https://github.com/smally84/go-i18n

```go
import (
	i18n "https://github.com/smally84/go-i18n"
)
var langPack map[string]map[interface{}]interface{}
func main(){
	i18n.LoadLangPack(langPack)
	fmt.Println(i18n.T("key"))
}
```

# 使用说明

```
// 加载语言包定义
LoadLangPack(langPack1, langPack2, langPack3)
// 简单字段翻译
T("zh-cn", "author")
// 嵌套key翻译，多级用"."
T("zh-cn", "author.name")
// 带有模板的翻译
T("zh-cn", "hello{name}", map[string]string{"name": "world"})
```

# 语言包定义示例

```
var langPack1 = map[string]map[interface{}]interface{}{
	"zh-cn": {
		1000: "输入有误",
	},
	"en-us": {
		1000: "Input error",
	},
}
var langPack2 = map[string]map[interface{}]interface{}{
	"zh-cn": {
		"author": "smally84",
		"user": map[string]string{
			"name": "姓名",
			"sex":  "性别",
		},
	},
	"en-us": {
		"user": map[string]string{
			"name": "name",
			"sex":  "sex",
		},
	},
}
var langPack3 = map[string]map[interface{}]interface{}{
	"zh-cn": {
		"hello{name}": "你好{name}",
	},
	"en-us": {
		"hello{name}": "hello{name}",
	},
```

# 注意事项
当键值为整数时，务必将其值转为int，因为接口类型的整数默认为int类型

# License

MIT License