package robot

import (
	"comics/tools"
	"github.com/tidwall/gjson"
	"net/http"
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

func GetProxy() string {
	content, code, _ := tools.HttpRequest("https://dvapi.doveproxy.net/cmapi.php?rq=distribute&user=yipinbao6688&token=eUkxbHhCSFZFcit1TS9XRWdxVy9mUT09&auth=0&geo=PH&city=208622&agreement=1&timeout=35&num=1&rtype=0",
		"GET", "", map[string]string{}, []*http.Cookie{})
	proxy := ""
	if code == 200 {
		res := gjson.Parse(content)
		proxy = "http://" + res.Get("data").Get("ip").String() + ":" + res.Get("data").Get("port").String()
	}
	return proxy
}
