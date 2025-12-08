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

// -------------------- GET ALL --------------------
func (h *Handler) GetAllKomoditas(c *gin.Context) {
    result := h.service.GetAllKomoditas(c.Request.Context())

    fx.Match(
        result,
        func(data []Komoditas) any {
            responses := fx.Map(data, ToResponse)
            c.JSON(http.StatusOK, gin.H{
                "success": true,
                "data":    responses,
                "count":   len(responses),
            })
            return nil
        },
        func(err error) any {
            c.JSON(http.StatusInternalServerError, gin.H{
                "success": false,
                "error":   err.Error(),
            })
            return nil
        },
    )
}

// -------------------- GET BY ID --------------------
func (h *Handler) GetKomoditasByID(c *gin.Context) {
    id, ok := parseID(c)
    if !ok {
        return
    }

    result := h.service.GetKomoditasByID(c.Request.Context(), id)

    fx.Match(
        result,
        func(data *Komoditas) any {
            if data == nil {
                c.JSON(http.StatusNotFound, gin.H{
                    "success": false,
                    "error":   "Komoditas not found",
                })
                return nil
            }
            c.JSON(http.StatusOK, gin.H{
                "success": true,
                "data":    ToResponse(*data),
            })
            return nil
        },
        func(err error) any {
            c.JSON(http.StatusNotFound, gin.H{
                "success": false,
                "error":   err.Error(),
            })
            return nil
        },
    )
}

// -------------------- CREATE --------------------
func (h *Handler) CreateKomoditas(c *gin.Context) {
    var req CreateKomoditasRequest
    if !bindJSON(c, &req) {
        return
    }

    result := h.service.CreateKomoditas(c.Request.Context(), req)

    fx.Match(
        result,
        func(data *Komoditas) any {
            if data == nil {
                c.JSON(http.StatusInternalServerError, gin.H{
                    "success": false,
                    "error":   "failed to create komoditas",
                })
                return nil
            }
            c.JSON(http.StatusCreated, gin.H{
                "success": true,
                "data":    ToResponse(*data),
            })
            return nil
        },
        func(err error) any {
            c.JSON(http.StatusBadRequest, gin.H{
                "success": false,
                "error":   err.Error(),
            })
            return nil
        },
    )
}

// -------------------- UPDATE --------------------
func (h *Handler) UpdateKomoditas(c *gin.Context) {
    id, ok := parseID(c)
    if !ok {
        return
    }

    var req UpdateKomoditasRequest
    if !bindJSON(c, &req) {
        return
    }

    result := h.service.UpdateKomoditas(c.Request.Context(), id, req)

    fx.Match(
        result,
        func(data *Komoditas) any {
            if data == nil {
                c.JSON(http.StatusNotFound, gin.H{
                    "success": false,
                    "error":   "Komoditas not found",
                })
                return nil
            }
            c.JSON(http.StatusOK, gin.H{
                "success": true,
                "data":    ToResponse(*data),
            })
            return nil
        },
        func(err error) any {
            c.JSON(http.StatusBadRequest, gin.H{
                "success": false,
                "error":   err.Error(),
            })
            return nil
        },
    )
}

// -------------------- DELETE --------------------
func (h *Handler) DeleteKomoditas(c *gin.Context) {
    id, ok := parseID(c)
    if !ok {
        return
    }

    result := h.service.DeleteKomoditas(c.Request.Context(), id)

    fx.Match(
        result,
        func(success bool) any {
            if !success {
                c.JSON(http.StatusInternalServerError, gin.H{
                    "success": false,
                    "error":   "failed to delete komoditas",
                })
                return nil
            }
            c.JSON(http.StatusOK, gin.H{
                "success": true,
                "message": "Komoditas deleted successfully",
            })
            return nil
        },
        func(err error) any {
            c.JSON(http.StatusBadRequest, gin.H{
                "success": false,
                "error":   err.Error(),
            })
            return nil
        },
    )
}

// -------------------- STATISTICS --------------------
func (h *Handler) GetKomoditasStats(c *gin.Context) {
    id, ok := parseID(c)
    if !ok {
        return
    }

    result := h.service.GetKomoditasWithStats(c.Request.Context(), id)

    fx.Match(
        result,
        func(data KomoditasWithStats) any {
            c.JSON(http.StatusOK, gin.H{
                "success": true,
                "data":    ToStatsResponse(data),
            })
            return nil
        },
        func(err error) any {
            c.JSON(http.StatusNotFound, gin.H{
                "success": false,
                "error":   err.Error(),
            })
            return nil
        },
    )
}

// -------------------- HELPERS --------------------
func parseID(c *gin.Context) (uint, bool) {
    id64, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "error":   "Invalid ID format",
        })
        return 0, false
    }
    return uint(id64), true
}

func bindJSON[T any](c *gin.Context, target *T) bool {
    if err := c.ShouldBindJSON(target); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "error":   err.Error(),
        })
        return false
    }
    return true
}
