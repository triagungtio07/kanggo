package product

import (
	"bytes"
	"context"
	"encoding/json"
	"kanggo/pkg/entity/model"
	"kanggo/pkg/mocks"
	"kanggo/utils"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestInsert(t *testing.T) {
	mockProductUsecase := new(mocks.ProductUsecase)

	t.Run("success", func(t *testing.T) {
		mockRequest := model.ProductRequest{
			Name:  "Produk 1",
			Price: 10000,
			Qty:   2,
		}

		mockProductUsecase.On("Insert", mock.Anything, mockRequest).Return(nil)

		body, err := json.Marshal(mockRequest)
		assert.Nil(t, err)

		httpReq, err := http.NewRequest(http.MethodPost, "/api/v1/product", bytes.NewReader(body))
		httpReq.Header.Set("X-Custom-Header", "myvalue")
		httpReq.Header.Set("Content-Type", "application/json")
		assert.Nil(t, err)

		r := gin.Default()
		rr := httptest.NewRecorder()

		h := NewProductHandler(mockProductUsecase)

		r.POST("/api/v1/product", h.Insert)
		r.ServeHTTP(rr, httpReq)

		var resp utils.Respond
		err = json.Unmarshal(rr.Body.Bytes(), &resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.EqualValues(t, http.StatusCreated, rr.Code)
		assert.EqualValues(t, 201, resp.Status)
		assert.EqualValues(t, "success insert product", resp.Message)
		mockProductUsecase.AssertExpectations(t)
	})
}

func TestGetAll(t *testing.T) {
	mockProductUsecase := new(mocks.ProductUsecase)

	t.Run("success", func(t *testing.T) {
		mockProductList := []model.ProductResponse{
			{
				Id:        1,
				Name:      "product 1",
				Price:     10000,
				Qty:       10,
				CreatedAt: "2021-11-05 15:52:59 +0700 WIB",
			},
			{
				Id:        2,
				Name:      "product 2",
				Price:     100000,
				Qty:       8,
				CreatedAt: "2021-11-06 15:52:59 +0700 WIB",
			},
		}

		mockProductUsecase.On("GetAll", mock.Anything).Return(mockProductList, nil)

		httpReq, err := http.NewRequest(http.MethodGet, "/api/v1/product", nil)
		assert.Nil(t, err)

		r := gin.Default()
		rr := httptest.NewRecorder()

		h := NewProductHandler(mockProductUsecase)

		r.GET("/api/v1/product", h.GetAll)
		r.ServeHTTP(rr, httpReq)

		var resp utils.Respond
		err = json.Unmarshal(rr.Body.Bytes(), &resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.EqualValues(t, http.StatusOK, rr.Code)
		assert.EqualValues(t, 200, resp.Status)
		assert.EqualValues(t, "success", resp.Message)
		assert.Equal(t, mockProductList[1].Name, "product 2")
		mockProductUsecase.AssertExpectations(t)
	})
}

func TestGetById(t *testing.T) {
	mockProductUsecase := new(mocks.ProductUsecase)

	t.Run("success", func(t *testing.T) {
		mockResponse := model.ProductResponse{
			Id:        1,
			Name:      "product 1",
			Price:     10000,
			Qty:       10,
			CreatedAt: "2021-11-05 15:52:59 +0700 WIB",
		}

		mockProductUsecase.On("GetById", mock.Anything, mock.AnythingOfType("int64")).Return(&mockResponse, nil)

		id := "1"
		httpReq, err := http.NewRequest(http.MethodGet, "/api/v1/products/"+id, nil)
		assert.Nil(t, err)

		r := gin.Default()
		rr := httptest.NewRecorder()

		h := NewProductHandler(mockProductUsecase)

		r.GET("/api/v1/products/:id", h.GetById)
		r.ServeHTTP(rr, httpReq)

		var resp utils.Respond
		err = json.Unmarshal(rr.Body.Bytes(), &resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.EqualValues(t, http.StatusOK, rr.Code)
		assert.EqualValues(t, 200, resp.Status)
		assert.EqualValues(t, "success", resp.Message)
		assert.Equal(t, mockResponse.Price, float64(10000))
		mockProductUsecase.AssertExpectations(t)
	})
}

func TestUpdate(t *testing.T) {
	mockProductUsecase := new(mocks.ProductUsecase)

	t.Run("success", func(t *testing.T) {
		mockRequest := model.ProductRequest{
			Name:  "Produk 1",
			Price: 10000,
			Qty:   2,
		}

		ctx := context.Background()

		mockProductUsecase.On("Update", ctx, mock.AnythingOfType("uint"), mockRequest).Return(nil)

		body, err := json.Marshal(mockRequest)
		assert.Nil(t, err)

		id := "1"
		httpReq, err := http.NewRequest(http.MethodPut, "/api/v1/product/"+id, bytes.NewReader(body))
		httpReq.Header.Set("X-Custom-Header", "myvalue")
		httpReq.Header.Set("Content-Type", "application/json")
		assert.Nil(t, err)

		r := gin.Default()
		rr := httptest.NewRecorder()

		h := NewProductHandler(mockProductUsecase)

		r.PUT("/api/v1/product/:id", h.Update)
		r.ServeHTTP(rr, httpReq)

		var resp utils.Respond
		err = json.Unmarshal(rr.Body.Bytes(), &resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.EqualValues(t, http.StatusOK, rr.Code)
		assert.EqualValues(t, 200, resp.Status)
		assert.EqualValues(t, "success update product", resp.Message)
		mockProductUsecase.AssertExpectations(t)
	})
}

func TestDelete(t *testing.T) {
	mockProductUsecase := new(mocks.ProductUsecase)

	t.Run("success", func(t *testing.T) {
		mockProductUsecase.On("Delete", mock.Anything, mock.AnythingOfType("int64")).Return(nil)

		id := "1"
		httpReq, err := http.NewRequest(http.MethodDelete, "/api/v1/product/"+id, nil)
		assert.Nil(t, err)

		r := gin.Default()
		rr := httptest.NewRecorder()

		h := NewProductHandler(mockProductUsecase)

		r.DELETE("/api/v1/product/:id", h.Delete)
		r.ServeHTTP(rr, httpReq)

		var resp utils.Respond
		err = json.Unmarshal(rr.Body.Bytes(), &resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.EqualValues(t, http.StatusOK, rr.Code)
		assert.EqualValues(t, 200, resp.Status)
		assert.EqualValues(t, "success delete product", resp.Message)
		mockProductUsecase.AssertExpectations(t)
	})
}
