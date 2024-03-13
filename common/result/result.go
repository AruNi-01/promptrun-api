package result

import "promptrun-api/common/errs"

// Response 基础响应信息
type Response struct {
	ErrCode int         `json:"errCode"`
	ErrMsg  string      `json:"errMsg"`
	Data    interface{} `json:"data,omitempty"`
}

// Succ 成功返回
func Succ(data interface{}) Response {
	return Response{ErrCode: errs.Success, ErrMsg: "", Data: data}
}

// Err 错误返回
func Err(errCode int, errMsg string) Response {
	return Response{
		ErrCode: errCode,
		ErrMsg:  errMsg,
	}
}
