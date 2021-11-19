package user

import (
	"context"
	"kanggo/pkg/entity/model"
	"kanggo/pkg/usecase/user"
	"kanggo/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

type UserHandler struct {
	userUsecase user.UserUsecase
}

func NewUserhandler(userUsecase user.UserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
	}
}

func (h *UserHandler) Route(app *gin.Engine) {
	v1 := app.Group("api/v1")
	{
		v1.POST("/register", h.Insert)
		v1.POST("/login", h.Login)
	}
}

func (h *UserHandler) Insert(c *gin.Context) {
	validate = validator.New()
	user := model.RegisterRequest{}

	err := c.ShouldBind(&user)
	if err != nil {
		utils.Response(c, 400, err.Error(), nil)
		return
	}

	err = validate.Struct(user)
	if err != nil {
		utils.Response(c, 400, err.Error(), nil)
		return
	}

	ctx := c.Request.Context()
	err = h.userUsecase.Insert(ctx, user)
	if err != nil {
		utils.Response(c, 500, err.Error(), nil)
		return
	}

	utils.Response(c, 201, "success insert user", nil)

}

func (h *UserHandler) Login(c *gin.Context) {
	ctx := c.Request.Context()
	validate = validator.New()
	login := model.LoginRequest{}
	val := model.ValidateResponse{}

	if err := c.ShouldBind(&login); err != nil {
		utils.Response(c, 400, err.Error(), nil)
		return
	}

	if err := validate.Struct(login); err != nil {
		utils.Response(c, 400, err.Error(), nil)
		return
	}

	var admin bool
	switch {
	case login.Email == "admin@gmail.com" && login.Password == "admin123":
		admin = true
	default:
		res, err := h.ValidateUser(ctx, login.Email, login.Password)
		if err != nil {
			utils.Response(c, 500, err.Error(), nil)
			return
		}
		if !res.Status {
			utils.Response(c, 400, "invalid password or username", nil)
			return
		}
		val = *res
	}

	expired := time.Now().Local().Add(time.Minute * time.Duration(300)).Unix()
	token, err := utils.GenerateToken(int64(val.Id), expired, admin)
	if err != nil {
		utils.Response(c, 400, err.Error(), nil)
	}

	result := &model.LoginResponse{
		Token:   token,
		Expired: expired,
	}

	utils.Response(c, 200, "success", result)
}

func (h *UserHandler) ValidateUser(ctx context.Context, email, pass string) (*model.ValidateResponse, error) {

	res, err := h.userUsecase.GetByEmail(ctx, email)
	if err != nil {
		return &model.ValidateResponse{UserResponse: *res, Status: false}, err
	}

	if !utils.CheckPasswordHash(pass, res.Password) {
		return &model.ValidateResponse{UserResponse: *res, Status: false}, nil
	}

	return &model.ValidateResponse{UserResponse: *res, Status: true}, nil
}
