package mocks

import "github.com/sgoldenf/a-place-for-your-thoughts/internal/models"

type UserModel struct{}

func (m *UserModel) CreateUser(name, email, password string) error {
	switch email {
	case "duplicate@example.com":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	if email == "login@example.com" && password == "password" {
		return 1, nil
	} else {
		return 0, models.ErrInvalidCredentials
	}
}

func (m *UserModel) Exists(id int) (bool, error) {
	switch id {
	case 1:
		return true, nil
	default:
		return false, nil
	}
}
