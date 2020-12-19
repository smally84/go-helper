package httpclient

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

// Client httpclient
type Client struct {
	url         string
	method      string
	header      map[string]string
	cookies     []http.Cookie
	timeout     time.Duration
	request     *http.Request
	queryParams string
	body        io.Reader
}

// New 新的client实例
func New(URL string, method string) *Client {
	client := &Client{}
	// 设置默认的超时时间
	client.timeout = time.Second * 10
	return client
}

// WithHeaders 设置请求头信息
func (c *Client) WithHeaders(header map[string]string) *Client {
	for k, v := range header {
		c.header[k] = v
	}
	return c
}

// WithCookies 设置cookie信息
func (c *Client) WithCookies(cookies []http.Cookie) *Client {
	for _, cookie := range cookies {
		c.cookies = append(c.cookies, cookie)
	}
	return c
}

// WithTimeout 设置超时时间
func (c *Client) WithTimeout(t time.Duration) *Client {
	c.timeout = t
	return c
}

// WithQueryParams 携带查询参数
func (c *Client) WithQueryParams(queryParams map[string]string) *Client {
	if queryParams == nil {
		return c
	}
	if c.queryParams != "" {
		c.queryParams = c.queryParams + "&"
	}
	for k, v := range queryParams {
		if c.queryParams == "" {
			c.queryParams = k + "=" + v
			continue
		}
		c.queryParams = c.queryParams + "&" + k + "=" + v
	}
	return c
}

// WithFormData 表单数据
func (c *Client) WithFormData(data url.Values) {
	c.request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	c.method = "POST"
	c.body = bytes.NewBuffer([]byte(data.Encode()))
}

// WithFileData 文件上传数据
func (c *Client) WithFileData(srcFile string, fileName string) *Client {
	// 文件上传只允许POST方法
	c.method = "POST"
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	// the Part headers to the buffer
	_, err := writer.CreateFormFile(fileName, srcFile)
	if err != nil {
		panic(err)
	}
	// the file data will be the second part of the body
	fh, err := os.Open(srcFile)
	if err != nil {
		panic(err)
	}
	boundary := writer.Boundary()
	closeBuf := bytes.NewBufferString(fmt.Sprintf("\r\n--%s--\r\n", boundary))
	c.body = io.MultiReader(body, fh, closeBuf)
	fi, err := fh.Stat()
	if err != nil {
		return c
	}
	// Set headers for multipart, and Content Length
	c.request.Header.Add("Content-Type", "multipart/form-data; boundary="+boundary)
	c.request.ContentLength = fi.Size() + int64(body.Len()) + int64(closeBuf.Len())
	return c
}

// WithRawData 发送的数据
func (c *Client) WithRawData(data []byte) *Client {
	c.body = bytes.NewBuffer(data)
	return c
}

// Do 执行请求
func (c *Client) Do() ([]byte, error) {
	c.request, _ = http.NewRequest(c.method, c.url, c.body)
	if c.request == nil {
		return nil, errors.New("request is nil,should be init first")
	}
	// 设置请求对象的header信息
	for k, v := range c.header {
		c.request.Header.Add(k, v)
	}

	// 设置cookie
	for _, cookie := range c.cookies {
		c.request.AddCookie(&cookie)
	}

	// 设置超时时间
	client := &http.Client{}
	client.Timeout = c.timeout

	// 执行请求
	resp, err := client.Do(c.request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 获取请求结果
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respData, nil
}

// DownloadFile 下载文件
// dstURL 文件的保存地址
func (c *Client) DownloadFile(dstURL string) error {
	resData, _ := c.Do()
	// 创建目标文件的文件夹
	dstDir := filepath.Dir(dstURL)
	_, err := os.Stat(dstDir)
	if err != nil {
		err := os.MkdirAll(dstDir, os.ModePerm)
		if err != nil {
			return err
		}
	}
	// 创建一个文件用于保存
	out, err := os.Create(dstURL)
	if err != nil {
		return err
	}
	defer out.Close()

	// 然后将响应流和文件流对接起来
	_, err = io.Copy(out, bytes.NewBuffer(resData))
	if err != nil {
		return err
	}
	return nil
}
