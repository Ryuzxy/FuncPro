package komoditas

import (
	"context"
	"fmt"

	"github.com/ryuzxy/FuncPro/pkg/fx"
)

type Service interface {
	GetAllKomoditas(ctx context.Context) fx.Result[[]Komoditas]
	GetKomoditasByID(ctx context.Context, id uint) fx.Result[*Komoditas]
	CreateKomoditas(ctx context.Context, req CreateKomoditasRequest) fx.Result[*Komoditas]
	UpdateKomoditas(ctx context.Context, id uint, req UpdateKomoditasRequest) fx.Result[*Komoditas]
	DeleteKomoditas(ctx context.Context, id uint) fx.Result[bool]
	GetKomoditasWithStats(ctx context.Context, id uint) fx.Result[KomoditasWithStats]
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetAllKomoditas(ctx context.Context) fx.Result[[]Komoditas] {
	return s.repo.GetAll(ctx)
}

func (s *service) GetKomoditasByID(ctx context.Context, id uint) fx.Result[*Komoditas] {
	return s.repo.GetByID(ctx, id)
}

func (s *service) CreateKomoditas(ctx context.Context, req CreateKomoditasRequest) fx.Result[*Komoditas] {
	kom := &Komoditas{
		Name: req.Name,
		Type: req.Type,
	}
	return s.repo.Create(ctx, kom)
}

func (s *service) UpdateKomoditas(ctx context.Context, id uint, req UpdateKomoditasRequest) fx.Result[*Komoditas] {

	existing, err := s.repo.GetByID(ctx, id).Unwrap()
	if err != nil {
		return fx.Err[*Komoditas](fmt.Errorf("komoditas not found: %w", err))
	}

	if req.Name != "" {
		existing.Name = req.Name
	}
	if req.Type != "" {
		existing.Type = req.Type
	}

	return s.repo.Update(ctx, id, existing)
}

func (s *service) DeleteKomoditas(ctx context.Context, id uint) fx.Result[bool] {
	return s.repo.Delete(ctx, id)
}

func (s *service) GetKomoditasWithStats(ctx context.Context, id uint) fx.Result[KomoditasWithStats] {
	kom, err := s.repo.GetByID(ctx, id).Unwrap()
	if err != nil {
		return fx.Err[KomoditasWithStats](fmt.Errorf("komoditas not found: %w", err))
	}

	// Placeholder sementara sampai tabel price tersedia
	stats := PriceStats{
		Average: 0,
		Min:     0,
		Max:     0,
		Count:   0,
		Trend:   "stable",
	}

	return fx.Ok(KomoditasWithStats{
		Komoditas: *kom,
		Stats:     stats,
	})
}
