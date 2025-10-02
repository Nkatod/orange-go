package errors

import (
	"errors"
	"fmt"

	"github.com/labstack/echo/v4"

	"orange-go/internal/model"
	"orange-go/internal/utils"
)

type CommonError struct {
	msg  string
	code Code
}

func NewCommonError(msg string, code Code) *CommonError {
	return &CommonError{msg, code}
}

func (r *CommonError) Error() string {
	return r.msg
}

func (r *CommonError) Code() Code {
	return r.code
}

func IsCommonError(err error) bool {
	var ce *CommonError
	return errors.As(err, &ce)
}

func GetCommonError(err error) *CommonError {
	var ce *CommonError
	if !errors.As(err, &ce) {
		return nil
	}

	return ce
}

// SendErrorJSON makes {error: blah, details: blah} json body and responds with error code
func SendErrorJSON(c echo.Context, httpCode int, err error, debug string) error {
	if IsCommonError(err) {
		cerr := GetCommonError(err)
		errorResponse := model.BaseResponse{
			Data: nil,
			Error: &model.BaseError{
				Code:    utils.StringPtr(fmt.Sprintf("%d", cerr.Code())),
				Debug:   utils.StringPtr(debug),
				Message: utils.StringPtr(err.Error()),
			},
		}
		return c.JSON(httpCode, errorResponse)
	}

	errorResponse := model.BaseResponse{
		Data: nil,
		Error: &model.BaseError{
			Code:    utils.StringPtr(fmt.Sprintf("%d", Unknown)),
			Debug:   utils.StringPtr(debug),
			Message: utils.StringPtr(err.Error()),
		},
	}
	return c.JSON(httpCode, errorResponse)
}
