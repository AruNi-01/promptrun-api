package errs

// Errs 错误处理，service 层收集
type Errs struct {
	ErrCode int
	Err     error
}

func NewErrs(errCode int, err error) *Errs {
	return &Errs{
		ErrCode: errCode,
		Err:     err,
	}
}
