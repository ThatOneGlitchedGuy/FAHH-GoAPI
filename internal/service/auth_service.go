package service

import (
	"errors"
	"golang-app/internal/config"
	"golang-app/internal/models"
	"golang-app/internal/repository"
	"golang-app/internal/schemas"
	"golang-app/internal/security"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	RegisterUser(input *schemas.UserCreate) (*models.User, error)
	LoginUser(input *schemas.UserLogin) (*models.User, error)
	GenerateJWT(user *models.User) (string, error)
}

type authService struct {
	userRepo repository.UserRepository
	config   *config.Settings
}

func NewAuthService(userRepo repository.UserRepository, cfg *config.Settings) AuthService {
	return &authService{userRepo: userRepo, config: cfg}
}

func (s *authService) RegisterUser(input *schemas.UserCreate) (*models.User, error) {
	_, err := s.userRepo.FindByEmail(input.Email)
	if err == nil {
		return nil, errors.New("user with this email already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:          input.Email,
		HashedPassword: string(hashedPassword),
		Role:           input.Role,
		IsActive:       true,
	}

	if input.FullName != nil {
		encrypted, err := security.Encrypt([]byte(*input.FullName), s.config.FernetKey)
		if err != nil {
			return nil, err
		}
		user.FullNameEncrypted = encrypted
	}
	if input.Address != nil {
		encrypted, err := security.Encrypt([]byte(*input.Address), s.config.FernetKey)
		if err != nil {
			return nil, err
		}
		user.AddressEncrypted = encrypted
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) LoginUser(input *schemas.UserLogin) (*models.User, error) {
	user, err := s.userRepo.FindByEmail(input.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(input.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func (s *authService) GenerateJWT(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":  user.ID,
		"role": user.Role,
		"exp":  time.Now().Add(time.Minute * time.Duration(s.config.AccessTokenExpireMinutes)).Unix(),
		"iat":  time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.SecretKey))
}
