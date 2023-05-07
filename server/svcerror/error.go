package svcerror

// 自定义错误字段
type SvcErr struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (e SvcErr) Error() string {
	return e.Msg
}

func New(code int, err error) *SvcErr {
	return &SvcErr{
		Msg:  err.Error(),
		Code: code,
	}
}

func AppendData(e *SvcErr, data interface{}) *SvcErr {
	return &SvcErr{
		Msg:  e.Msg,
		Code: e.Code,
		Data: data,
	}
}
