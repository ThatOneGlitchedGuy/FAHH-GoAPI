package errors

import (
	"fmt"
	"net/http"
)

type ErrorCode string
type ErrorSeverity string

const (
	ErrValidation       ErrorCode = "VALIDATION_ERROR"
	ErrNotFound         ErrorCode = "NOT_FOUND"
	ErrUnauthorized     ErrorCode = "UNAUTHORIZED"
	ErrForbidden        ErrorCode = "FORBIDDEN"
	ErrConflict         ErrorCode = "CONFLICT"
	ErrInternalServer   ErrorCode = "INTERNAL_SERVER_ERROR"
	ErrServiceUnavail   ErrorCode = "SERVICE_UNAVAILABLE"
	ErrBadRequest       ErrorCode = "BAD_REQUEST"
	ErrRateLimit        ErrorCode = "RATE_LIMIT_EXCEEDED"

	SeverityLow    ErrorSeverity = "LOW"
	SeverityMedium ErrorSeverity = "MEDIUM"
	SeverityHigh   ErrorSeverity = "HIGH"
	SeverityCritical ErrorSeverity = "CRITICAL"
)

type AppError struct {
	Code       ErrorCode      `json:"code"`
	Message    string         `json:"message"`
	Details    string         `json:"details,omitempty"`
	Severity   ErrorSeverity  `json:"severity"`
	StatusCode int            `json:"status_code"`
	RequestID  string         `json:"request_id,omitempty"`
	Timestamp  int64          `json:"timestamp"`
	Path       string         `json:"path,omitempty"`
}

func NewAppError(code ErrorCode, message string, statusCode int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
		Severity:   SeverityMedium,
	}
}

func (e *AppError) WithDetails(details string) *AppError {
	e.Details = details
	return e
}

func (e *AppError) WithSeverity(sev ErrorSeverity) *AppError {
	e.Severity = sev
	return e
}

func (e *AppError) WithRequestID(id string) *AppError {
	e.RequestID = id
	return e
}

func (e *AppError) WithPath(path string) *AppError {
	e.Path = path
	return e
}

func (e *AppError) WithTimestamp(ts int64) *AppError {
	e.Timestamp = ts
	return e
}

func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s: %s", e.Code, e.Message, e.Details)
}

func ValidationError(message string) *AppError {
	return NewAppError(ErrValidation, message, http.StatusBadRequest).WithSeverity(SeverityLow)
}

func NotFoundError(resource string) *AppError {
	return NewAppError(ErrNotFound, fmt.Sprintf("%s not found", resource), http.StatusNotFound).WithSeverity(SeverityLow)
}

func UnauthorizedError(message string) *AppError {
	return NewAppError(ErrUnauthorized, message, http.StatusUnauthorized).WithSeverity(SeverityMedium)
}

func ForbiddenError(message string) *AppError {
	return NewAppError(ErrForbidden, message, http.StatusForbidden).WithSeverity(SeverityMedium)
}

func ConflictError(resource string) *AppError {
	return NewAppError(ErrConflict, fmt.Sprintf("%s already exists", resource), http.StatusConflict).WithSeverity(SeverityLow)
}

func InternalServerError(message string) *AppError {
	return NewAppError(ErrInternalServer, message, http.StatusInternalServerError).WithSeverity(SeverityHigh)
}

func ServiceUnavailableError() *AppError {
	return NewAppError(ErrServiceUnavail, "Service temporarily unavailable", http.StatusServiceUnavailable).WithSeverity(SeverityCritical)
}

func RateLimitError() *AppError {
	return NewAppError(ErrRateLimit, "Rate limit exceeded", http.StatusTooManyRequests).WithSeverity(SeverityMedium)
}

func BadRequestError(message string) *AppError {
	return NewAppError(ErrBadRequest, message, http.StatusBadRequest).WithSeverity(SeverityLow)
}
