package middleware

import (
	"bytes"
	"io"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RequestLog struct {
	RequestID      string        `json:"request_id"`
	Method         string        `json:"method"`
	Path           string        `json:"path"`
	Query          string        `json:"query,omitempty"`
	RemoteAddr     string        `json:"remote_addr"`
	UserAgent      string        `json:"user_agent"`
	RequestBody    string        `json:"request_body,omitempty"`
	ResponseStatus int           `json:"response_status"`
	ResponseBody   string        `json:"response_body,omitempty"`
	Duration       time.Duration `json:"duration_ms"`
	Timestamp      time.Time     `json:"timestamp"`
}

func RequestLoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := uuid.New().String()
		c.Set("requestID", requestID)

		startTime := time.Now()

		var requestBody string
		if c.Request.Body != nil {
			bodyBytes, _ := io.ReadAll(c.Request.Body)
			requestBody = string(bodyBytes)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		writer := &responseWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = writer

		c.Next()

		duration := time.Since(startTime)

		responseBody := writer.body.String()
		if len(responseBody) > 5000 {
			responseBody = responseBody[:5000] + "..."
		}
		if len(requestBody) > 5000 {
			requestBody = requestBody[:5000] + "..."
		}

		log.Printf("[%s] %s %s %s - Status: %d - Duration: %dms",
			requestID,
			c.Request.Method,
			c.Request.RequestURI,
			requestBody,
			c.Writer.Status(),
			duration.Milliseconds(),
		)
	}
}

type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
