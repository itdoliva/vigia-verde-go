package repoUser

import (
	"context"
	"fmt"
	User "vigia-verde-go/internal/domain/user"

	"cloud.google.com/go/firestore"
)

type UserRepository struct {
	client *firestore.Client
}

func NewRepository(client *firestore.Client) User.Repository {
	return &UserRepository{client: client}
}

func (r *UserRepository) Save(ctx context.Context, User *User.User) error {
	docRef := r.client.Collection("treeUser").Doc(User.Id)

	if _, err := docRef.Set(ctx, User); err != nil {
		return fmt.Errorf("falha ao persistir usuário no Firestore: %v", err)
	}

	return nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*User.User, error) {
	docs, err := r.client.
		Collection("treeUser").
		Where("email", "==", email).
		Limit(1).
		Documents(ctx).
		GetAll()
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar usuário por email: %w", err)
	}

	if len(docs) == 0 {
		return nil, fmt.Errorf("usuário não existe")
	}

	var user User.User
	if err := docs[0].DataTo(&user); err != nil {
		return nil, fmt.Errorf("erro ao converter documento: %w", err)
	}

	return &user, nil
}
func (r *UserRepository) GetByPhone(ctx context.Context, phone string) (*User.User, error) {
	docs, err := r.client.
		Collection("treeUser").
		Where("phone", "==", phone).
		Limit(1).
		Documents(ctx).
		GetAll()
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar usuário por telefone: %w", err)
	}

	if len(docs) == 0 {
		return nil, fmt.Errorf("usuário não existe")
	}

	var user User.User
	if err := docs[0].DataTo(&user); err != nil {
		return nil, fmt.Errorf("erro ao converter documento: %w", err)
	}

	return &user, nil
}
func (r *UserRepository) GetById(ctx context.Context, id string) (*User.User, error) {
	doc, err := r.client.Collection("treeUser").Doc(id).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("Usuario não existe: %v", err)
	}
	var user User.User
	if err := doc.DataTo(&user); err != nil {
		return nil, fmt.Errorf("erro ao converter documento: %v", err)
	}
	return &user, nil
}
