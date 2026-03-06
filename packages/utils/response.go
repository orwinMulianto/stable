package utils

import (
	"strconv"
	"time"
)

type Response struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Error   any    `json:"error,omitempty"`
	Data    any    `json:"data,omitempty"`
	Meta    any    `json:"meta,omitempty"`
}

type EmptyObj struct{}

func BuildResponseSuccess(message string, data any) Response {
	res := Response{
		Status:  true,
		Message: message,
		Data:    data,
	}
	return res
}

func BuildResponseFailed(message string, err string, data any) Response {
	res := Response{
		Status:  false,
		Message: message,
		Error:   err,
		Data:    data,
	}
	return res
}

// ValidationError represents a single field validation error
type ValidationError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

// BuildValidationErrorResponse creates a response with structured validation errors
func BuildValidationErrorResponse(message string, errors []ValidationError) Response {
	return Response{
		Status:  false,
		Message: message,
		Error:   errors,
		Data:    nil,
	}
}

func StringToInt(str string) (int, error) {
	result, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}
	return int(result), nil
}

type Meta struct {
	Timestamp string `json:"timestamp"`
	Version   string `json:"version"`
	RequestID string `json:"request_id,omitempty"`
	Page      *int   `json:"page,omitempty"`
	Limit     *int   `json:"limit,omitempty"`
	Total     *int64 `json:"total,omitempty"`
}

func DefaultMeta() Meta {
	return Meta{
		Timestamp: time.Now().Format(time.RFC3339),
		Version:   "v1",
	}
}

func BuildResponseSuccessWithMeta(
	message string,
	data any,
	meta Meta,
) Response {
	res := Response{
		Status:  true,
		Message: message,
		Data:    data,
		Meta:    meta,
	}
	return res
}
