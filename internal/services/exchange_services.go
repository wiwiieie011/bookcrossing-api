package services

import (
	"errors"
	"log/slog"

	"github.com/dasler-fw/bookcrossing/internal/dto"
	"github.com/dasler-fw/bookcrossing/internal/models"
	"github.com/dasler-fw/bookcrossing/internal/repository"
)

type ExchangeService interface {
	CreateExchange(req *dto.CreateExchangeRequest, actingUserID uint) (*models.Exchange, error)
	AcceptExchange(exchangeID uint, actingUserID uint) error
	CompleteExchange(exchangeID uint, actingUserID uint) error
	CancelExchange(exchangeID uint, actingUserID uint) error
	GetByID(exchangeID uint) (*models.Exchange, error)
	GetAll() ([]models.Exchange, error)
}

type exchangeService struct {
	exchangeRepo repository.ExchangeRepository
	bookRepo     repository.BookRepository
	log          *slog.Logger
}

func NewExchangeService(exchangeRepo repository.ExchangeRepository, bookRepo repository.BookRepository, log *slog.Logger) ExchangeService {
	return &exchangeService{exchangeRepo: exchangeRepo, bookRepo: bookRepo, log: log}
}

func (s *exchangeService) CancelExchange(exchangeID uint, actingUserID uint) error {
	if exchangeID == 0 {
		s.log.Error("error in CancelExchange function exchange_services.go")
		return dto.ErrExchangeInvalidID
	}

	exchange, err := s.exchangeRepo.GetByID(exchangeID)
	if err != nil {
		s.log.Error("error in CancelExchange function exchange_services.go", "error", err)
		return err
	}

	// todo: check if the current user is the initiator, after the auth handler is implemented

	if exchange.Status != "pending" {
		s.log.Error("error in CancelExchange function exchange_services.go", "error", errors.New("exchange is not pending"))
		return dto.ErrExchangeNotPending
	}

	// Только инициатор может отменять pending обмен
	if exchange.InitiatorID != actingUserID {
		return errors.New("forbidden")
	}

	return s.exchangeRepo.CancelExchange(exchange)
}

func (s *exchangeService) CompleteExchange(exchangeID uint, actingUserID uint) error {
	if exchangeID == 0 {
		s.log.Error("error in CompleteExchange function exchange_services.go")
		return dto.ErrExchangeInvalidID
	}

	exchange, err := s.exchangeRepo.GetByID(exchangeID)
	if err != nil {
		s.log.Error("error in CompleteExchange function exchange_services.go", "error", err)
		return err
	}

	// todo: check if the current user is the initiator or the recipient, after the auth handler is implemented

	if exchange.Status != "accepted" {
		s.log.Error("error in CompleteExchange function exchange_services.go", "error", errors.New("exchange is not accepted"))
		return dto.ErrExchangeNotAccepted
	}

	// Завершить может любая сторона обмена (инициатор или получатель)
	if actingUserID != exchange.InitiatorID && actingUserID != exchange.RecipientID {
		return errors.New("forbidden")
	}

	return s.exchangeRepo.CompleteExchange(exchange)
}

func (s *exchangeService) AcceptExchange(exchangeID uint, actingUserID uint) error {
	if exchangeID == 0 {
		s.log.Error("error in AcceptExchange function exchange_services.go")
		return dto.ErrExchangeInvalidID
	}

	exchange, err := s.exchangeRepo.GetByID(exchangeID)
	if err != nil {
		s.log.Error("error in AcceptExchange function exchange_services.go", "error", err)
		return err
	}

	// todo: check if the current user is the recipient, after the auth handler is implemented

	if exchange.Status != "pending" {
		s.log.Error("error in AcceptExchange function exchange_services.go", "error", errors.New("exchange is not pending"))
		return dto.ErrExchangeNotPending
	}

	// Принять может только получатель
	if exchange.RecipientID != actingUserID {
		return errors.New("forbidden")
	}

	exchange.Status = "accepted"
	return s.exchangeRepo.Update(exchange)
}

func (s *exchangeService) CreateExchange(req *dto.CreateExchangeRequest, actingUserID uint) (*models.Exchange, error) {
	if req == nil {
		s.log.Error("error in CreateExchange function exchange_services.go")
		return nil, dto.ErrExchangeInvalidID
	}

	// Инициатором может быть только текущий пользователь
	exchange := &models.Exchange{
		InitiatorID:     actingUserID,
		RecipientID:     req.RecipientID,
		InitiatorBookID: req.InitiatorBookID,
		RecipientBookID: req.RecipientBookID,
		Status:          "pending",
	}

	initiatorBook, err := s.bookRepo.GetByID(req.InitiatorBookID)
	if err != nil {
		s.log.Error("error in CreateExchange function exchange_services.go", "error", err)
		return nil, err
	}

	recipientBook, err := s.bookRepo.GetByID(req.RecipientBookID)
	if err != nil {
		s.log.Error("error in CreateExchange function exchange_services.go", "error", err)
		return nil, err
	}

	if err := s.CheckIsTheSameUser(initiatorBook.UserID, recipientBook.UserID); err != nil {
		s.log.Error("error in CreateExchange function exchange_services.go", "error", err)
		return nil, err
	}

	if err := s.CheckInitiatorOwnsBook(actingUserID, initiatorBook); err != nil {
		s.log.Error("error in CreateExchange function exchange_services.go", "error", err)
		return nil, err
	}

	if err := s.CheckRecipientOwnsBook(req.RecipientID, recipientBook); err != nil {
		s.log.Error("error in CreateExchange function exchange_services.go", "error", err)
		return nil, err
	}

	if err := s.CheckIsAvailable(initiatorBook, recipientBook); err != nil {
		s.log.Error("error in CreateExchange function exchange_services.go", "error", err)
		return nil, err
	}
	if err := s.exchangeRepo.CreateExchange(exchange); err != nil {
		return nil, err
	}

	return exchange, nil
}

func (s *exchangeService) CheckInitiatorOwnsBook(initiatorID uint, initiatorBook *models.Book) error {
	if initiatorID != initiatorBook.UserID {
		s.log.Error("error in CreateExchange function exchange_services.go", "error", errors.New("initiator does not own the book"))
		return dto.ErrInitiatorNotOwner
	}

	return nil
}

func (s *exchangeService) CheckRecipientOwnsBook(recipientID uint, recipientBook *models.Book) error {
	if recipientID != recipientBook.UserID {
		s.log.Error("error in CreateExchange function exchange_services.go", "error", errors.New("recipient does not own the book"))
		return dto.ErrRecipientNotOwner
	}

	return nil
}
func (s *exchangeService) CheckIsTheSameUser(initiatorID uint, recipientID uint) error {
	if initiatorID == recipientID {
		s.log.Error("error in CreateExchange function exchange_services.go", "error", errors.New("initiator and recipient book cannot be the same user"))
		return errors.New("initiator and recipient book cannot be the same user")
	}

	return nil
}

func (s *exchangeService) CheckIsAvailable(initiatorBook *models.Book, recipientBook *models.Book) error {
	if initiatorBook.Status != "available" {
		s.log.Error("error in CreateExchange function exchange_services.go", "error", errors.New("initiator book is unavailable"))
		return dto.ErrUnavailable
	}

	if recipientBook.Status != "available" {
		s.log.Error("error in CreateExchange function exchange_services.go", "error", errors.New("recipient book is unavailable"))
		return dto.ErrRUnavailable
	}

	return nil
}

func (s *exchangeService) GetByID(exchangeID uint) (*models.Exchange, error) {
	return s.exchangeRepo.GetByID(exchangeID)
}

func (s *exchangeService) GetAll() ([]models.Exchange, error) {
	return s.exchangeRepo.GetAll()
}
