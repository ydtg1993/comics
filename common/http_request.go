package common

import (
	"bytes"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/logs"
)

/**
* 将复杂请求，因为网络的不稳定，将常规请求变成5次连续请求
 */
func HttpRequestByHeaderFor5(sUrl, method, sParame string, mHeader map[string]string) (string, int64, int) {
	res, l, n := HttpContentByHeader(sUrl, sParame, mHeader, method)
	if n != 200 {
		for i := 0; i < 5; i++ {
			res, l, n = HttpContentByHeader(sUrl, sParame, mHeader, method)
			if n == 200 {
				break
			}
		}
	}
	return res, l, n
}

/**
* 将复杂请求带cookie，因为网络的不稳定，将常规请求变成5次连续请求
 */
func HttpRequestByCookieFor5(sUrl, method, sParame string, mHeader map[string]string, cook []*http.Cookie) (string, int, []*http.Cookie) {
	res, n, cookie := httpRequestByCookie(sUrl, method, sParame, mHeader, cook)
	if n != 200 {
		for i := 0; i < 5; i++ {
			res, n, cookie = httpRequestByCookie(sUrl, method, sParame, mHeader, cook)
			if n == 200 {
				break
			}
		}
	}
	return res, n, cookie
}

/**
* 将复杂请求，因为网络的不稳定，将常规请求变成5次连续请求
 */
func DownFileFor5(sUrl, sPath, sName string) (string, bool) {
	res, b := DownFile(sUrl, sPath, sName)
	if b == false {
		for i := 0; i < 5; i++ {
			res, b = DownFile(sUrl, sPath, sName)
			if b == true {
				break
			}
		}
	}
	return res, b
}

/**
*  简单版http请求，适用于没有特别要求的
* httpUrl	请求的网址
* method	网络请求方式，一般为POST或者GET
* sParam	需要传递的参数
* mHeader	http的头部
 */
func httpRequest(httpUrl, method, sParam string, mHeader map[string]string) (string, error) {
	client := &http.Client{}

	req, er := http.NewRequest(method, httpUrl, bytes.NewReader([]byte(sParam)))
	if er != nil {
		req, er = http.NewRequest(method, httpUrl, strings.NewReader(sParam))
		if er != nil {
			//两次连接都失败了，需要返回一个空
			return "", er
		}
	}
	req.Close = true

	defer req.Body.Close()

	for key, val := range mHeader {
		req.Header.Add(key, val)
	}

	var body []byte
	resp, err := client.Do(req)
	if err == nil {
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			httpRequest(req.RequestURI, req.Method, "", mHeader)
		}

		body, _ = ioutil.ReadAll(resp.Body)
	}

	return string(body), err
}

/**
* 复杂版http，请求携带cookie
* httpUrl	请求的网址
* method	网络请求方式，一般为POST或者GET
* sParam	需要传递的参数
* mHeader	http的头部
* setCookie	传递的cookie
 */
func httpRequestByCookie(httpUrl, method, sParam string, mHeader map[string]string, setCookie []*http.Cookie) (string, int, []*http.Cookie) {
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

/**
* 复杂版http请求
* @param 	string 	sUrl	请求的地址
* @param	string	params	带入的参数
* @param	map	mHeader		head头
* @param	string	method		http方法
* 返回 字符串，长度，状态码
 */
func HttpContentByHeader(sUrl, params string, mHeader map[string]string, method string) (string, int64, int) {
	var strLen int64 = 0
	src := ""
	httpStart := true
	statusCode := 101
	req, err := http.NewRequest(method, sUrl, strings.NewReader(params))
	if err != nil {

		logs.Warning(sUrl, err.Error())

		req, err = http.NewRequest(method, sUrl, strings.NewReader(params))
		if err != nil {
			httpStart = false
			//两次连接都失败了，需要返回一个空
			return "", strLen, statusCode
		}
	}

	if req != nil && req.Body != nil {
		defer req.Body.Close()
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
		resp, err := client.Do(req)

		if err == nil {
			statusCode = resp.StatusCode
			contents, err := ioutil.ReadAll(resp.Body)
			if err == nil {
				src = string(contents)
			}

		}
		if resp != nil && resp.Body != nil {
			strLen = resp.ContentLength
			defer resp.Body.Close()
		}
	}

	return src, strLen, statusCode
}

/**
* 下载远程http的文件到本地
* @param 	string		远程地址
* @param	string		本地路径
* @param	string		文件名称
 */
func DownFile(sUrl, filepath, fileName string) (string, bool) {
	res := false
	//拼接完整地址
	allPathName := filepath + "/" + fileName

	//建立远程连接
	sParam := ""
	client := &http.Client{}
	client.Timeout = time.Minute * 1
	req, er := http.NewRequest(http.MethodGet, sUrl, bytes.NewReader([]byte(sParam)))
	if er != nil {
		logs.Warning("连接请求失败 error->", sUrl, er.Error())
		return "", false
	}

	if req != nil && req.Body != nil {
		defer req.Body.Close()
	} else {
		return "", false
	}

	//解析url
	u, uErr := url.Parse(sUrl)
	if uErr != nil {
		logs.Warning("url地址无法解析", sUrl, uErr.Error())
		return "", false
	}

	//配置参数
	req.Header.Set("Host", u.Host)
	req.Header.Set("Accept-Encoding", "identity") //强制服务器不走压缩，不然会得不到contentLength
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.54 Safari/537.36")
	req.Header.Set("Connection", "Close")

	resp, err := client.Do(req)
	//针对请求进行多次尝试

	if err != nil {
		logs.Warning("读取远程文件失败->", sUrl, err.Error())
		return allPathName, false
	}

	if resp != nil && resp.Body != nil {

		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			logs.Warning("网络状态异常", sUrl, resp.StatusCode)
			return allPathName, false
		}

		//创建多层目录,最后的是文件名不是目录
		_, cerr := CreateFile(filepath)
		if cerr != nil {
			return allPathName, res
		}

		//创建文件
		out, er := os.Create(allPathName)
		defer out.Close()

		if er != nil {
			logs.Warning("创建文件失败", allPathName, er.Error())
			return allPathName, res
		}

		//分块读取字节流
		/*var body []byte
		aw := 0

		bs := make([]byte, 0, 1024) //建立缓冲区块
		for {
			if len(bs) == cap(bs) {
				//申请内存空间
				bs = append(bs, 0)[:len(bs)]
			}
			//分块读取wenjin
			n, err := resp.Body.Read(bs[len(bs):cap(bs)])
			bs = bs[:len(bs)+n]
			//计算读取完成的字节流数量
			aw += n

			if err != nil {
				if err == io.EOF {
					//读完数据，跳出去
					break
				}
			}
		}
		body = bs

		logs.Debug("response size->", "contentLength=", resp.ContentLength, "body=", len(bs), "write=", aw)

		//如果读取不到文件，这里直接返回失败
		if len(body) < 1 {
			logs.Warning("获取文件大小失败", sUrl, resp.ContentLength, len(body))
			//删除空文件
			os.RemoveAll(allPathName)
			return allPathName, false
		}

		//检验文件大小
		rSize := resp.ContentLength
		if rSize < 1 {
			rSize = int64(len(bs))
		}
		if rSize != int64(aw) {
			logs.Warning("文件下载不完整", sUrl, rSize, aw)
			//删除空文件
			os.RemoveAll(allPathName)
			return allPathName, false
		}

		_, err = io.Copy(out, bytes.NewReader(body))*/

		_, err = io.Copy(out, resp.Body)

		aw, aErr := GetFileSize(allPathName)

		logs.Debug("response size->", sUrl, allPathName, "contentLength=", resp.ContentLength, "write=", aw)

		if err != nil {
			logs.Warning("写入文件失败", sUrl, err.Error())
			return allPathName, false
		}
		if aErr != nil {
			logs.Warning("保存本地文件失败", sUrl, err.Error())
			return allPathName, false
		}
		if resp.ContentLength > 0 && resp.ContentLength != aw {
			logs.Warning("文件下载不完整", sUrl, err.Error())
			//删除空文件
			os.RemoveAll(allPathName)
			return allPathName, false
		}

		return allPathName, true
	} else {
		return allPathName, false
	}

}
