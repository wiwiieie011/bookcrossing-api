package mocks

import (
	"github.com/dasler-fw/bookcrossing/internal/dto"
	"github.com/dasler-fw/bookcrossing/internal/models"
	"github.com/stretchr/testify/mock"
)

type ExchangeServiceMock struct {
	mock.Mock
}

func (m *ExchangeRepositoryMock) CreateExchangeService(req *dto.CreateExchangeRequest, actingUserID uint) (*models.Exchange, error) {
	args := m.Called(req, actingUserID)
	var exc *models.Exchange
	if args.Get(0) != nil {
		exc = args.Get(0).(*models.Exchange)
	}

	return exc, args.Error(1)
}

func (m *ExchangeRepositoryMock) AcceptExchangeService(exchangeID uint, actingUserID uint) error {
	args := m.Called(exchangeID, actingUserID)
	return args.Error(0)
}

func (m *ExchangeRepositoryMock) CompleteExchangeService(exchangeID uint, actingUserID uint) error {
	args := m.Called(exchangeID, actingUserID)
	return args.Error(0)
}

func (m *ExchangeRepositoryMock) CancelExchangeService(exchangeID uint, actingUserID uint) error {
	args := m.Called(exchangeID, actingUserID)
	return args.Error(0)
}

func (m *ExchangeRepositoryMock) GetByIDService(exchangeID uint) (*models.Exchange, error) {
	args := m.Called(exchangeID)

	var exc *models.Exchange
	if args.Get(0) != nil {
		exc = args.Get(0).(*models.Exchange)
	}

	return exc, args.Error(1)
}

func (m *ExchangeRepositoryMock) GetAllService() ([]models.Exchange, error) {
	args := m.Called()
	var exc []models.Exchange
	if args.Get(0) != nil {
		exc = args.Get(0).([]models.Exchange)
	}

	return exc, args.Error(1)
}
