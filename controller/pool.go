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

func TaskComic(source *SourceStrategy) {
	t := time.NewTicker(time.Hour * 36)
	defer t.Stop()
	rd.RPush(common.TaskStepRecord, fmt.Sprintf("漫画-进程开始 %s %s", config.Spe.SourceUrl, time.Now().String()))
	source.ComicPaw()

	for {
		<-t.C
		rd.Delete(common.TaskStepRecord)
		rd.RPush(common.TaskStepRecord, fmt.Sprintf("漫画更新-进程开始 %s %s", config.Spe.SourceUrl, time.Now().String()))
		source.ComicUpdate()
	}
}

func TaskChapter(source *SourceStrategy) {
	t := time.NewTicker(time.Minute * 15)
	defer t.Stop()
	threads := 2
	for {
		<-t.C
		wg := sync.WaitGroup{}
		wg.Add(threads)
		rd.RPush(common.TaskStepRecord, fmt.Sprintf("章节-进程开始 %s %s", config.Spe.SourceUrl, time.Now().String()))
		for i := 0; i < threads; i++ {
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
		rd.RPush(common.TaskStepRecord, fmt.Sprintf("连载漫画更新-进程开始 %s %s", config.Spe.SourceUrl, time.Now().String()))

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

func TaskImage(source *SourceStrategy, doneImageSignal chan struct{}) {
	threads := 3
	for {
		<-doneImageSignal
		wg := sync.WaitGroup{}
		wg.Add(threads)
		rd.Set(common.SourceImageCapture,
			fmt.Sprintf("图片抓取 %s", time.Now().String()),
			time.Hour*1)
		for i := 0; i < threads; i++ {
			go func() {
				source.ImagePaw()
				wg.Done()
			}()
		}
		wg.Wait()
	}
}

func TaskDownImage(doneImageSignal chan struct{}) {
	threads := 6
	for {
		wg := sync.WaitGroup{}
		wg.Add(threads)
		rd.Set(common.SourceImageDownload,
			fmt.Sprintf("图片下载 %s", time.Now().String()),
			time.Hour*1)
		for i := 0; i < threads; i++ {
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
			t := time.NewTicker(time.Second * 40)
			<-t.C
		}
		wg.Wait()
		doneImageSignal <- struct{}{}
	}
}
