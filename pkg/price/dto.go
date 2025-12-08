package price

import "time"

type CreatePriceRequest struct {
    KomoditasID uint      `json:"komoditas_id" binding:"required"`
    Value       float64   `json:"value" binding:"required,gt=0"`
    Date        time.Time `json:"date" binding:"required"`
    Market      string    `json:"market" binding:"max=100"`
}

type PriceResponse struct {
    ID          uint      `json:"id"`
    KomoditasID uint      `json:"komoditas_id"`
    Value       float64   `json:"value"`
    Date        time.Time `json:"date"`
    Market      string    `json:"market"`
    CreatedAt   time.Time `json:"created_at"`
}

type PriceAnalysisResponse struct {
    Current    float64 `json:"current"`
    Previous   float64 `json:"previous"`
    Change     float64 `json:"change"`
    ChangePct  float64 `json:"change_percentage"`
    Trend      string  `json:"trend"`
    Volatility float64 `json:"volatility"`
}

func ToResponse(p Price) PriceResponse {
    return PriceResponse{
        ID:          p.ID,
        KomoditasID: p.KomoditasID,
        Value:       p.Value,
        Date:        p.Date,
        Market:      p.Market,
        CreatedAt:   p.CreatedAt,
    }
}

func ToAnalysisResponse(a PriceAnalysis) PriceAnalysisResponse {
    r := PriceAnalysisResponse{}
    r.Current = a.Current
    r.Previous = a.Previous
    r.Change = a.Change
    r.ChangePct = a.ChangePct
    r.Trend = a.Trend
    r.Volatility = a.Volatility
    return r
}
