package price

import (
    "context"
    "fmt"
    "time"
    
    "gorm.io/gorm"
    "github.com/ryuzxy/FuncPro/pkg/fx"
)

type PriceRepository interface {
    Create(ctx context.Context, price Price) fx.Result[Price]
    GetByKomoditasID(ctx context.Context, komoditasID uint) fx.Result[[]Price]
    GetByKomoditasIDAndDateRange(ctx context.Context, komoditasID uint, start, end time.Time) fx.Result[[]Price]
    GetLatestByKomoditasID(ctx context.Context, komoditasID uint) fx.Result[Price]
    BulkCreate(ctx context.Context, prices []Price) fx.Result[[]Price]
    Delete(ctx context.Context, id uint) fx.Result[bool]
}

type priceRepository struct {
    db *gorm.DB
}

func NewPriceRepository(db *gorm.DB) PriceRepository {
    return &priceRepository{db: db}
}

func (r *priceRepository) Create(ctx context.Context, price Price) fx.Result[Price] {
    err := r.db.WithContext(ctx).Create(&price).Error
    if err != nil {
        return fx.Err[Price](fmt.Errorf("failed to create price: %w", err))
    }
    return fx.Ok(price)
}

func (r *priceRepository) GetByKomoditasID(ctx context.Context, komoditasID uint) fx.Result[[]Price] {
    var prices []Price
    err := r.db.WithContext(ctx).Where("komoditas_id = ?", komoditasID).Order("date desc").Find(&prices).Error
    if err != nil {
        return fx.Err[[]Price](fmt.Errorf("failed to get prices: %w", err))
    }
    return fx.Ok(prices)
}

func (r *priceRepository) GetByKomoditasIDAndDateRange(ctx context.Context, komoditasID uint, start, end time.Time) fx.Result[[]Price] {
    var prices []Price
    err := r.db.WithContext(ctx).
        Where("komoditas_id = ? AND date BETWEEN ? AND ?", komoditasID, start, end).
        Order("date asc").
        Find(&prices).Error
    if err != nil {
        return fx.Err[[]Price](fmt.Errorf("failed to get prices: %w", err))
    }
    return fx.Ok(prices)
}

func (r *priceRepository) GetLatestByKomoditasID(ctx context.Context, komoditasID uint) fx.Result[Price] {
    var price Price
    err := r.db.WithContext(ctx).
        Where("komoditas_id = ?", komoditasID).
        Order("date desc").
        First(&price).Error
    if err != nil {
        return fx.Err[Price](fmt.Errorf("failed to get latest price: %w", err))
    }
    return fx.Ok(price)
}

func (r *priceRepository) BulkCreate(ctx context.Context, prices []Price) fx.Result[[]Price] {
    err := r.db.WithContext(ctx).CreateInBatches(&prices, 100).Error
    if err != nil {
        return fx.Err[[]Price](fmt.Errorf("failed to bulk create prices: %w", err))
    }
    return fx.Ok(prices)
}

func (r *priceRepository) Delete(ctx context.Context, id uint) fx.Result[bool] {
    err := r.db.WithContext(ctx).Delete(&Price{}, id).Error
    if err != nil {
        return fx.Err[bool](fmt.Errorf("failed to delete price: %w", err))
    }
    return fx.Ok(true)
}