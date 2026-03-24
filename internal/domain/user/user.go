package User

import (
	"context"
	"errors"
	"time"
)

type Repository interface {
	Register(ctx context.Context, User *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByPhone(ctx context.Context, phone string) (*User, error)
	GetById(ctx context.Context, id string) (*User, error)
}

type User struct {
	Id                  string       `firestore:"id"`
	FullName            string       `firestore:"full_name"`
	Email               string       `firestore:"email"`
	Phone               string       `firestore:"phone"`
	PassHash            string       `firestore:"pass_hash"`
	Status              Status       `firestore:"status"`
	Emoji               string       `firestore:"emoji"`
	VerificationMethods verification `firestore:"verification_methods"`
	CreateAt            time.Time    `firestore:"create_at"`
}

type Status string

const (
	active   Status = "active"
	pending  Status = "pending"
	disabled Status = "disabled"
)

func (s Status) Validate() error {
	switch s {
	case active, pending, disabled:
		return nil
	default:
		return errors.New("invalid status")
	}
}

type verification struct {
	value      string    `firestore:"value"`
	verified   bool      `firestore:"verified"`
	verifiedAt time.Time `firestore:"verified_at"`
}

func (v verification) Validate() error {
	if v.value != "email" && v.value != "wpp" {
		return errors.New("invalid verification value")
	}
	return nil
}
