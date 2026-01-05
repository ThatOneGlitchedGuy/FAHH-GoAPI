package audit

import (
	"log"
	"sync"
	"time"
)

type AuditAction string
type ResourceType string

const (
	ActionCreate AuditAction = "CREATE"
	ActionUpdate AuditAction = "UPDATE"
	ActionDelete AuditAction = "DELETE"
	ActionRead   AuditAction = "READ"
	ActionLogin  AuditAction = "LOGIN"
	ActionLogout AuditAction = "LOGOUT"
)

const (
	ResourceUser    ResourceType = "USER"
	ResourceProduct ResourceType = "PRODUCT"
	ResourceOrder   ResourceType = "ORDER"
	ResourceReview  ResourceType = "REVIEW"
)

type AuditLog struct {
	ID           uint          `json:"id"`
	UserID       uint          `json:"user_id"`
	Action       AuditAction   `json:"action"`
	ResourceType ResourceType  `json:"resource_type"`
	ResourceID   uint          `json:"resource_id"`
	OldValues    string        `json:"old_values,omitempty"`
	NewValues    string        `json:"new_values,omitempty"`
	IPAddress    string        `json:"ip_address"`
	UserAgent    string        `json:"user_agent,omitempty"`
	Changes      string        `json:"changes,omitempty"`
	Status       string        `json:"status"`
	ErrorMsg     string        `json:"error_msg,omitempty"`
	Timestamp    time.Time     `json:"timestamp"`
}

type AuditLogger struct {
	logs chan *AuditLog
	mu   sync.RWMutex
	stop chan bool
}

func NewAuditLogger(bufferSize int) *AuditLogger {
	al := &AuditLogger{
		logs: make(chan *AuditLog, bufferSize),
		stop: make(chan bool),
	}
	go al.processlogs()
	return al
}

func (al *AuditLogger) Log(userID uint, action AuditAction, resourceType ResourceType, resourceID uint, ipAddress string) *AuditLog {
	return &AuditLog{
		UserID:       userID,
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		IPAddress:    ipAddress,
		Status:       "SUCCESS",
		Timestamp:    time.Now(),
	}
}

func (al *AuditLogger) LogWithChanges(userID uint, action AuditAction, resourceType ResourceType, resourceID uint, ipAddress string, oldValues, newValues string) *AuditLog {
	auditLog := al.Log(userID, action, resourceType, resourceID, ipAddress)
	auditLog.OldValues = oldValues
	auditLog.NewValues = newValues
	return auditLog
}

func (al *AuditLogger) LogError(userID uint, action AuditAction, resourceType ResourceType, resourceID uint, ipAddress string, err string) *AuditLog {
	auditLog := al.Log(userID, action, resourceType, resourceID, ipAddress)
	auditLog.Status = "FAILED"
	auditLog.ErrorMsg = err
	return auditLog
}

func (al *AuditLogger) Submit(log *AuditLog) {
	select {
	case al.logs <- log:
	case <-al.stop:
		return
	default:
		log.Printf("Audit log channel full, dropping log for user %d", log.UserID)
	}
}

func (al *AuditLogger) processlogs() {
	for {
		select {
		case log := <-al.logs:
			al.persistLog(log)
		case <-al.stop:
			return
		}
	}
}

func (al *AuditLogger) persistLog(auditLog *AuditLog) {
	log.Printf("[AUDIT] UserID: %d | Action: %s | Resource: %s#%d | Status: %s | IP: %s | Time: %s",
		auditLog.UserID,
		auditLog.Action,
		auditLog.ResourceType,
		auditLog.ResourceID,
		auditLog.Status,
		auditLog.IPAddress,
		auditLog.Timestamp.Format(time.RFC3339),
	)
}

func (al *AuditLogger) Stop() {
	close(al.stop)
}
