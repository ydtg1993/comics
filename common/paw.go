package common

import (
	"bytes"
	"comics/tools"
	"comics/tools/config"
	"crypto/tls"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

func RequestApi(url, method, param string, timeout int) (gjson.Result, error) {
	header := map[string]string{
		"User-Agent": config.Spe.UserAgent,
		"Referer":    config.Spe.SourceUrl}
	var cookie []*http.Cookie
	content, code, _ := tools.HttpRequest(url, method, param, header, cookie)
	if code != 200 {
		logs.Error("无法抓取目标页 接口:" + url)
		return gjson.Parse(""), fmt.Errorf("无法抓取目标页 接口:" + url)
	}
	t := time.NewTicker(time.Second * time.Duration(timeout))
	<-t.C
	return gjson.Parse(content), nil
}

/**
* 下载远程http的文件到本地
* @param 	string		远程地址
* @param	string		本地路径
* @param	string		文件名称
 */
func DownFile(sUrl, filepath, fileName, proxy string, cookies map[string]string) string {
	//拼接完整地址
	allPathName := filepath + "/" + fileName
	//建立远程连接
	sParam := ""
	var client *http.Client
	if proxy != "" {
		proxy, _ := url.Parse(proxy)
		tr := &http.Transport{
			Proxy:           http.ProxyURL(proxy),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		client = &http.Client{
			Transport: tr,
			Timeout:   time.Second * 30,
		}
	} else {
		client = &http.Client{
			Timeout: time.Second * 30,
		}
	}
	req, er := http.NewRequest(http.MethodGet, sUrl, bytes.NewReader([]byte(sParam)))
	if er != nil {
		logs.Warning("连接请求失败 error->", sUrl, er.Error())
		return ""
	}
	//配置参数
	req.Header.Set("Host", config.Spe.SourceUrl)
	req.Header.Set("referer", "https://"+config.Spe.SourceUrl+"/")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("User-Agent", config.Spe.UserAgent)
	req.Header.Set("Connection", "Close")

	ck := new(bytes.Buffer)
	for key, value := range cookies {
		fmt.Fprintf(ck, "%s=\"%s\";", key, value)
	}
	req.Header.Set("Cookie", ck.String())

	resp, err := client.Do(req)
	if err != nil {
		logs.Warning("读取远程文件失败->", sUrl, err.Error())
		return ""
	}

	if resp != nil && resp.Body != nil {

		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			logs.Warning("网络状态异常", sUrl, resp.StatusCode)
			return ""
		}

		//创建多层目录,最后的是文件名不是目录
		_, cerr := tools.CreateFile(filepath)
		if cerr != nil {
			return allPathName
		}

		//创建文件
		out, er := os.Create(allPathName)
		defer out.Close()

		if er != nil {
			logs.Warning("创建文件失败", allPathName, er.Error())
			return ""
		}

		_, err = io.Copy(out, resp.Body)
		aw, aErr := tools.GetFileSize(allPathName)

		if err != nil {
			logs.Warning("写入文件失败", sUrl, err.Error())
			return ""
		}
		if aErr != nil {
			logs.Warning("保存本地文件失败", sUrl, err.Error())
			return ""
		}
		if resp.ContentLength > 0 && resp.ContentLength != aw {
			logs.Warning("文件下载不完整", sUrl, err.Error())
			//删除空文件
			os.RemoveAll(allPathName)
			return ""
		}

		return allPathName
	} else {
		return ""
	}
}
