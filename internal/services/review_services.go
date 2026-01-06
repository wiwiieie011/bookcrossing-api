package services

import (
	"errors"
	"strings"

	"github.com/dasler-fw/bookcrossing/internal/dto"
	"github.com/dasler-fw/bookcrossing/internal/models"
	"github.com/dasler-fw/bookcrossing/internal/repository"
)

type ReviewService interface {
	Create(authorID uint, req dto.CreateReviewRequest) (*models.Review, error)
	GetByUserID(userID uint) ([]models.Review, error)
	GetByBookID(bookID uint) ([]models.Review, error)
	Delete(reviewID uint, authorID uint) error
}

type reviewService struct {
	repo repository.ReviewRepository
}

func NewReviewService(repo repository.ReviewRepository) ReviewService {
	return &reviewService{repo: repo}
}

func (s *reviewService) Create(authorID uint, req dto.CreateReviewRequest) (*models.Review, error) {
	trimmedText := strings.TrimSpace(req.Text)

	length := len([]rune(trimmedText))
	if length < 10 || length > 150 {
		return nil, dto.ErrReviewTextLength
	}

	if strings.TrimSpace(req.Text) == "" {
		return nil,errors.New("review text is request")
	}

	if req.Rating < 1 || req.Rating > 5 {
		return nil, dto.ErrInvalidRating
	}

	if req.TargetUserID == authorID {
		return nil,dto.ErrSelfReviewForbidden
	}

	review := &models.Review{
		AuthorID:     authorID,
		TargetUserID: req.TargetUserID,
		TargetBookID: req.TargetBookID,
		Text:         req.Text,
		Rating:       req.Rating,
	}

	if err:= s.repo.Create(review); err!= nil {
		return  nil, err
	}

	return  review, nil
}

func (s *reviewService) GetByUserID(userID uint) ([]models.Review, error) {
	return s.repo.GetByTargetUserID(userID)
}

func (s *reviewService) GetByBookID(bookID uint) ([]models.Review, error) {
	return s.repo.GetByTargetBookID(bookID)
}

func (s *reviewService) Delete(reviewID uint, authorID uint) error {
	review, err := s.repo.GetByID(reviewID)
	if err != nil {
		return err
	}

	if review.AuthorID != authorID {
		return dto.ErrReviewDeleteForbidden
	}
	return s.repo.Delete(reviewID)
}
