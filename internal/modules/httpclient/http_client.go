package httpclient

// http-client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type ResponseWrapper struct {
	StatusCode int
	Body       string
	Header     http.Header
}

//多了个传参headers []string
func Get(url string, timeout int, headers []string) ResponseWrapper {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return createRequestError(err)
	}
	return request(req, timeout, headers)
}

//多了个传参headers []string
func PostParams(url string, params string, timeout int, headers []string) ResponseWrapper {
	buf := bytes.NewBufferString(params)
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return createRequestError(err)
	}
	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	return request(req, timeout, headers)
}

//多了个传参headers []string
func PostJson(url string, body string, timeout int, headers []string) ResponseWrapper {
	buf := bytes.NewBufferString(body)
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return createRequestError(err)
	}
	req.Header.Set("Content-type", "application/json")
	return request(req, timeout, headers)
}

//多了个传参headers []string
func request(req *http.Request, timeout int, headers []string) ResponseWrapper {
	wrapper := ResponseWrapper{StatusCode: 0, Body: "", Header: make(http.Header)}
	client := &http.Client{}
	if timeout > 0 {
		client.Timeout = time.Duration(timeout) * time.Second
	}
	setRequestHeader(req, headers)
	resp, err := client.Do(req)
	if err != nil {
		wrapper.Body = fmt.Sprintf("执行HTTP请求错误-%s", err.Error())
		return wrapper
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		wrapper.Body = fmt.Sprintf("读取HTTP请求返回值失败-%s", err.Error())
		return wrapper
	}
	wrapper.StatusCode = resp.StatusCode
	wrapper.Body = string(body)
	wrapper.Header = resp.Header
	return wrapper
}

func setRequestHeader(req *http.Request, headers []string) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")
	req.Header.Set("Cookie", "test")
	//增加了
	for _, value := range headers {
		value = strings.Replace(value, "\"", "", -1)
		value = strings.Replace(value, ",", "", -1)
		headers_value := strings.Split(value, ":")
		if len(headers_value) > 1 {
			req.Header.Set(headers_value[0], headers_value[1])
		}
	}
	//增加结束
}

func createRequestError(err error) ResponseWrapper {
	errorMessage := fmt.Sprintf("创建HTTP请求错误-%s", err.Error())
	return ResponseWrapper{0, errorMessage, make(http.Header)}
}
