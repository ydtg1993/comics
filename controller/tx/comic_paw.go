package tx

import (
	"comics/common"
	"comics/global/orm"
	"comics/model"
	"comics/tools"
	"comics/tools/config"
	"comics/tools/rd"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"math"
	"regexp"
	"strconv"
	"time"
)

func ComicPaw() {
	tags := map[string]int{
		"恋爱": 105,
		"玄幻": 101,
		"异能": 103,
		"恐怖": 110,
		"剧情": 106,
		"科幻": 108,
		"悬疑": 112,
		"奇幻": 102,
		"冒险": 104,
		"犯罪": 111,
		"动作": 109,
		"日常": 113,
		"竞技": 114,
		"武侠": 115,
		"历史": 116,
		"战争": 117,
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
				tx := common.Kind{
					Tag:   common.Kv{Name: tag, Val: tagId},
					Pay:   common.Kv{Name: pay, Val: payId},
					State: common.Kv{Name: state, Val: stateId},
				}
				category(tx)
			}
		}
	}
}

func ComicUpdate() {
	bot := colly.NewCollector(
		colly.AllowedDomains(config.Spe.SourceUrl),
	)
	extensions.RandomUserAgent(bot)
	extensions.Referer(bot)

	for page := 1; page < 13; page++ {
		url := fmt.Sprintf("https://"+config.Spe.SourceUrl+"/Comic/all/search/time/page/%d",
			page)

		bot.OnHTML("li.ret-search-item", func(e *colly.HTMLElement) {
			info := e.DOM.Find(".ret-works-info")
			title := info.Find(".ret-works-title>a").Text()
			url, _ := info.Find(".ret-works-title>a").Attr("href")
			id := tools.FindStringNumber(url)
			author := info.Find(".ret-works-author").Text()
			cover, _ := e.DOM.Find(".ret-works-cover img.lazy").Attr("data-original")
			popularity := e.DOM.Find(".ret-works-tags span").Last().Find("em").Text()

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
			sourceComic.Label = model.Label{}
			sourceComic.LikeCount = ""
			sourceComic.Popularity = popularity

			var cookies map[string]string
			dir := fmt.Sprintf(config.Spe.DownloadPath+"comic/%d", id%10)
			downCover := common.DownFile(cover, dir, tools.RandStr(9)+".jpg", cookies)
			if downCover != "" {
				sourceComic.Cover = downCover
			}
			err := orm.Eloquent.Create(&sourceComic).Error
			if err != nil {
				msg := fmt.Sprintf("漫画入库失败 source = %d source_id = %d", config.Spe.SourceId, id)
				model.RecordFail(url, msg, "漫画入库", 1)
			} else {
				rd.RPush(common.SourceComicTASK, sourceComic.Id)
			}
		})

		for i := 0; i < 3; i++ {
			err := bot.Visit(url)
			if err != nil && i == 3 {
				model.RecordFail(url, "无法抓取分类列表页信息 :"+url, "列表错误", 0)
			} else {
				break
			}
		}
	}
}

func category(tx common.Kind) {
	bot := colly.NewCollector(
		colly.AllowedDomains(config.Spe.SourceUrl),
	)
	extensions.RandomUserAgent(bot)
	extensions.Referer(bot)

	url := fmt.Sprintf("https://"+config.Spe.SourceUrl+"/Comic/all/theme/%d/finish/%d/search/time/vip/%d/page/1",
		tx.Tag.Val, tx.State.Val, tx.Pay.Val)

	page := 1
	bot.OnResponse(func(r *colly.Response) {
		regexp := regexp.MustCompile(`var totalNum = "(\d+)";`)
		params := regexp.FindStringSubmatch(string(r.Body))
		if len(params) >= 2 {
			total, _ := strconv.Atoi(params[1])
			page = int(math.Ceil(float64(total) / float64(12)))
			for {
				if page < 1 || page > 5 {
					break
				}
				paw(bot, tx, page)
				t := time.NewTicker(time.Second * 2)
				<-t.C
				page--
			}
		}
	})
	err := bot.Visit(url)
	t := time.NewTicker(time.Second * 2)
	<-t.C
	if err != nil {
		model.RecordFail(url, "无法抓取分类列表页信息 :"+url, "列表错误", 0)
	}
}

func paw(bot *colly.Collector, tx common.Kind, page int) {
	url := fmt.Sprintf("https://"+config.Spe.SourceUrl+"/Comic/all/theme/%d/finish/%d/search/time/vip/%d/page/%d",
		tx.Tag.Val, tx.State.Val, tx.Pay.Val, page)

	bot.OnHTML("li.ret-search-item", func(e *colly.HTMLElement) {
		info := e.DOM.Find(".ret-works-info")
		title := info.Find(".ret-works-title>a").Text()
		url, _ := info.Find(".ret-works-title>a").Attr("href")
		id := tools.FindStringNumber(url)
		author := info.Find(".ret-works-author").Text()
		cover, _ := e.DOM.Find(".ret-works-cover img.lazy").Attr("data-original")
		popularity := e.DOM.Find(".ret-works-tags span").Last().Find("em").Text()

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
		sourceComic.Label = model.Label{tx.Tag.Name}
		sourceComic.Category = tx.Tag.Name
		sourceComic.LikeCount = ""
		sourceComic.Popularity = popularity
		if tx.State.Val == 2 {
			sourceComic.IsFinish = 1
		}
		var cookies map[string]string
		dir := fmt.Sprintf(config.Spe.DownloadPath+"comic/%d/%d", config.Spe.SourceId, id%128)
		downCover := common.DownFile(cover, dir, tools.RandStr(9)+".jpg", cookies)
		if downCover != "" {
			sourceComic.Cover = downCover
		}
		err := orm.Eloquent.Create(&sourceComic).Error
		if err != nil {
			msg := fmt.Sprintf("漫画入库失败 source = %d source_id = %d err = %s", config.Spe.SourceId, id, err.Error())
			model.RecordFail(url, msg, "漫画入库", 1)
		} else {
			rd.RPush(common.SourceComicTASK, sourceComic.Id)
		}
	})

	for i := 0; i < 3; i++ {
		err := bot.Visit(url)
		if err != nil && i == 3 {
			model.RecordFail(url, "无法抓取分类列表页信息 :"+url, "列表错误", 0)
		} else {
			break
		}
	}
}
