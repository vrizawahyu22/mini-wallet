package utils

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

const (
	shouldRetryTx = "should retry transaction"
)

var Logger *zap.Logger

type AppError struct {
	Message    string
	StatusCode int
}

type ValidationError struct {
	Key     string
	Message string
}

type ValidationErrors struct {
	Errors     []ValidationError
	StatusCode int
}

type SuccessResponse[T any] struct {
	Status string `json:"status"`
	Data   T      `json:"data"`
}

type DataError[T any] struct {
	Error T `json:"error"`
}

type FailedResponse[T any] struct {
	Status string       `json:"status"`
	Data   DataError[T] `json:"data"`
}

func InitLogger() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(fmt.Sprintf("Error initializing logger: %v", err))
	}

	Logger = logger
}

func LogAndPanicIfError(err error, message string) {
	if err != nil {
		errMsg := fmt.Sprintf("%s :%v", message, err)
		LogError(errMsg, zap.Error(err))
		panic(err)
	}
}

func LogError(message string, fields ...zap.Field) {
	Logger.Error(message, fields...)
}

func LogWarning(message string, fields ...zap.Field) {
	Logger.Warn(message, fields...)
}

func LogIfError(err error, message ...string) {
	if err != nil {
		if len(message) > 0 && message[0] != "" {
			LogError(message[0], zap.Error(err))
		} else {
			LogError("error occurred", zap.Error(err))
		}
	}
}

func LogInfo(message string, fields ...zap.Field) {
	Logger.Info(message, fields...)
}

func CustomErrorWithTrace(err error, message string, statusCode int, isShouldRetry ...bool) error {
	if len(isShouldRetry) > 0 && isShouldRetry[0] {
		return fmt.Errorf("%s|%s<->%d|%s", message, message, statusCode, shouldRetryTx)
	}
	return fmt.Errorf("%s|%s<->%d|", err.Error(), message, statusCode)
}

func PanicIfError(err error) {
	if err != nil {
		customError := strings.Split(err.Error(), "<->")
		message := customError[0]
		statusCode := 500

		if len(customError) > 1 {
			statusCode, _ = strconv.Atoi(strings.Split(customError[1], "|")[0])
		}

		appErr := AppError{
			Message:    message,
			StatusCode: statusCode,
		}
		panic(appErr)
	}
}

func PanicValidationError(errors []ValidationError, statusCode int) {
	validationErrors := ValidationErrors{
		Errors:     errors,
		StatusCode: statusCode,
	}
	panic(validationErrors)
}

func PanicIfAppError(err error, message string, statusCode int) {
	if err != nil {
		customErr := CustomErrorWithTrace(err, message, statusCode)
		PanicIfError(customErr)
	}
}

func GenerateSuccessResp[T any](w http.ResponseWriter, data T, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := SuccessResponse[T]{
		Status: "success",
		Data:   data,
	}

	responseEncode, err := Marshal(response)
	PanicIfAppError(err, "failed when marshal response", 500)

	_, err = w.Write(responseEncode)
	PanicIfAppError(err, "failed when write success response", 500)
}

func GenerateErrorResp[T any](w http.ResponseWriter, data T, statusCode int, errorCode ...int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := FailedResponse[T]{
		Status: "fail",
		Data: DataError[T]{
			Error: data,
		},
	}

	responseEncode, err := Marshal(response)
	PanicIfAppError(err, "failed when marshar response", 500)

	_, err = w.Write(responseEncode)
	PanicIfAppError(err, "failed when write success response", 500)
}

func CustomError(message string, statusCode int, isShouldRetry ...bool) error {
	if len(isShouldRetry) > 0 && isShouldRetry[0] {
		return fmt.Errorf("%s|%s<->%d|%s", message, message, statusCode, shouldRetryTx)
	}
	return fmt.Errorf("%s|%s<->%d|", message, message, statusCode)
}
