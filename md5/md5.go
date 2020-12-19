package md5

import (
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"os"
)

// StringMd5 字符串md5
func StringMd5(str string) string {
	plain := md5.New()
	plain.Write([]byte(str))
	return hex.EncodeToString(plain.Sum(nil))
}

// FileMd5 获取文件的md5值
func FileMd5(url string) string {
	finfo, err := os.Stat(url)
	if err == nil && !finfo.IsDir() {
		binData, err := ioutil.ReadFile(url)
		if err != nil {
			return ""
		}
		plain := md5.New()
		plain.Write(binData)
		md5 := hex.EncodeToString(plain.Sum(nil))
		return md5
	}
	return ""
}
