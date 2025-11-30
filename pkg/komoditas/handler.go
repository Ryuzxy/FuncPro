package komoditas

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

// GetAllKomoditas returns all komoditas
func (h *Handler) GetAllKomoditas(c *gin.Context) {
    result := h.service.GetAllKomoditas(c.Request.Context())
    
    result.Match(
        func(komoditas []Komoditas) interface{} {
            responses := fx.Map(komoditas, ToResponse)
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

// GetKomoditasByID returns komoditas by ID
func (h *Handler) GetKomoditasByID(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "error":   "Invalid ID format",
        })
        return
    }
    
    result := h.service.GetKomoditasByID(c.Request.Context(), uint(id))
    
    result.Match(
        func(komoditas Komoditas) interface{} {
            c.JSON(http.StatusOK, gin.H{
                "success": true,
                "data":    ToResponse(komoditas),
            })
            return nil
        },
        func(err error) interface{} {
            c.JSON(http.StatusNotFound, gin.H{
                "success": false,
                "error":   err.Error(),
            })
            return nil
        },
    )
}

// CreateKomoditas creates new komoditas
func (h *Handler) CreateKomoditas(c *gin.Context) {
    var req CreateKomoditasRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "error":   err.Error(),
        })
        return
    }
    
    result := h.service.CreateKomoditas(c.Request.Context(), req)
    
    result.Match(
        func(komoditas Komoditas) interface{} {
            c.JSON(http.StatusCreated, gin.H{
                "success": true,
                "data":    ToResponse(komoditas),
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

// UpdateKomoditas updates existing komoditas
func (h *Handler) UpdateKomoditas(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "error":   "Invalid ID format",
        })
        return
    }
    
    var req UpdateKomoditasRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "error":   err.Error(),
        })
        return
    }
    
    result := h.service.UpdateKomoditas(c.Request.Context(), uint(id), req)
    
    result.Match(
        func(komoditas Komoditas) interface{} {
            c.JSON(http.StatusOK, gin.H{
                "success": true,
                "data":    ToResponse(komoditas),
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

// DeleteKomoditas deletes komoditas
func (h *Handler) DeleteKomoditas(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "error":   "Invalid ID format",
        })
        return
    }
    
    result := h.service.DeleteKomoditas(c.Request.Context(), uint(id))
    
    result.Match(
        func(success bool) interface{} {
            c.JSON(http.StatusOK, gin.H{
                "success": true,
                "message": "Komoditas deleted successfully",
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

// GetKomoditasStats returns komoditas with price statistics
func (h *Handler) GetKomoditasStats(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "error":   "Invalid ID format",
        })
        return
    }
    
    result := h.service.GetKomoditasWithStats(c.Request.Context(), uint(id))
    
    result.Match(
        func(komoditasWithStats KomoditasWithStats) interface{} {
            c.JSON(http.StatusOK, gin.H{
                "success": true,
                "data":    ToStatsResponse(komoditasWithStats),
            })
            return nil
        },
        func(err error) interface{} {
            c.JSON(http.StatusNotFound, gin.H{
                "success": false,
                "error":   err.Error(),
            })
            return nil
        },
    )
}