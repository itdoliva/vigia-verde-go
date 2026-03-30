package User

import (
	"context"
	"time"
)

type Repository interface {
	Save(ctx context.Context, User *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByPhone(ctx context.Context, phone string) (*User, error)
	GetById(ctx context.Context, id string) (*User, error)
}

type User struct {
	Id         string    `firestore:"id"`
	FullName   string    `firestore:"full_name"`
	Email      string    `firestore:"email"`
	Phone      string    `firestore:"phone"`
	Emoji      string    `firestore:"emoji"`
	IsVerified bool      `firestore:"is_verified"`
	CreateAt   time.Time `firestore:"create_at"`
}
