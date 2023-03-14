package product

import (
	"context"
	"errors"
	"github.com/ServiceWeaver/weaver"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"time"
)

type Product struct {
	weaver.AutoMarshal
	ID          int64     `gorm:"column:id"`
	Name        string    `gorm:"column:name"`
	Description string    `gorm:"column:description"`
	Price       float64   `gorm:"column:price"`
	CategoryId  int64     `gorm:"column:category_id"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

type T interface {
	List(ctx context.Context, ids []int64) ([]Product, error)
	Create(ctx context.Context, product Product) (int64, error)
	Update(ctx context.Context, id int64, product Product) error
	Delete(ctx context.Context, id int64) error
}

type impl struct {
	weaver.Implements[T]
	weaver.WithConfig[config]
	db *gorm.DB
}

type config struct {
	Dsn string `toml:"dsn"`
}

func (p *Product) TableName() string {
	return "product_info"
}

func (s *impl) Init(_ context.Context) error {
	cfg := s.Config()
	db, err := gorm.Open(sqlite.Open(cfg.Dsn), &gorm.Config{})
	if err != nil {
		return errors.New("failed to connect database")
	}
	// Migrate the schema
	if err = db.AutoMigrate(&Product{}); err != nil {
		return err
	}
	s.db = db
	return nil
}

func (s *impl) List(ctx context.Context, ids []int64) (list []Product, err error) {
	if err = s.db.WithContext(ctx).Find(&list, ids).Error; err != nil {
		s.Logger().Error("db Find error", err, "ids", ids)
	}
	return
}

func (s *impl) Create(ctx context.Context, product Product) (id int64, err error) {
	if err = s.db.WithContext(ctx).Create(&product).Error; err != nil {
		s.Logger().Error("db Create error", err, "product", product)
	}
	return product.ID, err
}

func (s *impl) Update(ctx context.Context, id int64, product Product) error {
	if err := s.db.WithContext(ctx).Model(&Product{}).Where(`id = ?`, id).Updates(&product).Error; err != nil {
		s.Logger().Error("db Updates error", err, "id", id)
	}
	return nil
}

func (s *impl) Delete(ctx context.Context, id int64) error {
	if err := s.db.WithContext(ctx).Delete(&Product{}, id).Error; err != nil {
		s.Logger().Error("db Delete error", err, "id", id)
	}
	return nil
}
