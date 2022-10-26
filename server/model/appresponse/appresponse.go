package appresponse

import "time"

type Ok struct {
	T       int64 `json:"t"`
	Success bool  `json:"success"`
	Result  any   `json:"result"`
}

func New(result any) *Ok {
	return &Ok{
		T:       time.Now().UnixMilli(),
		Success: true,
		Result:  result,
	}
}
