package i18n

import (
	"fmt"
	"testing"
)

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
}

func TestInit(t *testing.T) {
	LoadLangPack(langPack1, langPack2, langPack3)
	fmt.Println(T("zh-cn", "author"))
	fmt.Println(T("zh-cn", "user.name"))
	fmt.Println(T("en-us", "user.name"))
	fmt.Println(T("zh-cn", "hello{name}", map[string]string{"name": "world"}))
}
