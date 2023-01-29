package tools

import (
	"bytes"
	"github.com/beego/beego/v2/core/logs"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

/**
* 复杂版http，请求携带cookie
* httpUrl	请求的网址
* method	网络请求方式，一般为POST或者GET
* sParam	需要传递的参数
* mHeader	http的头部
* setCookie	传递的cookie
 */
func HttpRequest(httpUrl, method, sParam string, mHeader map[string]string, setCookie []*http.Cookie) (string, int, []*http.Cookie) {
	src := ""
	httpStart := true
	statusCode := 101

	cook := []*http.Cookie{}

	req, er := http.NewRequest(method, httpUrl, bytes.NewReader([]byte(sParam)))
	if er != nil {

		logs.Warning("http request error->", httpUrl, er.Error())

		req, er = http.NewRequest(method, httpUrl, strings.NewReader(sParam))
		if er != nil {
			httpStart = false
			//两次连接都失败了，需要返回一个空
			return "", statusCode, setCookie
		}
	}

	for key, val := range mHeader {
		req.Header.Add(key, val)
	}

	if setCookie != nil && len(setCookie) > 0 {
		for _, v := range setCookie {
			req.AddCookie(v)
		}
	}

	//只有连接成功后，才会写入头的读取字节流
	if httpStart == true {
		if len(mHeader) > 0 {
			//mid := ""
			for key, val := range mHeader {
				//mid = key + ":" + val
				req.Header.Set(key, val)
			}
		}

		req.Header.Set("Accept-Charset", "utf-8")
		req.Header.Set("Connection", "Close")

		tr := &http.Transport{DisableKeepAlives: true,
			Dial: func(netw, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(netw, addr, time.Second*30) //设置建立连接超时
				if err != nil {
					return nil, err
				}
				c.SetDeadline(time.Now().Add(30 * time.Second)) //设置发送接收数据超时
				return c, nil
			}}
		client := &http.Client{Transport: tr}
		req2 := req
		resp, err := client.Do(req2)

		if err != nil {
			logs.Warning("连接失败", httpUrl, err.Error())
			//两次连接都失败了，需要返回一个空
			return "", statusCode, setCookie
		} else {
			defer resp.Body.Close()

			statusCode = resp.StatusCode
			cook = resp.Cookies()
			contents, _ := ioutil.ReadAll(resp.Body)
			src = string(contents)
		}
	}

	defer req.Body.Close()

	return src, statusCode, cook
}
