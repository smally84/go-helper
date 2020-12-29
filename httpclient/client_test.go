package httpclient

import (
	"fmt"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	res, err := New("https://www.baidu.com", "GET").
		WithHeaders(map[string]string{}).
		WithTimeout(time.Second * 10).
		Request()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(res))
	}

}
