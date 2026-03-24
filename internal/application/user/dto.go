package appUser

import (
	"errors"
)

type RegisterReq struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	Status   string `json:"status"`
	Emoji    string `json:"emoji"`
}

func (r *RegisterReq) Validate() error {
	if r.Email == "" && r.Phone == "" {
		return errors.New("Necessario informar email ou telefone")

	}
	return nil
}

type LoginReq struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

func (lr *LoginReq) Validate() error {
	if lr.Email == "" && lr.Phone == "" {
		return errors.New("Necessario informar email ou telefone")

	}
	return nil
}
