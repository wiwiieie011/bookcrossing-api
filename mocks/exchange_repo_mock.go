package mocks

import (
	"github.com/dasler-fw/bookcrossing/internal/models"
	"github.com/stretchr/testify/mock"
)

type ExchangeRepositoryMock struct {
	mock.Mock
}

func (m *ExchangeRepositoryMock) CreateExchange(req *models.Exchange) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *ExchangeRepositoryMock) CompleteExchange(req *models.Exchange) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *ExchangeRepositoryMock) CancelExchange(req *models.Exchange) error {
	args := m.Called(req)
	return args.Error(0)
}
func (m *ExchangeRepositoryMock) Update(req *models.Exchange) error {
	args := m.Called(req)
	return args.Error(0)
}
func (m *ExchangeRepositoryMock) GetByID(id uint) (*models.Exchange, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Exchange), args.Error(1)
}

func (m *ExchangeRepositoryMock) GetAll() ([]models.Exchange, error) {
	args := m.Called()

	var exchs []models.Exchange
	if args.Get(0) != nil {
		exchs = args.Get(0).([]models.Exchange)
	}

	return exchs, args.Error(1)
}