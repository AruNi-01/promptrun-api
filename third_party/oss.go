package third_party

import (
	"bytes"
	"encoding/base64"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"os"
	"promptrun-api/utils"
	"strings"
)

var (
	OSS *oss.Bucket
)

const (
	OSSPrefixHeaderImg = "header_img/"
	OSSPrefixPromptImg = "prompt_img/"
)

func OSSInit() {
	client, err := oss.New(os.Getenv("OSS_ENDPOINT"), os.Getenv("OSS_ACCESS_KEY"), os.Getenv("OSS_ACCESS_SECRET"))
	if err != nil {
		utils.Log().Panic("", "OSS 初始化失败，errMsg: %s", err.Error())
		panic(err)
	}
	bucket, err := client.Bucket(os.Getenv("OSS_BUCKET"))
	if err != nil {
		utils.Log().Panic("", "OSS 初始化失败，errMsg: %s", err.Error())
		panic(err)
	}
	OSS = bucket
}

func UploadBase64ImgToOSS(objectName, base64Img string) (string, error) {
	// 去掉 base64 图片前缀（部分）
	content := strings.TrimPrefix(base64Img, "data:image/")
	// 获取图片类型（后缀 png、jpeg、jpg...）
	imgType := strings.Split(content, ";")[0]
	// 去掉 base64 图片前缀（剩余）
	content = strings.TrimPrefix(strings.Split(content, ";")[1], "base64,")
	// 解码 base64 图片，得到 byte 数组
	byteData, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return "", err
	}

	objectName = objectName + "." + imgType
	err = OSS.PutObject(objectName, bytes.NewReader(byteData))
	if err != nil {
		return "", err
	}

	return "https://" + os.Getenv("OSS_BUCKET") + "." + os.Getenv("OSS_ENDPOINT") + "/" + objectName, nil
}

func DeleteOSSImgByUrl(path string) error {
	path = strings.TrimPrefix(path, "https://"+os.Getenv("OSS_BUCKET")+"."+os.Getenv("OSS_ENDPOINT")+"/")
	return OSS.DeleteObject(path)
}
