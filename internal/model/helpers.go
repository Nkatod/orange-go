package model

func NewErr(code, msg, debug string) *BaseError {
	return &BaseError{
		Code:    StrPtr(code),
		Message: StrPtr(msg),
		Debug:   StrPtr(debug),
	}
}

func StrPtr(s string) *string { return &s }
