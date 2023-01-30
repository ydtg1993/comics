package robot

import (
	"comics/tools/config"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func (Robot *Robot) TapIn(method, path string) bool {
	for i := 0; i < 3; i++ {
		dom, err := Robot.WebDriver.FindElement(method, path)
		if err != nil {
			continue
		}
		dom.Click()
		time.Sleep(1 * time.Second)
		//检查dom 判断跳转
		_, err = Robot.WebDriver.FindElement(method, path)
		if err == nil {
			return true
		}
	}
	return false
}

func (Robot *Robot) CatchImage(url string) []byte {
	payload := strings.NewReader("")
	req, _ := http.NewRequest("GET", url, payload)
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")
	current_url, _ := Robot.WebDriver.CurrentURL()
	req.Header.Add("Referer", current_url)
	req.Header.Add("Origin", current_url)
	req.Header.Add("Host", current_url)
	req.Header.Add("User-Agent", config.Spe.UserAgent)
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Cache-Control", "no-cache")

	res, err := http.DefaultClient.Do(req)
	defer res.Body.Close()

	var response []byte
	if err != nil {
		fmt.Println(err.Error())
		return response
	}

	reader, err := gzip.NewReader(res.Body)
	if err != nil {
		fmt.Println(err.Error())
		return response
	}
	response, _ = ioutil.ReadAll(reader)

	return response
}
