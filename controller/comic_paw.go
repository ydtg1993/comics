package controller

import (
	"fmt"
	"os"
)

func ComicPaw() {
	sourceUrl := os.Getenv("SOURCE_URL")
	fmt.Println(sourceUrl)
}
