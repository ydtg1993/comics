package tools

import (
	"math/rand"
	"regexp"
	"strconv"
	"time"
)

func StringToInt64(e string) (int64, error) {
	return strconv.ParseInt(e, 10, 64)
}

func StringToInt(e string) (int, error) {
	return strconv.Atoi(e)
}

func GetCurrentTimeStr() string {
	return time.Now().Format("2006/01/02 15:04:05")
}

func GetCurrentTime() time.Time {
	return time.Now()
}

func FindStringNumber(s string) int {
	pattern := regexp.MustCompile(`\d+`)
	numberStrings := pattern.FindAllStringSubmatch(s, -1)

	result := ""
	for _, number := range numberStrings[0] {
		result += number
	}
	res, err := strconv.Atoi(result)
	if err != nil {
		return 0
	}
	return res
}

func RandStr(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	var result []byte
	rand.Seed(time.Now().UnixNano() + int64(rand.Intn(100)))
	for i := 0; i < length; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}
	return string(result)
}
