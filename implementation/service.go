package implementation

import (
	"context"
	"strings"

	"pharmacy/medicine"
)

type Repository interface {
	Create(medicine.Medicine) (string, error)
	GetByID(string) (medicine.Medicine, error)
	Sell(string, int) error
}

type Service interface {
	Create(ctx context.Context, medicine medicine.Medicine) (string, error)
	GetByID(ctx context.Context, id string) (medicine.Medicine, error)
	Sell(ctx context.Context, id string, count int) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(_ context.Context, m medicine.Medicine) (string, error) {
	if strings.TrimSpace(m.Name) == "" ||
		strings.TrimSpace(m.Manufacturer) == "" ||
		m.Price < 0 || m.Quantity < 0 {
		return "", medicine.ErrInvalidMedicine
	}
	return s.repo.Create(m)
}

func (s *service) GetByID(_ context.Context, id string) (medicine.Medicine, error) {
	if strings.TrimSpace(id) == "" {
		return medicine.Medicine{}, medicine.ErrNotFound
	}
	return s.repo.GetByID(id)
}

func (s *service) Sell(_ context.Context, id string, count int) error {
	if strings.TrimSpace(id) == "" {
		return medicine.ErrNotFound
	}
	if count <= 0 {
		return medicine.ErrInvalidCount
	}
	return s.repo.Sell(id, count)
}
