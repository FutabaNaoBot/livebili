package request

import (
	"io"
	"net/http"
)

func setHeader(req *http.Request, cookies string) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")
	req.Header.Set("Referer", "https://www.bilibili.com")
	req.Header.Set("Cookie", cookies)
}

func DoGet(url string, cookies string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	setHeader(req, cookies)

	// 发送请求
	client := &http.Client{}
	return client.Do(req)
}

func DoPost(url, contentType string, body io.Reader, cookies string) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	setHeader(req, cookies)
	client := &http.Client{}
	return client.Do(req)
}
