package errs

// Success errCode 除了 0，其他都是有错误
const (
	Success = 0
)

// 三位数错误编码为复用 http 原本含义，五位数错误编码为应用自定义错误
// 5 开头的五位数错误编码为服务器端错误，4 开头的五位数错误编码为客户端错误

const (
	ErrNotLogin = 401
)

const (
	ErrParam = 40001 + iota

	ErrConfirmPasswordDiff
	ErrEmailExist
	ErrUserNotExist
	ErrWrongPassword

	ErrRecordNotFound

	ErrLikeIntervalTooShort

	ErrBalanceNotEnough
)

const (
	ErrDBError          = 50001 + iota // CodeDBError 数据库操作失败
	ErrEncryptError                    // CodeEncryptError 加密失败
	ErrJsonConvertError                // JsonConvertError json 转换异常

	ErrUploadImgToOSS // 上传图片到 OSS 失败

	ErrPayFacadeError     // 支付接口调用失败
	ErrPayOrderQueryError // 支付订单查询失败

)
