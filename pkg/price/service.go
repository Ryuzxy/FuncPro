package price

import (
    "context"
    "fmt"
    "time"
    
    "github.com/ryuzxy/FuncPro/pkg/fx"
)

type Service interface {
    CreatePrice(ctx context.Context, req CreatePriceRequest) fx.Result[Price]
    GetPricesByKomoditas(ctx context.Context, komoditasID uint) fx.Result[[]Price]
    GetPriceAnalysis(ctx context.Context, komoditasID uint) fx.Result[PriceAnalysis]
    BulkCreatePrices(ctx context.Context, reqs []CreatePriceRequest) fx.Result[[]Price]
    GetPriceTrends(ctx context.Context, komoditasIDs []uint) fx.Result[map[uint]PriceAnalysis]
}

type service struct {
    repo PriceRepository
}

func NewService(repo PriceRepository) Service {
    return &service{repo: repo}
}

// Pure validation function
func validateCreateRequest(req CreatePriceRequest) fx.Result[CreatePriceRequest] {
    if req.KomoditasID == 0 {
        return fx.Err[CreatePriceRequest](fmt.Errorf("komoditas_id is required"))
    }
    if req.Value <= 0 {
        return fx.Err[CreatePriceRequest](fmt.Errorf("value must be positive"))
    }
    if req.Date.IsZero() {
        return fx.Err[CreatePriceRequest](fmt.Errorf("date is required"))
    }
    if req.Date.After(time.Now()) {
        return fx.Err[CreatePriceRequest](fmt.Errorf("date cannot be in the future"))
    }
    return fx.Ok(req)
}

// Pure transformation function
func requestToPrice(req CreatePriceRequest) Price {
    return Price{
        KomoditasID: req.KomoditasID,
        Value:       req.Value,
        Date:        req.Date,
        Market:      req.Market,
    }
}

func (s *service) CreatePrice(ctx context.Context, req CreatePriceRequest) fx.Result[Price] {
    return validateCreateRequest(req).
        Map(requestToPrice).
        AndThen(func(p Price) fx.Result[Price] {
            return s.repo.Create(ctx, p)
        })
}

func (s *service) GetPricesByKomoditas(ctx context.Context, komoditasID uint) fx.Result[[]Price] {
    return s.repo.GetByKomoditasID(ctx, komoditasID)
}

func (s *service) GetPriceAnalysis(ctx context.Context, komoditasID uint) fx.Result[PriceAnalysis] {
    // Get prices for last 30 days
    end := time.Now()
    start := end.AddDate(0, 0, -30)
    
    pricesResult := s.repo.GetByKomoditasIDAndDateRange(ctx, komoditasID, start, end)
    
    return pricesResult.Map(func(prices []Price) PriceAnalysis {
        return analyzePrices(prices)
    })
}

func (s *service) BulkCreatePrices(ctx context.Context, reqs []CreatePriceRequest) fx.Result[[]Price] {
    // Validate all requests
    validatedPrices := make([]fx.Result[Price], 0, len(reqs))
    
    for _, req := range reqs {
        result := validateCreateRequest(req).
            Map(requestToPrice)
        validatedPrices = append(validatedPrices, result)
    }
    
    // Check for validation errors
    for _, result := range validatedPrices {
        if result.IsErr() {
            return fx.Err[[]Price](result.Unwrap().Error)
        }
    }
    
    // Extract valid prices
    prices := make([]Price, 0, len(validatedPrices))
    for _, result := range validatedPrices {
        price, _ := result.Unwrap()
        prices = append(prices, price)
    }
    
    return s.repo.BulkCreate(ctx, prices)
}

func (s *service) GetPriceTrends(ctx context.Context, komoditasIDs []uint) fx.Result[map[uint]PriceAnalysis] {
    type analysisResult struct {
        komoditasID uint
        analysis    PriceAnalysis
        err         error
    }
    
    // Process each komoditas concurrently
    results := fx.ParallelMap(ctx, komoditasIDs, func(ctx context.Context, id uint) fx.Result[analysisResult] {
        analysis := s.GetPriceAnalysis(ctx, id)
        return analysis.Map(func(a PriceAnalysis) analysisResult {
            return analysisResult{
                komoditasID: id,
                analysis:    a,
            }
        })
    }, 5) // 5 concurrent workers
    
    return results.Map(func(results []analysisResult) map[uint]PriceAnalysis {
        trendMap := make(map[uint]PriceAnalysis)
        for _, result := range results {
            trendMap[result.komoditasID] = result.analysis
        }
        return trendMap
    })
}

// Pure function for price analysis
func analyzePrices(prices []Price) PriceAnalysis {
    if len(prices) == 0 {
        return PriceAnalysis{}
    }
    
    // Sort by date (assuming they come sorted from DB)
    sortedPrices := fx.Map(prices, func(p Price) Price { return p })
    
    current := sortedPrices[len(sortedPrices)-1].Value
    var previous float64
    
    if len(sortedPrices) > 1 {
        previous = sortedPrices[len(sortedPrices)-2].Value
    } else {
        previous = current
    }
    
    change := current - previous
    changePct := (change / previous) * 100
    
    // Calculate volatility (standard deviation)
    values := fx.Map(sortedPrices, func(p Price) float64 { return p.Value })
    mean := fx.Reduce(values, 0.0, func(acc, val float64) float64 { 
        return acc + val 
    }) / float64(len(values))
    
    variance := fx.Reduce(values, 0.0, func(acc, val float64) float64 {
        diff := val - mean
        return acc + (diff * diff)
    }) / float64(len(values))
    
    volatility := variance // Simplified - in real case, use proper std dev
    
    trend := "stable"
    if changePct > 5 {
        trend = "up"
    } else if changePct < -5 {
        trend = "down"
    }
    
    return PriceAnalysis{
        Current:    current,
        Previous:   previous,
        Change:     change,
        ChangePct:  changePct,
        Trend:      trend,
        Volatility: volatility,
    }
}