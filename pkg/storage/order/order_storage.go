package order

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"kanggo/pkg/entity/model"
	"kanggo/pkg/entity/schema"

	"gorm.io/gorm"
)

//go:generate mockery --name OrderStorage --case snake --output ../../mocks --disable-version-string

type (
	OrderStorage interface {
		InsertOrder(ctx context.Context, data schema.Order, quantity int64) error
		GetAllOrder(ctx context.Context) ([]model.OrderResponse, error)
		GetAllOrderPerUser(ctx context.Context, userId uint64) ([]model.OrderResponse, error)
		GetOrderById(ctx context.Context, orderId int64, userId uint64) (*model.OrderResponse, error)
		UpdatePayment(ctx context.Context, data schema.Order) error
	}

	orderStorage struct {
		Native *sql.DB
		Gorm   *gorm.DB
	}
)

func NewOrderStorage(native *sql.DB, gorm *gorm.DB) OrderStorage {
	return &orderStorage{
		Native: native,
		Gorm:   gorm,
	}
}

func (o *orderStorage) InsertOrder(ctx context.Context, data schema.Order, quantity int64) error {
	var product schema.Product
	var qty int64
	var amount float64
	var newQty int64

	type id struct {
		UserId    int64
		ProductId int64
	}

	var ids id

	tx := o.Gorm.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	_ = tx.WithContext(ctx).Where("user_id = ? and product_id=? and status <> 'paid' ", data.UserId, data.ProductId).
		Select("user_id", "product_id").First(&data).Scan(&ids).Error

	if ids.ProductId == 0 && ids.UserId == 0 {
		if err := tx.WithContext(ctx).Create(&data).Error; err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.WithContext(ctx).Where("id = ?", data.ProductId).Select("qty").
			First(&product).Scan(&qty).Error; err != nil {
			tx.Rollback()
			return err
		}

		newQty = qty - quantity

		if err := tx.WithContext(ctx).Model(&product).Where("id=?", data.ProductId).
			Update("qty", newQty).Error; err != nil {
			tx.Rollback()
			return err
		}

	} else {
		addAmount := data.Amount
		if err := tx.WithContext(ctx).Where("user_id = ? and product_id=? and status <> 'paid'", ids.UserId, ids.ProductId).Select("amount").
			First(&data).Scan(&amount).Error; err != nil {
			tx.Rollback()
			return err
		}

		newAmount := amount + addAmount
		data.Amount = newAmount

		if err := tx.WithContext(ctx).Model(&data).Where("user_id = ? and product_id=? and status <> 'paid'", ids.UserId, ids.ProductId).Update("amount", newAmount).Error; err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.WithContext(ctx).Where("id = ?", data.ProductId).Select("qty").
			First(&product).Scan(&qty).Error; err != nil {
			tx.Rollback()
			return err
		}

		newQty = qty - quantity

		if err := tx.WithContext(ctx).Model(&product).Where("id=?", data.ProductId).
			Update("qty", newQty).Error; err != nil {
			tx.Rollback()
			return err
		}

	}

	return tx.Commit().Error
}

func (o *orderStorage) GetAllOrder(ctx context.Context) ([]model.OrderResponse, error) {
	qry := `SELECT o.id, COALESCE(o.user_id,0), COALESCE(u.name,""), COALESCE(o.product_id,0),
	COALESCE(p.name,""), COALESCE(o.amount,0), COALESCE(o.status,"")
	FROM orders as o 
	LEFT JOIN products as p ON p.id = o.product_id
	LEFT JOIN users as u ON u.id = o.user_id
	`

	rows, err := o.Native.QueryContext(ctx, qry)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	products := []model.OrderResponse{}
	for rows.Next() {
		var res model.OrderResponse
		if err := rows.Scan(&res.OrderId, &res.UserId, &res.UserName,
			&res.ProductId, &res.ProductName, &res.Amount, &res.Status); err != nil {
			return nil, err
		}
		products = append(products, res)
	}

	if len(products) == 0 {
		return nil, errors.New("data not found")
	}
	return products, nil
}

func (o *orderStorage) GetAllOrderPerUser(ctx context.Context, userId uint64) ([]model.OrderResponse, error) {
	qry := `SELECT o.id, COALESCE(o.user_id,0), COALESCE(u.name,""), COALESCE(o.product_id,0),
	COALESCE(p.name,""), COALESCE(o.amount,0), COALESCE(o.status,"")
	FROM orders as o 
	LEFT JOIN products as p ON p.id = o.product_id
	LEFT JOIN users as u ON u.id = o.user_id
	WHERE o.user_id = ?
	`

	rows, err := o.Native.QueryContext(ctx, qry, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	products := []model.OrderResponse{}
	for rows.Next() {
		var res model.OrderResponse
		if err := rows.Scan(&res.OrderId, &res.UserId, &res.UserName,
			&res.ProductId, &res.ProductName, &res.Amount, &res.Status); err != nil {
			return nil, err
		}
		products = append(products, res)
	}

	if len(products) == 0 {
		return nil, errors.New("data not found")
	}
	return products, nil
}

func (o *orderStorage) GetOrderById(ctx context.Context, orderId int64, userId uint64) (*model.OrderResponse, error) {
	var result model.OrderResponse

	qry := `SELECT o.id, COALESCE(o.user_id,0), COALESCE(u.name,""), COALESCE(o.product_id,0),
	COALESCE(p.name,""), COALESCE(o.amount,0), COALESCE(o.status,"")
	FROM orders as o 
	LEFT JOIN products as p ON p.id = o.product_id
	LEFT JOIN users as u ON u.id = o.user_id
	WHERE o.id = ? AND o.user_id = ?
	`

	res := o.Native.QueryRowContext(ctx, qry, orderId, userId)
	if err := res.Scan(&result.OrderId, &result.UserId, &result.UserName,
		&result.ProductId, &result.ProductName, &result.Amount, &result.Status); err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	return &result, nil
}

func (o *orderStorage) UpdatePayment(ctx context.Context, data schema.Order) error {
	var amount float64
	data1 := fmt.Sprintf("%.2f", data.Amount)

	tx := o.Gorm.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.WithContext(ctx).Where("user_id = ? and product_id=? and status <> 'paid'", data.UserId, data.ProductId).Select("amount").
		First(&data).Scan(&amount).Error; err != nil {
		tx.Rollback()
		return err
	}
	data2 := fmt.Sprintf("%.2f", amount)

	if data1 != data2 {
		tx.Rollback()
		return errors.New("payment amount does not match")
	}

	if err := tx.WithContext(ctx).Model(&data).Where("user_id = ? and product_id=? and status <> 'paid'", data.UserId, data.ProductId).
		Update("status", "paid").Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
