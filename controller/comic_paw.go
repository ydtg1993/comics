package controller

import (
	"comics/model"
	"fmt"
)

func ComicPaw() {
	fmt.Println("hk")
	comic := new(model.Comic)
	var total int64
	model.GetGormDb().Model(comic).Count(&total)
	fmt.Println(total)
}
