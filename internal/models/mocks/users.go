package mocks

import "snippet.darieldejesus.com/internal/models"

type UserModel struct {
	Err error
}

func (m *UserModel) Insert(name, email, password string) error {
	if m.Err != nil {
		return m.Err
	}
	switch email {
	case "user@darieldejesus.com":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	if m.Err != nil {
		return 0, m.Err
	}
	if email == "user@darieldejesus.com" && password == "Pa$$word123" {
		return 7, nil
	}

	return 0, models.ErrInvalidCredentials
}

func (m *UserModel) Exists(id int) (bool, error) {
	if m.Err != nil {
		return false, m.Err
	}
	switch id {
	case 7:
		return true, nil
	default:
		return false, nil
	}
}
