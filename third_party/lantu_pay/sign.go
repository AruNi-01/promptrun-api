package lantu_pay

import (
	"crypto/md5"
	"encoding/hex"
	"os"
	"sort"
	"strings"
)

func packageSign(params map[string]string) string {
	// 先将参数以其参数名的字典序升序进行排序，用切片保存有序的 key
	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// 拼接参数键值对
	var builder strings.Builder
	for _, key := range keys {
		value := params[key]
		if value != "" {
			builder.WriteString(key + "=" + value + "&")
		}
	}
	return builder.String()
}

func genSign(params map[string]string) string {
	// 生成签名前先去除sign
	delete(params, "sign")
	stringA := packageSign(params)
	stringSignTemp := stringA + "key=" + os.Getenv("SECRET_KEY")

	// 计算 MD5 哈希值
	hash := md5.Sum([]byte(stringSignTemp))
	sign := hex.EncodeToString(hash[:])

	return strings.ToUpper(sign)
}
