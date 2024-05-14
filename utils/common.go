package utils

import (
	"github.com/bwmarrin/snowflake"
	"github.com/google/uuid"
	"strings"
	"time"
)

// GenUUID 生成 32bit 纯字符串的随机串
func GenUUID() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

// GenSnowflakeId 生成雪花 ID
func GenSnowflakeId() int64 {
	node, err := snowflake.NewNode(1)
	if err != nil {
		Log().Error("", "snowflake.NewNode error: %s", err.Error())
	}

	// 截取低 16 位作为唯一 ID
	return node.Generate().Int64() / 1000
}

func GetShanghaiLocation() *time.Location {
	location, _ := time.LoadLocation("Asia/Shanghai")
	return location
}
