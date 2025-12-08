package price

import (
    "context"
    "fmt"
    "time"

    "github.com/ryuzxy/FuncPro/pkg/fx"
)

type Service interface {
    CreatePrice(ctx context.Context, req CreatePriceRequest) fx.Result[Price]
    GetPricesByKomoditas(ctx context.Context, id uint) fx.Result[[]Price]
    GetPriceAnalysis(ctx context.Context, id uint) fx.Result[PriceAnalysis]
    BulkCreatePrices(ctx context.Context, reqs []CreatePriceRequest) fx.Result[[]Price]
    GetPriceTrends(ctx context.Context, ids []uint) fx.Result[map[uint]PriceAnalysis]
}

type service struct {
    repo PriceRepository
}

func NewService(repo PriceRepository) Service {
    return &service{repo: repo}
}

func validateCreateRequest(req CreatePriceRequest) error {
    if req.KomoditasID == 0 {
        return fmt.Errorf("komoditas_id required")
    }
    if req.Value <= 0 {
        return fmt.Errorf("value must > 0")
    }
    if req.Date.IsZero() {
        return fmt.Errorf("date required")
    }
    if req.Date.After(time.Now()) {
        return fmt.Errorf("date cannot be future")
    }
    return nil
}

func requestToPrice(req CreatePriceRequest) Price {
    return Price{
        KomoditasID: req.KomoditasID,
        Value:       req.Value,
        Date:        req.Date,
        Market:      req.Market,
    }
}

func (s *service) CreatePrice(ctx context.Context, req CreatePriceRequest) fx.Result[Price] {
    if err := validateCreateRequest(req); err != nil {
        return fx.Err[Price](err)
    }

    p := requestToPrice(req)
    return s.repo.Create(ctx, p)
}

func (s *service) GetPricesByKomoditas(ctx context.Context, id uint) fx.Result[[]Price] {
    return s.repo.GetByKomoditasID(ctx, id)
}

func (s *service) GetPriceAnalysis(ctx context.Context, id uint) fx.Result[PriceAnalysis] {
    end := time.Now()
    start := end.AddDate(0, 0, -30)

    prices, err := s.repo.GetByKomoditasIDAndDateRange(ctx, id, start, end).Unwrap()
    if err != nil {
        return fx.Err[PriceAnalysis](err)
    }

    analysis := analyzePrices(prices)
    return fx.Ok(analysis)
}

func (s *service) BulkCreatePrices(ctx context.Context, reqs []CreatePriceRequest) fx.Result[[]Price] {
    prices := make([]Price, 0, len(reqs))

    for _, req := range reqs {
        if err := validateCreateRequest(req); err != nil {
            return fx.Err[[]Price](err)
        }
        prices = append(prices, requestToPrice(req))
    }

    return s.repo.BulkCreate(ctx, prices)
}

func (s *service) GetPriceTrends(ctx context.Context, ids []uint) fx.Result[map[uint]PriceAnalysis] {
    trends := make(map[uint]PriceAnalysis)

    for _, id := range ids {
        analysis, err := s.GetPriceAnalysis(ctx, id).Unwrap()
        if err != nil {
            return fx.Err[map[uint]PriceAnalysis](err)
        }
        trends[id] = analysis
    }

    return fx.Ok(trends)
}

func analyzePrices(prices []Price) PriceAnalysis {
    if len(prices) == 0 {
        return PriceAnalysis{}
    }

    current := prices[len(prices)-1].Value
    previous := current
    if len(prices) > 1 {
        previous = prices[len(prices)-2].Value
    }

    change := current - previous
    changePct := 0.0
    if previous != 0 {
        changePct = (change / previous) * 100
    }

    values := make([]float64, len(prices))
    for i, p := range prices {
        values[i] = p.Value
    }

    mean := AveragePrice(values)
    variance := 0.0
    for _, v := range values {
        diff := v - mean
        variance += diff * diff
    }
    variance /= float64(len(values))

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
        Volatility: variance,
    }
}
