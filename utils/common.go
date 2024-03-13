package utils

import (
	"github.com/google/uuid"
	"strings"
)

// GenUUID 生成 32bit 纯字符串的随机串
func GenUUID() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}
