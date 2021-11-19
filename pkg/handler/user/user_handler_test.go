package user

import (
	"bytes"
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

func TestRegister(t *testing.T) {
	mockUserUsecase := new(mocks.UserUsecase)

	t.Run("success", func(t *testing.T) {
		mockRequest := model.RegisterRequest{
			Name:     "Agung",
			Email:    "agung@gmail.com",
			Password: "agung123",
		}

		mockUserUsecase.On("Insert", mock.Anything, mockRequest).Return(nil)

		body, err := json.Marshal(mockRequest)
		assert.Nil(t, err)

		httpReq, err := http.NewRequest(http.MethodPost, "/api/v1/register", bytes.NewReader(body))
		httpReq.Header.Set("X-Custom-Header", "myvalue")
		httpReq.Header.Set("Content-Type", "application/json")
		assert.Nil(t, err)

		r := gin.Default()
		rr := httptest.NewRecorder()

		h := NewUserhandler(mockUserUsecase)

		r.POST("/api/v1/register", h.Insert)

		r.ServeHTTP(rr, httpReq)

		var resp utils.Respond
		err = json.Unmarshal(rr.Body.Bytes(), &resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.EqualValues(t, http.StatusCreated, rr.Code)
		assert.EqualValues(t, 201, resp.Status)
		assert.EqualValues(t, "success insert user", resp.Message)
		mockUserUsecase.AssertExpectations(t)
	})
}

func TestLogin(t *testing.T) {
	mockUserUsecase := new(mocks.UserUsecase)

	t.Run("success", func(t *testing.T) {
		mockRequest := model.LoginRequest{
			Email:    "agung@gmail.com",
			Password: "agung123",
		}

		mockResponse := model.UserResponse{
			Id:       1,
			Name:     "agung",
			Email:    "agung@gmail.com",
			Password: "agung123",
		}

		mockUserUsecase.On("GetByEmail", mock.Anything, mockRequest.Email).Return(&mockResponse, nil)

		body, err := json.Marshal(mockRequest)
		assert.Nil(t, err)

		httpReq, err := http.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewReader(body))
		httpReq.Header.Set("X-Custom-Header", "myvalue")
		httpReq.Header.Set("Content-Type", "application/json")
		assert.Nil(t, err)

		r := gin.Default()
		rr := httptest.NewRecorder()

		h := NewUserhandler(mockUserUsecase)

		r.POST("/api/v1/login", h.Login)
		r.ServeHTTP(rr, httpReq)

		var resp utils.Respond
		err = json.Unmarshal(rr.Body.Bytes(), &resp)

		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.EqualValues(t, http.StatusBadRequest, rr.Code)
		assert.EqualValues(t, 400, resp.Status)
		assert.EqualValues(t, "invalid password or username", resp.Message)
		mockUserUsecase.AssertExpectations(t)
	})
}
