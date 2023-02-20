package robot

import (
	"comics/tools/config"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
)

func GetColly() *colly.Collector {
	bot := colly.NewCollector(
		colly.AllowedDomains(config.Spe.SourceUrl),
	)
	extensions.RandomUserAgent(bot)
	extensions.Referer(bot)
	proxy := GetProxy()
	if proxy != "" && config.Spe.AppDebug == false {
		bot.SetProxy(proxy)
	}
	return bot
}
