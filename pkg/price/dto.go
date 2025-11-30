package price

import "time"

// CreatePriceRequest DTO for creating price
type CreatePriceRequest struct {
    KomoditasID uint      `json:"komoditas_id" binding:"required"`
    Value       float64   `json:"value" binding:"required,gt=0"`
    Date        time.Time `json:"date" binding:"required"`
    Market      string    `json:"market" binding:"max=100"`
}

// PriceResponse DTO for price response
type PriceResponse struct {
    ID          uint      `json:"id"`
    KomoditasID uint      `json:"komoditas_id"`
    Value       float64   `json:"value"`
    Date        time.Time `json:"date"`
    Market      string    `json:"market"`
    CreatedAt   time.Time `json:"created_at"`
}

// PriceAnalysisResponse DTO for price analysis
type PriceAnalysisResponse struct {
    Current    float64 `json:"current"`
    Previous   float64 `json:"previous"`
    Change     float64 `json:"change"`
    ChangePct  float64 `json:"change_percentage"`
    Trend      string  `json:"trend"`
    Volatility float64 `json:"volatility"`
}

// ToResponse converts Price to response DTO
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

// ToAnalysisResponse converts PriceAnalysis to response DTO
func ToAnalysisResponse(analysis PriceAnalysis) PriceAnalysisResponse {
    return PriceAnalysisResponse{
        Current:    analysis.Current,
        Previous:   analysis.Previous,
        Change:     analysis.Change,
        ChangePct:  analysis.ChangePct,
        Trend:      analysis.Trend,
        Volatility: analysis.Volatility,
    }
}