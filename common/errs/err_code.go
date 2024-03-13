package errs

// Success errCode 除了 0，其他都是有错误
const (
	Success = 0
)

// 三位数错误编码为复用 http 原本含义，五位数错误编码为应用自定义错误
// 5 开头的五位数错误编码为服务器端错误，4 开头的五位数错误编码为客户端错误

const (
	CodeCheckLogin = 401 // CodeCheckLogin 未登录
	CodeNoRightErr = 403 // CodeNoRightErr 未授权访问
	CodeParamErr   = 500 // CodeParamErr 参数错误
)

const (
	ErrParam = 40001 + iota
	ErrConfirmPasswordDiff
	ErrEmailExist
	ErrUserNotExist
	ErrWrongPassword
)

const (
	ErrDBError          = 50001 + iota // CodeDBError 数据库操作失败
	ErrEncryptError                    // CodeEncryptError 加密失败
	ErrJsonConvertError                // JsonConvertError json 转换异常

)
