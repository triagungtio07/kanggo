package order

import (
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

func TestGetAllOrder(t *testing.T) {
	mockOrderUsecase := new(mocks.OrderUsecase)

	t.Run("success", func(t *testing.T) {
		mockOrderList := []model.OrderResponse{
			{
				OrderId:     1,
				UserId:      1,
				UserName:    "Agung",
				ProductId:   1,
				ProductName: "Produk 1",
				Amount:      50000,
				Status:      "paid",
			},
			{
				OrderId:     2,
				UserId:      1,
				UserName:    "Agung",
				ProductId:   2,
				ProductName: "Produk 2",
				Amount:      50000,
				Status:      "unpaid",
			},
		}

		mockOrderUsecase.On("GetAllOrder", mock.Anything).Return(mockOrderList, nil)

		httpReq, err := http.NewRequest(http.MethodGet, "/api/v1/order", nil)
		assert.Nil(t, err)

		r := gin.Default()
		rr := httptest.NewRecorder()

		h := NewOrderHandler(mockOrderUsecase)

		r.GET("/api/v1/order", h.GetAllOrder)
		r.ServeHTTP(rr, httpReq)

		var resp utils.Respond
		err = json.Unmarshal(rr.Body.Bytes(), &resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.EqualValues(t, http.StatusOK, rr.Code)
		assert.EqualValues(t, 200, resp.Status)
		assert.EqualValues(t, "success", resp.Message)
		assert.Equal(t, mockOrderList[1].UserName, "Agung")
		mockOrderUsecase.AssertExpectations(t)
	})
}
