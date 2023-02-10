package controller

import (
	"comics/common"
	"comics/global/orm"
	"comics/model"
	"comics/tools/config"
	"comics/tools/rd"
	"fmt"
	"sync"
	"time"
)

var TaskStepRecord = "task:step:record:"

func TaskComic(source *SourceStrategy) {
	t := time.NewTicker(time.Hour * 36)
	defer t.Stop()
	rd.RPush(TaskStepRecord, fmt.Sprintf("漫画-进程开始 %s %s", config.Spe.SourceUrl, time.Now().String()))
	source.ComicPaw()

	for {
		<-t.C
		rd.Delete(TaskStepRecord)
		rd.RPush(TaskStepRecord, fmt.Sprintf("漫画更新-进程开始 %s %s", config.Spe.SourceUrl, time.Now().String()))
		source.ComicUpdate()
	}
}

func TaskChapter(source *SourceStrategy) {
	t := time.NewTicker(time.Minute * 7)
	defer t.Stop()
	for {
		<-t.C
		wg := sync.WaitGroup{}
		wg.Add(2)
		rd.RPush(TaskStepRecord, fmt.Sprintf("章节-进程开始 %s %s", config.Spe.SourceUrl, time.Now().String()))
		for i := 0; i < config.Spe.Maxthreads; i++ {
			go func() {
				source.ChapterPaw()
				wg.Done()
			}()
		}
		wg.Wait()
	}
}

func TaskChapterUpdate() {
	t := time.NewTicker(time.Hour * 12)
	defer t.Stop()
	for {
		<-t.C
		rd.RPush(TaskStepRecord, fmt.Sprintf("连载漫画更新-进程开始 %s %s", config.Spe.SourceUrl, time.Now().String()))

		page := 0
		limit := 500
		for {
			var sourceComics []model.SourceComic
			orm.Eloquent.Offset(page*limit).Limit(limit).Where("source = ? and is_finish = 0", config.Spe.SourceId).Find(&sourceComics)
			if len(sourceComics) == 0 {
				break
			}
			page = page + 1
			for _, sourceComic := range sourceComics {
				rd.RPush(common.SourceComicTASK, sourceComic.Id)
			}
		}
	}
}

func TaskImage(source *SourceStrategy) {
	t := time.NewTicker(time.Minute * 9)
	defer t.Stop()
	for {
		<-t.C
		wg := sync.WaitGroup{}
		wg.Add(3)
		rd.RPush(TaskStepRecord, fmt.Sprintf("图片-进程开始 %s %s", config.Spe.SourceUrl, time.Now().String()))
		for i := 0; i < config.Spe.Maxthreads; i++ {
			go func() {
				source.ImagePaw()
				wg.Done()
			}()
		}
		wg.Wait()
	}
}

func TaskDownImage() {
	t := time.NewTicker(time.Minute * 5)
	defer t.Stop()
	for {
		<-t.C
		wg := sync.WaitGroup{}
		wg.Add(3)
		rd.RPush(TaskStepRecord, fmt.Sprintf("图片下载-进程开始 %s %s", config.Spe.SourceUrl, time.Now().String()))
		for i := 0; i < config.Spe.Maxthreads; i++ {
			go func() {
				var ext string
				if config.Spe.SourceId == 1 {
					ext = "webp"
				} else {
					ext = "jpg"
				}
				DownImage(ext)
				wg.Done()
			}()
		}
		wg.Wait()
	}
}
