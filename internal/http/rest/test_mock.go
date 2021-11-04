package rest

import "github.com/stretchr/testify/mock"

type MockProxyService struct {
	mock.Mock
}

func (m *MockProxyService) VatValidator() (string, error) {
	args := m.Called()
	return "", args.Error(1)
}
