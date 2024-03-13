package api

import "promptrun-api/common/result"

// SuccessResponse 返回成功
func SuccessResponse(data interface{}) result.Response {
	return result.Succ(data)
}

// ErrorResponse 返回错误消息
func ErrorResponse(errCode int, errMsg string) result.Response {
	return result.Err(errCode, errMsg)
}
