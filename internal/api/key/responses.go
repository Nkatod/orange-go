package key

import "orange-go/internal/model"

type KeyValuePutData struct {
	Key    string `json:"key"`
	Status string `json:"status"`
}

type KeyValuePutResponse struct {
	Data  *KeyValuePutData `json:"data"`
	Error *model.BaseError `json:"error,omitempty"`
}

type KeyValueGetResponse struct {
	Data  string           `json:"data"`
	Error *model.BaseError `json:"error,omitempty"`
}
