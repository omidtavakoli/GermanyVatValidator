package validator

import "github.com/stretchr/testify/mock"

type MockProxyService struct {
	mock.Mock
}

func (m *MockProxyService) VatValidator(vatNum string) (bool, error) {
	args := m.Called(vatNum)
	return true, args.Error(1)
}
