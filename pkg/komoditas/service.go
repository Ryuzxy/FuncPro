package Komoditas

import (
    "context"
    "fmt"
    
    "github.com/ryuzxy/fucpro/pkg/fx"
)

type Service interface {
    GetAllKomoditas(ctx context.Context) fx.Result[[]Komoditas]
    GetKomoditasByID(ctx context.Context, id uint) fx.Result[Komoditas]
    CreateKomoditas(ctx context.Context, req CreateKomoditasRequest) fx.Result[Komoditas]
    UpdateKomoditas(ctx context.Context, id uint, req UpdateKomoditasRequest) fx.Result[Komoditas]
    DeleteKomoditas(ctx context.Context, id uint) fx.Result[bool]
    GetKomoditasWithStats(ctx context.Context, id uint) fx.Result[KomoditasWithStats]
    BulkCreateKomoditas(ctx context.Context, reqs []CreateKomoditasRequest) fx.Result[[]Komoditas]
}

type service struct {
    repo      Repository
    priceRepo PriceRepository
}

func NewService(repo Repository, priceRepo PriceRepository) Service {
    return &service{repo: repo, priceRepo: priceRepo}
}

// Pure validation function
func validateCreateRequest(req CreateKomoditasRequest) fx.Result[CreateKomoditasRequest] {
    if req.Name == "" {
        return fx.Err[CreateKomoditasRequest](fmt.Errorf("name is required"))
    }
    if req.Type == "" {
        return fx.Err[CreateKomoditasRequest](fmt.Errorf("type is required"))
    }
    if len(req.Name) > 100 {
        return fx.Err[CreateKomoditasRequest](fmt.Errorf("name too long"))
    }
    return fx.Ok(req)
}

// Pure transformation function  
func requestToKomoditas(req CreateKomoditasRequest) Komoditas {
    return Komoditas{
        Name: req.Name,
        Type: req.Type,
    }
}

func (s *service) GetAllKomoditas(ctx context.Context) fx.Result[[]Komoditas] {
    return s.repo.GetAll(ctx)
}

func (s *service) GetKomoditasByID(ctx context.Context, id uint) fx.Result[Komoditas] {
    return s.repo.GetByID(ctx, id)
}

func (s *service) CreateKomoditas(ctx context.Context, req CreateKomoditasRequest) fx.Result[Komoditas] {
    return validateCreateRequest(req).
        Map(requestToKomoditas).
        AndThen(func(k Komoditas) fx.Result[Komoditas] {
            return s.repo.Create(ctx, k)
        })
}

func (s *service) UpdateKomoditas(ctx context.Context, id uint, req UpdateKomoditasRequest) fx.Result[Komoditas] {
    // First get existing komoditas
    return s.repo.GetByID(ctx, id).
        AndThen(func(existing Komoditas) fx.Result[Komoditas] {
            // Pure update function
            updated := applyUpdate(existing, req)
            return s.repo.Update(ctx, id, updated)
        })
}

// Pure function to apply updates
func applyUpdate(existing Komoditas, req UpdateKomoditasRequest) Komoditas {
    if req.Name != "" {
        existing.Name = req.Name
    }
    if req.Type != "" {
        existing.Type = req.Type
    }
    return existing
}

func (s *service) DeleteKomoditas(ctx context.Context, id uint) fx.Result[bool] {
    return s.repo.Delete(ctx, id)
}

func (s *service) GetKomoditasWithStats(ctx context.Context, id uint) fx.Result[KomoditasWithStats] {
    komoditasResult := s.repo.GetByID(ctx, id)
    
    return komoditasResult.AndThen(func(k Komoditas) fx.Result[KomoditasWithStats] {
        pricesResult := s.priceRepo.GetByKomoditasID(ctx, k.ID)
        
        return pricesResult.Map(func(prices []Price) KomoditasWithStats {
            stats := calculatePriceStats(prices)
            
            return KomoditasWithStats{
                Komoditas: k,
                Stats:     stats,
            }
        })
    })
}

func (s *service) BulkCreateKomoditas(ctx context.Context, reqs []CreateKomoditasRequest) fx.Result[[]Komoditas] {
    // Validate all requests first
    validatedReqs := make([]fx.Result[Komoditas], 0, len(reqs))
    
    for _, req := range reqs {
        result := validateCreateRequest(req).
            Map(requestToKomoditas)
        validatedReqs = append(validatedReqs, result)
    }
    
    // Check for any validation errors
    for _, result := range validatedReqs {
        if result.IsErr() {
            return fx.Err[[]Komoditas](result.Unwrap().Error)
        }
    }
    
    // Create all komoditas
    results := make([]Komoditas, 0, len(validatedReqs))
    for _, result := range validatedReqs {
        komoditas, _ := result.Unwrap()
        createResult := s.repo.Create(ctx, komoditas)
        if createResult.IsErr() {
            return fx.Err[[]Komoditas](createResult.Unwrap().Error)
        }
        results = append(results, createResult.Unwrap().Value)
    }
    
    return fx.Ok(results)
}

// Pure function for price statistics
func calculatePriceStats(prices []Price) PriceStats {
    if len(prices) == 0 {
        return PriceStats{
            Average: 0,
            Min:     0,
            Max:     0,
            Count:   0,
            Trend:   "stable",
        }
    }
    
    values := fx.Map(prices, func(p Price) float64 { return p.Value })
    
    sum := fx.Reduce(values, 0.0, func(acc, val float64) float64 { 
        return acc + val 
    })
    
    min := fx.Reduce(values, values[0], func(acc, val float64) float64 {
        if val < acc { return val }
        return acc
    })
    
    max := fx.Reduce(values, values[0], func(acc, val float64) float64 {
        if val > acc { return val } 
        return acc
    })
    
    trend := calculateTrend(prices)
    
    return PriceStats{
        Average: sum / float64(len(values)),
        Min:     min,
        Max:     max,
        Count:   len(prices),
        Trend:   trend,
    }
}

// Pure function for trend calculation
func calculateTrend(prices []Price) string {
    if len(prices) < 2 {
        return "stable"
    }
    
    // Sort by date (assuming prices are already sorted by date)
    first, last := prices[0].Value, prices[len(prices)-1].Value
    
    if last > first * 1.05 { // 5% increase
        return "up"
    } else if last < first * 0.95 { // 5% decrease
        return "down"
    }
    return "stable"
}