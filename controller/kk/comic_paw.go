package kk

import (
	"comics/common"
	"comics/global/orm"
	"comics/model"
	"comics/tools"
	"comics/tools/config"
	"comics/tools/rd"
	"fmt"
	"github.com/tidwall/gjson"
	"math"
	"path/filepath"
	"strconv"
	"strings"
)

func ComicPaw() {
	tags := map[string]int{
		"恋爱":  20,
		"古风":  46,
		"穿越":  80,
		"大女主": 77,
		"青春":  47,
		"非人类": 92,
		"奇幻":  22,
		"都市":  48,
		"总裁":  52,
		"强剧情": 82,
		"玄幻":  63,
		"系统":  86,
		"悬疑":  65,
		"末世":  91,
		"热血":  67,
		"萌系":  62,
		"搞笑":  71,
		"重生":  89,
		"异能":  68,
		"冒险":  93,
		"武侠":  85,
		"竞技":  72,
		"正能量": 54,
	}
	regions := map[string]int{
		"国漫": 2,
		"韩漫": 3,
		"日漫": 4,
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
		for region, regionId := range regions {
			for pay, payId := range pays {
				for state, stateId := range states {
					fmt.Printf("%s %s %s %s \n", tag, region, pay, state)
					kk := common.Kind{
						Tag:    common.Kv{Name: tag, Val: tagId},
						Region: common.Kv{Name: region, Val: regionId},
						Pay:    common.Kv{Name: pay, Val: payId},
						State:  common.Kv{Name: state, Val: stateId},
					}
					category(kk, 1)
				}
			}
		}
	}
}

func ComicUpdate() {
	kk := common.Kind{
		Tag:    common.Kv{Name: "", Val: 0},
		Region: common.Kv{Name: "", Val: 1},
		Pay:    common.Kv{Name: "", Val: 0},
		State:  common.Kv{Name: "", Val: 0},
	}
	category(kk, 3)
}

func category(kk common.Kind, sort int) {
	url := fmt.Sprintf("https://"+config.Spe.SourceUrl+"/search/mini/topic/multi_filter?tag_id=%d&label_dimension_origin=%d&pay_status=%d&update_status=%d&sort=%d&page=%d&size=48",
		kk.Tag.Val, kk.Region.Val, kk.Pay.Val, kk.State.Val, sort, 1)
	content, err := common.RequestApi(url, "GET", "", 7)
	if err != nil {
		return
	}
	total := tools.FindStringNumber(content.Get("total").String())
	page := int(math.Ceil(float64(total) / float64(48)))
	for {
		if page < 1 {
			break
		}
		paw(kk, sort, page)
		page--
	}
}

func paw(kk common.Kind, sort, page int) {
	url := fmt.Sprintf("https://"+config.Spe.SourceUrl+"/search/mini/topic/multi_filter?tag_id=%d&label_dimension_origin=%d&pay_status=%d&update_status=%d&sort=%d&page=%d&size=48",
		kk.Tag.Val, kk.Region.Val, kk.Pay.Val, kk.State.Val, sort, page)
	content, err := common.RequestApi(url, "GET", "", 3)
	if err != nil {
		return
	}

	list := content.Get("hits.topicMessageList")
	list.ForEach(func(key, value gjson.Result) bool {
		id, _ := strconv.Atoi(value.Get("id").String())
		var exists bool
		orm.Eloquent.Model(model.SourceComic{}).Select("count(*) > 0").
			Where("source = ? and source_id = ?", config.Spe.SourceId, id).Find(&exists)
		if exists == true {
			return true
		}
		sourceComic := new(model.SourceComic)
		sourceComic.Source = 1
		sourceComic.SourceId = id

		coverUrl := strings.TrimSuffix(value.Get("cover_image_url").String(), "-t.w207.webp.h")
		var cookies map[string]string
		dir := fmt.Sprintf(config.Spe.DownloadPath+"comic/%d", id%10)
		cover := common.DownFile(coverUrl, dir, filepath.Base(coverUrl)+".webp", cookies)
		if cover != "" {
			sourceComic.Cover = cover
		}
		sourceComic.SourceUrl = "https://" + config.Spe.SourceUrl + "/web/topic/" + value.Get("id").String()
		sourceComic.Title = value.Get("title").String()
		for _, v := range value.Get("category").Array() {
			sourceComic.Category = append(sourceComic.Category, v.Str)
		}
		sourceComic.Author = value.Get("author_name").String()
		sourceComic.IsFree = 0
		sourceComic.SourceData = value.String()
		if value.Get("is_free").Bool() == false {
			sourceComic.IsFree = 1
		}
		if kk.State.Val == 2 {
			sourceComic.IsFinish = 1
		}
		err := orm.Eloquent.Create(&sourceComic).Error
		if err != nil {
			msg := fmt.Sprintf("漫画入库失败 source = %d source_id = %d", config.Spe.SourceId, id)
			model.RecordFail(url, msg, "漫画入库", 1)
		}
		rd.RPush(model.SourceComicTASK, sourceComic.Id)
		return true
	})
}
