package appo

import (
	"fmt"
	"time"
)

func NewSuccess(result any) *SuccessResponse {
	return &SuccessResponse{
		BaseResponse: BaseResponse{
			T:       time.Now().UnixMilli(),
			Success: true,
		},
		Result: result,
	}
}

func NewFail(errCode, message string) *FailResponse {
	return &FailResponse{
		BaseResponse: BaseResponse{
			T:       time.Now().UnixMilli(),
			Success: false,
		},
		ErrCode: errCode,
		Message: message,
	}
}

func NewInternal() *FailResponse {
	return NewFail(Internal, "服务器内部错误，请稍后重试")
}

func NewServiceUnavailable() *FailResponse {
	return NewFail(ServiceUnavailable, "服务不可用或连接超时")
}

func NewNotFound(name, value string) *FailResponse {
	return NewFail(NotFound, fmt.Sprintf("resource: %v with value: %v not found", name, value))
}
