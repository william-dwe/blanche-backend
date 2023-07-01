package httperror

import (
	"net/http"
)

type AppError struct {
	StatusCode int    `json:"-"`
	Code       string `json:"code"`
	Message    string `json:"message"`
	Data       any    `json:"data"`
}

func (err AppError) Error() string {
	return err.Message
}

func BadRequestError(message string, code string) AppError {
	if code == "" {
		code = "BAD_REQUEST"
	}
	return AppError{
		Code:       code,
		Message:    message,
		StatusCode: http.StatusBadRequest,
		Data:       nil,
	}
}

func NotFoundError(message string) AppError {
	return AppError{
		Code:       "NOT_FOUND_ERROR",
		Message:    message,
		StatusCode: http.StatusNotFound,
		Data:       nil,
	}
}

func InternalServerError(message string) AppError {
	return AppError{
		Code:       "INTERNAL_SERVER_ERROR",
		Message:    message,
		StatusCode: http.StatusInternalServerError,
		Data:       nil,
	}
}

func UnauthorizedError() AppError {
	return AppError{
		Code:       "UNAUTHORIZED_ERROR",
		Message:    "Unauthorized",
		StatusCode: http.StatusUnauthorized,
		Data:       nil,
	}
}

func UnauthorizedErrorLogin() AppError {
	return AppError{
		Code:       "UNAUTHORIZED_ERROR",
		Message:    "Email or password is incorrect",
		StatusCode: http.StatusUnauthorized,
		Data:       nil,
	}
}

func ForbiddenError() AppError {
	return AppError{
		Code:       "FORBIDDEN_ERROR",
		Message:    "Don't have permission to access",
		StatusCode: http.StatusForbidden,
		Data:       nil,
	}
}

func ForbiddenErrorMsg(msg string) AppError {
	return AppError{
		Code:       "FORBIDDEN_ERROR",
		Message:    msg,
		StatusCode: http.StatusForbidden,
		Data:       nil,
	}
}

func TimeoutError() AppError {
	return AppError{
		Code:       "TIMEOUT_ERROR",
		Message:    "Timeout",
		StatusCode: http.StatusRequestTimeout,
		Data:       nil,
	}
}
