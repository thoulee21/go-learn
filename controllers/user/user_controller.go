package user

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	domainErrors "github.com/thoulee21/go-learn/errors"
	"github.com/thoulee21/go-learn/models"
	"gorm.io/gorm"
)

type UserController struct {
	DB          *gorm.DB
	UserService models.IUserService
}

// @Summary		创建用户
// @Description	创建一个新的用户
// @Accept			json
// @Produce		json
// @Param			user	body		models.User	true	"用户信息"
// @Success		200		{object}	models.User	"成功"
// @Failure		400		{object}	string		"请求错误"
// @Failure		500		{object}	string		"内部错误"
// @Router			/user [post]
func (c *UserController) NewUser(ctx *gin.Context) {
	var request models.User
	if err := ctx.BindJSON(&request); err != nil {
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	userModel, err := c.UserService.Create(&request)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	userResponse := userModel
	ctx.JSON(http.StatusOK, userResponse)
}

// @Summary		获取所有用户
// @Description	获取所有用户的信息
// @Produce		json
// @Success		200	{array}		models.User	"成功"
// @Failure		500	{object}	string		"内部错误"
// @Router			/user [get]
func (c *UserController) GetAllUsers(ctx *gin.Context) {
	users, err := c.UserService.GetAll()
	if err != nil {
		appError := domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		_ = ctx.Error(appError)
		return
	}
	ctx.JSON(http.StatusOK, users)
}

// @Summary		获取用户信息
// @Description	根据用户ID获取用户信息
// @Produce		json
// @Param			id	path		int			true	"用户ID"
// @Success		200	{object}	models.User	"成功"
// @Failure		400	{object}	string		"请求错误"
// @Failure		404	{object}	string		"用户未找到"
// @Failure		500	{object}	string		"内部错误"
// @Router			/user/{id} [get]
func (c *UserController) GetUserByID(ctx *gin.Context) {
	userID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		appError := domainErrors.NewAppError(errors.New("user id is invalid"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	user, err := c.UserService.GetByID(uint(userID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			appError := domainErrors.NewAppError(errors.New("user not found"), domainErrors.NotFound)
			_ = ctx.Error(appError)
			return
		}
		appError := domainErrors.NewAppError(err, domainErrors.UnknownError)
		_ = ctx.Error(appError)
		return
	}
	ctx.JSON(http.StatusOK, user)
}

// @Summary		更新用户信息
// @Description	根据用户ID更新用户信息
// @Accept			json
// @Produce		json
// @Param			id		path		int			true	"用户ID"
// @Param			user	body		models.User	true	"用户信息"
// @Success		200		{object}	models.User	"成功"
// @Failure		400		{object}	string		"请求错误"
// @Failure		404		{object}	string		"用户未找到"
// @Failure		500		{object}	string		"内部错误"
// @Router			/user/{id} [put]
func (c *UserController) UpdateUser(ctx *gin.Context) {
	userID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		appError := domainErrors.NewAppError(errors.New("param id is necessary"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	var requestMap models.User
	err = ctx.BindJSON(&requestMap)
	if err != nil {
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	err = updateValidation(requestMap)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	userUpdated, err := c.UserService.Update(uint(userID), &requestMap)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, userUpdated)
}

// @Summary		删除用户
// @Description	根据用户ID删除用户
// @Produce		json
// @Param			id	path		int		true	"用户ID"
// @Success		200	{object}	string	"成功"
// @Failure		400	{object}	string	"请求错误"
// @Failure		404	{object}	string	"用户未找到"
// @Failure		500	{object}	string	"内部错误"
// @Router			/user/{id} [delete]
func (c *UserController) DeleteUser(ctx *gin.Context) {
	userID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		appError := domainErrors.NewAppError(errors.New("param id is necessary"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	err = c.UserService.Delete(uint(userID))
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "resource deleted successfully"})
}
