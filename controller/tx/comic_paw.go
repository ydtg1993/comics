package tx

import (
	"comics/common"
	"comics/global/orm"
	"comics/model"
	"comics/tools"
	"comics/tools/config"
	"comics/tools/rd"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"math"
	"regexp"
	"strconv"
)

func ComicPaw() {
	tags := map[string]int{
		/*"恋爱": 105,
		"玄幻": 101,
		"异能": 103,
		"恐怖": 110,
		"剧情": 106,
		"科幻": 108,*/
		"悬疑": 112,
	}
	pays := map[string]int{
		"免费": 1,
		//"付费": 2,
	}
	states := map[string]int{
		"连载中": 1,
		"已完结": 2,
	}
	for tag, tagId := range tags {
		for pay, payId := range pays {
			for state, stateId := range states {
				fmt.Printf("%s %s %s \n", tag, pay, state)
				category(tagId, payId, stateId)
			}
		}
	}
}

func ComicUpdate() {

}

func category(tagId, payId, stateId int) {
	bot := colly.NewCollector(
		colly.AllowedDomains(config.Spe.SourceUrl),
	)
	extensions.RandomUserAgent(bot)
	extensions.Referer(bot)

	url := fmt.Sprintf("https://"+config.Spe.SourceUrl+"/Comic/all/theme/%d/finish/%d/search/time/vip/%d/page/1",
		tagId, stateId, payId)

	page := 1
	bot.OnResponse(func(r *colly.Response) {
		regexp := regexp.MustCompile(`var totalNum = "(\d+)";`)
		params := regexp.FindStringSubmatch(string(r.Body))
		if len(params) >= 2 {
			total, _ := strconv.Atoi(params[1])
			page = int(math.Ceil(float64(total) / float64(12)))
			for {
				if page <= 1 {
					break
				}
				paw(bot, tagId, payId, stateId, page)
				page--
			}
		}
	})
	err := bot.Visit(url)
	if err != nil {
		logs.Error("无法抓取分类列表页Dom:" + url)
	}
}

func paw(bot *colly.Collector, tagId, payId, stateId, page int) {
	url := fmt.Sprintf("https://"+config.Spe.SourceUrl+"/Comic/all/theme/%d/finish/%d/search/time/vip/%d/page/%d",
		tagId, stateId, payId, page)

	bot.OnHTML("li.ret-search-item", func(e *colly.HTMLElement) {
		info := e.DOM.Find(".ret-works-info")
		title := info.Find(".ret-works-title>a").Text()
		url, _ := info.Find(".ret-works-title>a").Attr("href")
		id := tools.FindStringNumber(url)
		author := info.Find(".ret-works-author").Text()
		description := info.Find(".ret-works-decs").Text()
		cover, _ := e.DOM.Find(".ret-works-cover img.lazy").Attr("data-original")

		var exists bool
		orm.Eloquent.Model(model.SourceComic{}).Select("count(*) > 0").
			Where("source = ? and source_id = ?", config.Spe.SourceId, id).Find(&exists)
		if exists == true {
			return
		}
		sourceComic := new(model.SourceComic)
		sourceComic.Source = config.Spe.SourceId
		sourceComic.SourceId = id
		sourceComic.SourceUrl = "https://" + config.Spe.SourceUrl + url
		sourceComic.Title = title
		sourceComic.Cover = cover
		sourceComic.Author = author
		sourceComic.Description = description
		sourceComic.LikeCount = ""
		sourceComic.Popularity = ""
		var cookies map[string]string
		dir := fmt.Sprintf(config.Spe.DownloadPath+"comic/%d", id%10)
		downCover := common.DownFile(cover, dir, tools.RandStr(9)+".jpg", cookies)
		if downCover != "" {
			sourceComic.Cover = downCover
		}
		err := orm.Eloquent.Create(&sourceComic).Error
		if err != nil {
			logs.Error(fmt.Sprintf("comic数据导入失败 source = %d source_id = %d", config.Spe.SourceId, id))
		} else {
			rd.RPush(model.SourceComicTASK, sourceComic.Id)
		}
	})

	err := bot.Visit(url)
	if err != nil {
		logs.Error("无法抓取分类列表页Dom:" + url)
	}
}
