package product

import (
	"kanggo/pkg/entity/model"
	"kanggo/pkg/middleware"
	"kanggo/pkg/usecase/product"
	"kanggo/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

type ProductHandler struct {
	productUsecase product.ProductUsecase
}

func NewProductHandler(productUsecase product.ProductUsecase) *ProductHandler {
	return &ProductHandler{
		productUsecase: productUsecase,
	}
}

func (h *ProductHandler) Route(app *gin.Engine) {
	v1 := app.Group("api/v1")
	{
		{
			v1.POST("/product", middleware.RoleAdmin(), h.Insert)
			v1.GET("/product", middleware.RoleUser(), h.GetAll)
			v1.GET("/product/:id", middleware.RoleUser(), h.GetById)
			v1.PUT("/product/:id", middleware.RoleAdmin(), h.Update)
			v1.DELETE("/product/:id", middleware.RoleAdmin(), h.Delete)

		}
	}

}

func (h *ProductHandler) Insert(c *gin.Context) {
	validate = validator.New()
	product := model.ProductRequest{}
	ctx := c.Request.Context()

	if err := c.ShouldBindJSON(&product); err != nil {
		utils.Response(c, 400, err.Error(), nil)
		return
	}

	if err := validate.Struct(product); err != nil {
		utils.Response(c, 400, err.Error(), nil)
		return
	}

	if err := h.productUsecase.Insert(ctx, product); err != nil {
		utils.Response(c, 500, err.Error(), nil)
		return
	}

	utils.Response(c, 201, "success insert product", nil)
}

func (h *ProductHandler) GetAll(c *gin.Context) {
	ctx := c.Request.Context()

	res, err := h.productUsecase.GetAll(ctx)
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

func (h *ProductHandler) GetById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	ctx := c.Request.Context()
	res, err := h.productUsecase.GetById(ctx, int64(id))
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

func (h *ProductHandler) Update(c *gin.Context) {
	validate = validator.New()
	id, _ := strconv.Atoi(c.Param("id"))
	product := model.ProductRequest{}
	ctx := c.Request.Context()

	if err := c.ShouldBind(&product); err != nil {
		utils.Response(c, 400, err.Error(), nil)
		return
	}

	if err := validate.Struct(product); err != nil {
		utils.Response(c, 400, err.Error(), nil)
		return
	}

	if err := h.productUsecase.Update(ctx, uint(id), product); err != nil {
		if err.Error() == "data not found" {
			utils.Response(c, 404, err.Error(), nil)
			return
		}
		utils.Response(c, 500, err.Error(), nil)
		return
	}

	utils.Response(c, 200, "success update product", nil)
}

func (h *ProductHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	ctx := c.Request.Context()

	if err := h.productUsecase.Delete(ctx, int64(id)); err != nil {
		if err.Error() == "data not found" {
			utils.Response(c, 404, err.Error(), nil)
			return
		}
		utils.Response(c, 500, err.Error(), nil)
		return
	}

	utils.Response(c, 200, "success delete product", nil)
}
