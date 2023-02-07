package controller

import (
	"comics/controller/kk"
	"comics/controller/tx"
)

type SourceStrategy struct {
	ComicPaw      func()
	ComicUpdate   func()
	ChapterPaw    func()
	ChapterUpdate func()
	ImagePaw      func()
}

func SourceOperate(source string) *SourceStrategy {
	switch source {
	case "www.kuaikanmanhua.com":
		return &SourceStrategy{
			ComicPaw:    kk.ComicPaw,
			ComicUpdate: kk.ComicUpdate,
			ChapterPaw:  kk.ChapterPaw,
			ImagePaw:    kk.ImagePaw,
		}
	case "ac.qq.com":
		return &SourceStrategy{
			ComicPaw:    tx.ComicPaw,
			ComicUpdate: tx.ComicUpdate,
			ChapterPaw:  tx.ChapterPaw,
			ImagePaw:    tx.ImagePaw,
		}
	}
	return &SourceStrategy{}
}
