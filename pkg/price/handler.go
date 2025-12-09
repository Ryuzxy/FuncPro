package price

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"

)

type Handler struct {
    service Service
}

func NewHandler(s Service) *Handler {
    return &Handler{service: s}
}

func (h *Handler) CreatePrice(c *gin.Context) {
    var req CreatePriceRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
        return
    }

    price, err := h.service.CreatePrice(c.Request.Context(), req).Unwrap()
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "success": true,
        "data":    ToResponse(price),
    })
}

func (h *Handler) GetPricesByKomoditas(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("komoditas_id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "invalid komoditas id"})
        return
    }

    prices, err := h.service.GetPricesByKomoditas(c.Request.Context(), uint(id)).Unwrap()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
        return
    }

    resp := make([]PriceResponse, 0, len(prices))
    for _, p := range prices {
        resp = append(resp, ToResponse(p))
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "data":    resp,
        "count":   len(resp),
    })
}

func (h *Handler) GetPriceAnalysis(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("komoditas_id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "invalid komoditas id"})
        return
    }

    analysis, err := h.service.GetPriceAnalysis(c.Request.Context(), uint(id)).Unwrap()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "data":    ToAnalysisResponse(analysis),
    })
}

func (h *Handler) BulkCreatePrices(c *gin.Context) {
    var reqs []CreatePriceRequest
    if err := c.ShouldBindJSON(&reqs); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
        return
    }

    prices, err := h.service.BulkCreatePrices(c.Request.Context(), reqs).Unwrap()
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
        return
    }

    resp := make([]PriceResponse, 0, len(prices))
    for _, p := range prices {
        resp = append(resp, ToResponse(p))
    }

    c.JSON(http.StatusCreated, gin.H{
        "success": true,
        "data":    resp,
        "count":   len(resp),
    })
}
