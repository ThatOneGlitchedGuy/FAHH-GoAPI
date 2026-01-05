package response

import (
	"net/http"
	"time"
)

type PageInfo struct {
	CurrentPage int   `json:"current_page"`
	PageSize    int   `json:"page_size"`
	TotalItems  int64 `json:"total_items"`
	TotalPages  int   `json:"total_pages"`
	HasNext     bool  `json:"has_next"`
	HasPrev     bool  `json:"has_prev"`
}

type PaginatedResponse struct {
	Success   bool        `json:"success"`
	Data      interface{} `json:"data"`
	PageInfo  *PageInfo   `json:"page_info,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
	RequestID string      `json:"request_id,omitempty"`
}

type DataResponse struct {
	Success   bool        `json:"success"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
	RequestID string      `json:"request_id,omitempty"`
}

type ErrorResponse struct {
	Success   bool        `json:"success"`
	Error     interface{} `json:"error"`
	Timestamp time.Time   `json:"timestamp"`
	RequestID string      `json:"request_id,omitempty"`
}

type MessageResponse struct {
	Success   bool      `json:"success"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	RequestID string    `json:"request_id,omitempty"`
}

func NewDataResponse(data interface{}, requestID string) *DataResponse {
	return &DataResponse{
		Success:   true,
		Data:      data,
		Timestamp: time.Now(),
		RequestID: requestID,
	}
}

func NewPaginatedResponse(data interface{}, pageInfo *PageInfo, requestID string) *PaginatedResponse {
	return &PaginatedResponse{
		Success:   true,
		Data:      data,
		PageInfo:  pageInfo,
		Timestamp: time.Now(),
		RequestID: requestID,
	}
}

func NewErrorResponse(err interface{}, requestID string) *ErrorResponse {
	return &ErrorResponse{
		Success:   false,
		Error:     err,
		Timestamp: time.Now(),
		RequestID: requestID,
	}
}

func NewMessageResponse(message string, requestID string) *MessageResponse {
	return &MessageResponse{
		Success:   true,
		Message:   message,
		Timestamp: time.Now(),
		RequestID: requestID,
	}
}

func CalculatePageInfo(currentPage, pageSize int, totalItems int64) *PageInfo {
	totalPages := int((totalItems + int64(pageSize) - 1) / int64(pageSize))
	if totalPages == 0 {
		totalPages = 1
	}

	return &PageInfo{
		CurrentPage: currentPage,
		PageSize:    pageSize,
		TotalItems:  totalItems,
		TotalPages:  totalPages,
		HasNext:     currentPage < totalPages,
		HasPrev:     currentPage > 1,
	}
}

type SuccessCode string

const (
	Created        SuccessCode = "CREATED"
	Updated        SuccessCode = "UPDATED"
	Deleted        SuccessCode = "DELETED"
	Retrieved      SuccessCode = "RETRIEVED"
	Retrieved_List SuccessCode = "RETRIEVED_LIST"
)

type DetailedResponse struct {
	Success   bool        `json:"success"`
	Code      SuccessCode `json:"code"`
	Data      interface{} `json:"data"`
	Message   string      `json:"message,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
	RequestID string      `json:"request_id,omitempty"`
	StatusCode int        `json:"status_code"`
}

func NewDetailedResponse(code SuccessCode, data interface{}, statusCode int, requestID string) *DetailedResponse {
	return &DetailedResponse{
		Success:    true,
		Code:       code,
		Data:       data,
		Timestamp:  time.Now(),
		RequestID:  requestID,
		StatusCode: statusCode,
	}
}

func (r *DetailedResponse) WithMessage(msg string) *DetailedResponse {
	r.Message = msg
	return r
}
