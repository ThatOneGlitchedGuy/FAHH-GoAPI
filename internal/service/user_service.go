package service

import (
	"golang-app/internal/config"
	"golang-app/internal/models"
	"golang-app/internal/repository"
	"golang-app/internal/schemas"
	"golang-app/internal/security"
)

type UserService interface {
	GetUserByID(id uint) (*schemas.UserOut, error)
	UpdateUser(id uint, data *schemas.UserBase) (*schemas.UserOut, error)
	GetUserProfile(id uint) (*schemas.UserProfileOut, error)
}

type userService struct {
	repo   repository.UserRepository
	config *config.Settings
}

func NewUserService(repo repository.UserRepository, cfg *config.Settings) UserService {
	return &userService{repo: repo, config: cfg}
}

func (s *userService) userToUserOut(user *models.User) (*schemas.UserOut, error) {
	userOut := &schemas.UserOut{
		ID:        user.ID,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		IsActive:  user.IsActive,
	}

	if user.FullNameEncrypted != nil {
		decrypted, err := security.Decrypt(user.FullNameEncrypted, s.config.FernetKey, 0)
		if err != nil {
			return nil, err
		}
		fullName := string(decrypted)
		userOut.FullName = &fullName
	}

	if user.AddressEncrypted != nil {
		decrypted, err := security.Decrypt(user.AddressEncrypted, s.config.FernetKey, 0)
		if err != nil {
			return nil, err
		}
		address := string(decrypted)
		userOut.Address = &address
	}

	return userOut, nil
}

func (s *userService) GetUserByID(id uint) (*schemas.UserOut, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return s.userToUserOut(user)
}

func (s *userService) UpdateUser(id uint, data *schemas.UserBase) (*schemas.UserOut, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if data.FullName != nil {
		encrypted, err := security.Encrypt([]byte(*data.FullName), s.config.FernetKey)
		if err != nil {
			return nil, err
		}
		user.FullNameEncrypted = encrypted
	}

	if data.Address != nil {
		encrypted, err := security.Encrypt([]byte(*data.Address), s.config.FernetKey)
		if err != nil {
			return nil, err
		}
		user.AddressEncrypted = encrypted
	}

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	return s.userToUserOut(user)
}

func (s *userService) GetUserProfile(id uint) (*schemas.UserProfileOut, error) {
	user, err := s.repo.FindUserWithProducts(id)
	if err != nil {
		return nil, err
	}

	if user.Role != models.UserRoleAgent {
		// Or return a different profile schema for non-agents
		return nil, nil // Or an error indicating not an agent
	}

	var fullName string
	if user.FullNameEncrypted != nil {
		decrypted, err := security.Decrypt(user.FullNameEncrypted, s.config.FernetKey, 0)
		if err == nil {
			fullName = string(decrypted)
		}
	}

	productsOut := make([]schemas.ProductOut, len(user.Products))
	for i, p := range user.Products {
		desc := p.Description
		productsOut[i] = schemas.ProductOut{
			ID:          p.ID,
			AgentID:     p.AgentID,
			Title:       p.Title,
			Description: &desc,
			Price:       p.Price,
			IsActive:    p.IsActive,
			CreatedAt:   p.CreatedAt,
		}
	}

	return &schemas.UserProfileOut{
		ID:        user.ID,
		FullName:  fullName,
		CreatedAt: user.CreatedAt,
		Products:  productsOut,
	}, nil
}
