package schemas

import (
	"time"

	"golang-app/internal/models"
	"github.com/go-playground/validator/v10"
)

type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type" default:"bearer"`
}

type TokenPayload struct {
	Sub string `json:"sub"`
	Exp int64  `json:"exp"`
}

type UserBase struct {
	Email    string          `json:"email" validate:"required,email"`
	Role     models.UserRole `json:"role" default:"CONSUMER"`
	FullName *string         `json:"full_name,omitempty"`
	Address  *string         `json:"address,omitempty"`
}

type UserCreate struct {
	UserBase
	Password string `json:"password" validate:"required,min=8"`
}

type UserLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserOut struct {
	ID        uint            `json:"id"`
	Email     string          `json:"email"`
	Role      models.UserRole `json:"role"`
	FullName  *string         `json:"full_name,omitempty"`
	Address   *string         `json:"address,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
	IsActive  bool            `json:"is_active"`
}

type UserProfileOut struct {
	ID        uint         `json:"id"`
	FullName  string       `json:"full_name"`
	CreatedAt time.Time    `json:"created_at"`
	Products  []ProductOut `json:"products"`
}

var Validate *validator.Validate

func init() {
	Validate = validator.New()
}