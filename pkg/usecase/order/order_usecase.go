package order

import (
	"context"
	"kanggo/pkg/entity/model"
	"kanggo/pkg/entity/schema"
	storage "kanggo/pkg/storage/order"
	productStorage "kanggo/pkg/storage/product"
)

//go:generate mockery --name OrderUsecase --case snake --output ../../mocks --disable-version-string

type (
	OrderUsecase interface {
		InsertOrder(ctx context.Context, data model.OrderRequest) error
		GetAllOrder(ctx context.Context) ([]model.OrderResponse, error)
		GetAllOrderPerUser(ctx context.Context, userId uint64) ([]model.OrderResponse, error)
		GetOrderById(ctx context.Context, orderId int64, userId uint64) (*model.OrderResponse, error)
		UpdatePayment(ctx context.Context, data model.PaymentRequest) error
	}

	orderUsecase struct {
		orderStorage   storage.OrderStorage
		productStorage productStorage.ProductStorage
	}
)

func NewOrderUsecase(orderStorage storage.OrderStorage, productStorage productStorage.ProductStorage) OrderUsecase {
	return &orderUsecase{
		orderStorage:   orderStorage,
		productStorage: productStorage,
	}
}

func (o *orderUsecase) InsertOrder(ctx context.Context, data model.OrderRequest) error {
	request := schema.Order{
		UserId:    data.UserId,
		ProductId: data.ProductId,
		Amount:    data.Amount,
	}

	status, err := o.productStorage.CheckQty(ctx, data.ProductId, data.Quantity)
	if err != nil {
		return err
	}

	if status {
		if err = o.orderStorage.InsertOrder(ctx, request, data.Quantity); err != nil {
			return err
		}

	}

	return nil

}

func (o *orderUsecase) GetAllOrder(ctx context.Context) ([]model.OrderResponse, error) {
	res, err := o.orderStorage.GetAllOrder(ctx)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (o *orderUsecase) GetAllOrderPerUser(ctx context.Context, userId uint64) ([]model.OrderResponse, error) {
	res, err := o.orderStorage.GetAllOrderPerUser(ctx, userId)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (o *orderUsecase) GetOrderById(ctx context.Context, orderId int64, userId uint64) (*model.OrderResponse, error) {

	res, err := o.orderStorage.GetOrderById(ctx, orderId, userId)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (o *orderUsecase) UpdatePayment(ctx context.Context, data model.PaymentRequest) error {
	request := schema.Order{
		UserId:    data.UserId,
		ProductId: data.ProductId,
		Amount:    data.Amount,
	}

	if err := o.orderStorage.UpdatePayment(ctx, request); err != nil {
		return err
	}

	return nil
}
