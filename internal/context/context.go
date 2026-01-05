package context

import (
	"github.com/gin-gonic/gin"
)

type ContextKey string

const (
	UserIDKey    ContextKey = "userID"
	RequestIDKey ContextKey = "requestID"
	RoleKey      ContextKey = "role"
	EmailKey     ContextKey = "email"
	IPAddressKey ContextKey = "ipAddress"
)

type RequestContext struct {
	UserID    uint
	RequestID string
	Role      string
	Email     string
	IPAddress string
}

func NewRequestContext(c *gin.Context) *RequestContext {
	rc := &RequestContext{
		IPAddress: c.ClientIP(),
	}

	if userID, exists := c.Get(string(UserIDKey)); exists {
		if id, ok := userID.(uint); ok {
			rc.UserID = id
		}
	}

	if requestID, exists := c.Get(string(RequestIDKey)); exists {
		if id, ok := requestID.(string); ok {
			rc.RequestID = id
		}
	}

	if role, exists := c.Get(string(RoleKey)); exists {
		if r, ok := role.(string); ok {
			rc.Role = r
		}
	}

	if email, exists := c.Get(string(EmailKey)); exists {
		if e, ok := email.(string); ok {
			rc.Email = e
		}
	}

	return rc
}

func SetUserID(c *gin.Context, userID uint) {
	c.Set(string(UserIDKey), userID)
}

func SetRequestID(c *gin.Context, requestID string) {
	c.Set(string(RequestIDKey), requestID)
}

func SetRole(c *gin.Context, role string) {
	c.Set(string(RoleKey), role)
}

func SetEmail(c *gin.Context, email string) {
	c.Set(string(EmailKey), email)
}

func GetUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get(string(UserIDKey))
	if !exists {
		return 0, false
	}
	id, ok := userID.(uint)
	return id, ok
}

func GetRequestID(c *gin.Context) string {
	requestID, exists := c.Get(string(RequestIDKey))
	if !exists {
		return ""
	}
	id, ok := requestID.(string)
	if !ok {
		return ""
	}
	return id
}

func GetRole(c *gin.Context) string {
	role, exists := c.Get(string(RoleKey))
	if !exists {
		return ""
	}
	r, ok := role.(string)
	if !ok {
		return ""
	}
	return r
}

func GetEmail(c *gin.Context) string {
	email, exists := c.Get(string(EmailKey))
	if !exists {
		return ""
	}
	e, ok := email.(string)
	if !ok {
		return ""
	}
	return e
}

func GetIPAddress(c *gin.Context) string {
	return c.ClientIP()
}
