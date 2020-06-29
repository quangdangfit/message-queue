package utils

import (
	"net/http"
)

const (
	StatusOK           = "OK"
	StatusError        = "ERROR"
	StatusBadRequest   = "BAD_REQUEST"
	StatusUnAuthorized = "UN_AUTHORIZED"
	StatusNotFound     = "NOT_FOUND"
	StatusDBError      = "BD_ERROR"
)

var Res = map[string]objResponse{
	StatusOK:           {HttpCode: http.StatusOK, Code: StatusOK, Message: "Thành công"},
	StatusError:        {HttpCode: http.StatusBadRequest, Code: StatusError, Message: "Có lỗi xảy ra, vui lòng thử lại"},
	StatusUnAuthorized: {HttpCode: http.StatusUnauthorized, Code: StatusUnAuthorized, Message: "Bạn không có quyền truy cập vào hệ thống. Vui lòng liên hệ quản lý để được phân quyền."},
	StatusBadRequest:   {HttpCode: http.StatusBadRequest, Code: StatusBadRequest, Message: "Thông tin gửi lên không hợp lệ, vui lòng kiểm tra lại"},
	StatusNotFound:     {HttpCode: http.StatusNotFound, Code: StatusNotFound, Message: "Không tìm thấy thông tin dữ liệu trên hệ thống, vui lòng kiểm tra lại"},
	StatusDBError:      {HttpCode: http.StatusBadRequest, Code: StatusDBError, Message: "Hệ thống xử lý dữ liệu không thành công, vui lòng thử lại"},
}

type objResponse struct {
	HttpCode        int    `json:"http_code"`
	Code            string `json:"code"`
	Message         string `json:"message"`
	OriginalMessage string `json:"original_message,omitempty"`
}

func MsgResponse(code string, internalErr *error) objResponse {
	if _, ok := Res[code]; !ok {
		code = StatusError
	}
	if internalErr != nil {
		return objResponse{
			Res[code].HttpCode,
			Res[code].Code,
			Res[code].Message,
			(*internalErr).Error()}
	}
	return Res[code]
}
