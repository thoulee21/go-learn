package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/thoulee21/go-learn/models"
	"gorm.io/gorm"
)

type TaskController struct {
	DB *gorm.DB
}

//	@Summary		获取所有任务
//	@Description	获取所有任务
//	@Produce		json
//	@Success		200	{array}		models.Task	"成功"
//	@Failure		500	{object}	string		"内部错误"
//	@Router			/api/v1/tasks [get]
func (tc *TaskController) GetTasks(c *gin.Context) {
	var tasks []models.Task
	if err := tc.DB.Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

//	@Summary		创建任务
//	@Description	创建任务
//	@Accept			json
//	@Produce		json
//	@Param			task	body		models.Task	true	"任务"
//	@Success		201		{object}	models.Task	"成功"
//	@Failure		400		{object}	string		"请求错误"
//	@Failure		500		{object}	string		"内部错误"
//	@Router			/api/v1/tasks [post]
func (tc *TaskController) CreateTask(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := tc.DB.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, task)
}

//	@Summary		获取任务
//	@Description	获取任务
//	@Produce		json
//	@Param			id	path		int			true	"任务ID"
//	@Success		200	{object}	models.Task	"成功"
//	@Failure		400	{object}	string		"请求错误"
//	@Failure		404	{object}	string		"任务不存在"
//	@Router			/api/v1/tasks/{id} [get]
func (tc *TaskController) GetTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task id"})
		return
	}

	var task models.Task
	if err := tc.DB.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	c.JSON(http.StatusOK, task)
}

//	@Summary		更新任务
//	@Description	更新任务
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int			true	"任务ID"
//	@Param			task	body		models.Task	true	"任务"
//	@Success		200		{object}	models.Task	"成功"
//	@Failure		400		{object}	string		"请求错误"
//	@Failure		404		{object}	string		"任务不存在"
//	@Failure		500		{object}	string		"内部错误"
//	@Router			/api/v1/tasks/{id} [put]
func (tc *TaskController) UpdateTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task id"})
		return
	}

	var task models.Task
	if err := tc.DB.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := tc.DB.Save(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, task)
}

//	@Summary		删除任务
//	@Description	删除任务
//	@Produce		json
//	@Param			id	path		int		true	"任务ID"
//	@Success		204	{object}	string	"成功"
//	@Failure		400	{object}	string	"请求错误"
//	@Failure		404	{object}	string	"任务不存在"
//	@Failure		500	{object}	string	"内部错误"
//	@Router			/api/v1/tasks/{id} [delete]
func (tc *TaskController) DeleteTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task id"})
		return
	}

	var task models.Task
	if err := tc.DB.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	if err := tc.DB.Delete(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
