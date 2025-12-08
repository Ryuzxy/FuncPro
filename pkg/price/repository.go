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
    if err := r.db.WithContext(ctx).Create(&price).Error; err != nil {
        return fx.Err[Price](fmt.Errorf("create failed: %w", err))
    }
    return fx.Ok(price)
}

func (r *priceRepository) GetByKomoditasID(ctx context.Context, komoditasID uint) fx.Result[[]Price] {
    var list []Price
    err := r.db.WithContext(ctx).
        Where("komoditas_id = ?", komoditasID).
        Order("date asc").
        Find(&list).Error

    if err != nil {
        return fx.Err[[]Price](fmt.Errorf("query failed: %w", err))
    }

    return fx.Ok(list)
}

func (r *priceRepository) GetByKomoditasIDAndDateRange(ctx context.Context, komoditasID uint, start, end time.Time) fx.Result[[]Price] {
    var list []Price
    err := r.db.WithContext(ctx).
        Where("komoditas_id = ? AND date BETWEEN ? AND ?", komoditasID, start, end).
        Order("date asc").
        Find(&list).Error

    if err != nil {
        return fx.Err[[]Price](fmt.Errorf("range query failed: %w", err))
    }

    return fx.Ok(list)
}

func (r *priceRepository) GetLatestByKomoditasID(ctx context.Context, komoditasID uint) fx.Result[Price] {
    var p Price
    err := r.db.WithContext(ctx).
        Where("komoditas_id = ?", komoditasID).
        Order("date desc").
        First(&p).Error

    if err != nil {
        return fx.Err[Price](fmt.Errorf("latest query failed: %w", err))
    }

    return fx.Ok(p)
}

func (r *priceRepository) BulkCreate(ctx context.Context, prices []Price) fx.Result[[]Price] {
    if err := r.db.WithContext(ctx).CreateInBatches(&prices, 100).Error; err != nil {
        return fx.Err[[]Price](fmt.Errorf("bulk insert failed: %w", err))
    }
    return fx.Ok(prices)
}

func (r *priceRepository) Delete(ctx context.Context, id uint) fx.Result[bool] {
    if err := r.db.WithContext(ctx).Delete(&Price{}, id).Error; err != nil {
        return fx.Err[bool](fmt.Errorf("delete failed: %w", err))
    }
    return fx.Ok(true)
}
