package price

import (
    "net/http"
    "strconv"
    
    "github.com/gin-gonic/gin"
    "github.com/ryuzxy/FuncPro/pkg/fx"
)

type Handler struct {
    service Service
}

func NewHandler(service Service) *Handler {
    return &Handler{service: service}
}

// CreatePrice creates new price entry
func (h *Handler) CreatePrice(c *gin.Context) {
    var req CreatePriceRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "error":   err.Error(),
        })
        return
    }
    
    result := h.service.CreatePrice(c.Request.Context(), req)
    
    result.Match(
        func(price Price) interface{} {
            c.JSON(http.StatusCreated, gin.H{
                "success": true,
                "data":    ToResponse(price),
            })
            return nil
        },
        func(err error) interface{} {
            c.JSON(http.StatusBadRequest, gin.H{
                "success": false,
                "error":   err.Error(),
            })
            return nil
        },
    )
}

// GetPricesByKomoditas returns prices for specific komoditas
func (h *Handler) GetPricesByKomoditas(c *gin.Context) {
    komoditasID, err := strconv.ParseUint(c.Param("komoditas_id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "error":   "Invalid komoditas ID format",
        })
        return
    }
    
    result := h.service.GetPricesByKomoditas(c.Request.Context(), uint(komoditasID))
    
    result.Match(
        func(prices []Price) interface{} {
            responses := fx.Map(prices, ToResponse)
            c.JSON(http.StatusOK, gin.H{
                "success": true,
                "data":    responses,
                "count":   len(responses),
            })
            return nil
        },
        func(err error) interface{} {
            c.JSON(http.StatusInternalServerError, gin.H{
                "success": false,
                "error":   err.Error(),
            })
            return nil
        },
    )
}

// GetPriceAnalysis returns price analysis for komoditas
func (h *Handler) GetPriceAnalysis(c *gin.Context) {
    komoditasID, err := strconv.ParseUint(c.Param("komoditas_id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "error":   "Invalid komoditas ID format",
        })
        return
    }
    
    result := h.service.GetPriceAnalysis(c.Request.Context(), uint(komoditasID))
    
    result.Match(
        func(analysis PriceAnalysis) interface{} {
            c.JSON(http.StatusOK, gin.H{
                "success": true,
                "data":    ToAnalysisResponse(analysis),
            })
            return nil
        },
        func(err error) interface{} {
            c.JSON(http.StatusInternalServerError, gin.H{
                "success": false,
                "error":   err.Error(),
            })
            return nil
        },
    )
}

// BulkCreatePrices creates multiple price entries
func (h *Handler) BulkCreatePrices(c *gin.Context) {
    var reqs []CreatePriceRequest
    if err := c.ShouldBindJSON(&reqs); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "error":   err.Error(),
        })
        return
    }
    
    result := h.service.BulkCreatePrices(c.Request.Context(), reqs)
    
    result.Match(
        func(prices []Price) interface{} {
            responses := fx.Map(prices, ToResponse)
            c.JSON(http.StatusCreated, gin.H{
                "success": true,
                "data":    responses,
                "count":   len(responses),
            })
            return nil
        },
        func(err error) interface{} {
            c.JSON(http.StatusBadRequest, gin.H{
                "success": false,
                "error":   err.Error(),
            })
            return nil
        },
    )
}