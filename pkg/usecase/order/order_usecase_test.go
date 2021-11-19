package order

import (
	"context"
	"testing"

	"kanggo/pkg/entity/model"
	"kanggo/pkg/entity/schema"
	"kanggo/pkg/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestInsert(t *testing.T) {
	mockProductStorage := new(mocks.ProductStorage)
	mockOrderStorage := new(mocks.OrderStorage)
	o := NewOrderUsecase(mockOrderStorage, mockProductStorage)
	ctx := context.Background()

	model := model.OrderRequest{
		UserId:    1,
		ProductId: 1,
		Amount:    100000,
		Quantity:  2,
	}
	schema := schema.Order{
		UserId:    1,
		ProductId: 1,
		Amount:    100000,
	}

	t.Run("success", func(t *testing.T) {
		mockProductStorage.On("CheckQty", ctx, model.ProductId, model.Quantity).Return(true, nil)
		mockOrderStorage.On("InsertOrder", ctx, schema, model.Quantity).Return(nil)

		err := o.InsertOrder(ctx, model)

		assert.Nil(t, err)
		assert.NoError(t, err)
		mockProductStorage.AssertExpectations(t)
	})
}

func TestGetAllOrder(t *testing.T) {
	mockProductStorage := new(mocks.ProductStorage)
	mockOrderStorage := new(mocks.OrderStorage)
	o := NewOrderUsecase(mockOrderStorage, mockProductStorage)
	ctx := context.Background()

	mockOrderList := []model.OrderResponse{
		{
			OrderId:     1,
			UserId:      1,
			UserName:    "Agung",
			ProductId:   1,
			ProductName: "Produk 1",
			Amount:      320000,
			Status:      "pending",
		},
		{
			OrderId:     2,
			UserId:      2,
			UserName:    "Bayu",
			ProductId:   1,
			ProductName: "Produk 1",
			Amount:      6000000,
			Status:      "paid",
		},
	}

	t.Run("success", func(t *testing.T) {
		mockOrderStorage.On("GetAllOrder", mock.Anything).Return(mockOrderList, nil)

		list, err := o.GetAllOrder(ctx)

		assert.NotNil(t, list)
		assert.Nil(t, err)
		assert.NoError(t, err)
		assert.Equal(t, mockOrderList[0].UserName, "Agung")
		assert.Len(t, list, len(mockOrderList))
		mockProductStorage.AssertExpectations(t)
	})
}

func TestGetAllOrderPerUser(t *testing.T) {
	mockProductStorage := new(mocks.ProductStorage)
	mockOrderStorage := new(mocks.OrderStorage)
	o := NewOrderUsecase(mockOrderStorage, mockProductStorage)
	ctx := context.Background()
	var userId uint64 = 1

	mockOrderList := []model.OrderResponse{
		{
			OrderId:     1,
			UserId:      1,
			UserName:    "Agung",
			ProductId:   1,
			ProductName: "Produk 1",
			Amount:      320000,
			Status:      "pending",
		},
		{
			OrderId:     4,
			UserId:      1,
			UserName:    "Agung",
			ProductId:   1,
			ProductName: "Produk 1",
			Amount:      6000000,
			Status:      "pending",
		},
	}

	t.Run("success", func(t *testing.T) {
		mockOrderStorage.On("GetAllOrderPerUser", mock.Anything, userId).Return(mockOrderList, nil)

		list, err := o.GetAllOrderPerUser(ctx, userId)

		assert.NotNil(t, list)
		assert.Nil(t, err)
		assert.NoError(t, err)
		assert.Equal(t, mockOrderList[0].UserName, "Agung")
		assert.Len(t, list, len(mockOrderList))
		mockProductStorage.AssertExpectations(t)
	})
}

func TestGetOrderById(t *testing.T) {
	mockProductStorage := new(mocks.ProductStorage)
	mockOrderStorage := new(mocks.OrderStorage)
	o := NewOrderUsecase(mockOrderStorage, mockProductStorage)
	ctx := context.Background()
	var userId uint64 = 1
	var orderId int64 = 1

	mockOrder := model.OrderResponse{
		OrderId:     1,
		UserId:      1,
		UserName:    "Agung",
		ProductId:   1,
		ProductName: "Produk 1",
		Amount:      320000,
		Status:      "pending",
	}

	t.Run("success", func(t *testing.T) {
		mockOrderStorage.On("GetOrderById", mock.Anything, orderId, userId).Return(&mockOrder, nil)

		list, err := o.GetOrderById(ctx, orderId, userId)

		assert.NotNil(t, list)
		assert.Nil(t, err)
		assert.NoError(t, err)
		assert.Equal(t, mockOrder.UserName, "Agung")
		mockProductStorage.AssertExpectations(t)
	})
}

func TestUpdatePayment(t *testing.T) {
	mockProductStorage := new(mocks.ProductStorage)
	mockOrderStorage := new(mocks.OrderStorage)
	o := NewOrderUsecase(mockOrderStorage, mockProductStorage)
	ctx := context.Background()

	model := model.PaymentRequest{
		UserId:    1,
		ProductId: 1,
		Amount:    5000,
	}
	schema := schema.Order{
		UserId:    1,
		ProductId: 1,
		Amount:    5000,
	}

	t.Run("success", func(t *testing.T) {
		mockOrderStorage.On("UpdatePayment", mock.Anything, schema).Return(nil)

		err := o.UpdatePayment(ctx, model)

		assert.Nil(t, err)
		assert.NoError(t, err)
		mockProductStorage.AssertExpectations(t)
	})
}
