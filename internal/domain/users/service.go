package users

import (
	"SpotSync/internal/auth"
	"SpotSync/internal/domain/users/dto"
	"errors"
	"strings"
)

type service struct {
	repo       Repository
	jwtService auth.JWTService
}

func NewService(repo Repository, jwtService auth.JWTService) *service {
	return &service{repo: repo, jwtService: jwtService}
}

func (s *service) CreateUser(req dto.RegisterRequest) (*dto.RegisterResponse, error) {
	// Normalize role: default to DRIVER
	role := UserRole(strings.ToUpper(req.Role))
	if role != RoleAdmin && role != RoleDriver {
		role = RoleDriver 
	}

	user := User{
		Email: req.Email,
		Name:  req.Name,
		Phone: req.Phone,
		Role:  role,      
	}

	if err := user.HashPassword(req.Password); err != nil {
		return nil, err
	}

	if err := s.repo.CreateUser(&user); err != nil {
		return nil, err
	}

	return &dto.RegisterResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone,
		Role:      string(user.Role),
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),	
	}, nil
}

func (s *service) LoginUser(req dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	// Fixed: was "return nil, err" (err is nil here, hides the real problem)
	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	if !user.CheckPassword(req.Password) {
		return nil, errors.New("invalid email or password")
	}

	token, err := s.jwtService.GenerateToken(user.ID, user.Name, user.Email, string(user.Role), user.Phone, user.CreatedAt.String())
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token: token,
		User: dto.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  string(user.Role),
		},
	}, nil
}
		