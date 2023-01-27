package tools

import (
	"io/ioutil"
	"os"

	"github.com/beego/beego/v2/core/logs"
)

/**
* 创建文件目录,包含多层目录
 */
func CreateFile(src string) (string, error) {
	//	src := dir + name + "/"
	if IsExist(src) {
		return src, nil
	}

	if err := os.MkdirAll(src, 0777); err != nil {
		if os.IsPermission(err) {
			logs.Error("创建文件失败，不够权限创建文件", src)
		}
		return "", err
	}

	return src, nil
}

/*
* 判断文件或者目录是否存在
 */
func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

/**
* 新建文件并写入内容
* 如果文件已存在,则覆盖以前内容
 */
func WriteFile(filePath, fileName, content string) (int, error) {
	_, err := CreateFile(filePath)
	if err != nil {
		return 0, err
	}
	src := filePath + "/" + fileName
	fs, e := os.Create(src)
	if e != nil {
		return 0, e
	}
	defer fs.Close()
	return fs.WriteString(content)
}

/**
* 获取文件大小,单位是Byte
 */
func GetFileSize(file string) (int64, error) {
	f, e := os.Stat(file)
	if e != nil {
		return 0, e
	}
	return f.Size(), nil
}

/**
* 读取文件内容并返回字符串
* @param  path  文件路径
 */
func ReadFileString(path string) string {
	fi, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)
	return string(fd)
}

/**
* 根据文件头，识别指定文件的后缀名
 */
func GetExtByMime(key string) string {
	val := ""
	switch key {
	case "image/jpeg":
		val = ".jpg"
	case "image/bmp":
		val = ".bmp"
	case "image/gif":
		val = ".gif"
	case "image/png":
		val = ".png"
	case "image/x-icon":
		val = ".ico"
	case "image/tiff":
		val = ".tif"
	case "application/octet-stream":
		val = ".*"
	case "application/x-jpg":
		val = ".jpg"
	case "video/mpeg4":
		val = ".mp4	"
	case "video/mpeg":
		val = ".mp2v"
	case "audio/x-ms-wma":
		val = ".wma"
	case "video/x-sgi-movie":
		val = ".movie"
	case "video/quicktime":
		val = ".avi"
	case "video/x-ms-wmv":
		val = ".wmv"
	case "video/x-flv":
		val = ".flv"
	case "video/x-msvideo":
		val = ".avi"
	case "application/x-mpegurl":
		val = ".m3u8"
	}
	return val
}
