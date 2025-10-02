package model

type SuccessData struct {
	Message string `json:"message"`
	Code    uint32 `json:"code"`
}

type Success struct {
	Data SuccessData `json:"data"`
}
