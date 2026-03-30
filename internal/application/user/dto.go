package appUser

import (
	"errors"
	"fmt"
)

type RegisterReq struct {
	UID        string `json:"uid"`
	FullName   string `json:"full_name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	IsVerified string `json:"is_verified"`
	Emoji      string `json:"emoji"`
}

func (r *RegisterReq) Validate() error {
	if r.UID == "" {
		return fmt.Errorf("o ID do usuário (UID do Firebase) é obrigatório")
	}
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
