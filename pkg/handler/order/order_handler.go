package order

import (
	"kanggo/pkg/entity/model"
	"kanggo/pkg/middleware"
	"kanggo/pkg/usecase/order"
	"kanggo/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

type OrderHandler struct {
	orderUsecase order.OrderUsecase
}

func NewOrderHandler(orderUsecase order.OrderUsecase) *OrderHandler {
	return &OrderHandler{
		orderUsecase: orderUsecase,
	}
}

func (o *OrderHandler) Route(app *gin.Engine) {
	v1 := app.Group("api/v1")
	{
		{
			v1.POST("/order", middleware.RoleUser(), o.InsertOrder)
			v1.GET("/order", middleware.RoleAdmin(), o.GetAllOrder)
			v1.GET("/order/user", middleware.RoleUser(), o.GetAllOrderPerUser)
			v1.GET("/order/:id", middleware.RoleUser(), o.GetOrderById)
			v1.PUT("/payment", middleware.RoleUser(), o.UpdatePayment)
		}
	}

}

func (o *OrderHandler) InsertOrder(c *gin.Context) {
	validate = validator.New()
	order := model.OrderRequest{}
	userId := c.MustGet("user_id").(uint64)
	order.UserId = int64(userId)
	ctx := c.Request.Context()

	err := c.ShouldBindJSON(&order)
	if err != nil {
		utils.Response(c, 400, err.Error(), nil)
		return
	}

	err = validate.Struct(order)
	if err != nil {
		utils.Response(c, 400, err.Error(), nil)
		return
	}

	err = o.orderUsecase.InsertOrder(ctx, order)
	if err != nil {
		if err.Error() == "not enough product quantity" {
			utils.Response(c, 400, "not enough product quantity", nil)
			return
		}
		utils.Response(c, 500, err.Error(), nil)
		return
	}

	utils.Response(c, 201, "success insert order", nil)
}

func (o *OrderHandler) GetAllOrder(c *gin.Context) {
	ctx := c.Request.Context()

	res, err := o.orderUsecase.GetAllOrder(ctx)
	if err != nil {
		if err.Error() == "data not found" {
			utils.Response(c, 404, err.Error(), nil)
			return
		}
		utils.Response(c, 500, err.Error(), nil)
		return
	}

	utils.Response(c, 200, "success", res)
}

func (o *OrderHandler) GetAllOrderPerUser(c *gin.Context) {
	ctx := c.Request.Context()

	userId := c.MustGet("user_id").(uint64)

	res, err := o.orderUsecase.GetAllOrderPerUser(ctx, userId)
	if err != nil {
		if err.Error() == "data not found" {
			utils.Response(c, 404, err.Error(), nil)
			return
		}
		utils.Response(c, 500, err.Error(), nil)
		return
	}

	utils.Response(c, 200, "success", res)
}

func (o *OrderHandler) GetOrderById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userId := c.MustGet("user_id").(uint64)
	ctx := c.Request.Context()
	res, err := o.orderUsecase.GetOrderById(ctx, int64(id), userId)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			utils.Response(c, 404, "data not found", nil)
			return
		}
		utils.Response(c, 500, err.Error(), nil)
		return
	}

	utils.Response(c, 200, "success", res)
}

func (o *OrderHandler) UpdatePayment(c *gin.Context) {
	validate = validator.New()
	payment := model.PaymentRequest{}
	userId := c.MustGet("user_id").(uint64)
	payment.UserId = int64(userId)
	ctx := c.Request.Context()

	err := c.ShouldBindJSON(&payment)
	if err != nil {
		utils.Response(c, 400, err.Error(), nil)
		return
	}

	err = validate.Struct(payment)
	if err != nil {
		utils.Response(c, 400, err.Error(), nil)
		return
	}

	err = o.orderUsecase.UpdatePayment(ctx, payment)
	if err != nil {
		if err.Error() == "payment amount does not match" {
			utils.Response(c, 400, "payment amount does not match", nil)
			return
		}
		if err.Error() == "record not found; record not found" {
			utils.Response(c, 404, "data not found", nil)
			return
		}
		utils.Response(c, 500, err.Error(), nil)
		return
	}

	utils.Response(c, 200, "success update payment", nil)
}
