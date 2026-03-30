package appUser

import (
	"context"
	"strconv"
	"time"
	User "vigia-verde-go/internal/domain/user"
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

	isVerified, err := strconv.ParseBool(dto.IsVerified)
	if err != nil {
		return err
	}

	user := &User.User{
		Id:         dto.UID,
		FullName:   dto.FullName,
		Email:      dto.Email,
		Phone:      dto.Phone,
		IsVerified: isVerified,
		Emoji:      dto.Emoji,
		CreateAt:   time.Now(),
	}
	if err := s.repo.Save(ctx, user); err != nil {
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
