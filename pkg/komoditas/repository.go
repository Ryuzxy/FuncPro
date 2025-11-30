package komoditas

import (
    "context"
    "fmt"
    
    "gorm.io/gorm"
    "github.com/ryuzxy/FuncPro/pkg/fx"
)

type Repository interface {
    GetAll(ctx context.Context) fx.Result[[]Komoditas]
    GetByID(ctx context.Context, id uint) fx.Result[Komoditas]
    Create(ctx context.Context, komoditas Komoditas) fx.Result[Komoditas]
    Update(ctx context.Context, id uint, komoditas Komoditas) fx.Result[Komoditas]
    Delete(ctx context.Context, id uint) fx.Result[bool]
    GetByName(ctx context.Context, name string) fx.Result[Komoditas]
}

type repository struct {
    db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
    return &repository{db: db}
}

func (r *repository) GetAll(ctx context.Context) fx.Result[[]Komoditas] {
    var komoditas []Komoditas
    err := r.db.WithContext(ctx).Find(&komoditas).Error
    if err != nil {
        return fx.Err[[]Komoditas](fmt.Errorf("failed to get komoditas: %w", err))
    }
    return fx.Ok(komoditas)
}

func (r *repository) GetByID(ctx context.Context, id uint) fx.Result[Komoditas] {
    var komoditas Komoditas
    err := r.db.WithContext(ctx).First(&komoditas, id).Error
    if err != nil {
        return fx.Err[Komoditas](fmt.Errorf("komoditas not found: %w", err))
    }
    return fx.Ok(komoditas)
}

func (r *repository) Create(ctx context.Context, komoditas Komoditas) fx.Result[Komoditas] {
    err := r.db.WithContext(ctx).Create(&komoditas).Error
    if err != nil {
        return fx.Err[Komoditas](fmt.Errorf("failed to create komoditas: %w", err))
    }
    return fx.Ok(komoditas)
}

func (r *repository) Update(ctx context.Context, id uint, komoditas Komoditas) fx.Result[Komoditas] {
    err := r.db.WithContext(ctx).Model(&Komoditas{}).Where("id = ?", id).Updates(komoditas).Error
    if err != nil {
        return fx.Err[Komoditas](fmt.Errorf("failed to update komoditas: %w", err))
    }
    
    return r.GetByID(ctx, id)
}

func (r *repository) Delete(ctx context.Context, id uint) fx.Result[bool] {
    err := r.db.WithContext(ctx).Delete(&Komoditas{}, id).Error
    if err != nil {
        return fx.Err[bool](fmt.Errorf("failed to delete komoditas: %w", err))
    }
    return fx.Ok(true)
}

func (r *repository) GetByName(ctx context.Context, name string) fx.Result[Komoditas] {
    var komoditas Komoditas
    err := r.db.WithContext(ctx).Where("name = ?", name).First(&komoditas).Error
    if err != nil {
        return fx.Err[Komoditas](fmt.Errorf("komoditas not found: %w", err))
    }
    return fx.Ok(komoditas)
}