package robot

import (
	"comics/tools"
	"comics/tools/config"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"github.com/tidwall/gjson"
	"net"
	"net/http"
	"time"
)

func GetColly() *colly.Collector {
	bot := colly.NewCollector(
		colly.AllowedDomains(config.Spe.SourceUrl),
	)
	extensions.RandomUserAgent(bot)
	extensions.Referer(bot)
	bot.WithTransport(&http.Transport{
		//Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          10,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	})
	content, code, _ := tools.HttpRequest("https://dvapi.doveproxy.net/cmapi.php?rq=distribute&user=yipinbao6688&token=eUkxbHhCSFZFcit1TS9XRWdxVy9mUT09&auth=0&geo=PH&city=208622&agreement=1&timeout=35&num=1&rtype=0",
		"GET", "", map[string]string{}, []*http.Cookie{})
	proxy := ""
	if code == 200 {
		res := gjson.Parse(content)
		proxy = "http://" + res.Get("data").Get("ip").String() + ":" + res.Get("data").Get("port").String()
	}
	if proxy != "" {
		//bot.SetProxy(proxy)
	}
	return bot
}
