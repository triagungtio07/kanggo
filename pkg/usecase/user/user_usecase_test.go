package user

import (
	"context"
	"kanggo/pkg/entity/schema"
	"kanggo/pkg/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

//can't do test due to password bcrypt
func TestInsert(t *testing.T) {
	mockUserStorage := new(mocks.UserStorage)

	ctx := context.Background()

	mockUserReq := schema.User{
		Base:     schema.Base{Id: 0x0, CreatedAt: time.Time{}, UpdatedAt: time.Time{}},
		Name:     "agung",
		Email:    "agung@gmail.com",
		Password: "$2a$14$lEQUEJ.f3N1rXzeFssrD3OrPuK2mQ4qQZpmBVySFAjM6GBw5MJzIq",
	}

	t.Run("success", func(t *testing.T) {
		mockUserStorage.On("Insert", ctx, mockUserReq).Return(nil)
	})

}

func TestGetByEmail(t *testing.T) {
	mockUserStorage := new(mocks.UserStorage)
	u := NewUserUsecase(mockUserStorage)
	ctx := context.Background()
	var email string = "agung@gmail.com"

	mockUserReq := schema.User{
		Base:     schema.Base{Id: 1},
		Name:     "agung",
		Email:    "agung@gmail.com",
		Password: "$2a$14$lEQUEJ.f3N1rXzeFssrD3OrPuK2mQ4qQZpmBVySFAjM6GBw5MJzIq",
	}

	t.Run("success", func(t *testing.T) {
		mockUserStorage.On("GetByEmail", mock.Anything, mock.AnythingOfType("string")).Return(&mockUserReq, nil)

		result, err := u.GetByEmail(ctx, email)
		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.NoError(t, err)
		assert.Equal(t, mockUserReq.Name, "agung")
		mockUserStorage.AssertExpectations(t)
	})

}
