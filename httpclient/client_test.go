package httpclient

import (
	"fmt"
	"testing"
)

func TestClient(t *testing.T) {
	client := New("https://www.piaoziyou.com", "GET")
	res, err := client.Do()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(res))
	}

}
