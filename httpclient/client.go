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
	URL         string
	Method      string
	Header      map[string]string
	Cookies     []http.Cookie
	Timeout     time.Duration
	Request     *http.Request
	QueryParams string
	Body        io.Reader
}

// New 新的client实例
func (c *Client) New(URL string, method string) *Client {
	client := &Client{}
	// 设置默认的超时时间
	client.Timeout = time.Second * 10
	return client
}

// WithHeaders 设置请求头信息
func (c *Client) WithHeaders(header map[string]string) *Client {
	for k, v := range header {
		c.Header[k] = v
	}
	return c
}

// WithCookies 设置cookie信息
func (c *Client) WithCookies(cookies []http.Cookie) *Client {
	for _, cookie := range cookies {
		c.Cookies = append(c.Cookies, cookie)
	}
	return c
}

// WithTimeout 设置超时时间
func (c *Client) WithTimeout(t time.Duration) *Client {
	c.Timeout = t
	return c
}

// WithQueryParams 携带查询参数
func (c *Client) WithQueryParams(queryParams map[string]string) *Client {
	if queryParams == nil {
		return c
	}
	if c.QueryParams != "" {
		c.QueryParams = c.QueryParams + "&"
	}
	for k, v := range queryParams {
		if c.QueryParams == "" {
			c.QueryParams = k + "=" + v
			continue
		}
		c.QueryParams = c.QueryParams + "&" + k + "=" + v
	}
	return c
}

// WithFormData 表单数据
func (c *Client) WithFormData(data url.Values) {
	c.Request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	c.Method = "POST"
	c.Body = bytes.NewBuffer([]byte(data.Encode()))
}

// WithFileData 文件上传数据
func (c *Client) WithFileData(srcFile string, fileName string) *Client {
	// 文件上传只允许POST方法
	c.Method = "POST"
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
	c.Body = io.MultiReader(body, fh, closeBuf)
	fi, err := fh.Stat()
	if err != nil {
		return c
	}
	// Set headers for multipart, and Content Length
	c.Request.Header.Add("Content-Type", "multipart/form-data; boundary="+boundary)
	c.Request.ContentLength = fi.Size() + int64(body.Len()) + int64(closeBuf.Len())
	return c
}

// WithRawData 发送的数据
func (c *Client) WithRawData(data []byte) *Client {
	c.Body = bytes.NewBuffer(data)
	return c
}

// Do 执行请求
func (c *Client) Do() ([]byte, error) {
	c.Request, _ = http.NewRequest(c.Method, c.URL, c.Body)
	if c.Request == nil {
		return nil, errors.New("request is nil,should be init first")
	}
	// 设置请求对象的header信息
	for k, v := range c.Header {
		c.Request.Header.Add(k, v)
	}

	// 设置cookie
	for _, cookie := range c.Cookies {
		c.Request.AddCookie(&cookie)
	}

	// 设置超时时间
	client := &http.Client{}
	client.Timeout = c.Timeout

	// 执行请求
	resp, err := client.Do(c.Request)
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
