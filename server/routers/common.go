package routers

type Response struct {
	Succeeded bool        `json:"succeeded"`
	ErrMsg    string      `json:"errMsg"`
	Data      interface{} `json:"data"`
}

func makeResponse(succeeded bool, err error, data interface{}) *Response {
	var errMsg string
	if err != nil {
		errMsg = err.Error()
	}
	return &Response{Succeeded: succeeded, ErrMsg: errMsg, Data: data}
}
