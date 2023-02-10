package robot

import (
	"comics/tools/config"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
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
	proxy := GetProxy()
	if proxy != "" && config.Spe.AppDebug == false {
		bot.SetProxy(proxy)
	}
	return bot
}
