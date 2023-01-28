package controller

import (
	"comics/tools"
	"comics/tools/config"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"github.com/tidwall/gjson"
	"time"
)

func requestHtml(url string, timeout int) (string, error) {
	header := map[string]string{
		"User-Agent": config.Spe.UserAgent,
		"Referer":    config.Spe.SourceUrl}
	content, _, code := tools.HttpRequestByHeaderFor5(url, "GET", "", header)
	if code != 200 {
		logs.Error("无法抓取目标页 分页:" + url)
		return "", fmt.Errorf("无法抓取目标页 分页:" + url)
	}
	t := time.NewTicker(time.Second * time.Duration(timeout))
	<-t.C
	return content, nil
}

func requestApi(url, method, param string, timeout int) (gjson.Result, error) {
	header := map[string]string{
		"User-Agent": config.Spe.UserAgent,
		"Referer":    config.Spe.SourceUrl}
	content, _, code := tools.HttpRequestByHeaderFor5(url, method, param, header)
	if code != 200 {
		logs.Error("无法抓取目标页 接口:" + url)
		return gjson.Parse(""), fmt.Errorf("无法抓取目标页 接口:" + url)
	}
	t := time.NewTicker(time.Second * time.Duration(timeout))
	<-t.C
	return gjson.Parse(content), nil
}
