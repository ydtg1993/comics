package controller

import (
	"comics/common"
	"comics/global/orm"
	"comics/model"
	"comics/tools/config"
	"comics/tools/rd"
	"fmt"
)

func DownImage(ext string) {
	taskLimit := 15
	for limit := 0; limit < taskLimit; limit++ {
		id, err := rd.LPop(common.SourceImageTASK)
		if err != nil || id == "" {
			return
		}

		sourceImage := new(model.SourceImage)
		if orm.Eloquent.Where("id = ?", id).First(&sourceImage); sourceImage.Id == 0 {
			continue
		}
		sourceChapter := new(model.SourceChapter)
		orm.Eloquent.Where("id = ?", sourceImage.ChapterId).First(&sourceChapter)

		dir := fmt.Sprintf(config.Spe.DownloadPath+"chapter/%d/%d/%d/%d",
			config.Spe.SourceId, sourceChapter.ComicId%128, sourceChapter.ComicId, sourceImage.ChapterId)
		for key, img := range sourceImage.SourceData {
			state := 0
			for i := 0; i < 3; i++ {
				file := common.DownFile(img, dir, fmt.Sprintf("%d.%s", key, ext), map[string]string{})
				if file != "" {
					state = 1
					sourceImage.Images = append(sourceImage.Images, file)
					break
				}
			}
			if state == 0 {
				sourceImage.Images = model.Images{}
				sourceImage.State = state
				break
			}
			sourceImage.State = state
		}
		orm.Eloquent.Save(&sourceImage)
	}
}
