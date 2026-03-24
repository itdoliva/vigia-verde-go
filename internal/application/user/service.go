package appUser

import (
	"context"
	User "vigia-verde-go/internal/domain/user"
	"vigia-verde-go/internal/infrastructure/security"
)

type UserService struct {
	repo User.Repository
}

func NewService(repo User.Repository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, dto RegisterReq) error {
	if err := dto.Validate(); err != nil {
		return err
	}

	hashed, err := security.HashedPassword(dto.Password)
	if err != nil {
		return err
	}

	user := &User.User{
		FullName: dto.FullName,
		Email:    dto.Email,
		Phone:    dto.Phone,
		PassHash: hashed,
		Status:   User.Status(dto.Status),
		Emoji:    dto.Emoji,
	}
	if err := s.repo.Register(ctx, user); err != nil {
		return err
	}
	return nil
}

func (s *UserService) GetByEmail(ctx context.Context, email string) (*User.User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (s *UserService) GetById(ctx context.Context, id string) (*User.User, error) {
	user, err := s.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetByPhone(ctx context.Context, phone string) (*User.User, error) {
	user, err := s.repo.GetByPhone(ctx, phone)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) Login(ctx context.Context, lr LoginReq) (string, error) {
	user, err := s.GetByPhone(ctx, lr.Phone)
	if err != nil {
		return "", err
	}

	if err := security.CheckPassword(lr.Password, user.PassHash); err != nil {
		return "", err
	}

	return "Usuario encontrado", nil
}
