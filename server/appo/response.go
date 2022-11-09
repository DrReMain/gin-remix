package appo

type BaseResponse struct {
	T       int64 `json:"t"`
	Success bool  `json:"success"`
}

type FailResponse struct {
	BaseResponse
	ErrCode string `json:"err_code"`
	Message string `json:"msg"`
}

type SuccessResponse struct {
	BaseResponse
	Result any `json:"result"`
}
