package product

import (
	"context"
	"database/sql"
	"errors"
	"kanggo/pkg/entity/schema"

	"gorm.io/gorm"
)

//go:generate mockery --name ProductStorage --case snake --output ../../mocks --disable-version-string

type (
	ProductStorage interface {
		Insert(ctx context.Context, data schema.Product) error
		Update(ctx context.Context, data schema.Product) error
		GetAll(ctx context.Context) ([]schema.Product, error)
		GetById(ctx context.Context, id int64) (*schema.Product, error)
		Delete(ctx context.Context, id int64) error
		CheckQty(ctx context.Context, productId, amount int64) (bool, error)
	}

	productStorage struct {
		Native *sql.DB
		Gorm   *gorm.DB
	}
)

func NewProductStorage(native *sql.DB, gorm *gorm.DB) ProductStorage {
	return &productStorage{
		Native: native,
		Gorm:   gorm,
	}
}

func (p *productStorage) Insert(ctx context.Context, data schema.Product) error {
	if err := p.Gorm.WithContext(ctx).Create(&data).Error; err != nil {
		return err
	}

	return nil
}

func (p *productStorage) Update(ctx context.Context, data schema.Product) error {
	result := p.Gorm.WithContext(ctx).Updates(data)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("data not found")
	}

	return nil
}

func (p *productStorage) GetAll(ctx context.Context) ([]schema.Product, error) {
	qry := `SELECT * FROM products`

	rows, err := p.Native.QueryContext(ctx, qry)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	products := []schema.Product{}
	for rows.Next() {
		var res schema.Product
		if err := rows.Scan(&res.Id, &res.CreatedAt, &res.UpdatedAt, &res.Name, &res.Price, &res.Qty); err != nil {
			return nil, err
		}
		products = append(products, res)
	}

	if len(products) == 0 {
		return nil, errors.New("data not found")
	}
	return products, nil
}

func (p *productStorage) GetById(ctx context.Context, id int64) (*schema.Product, error) {
	product := schema.Product{}
	qry := `SELECT * FROM products WHERE id = ?`

	res := p.Native.QueryRowContext(ctx, qry, id)
	if err := res.Scan(&product.Id, &product.CreatedAt, &product.UpdatedAt,
		&product.Name, &product.Price, &product.Qty); err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	return &product, nil
}

func (p *productStorage) Delete(ctx context.Context, id int64) error {
	result := p.Gorm.WithContext(ctx).Delete(schema.Product{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("data not found")
	}

	return nil
}

func (o *productStorage) CheckQty(ctx context.Context, productId, amount int64) (bool, error) {
	var id int64
	qry := `SELECT id FROM products WHERE id = ? AND qty > ?`

	res := o.Native.QueryRowContext(ctx, qry, productId, amount)
	if err := res.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return false, errors.New("not enough product quantity")
		}
		return false, err
	}

	return true, nil
}
