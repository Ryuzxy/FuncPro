package komoditas

import "time"

// CreateKomoditasRequest DTO for creating komoditas
type CreateKomoditasRequest struct {
    Name string `json:"name" binding:"required,min=1,max=100"`
    Type string `json:"type" binding:"required,min=1,max=50"`
}

// UpdateKomoditasRequest DTO for updating komoditas
type UpdateKomoditasRequest struct {
    Name string `json:"name" binding:"omitempty,min=1,max=100"`
    Type string `json:"type" binding:"omitempty,min=1,max=50"`
}

// KomoditasResponse DTO for komoditas response
type KomoditasResponse struct {
    ID        uint      `json:"id"`
    Name      string    `json:"name"`
    Type      string    `json:"type"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// KomoditasWithStatsResponse DTO for komoditas with stats
type KomoditasWithStatsResponse struct {
    KomoditasResponse
    Stats PriceStats `json:"stats"`
}

// ToResponse converts Komoditas to response DTO
func ToResponse(k Komoditas) KomoditasResponse {
    return KomoditasResponse{
        ID:        k.ID,
        Name:      k.Name,
        Type:      k.Type,
        CreatedAt: k.CreatedAt,
        UpdatedAt: k.UpdatedAt,
    }
}

// ToStatsResponse converts KomoditasWithStats to response DTO
func ToStatsResponse(k KomoditasWithStats) KomoditasWithStatsResponse {
    return KomoditasWithStatsResponse{
        KomoditasResponse: ToResponse(k.Komoditas),
        Stats:             k.Stats,
    }
}