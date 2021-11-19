package product

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
	p := NewProductUsecase(mockProductStorage)
	ctx := context.Background()

	model := model.ProductRequest{
		Name:  "Produk 1",
		Price: 10000,
		Qty:   10,
	}
	schema := schema.Product{
		Name:  "Produk 1",
		Price: 10000,
		Qty:   10,
	}

	t.Run("success", func(t *testing.T) {
		mockProductStorage.On("Insert", mock.Anything, schema).Return(nil)

		err := p.Insert(ctx, model)

		assert.Nil(t, err)
		assert.NoError(t, err)
		mockProductStorage.AssertExpectations(t)
	})
}

func TestUpdate(t *testing.T) {
	mockProductStorage := new(mocks.ProductStorage)
	p := NewProductUsecase(mockProductStorage)
	ctx := context.Background()
	var id uint = 1

	model := model.ProductRequest{
		Name:  "Produk 1",
		Price: 10000,
		Qty:   10,
	}
	schema := schema.Product{
		Base:  schema.Base{Id: 1},
		Name:  "Produk 1",
		Price: 10000,
		Qty:   10,
	}

	t.Run("success", func(t *testing.T) {
		mockProductStorage.On("Update", mock.Anything, schema).Return(nil)

		err := p.Update(ctx, id, model)

		assert.Nil(t, err)
		assert.NoError(t, err)
		mockProductStorage.AssertExpectations(t)
	})
}

func TestGetAll(t *testing.T) {
	mockProductStorage := new(mocks.ProductStorage)
	ctx := context.Background()

	mockProductList := []schema.Product{
		{
			Base:  schema.Base{Id: 1},
			Name:  "product 1",
			Price: 10000,
			Qty:   10,
		},
		{
			Base:  schema.Base{Id: 2},
			Name:  "product 2",
			Price: 20000,
			Qty:   5,
		},
	}

	t.Run("success", func(t *testing.T) {
		mockProductStorage.On("GetAll", mock.Anything).Return(mockProductList, nil)

		u := NewProductUsecase(mockProductStorage)
		list, err := u.GetAll(ctx)

		assert.NotNil(t, list)
		assert.Nil(t, err)
		assert.NoError(t, err)
		assert.Equal(t, mockProductList[0].Name, "product 1")
		assert.Len(t, list, len(mockProductList))
		mockProductStorage.AssertExpectations(t)
	})
}

func TestGetById(t *testing.T) {
	mockProductStorage := new(mocks.ProductStorage)
	ctx := context.Background()
	var idProduct int64 = 1

	mockProduct := schema.Product{
		Base:  schema.Base{Id: 1},
		Name:  "product 1",
		Price: 10000,
		Qty:   10,
	}

	t.Run("success", func(t *testing.T) {
		mockProductStorage.On("GetById", mock.Anything, mock.AnythingOfType("int64")).Return(&mockProduct, nil)

		u := NewProductUsecase(mockProductStorage)

		detail, err := u.GetById(ctx, idProduct)

		assert.NotNil(t, detail)
		assert.Nil(t, err)
		assert.NoError(t, err)
		assert.Equal(t, mockProduct.Name, "product 1")
		mockProductStorage.AssertExpectations(t)
	})
}

func TestDelete(t *testing.T) {
	mockProductStorage := new(mocks.ProductStorage)
	ctx := context.Background()
	var idProduct int64 = 1

	t.Run("success", func(t *testing.T) {
		mockProductStorage.On("Delete", mock.Anything, mock.AnythingOfType("int64")).Return(nil)

		u := NewProductUsecase(mockProductStorage)

		err := u.Delete(ctx, idProduct)

		assert.Nil(t, err)
		assert.NoError(t, err)
		mockProductStorage.AssertExpectations(t)
	})
}
