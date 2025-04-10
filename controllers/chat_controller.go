package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thoulee21/go-learn/models"
	"github.com/thoulee21/go-learn/services"
	"gorm.io/gorm"
)

type ChatController struct {
	DB        *gorm.DB
	AIService *services.AIService
}

//	@Summary		测试AI服务
//	@Description	测试AI服务是否正常工作
//	@Produce		json
//	@Param			msg	query		string	false	"测试消息"
//	@Success		200	{string}	string	"成功"
//	@Failure		500	{object}	string	"内部错误"
//	@Router			/test [get]
func (cc *ChatController) Test(c *gin.Context) {
	testMessage := c.Query("msg")

	responseText, err := cc.AIService.GenerateResponse(
		c.Request.Context(),
		[]services.ChatMessage{{
			Role:    "user",
			Content: testMessage,
		}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI服务错误: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, responseText)
}

//	@Summary		发送聊天消息
//	@Description	发送消息到AI并获取回复
//	@Accept			json
//	@Produce		json
//	@Param			request	body		models.ChatRequest	true	"聊天请求"
//	@Success		200		{object}	models.ChatResponse	"成功"
//	@Failure		400		{object}	string				"请求错误"
//	@Failure		500		{object}	string				"内部错误"
//	@Router			/chat [post]
func (cc *ChatController) Chat(c *gin.Context) {
	var request models.ChatRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 如果没有会话ID，创建一个新的
	if request.SessionID == "" {
		request.SessionID = uuid.New().String()
	}

	// 保存用户消息
	userMessage := models.ChatMessage{
		SessionID: request.SessionID,
		Role:      "user",
		Content:   request.Message,
	}
	if err := cc.DB.Create(&userMessage).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法保存消息"})
		return
	}

	// 获取历史消息（最多10条）
	var chatHistory []models.ChatMessage
	cc.DB.Where("session_id = ?", request.SessionID).Order("created_at desc").Limit(10).Find(&chatHistory)

	// 构造OpenAI消息格式
	var openAIMessages []services.ChatMessage
	for i := len(chatHistory) - 1; i >= 0; i-- {
		openAIMessages = append(openAIMessages, services.ChatMessage{
			Role:    chatHistory[i].Role,
			Content: chatHistory[i].Content,
		})
	}

	// 调用AI服务
	responseText, err := cc.AIService.GenerateResponse(c.Request.Context(), openAIMessages)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI服务错误: " + err.Error()})
		return
	}

	// 保存AI回复
	aiMessage := models.ChatMessage{
		SessionID: request.SessionID,
		Role:      "assistant",
		Content:   responseText,
	}
	if err := cc.DB.Create(&aiMessage).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法保存AI回复"})
		return
	}

	// 返回响应
	c.JSON(http.StatusOK, models.ChatResponse{
		SessionID: request.SessionID,
		Message:   responseText,
	})
}

//	@Summary		获取聊天历史
//	@Description	获取特定会话的聊天历史
//	@Produce		json
//	@Param			session_id	path		string				true	"会话ID"
//	@Success		200			{array}		models.ChatMessage	"成功"
//	@Failure		400			{object}	string				"请求错误"
//	@Failure		500			{object}	string				"内部错误"
//	@Router			/chat/history/{session_id} [get]
func (cc *ChatController) GetChatHistory(c *gin.Context) {
	sessionID := c.Param("session_id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "会话ID不能为空"})
		return
	}

	var messages []models.ChatMessage
	if err := cc.DB.Where("session_id = ?", sessionID).Order("created_at asc").Find(&messages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, messages)
}

//	@Summary		流式发送聊天消息
//	@Description	流式发送消息到AI并获取实时回复
//	@Accept			json
//	@Produce		text/event-stream
//	@Param			request	body		models.ChatRequest	true	"聊天请求"
//	@Success		200		{object}	string				"成功"
//	@Failure		400		{object}	string				"请求错误"
//	@Failure		500		{object}	string				"内部错误"
//	@Router			/chat/stream [post]
func (cc *ChatController) StreamChat(c *gin.Context) {
	var request models.ChatRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 如果没有会话ID，创建一个新的
	if request.SessionID == "" {
		request.SessionID = uuid.New().String()
	}

	// 保存用户消息
	userMessage := models.ChatMessage{
		SessionID: request.SessionID,
		Role:      "user",
		Content:   request.Message,
	}
	if err := cc.DB.Create(&userMessage).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法保存消息"})
		return
	}

	// 获取历史消息（最多10条）
	var chatHistory []models.ChatMessage
	cc.DB.Where("session_id = ?", request.SessionID).Order("created_at desc").Limit(10).Find(&chatHistory)

	// 构造OpenAI消息格式
	var openAIMessages []services.ChatMessage
	for i := len(chatHistory) - 1; i >= 0; i-- {
		openAIMessages = append(openAIMessages, services.ChatMessage{
			Role:    chatHistory[i].Role,
			Content: chatHistory[i].Content,
		})
	}

	// 设置SSE响应头
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")

	// 保存完整响应用于数据库存储
	fullResponse := ""

	// 创建回调函数处理流式响应
	callback := func(chunk string) {
		// 发送数据块
		c.Writer.Write([]byte("data: " + chunk + "\n\n"))
		c.Writer.Flush()
		fullResponse += chunk
	}

	// 调用AI服务的流式响应方法
	err := cc.AIService.GenerateStreamResponse(c.Request.Context(), openAIMessages, callback)
	if err != nil {
		// 尝试发送错误消息，但此时可能连接已关闭
		c.Writer.Write([]byte("data: {\"error\": \"" + err.Error() + "\"}\n\n"))
		c.Writer.Flush()
		return
	}

	// 发送结束信号
	c.Writer.Write([]byte("data: [DONE]\n\n"))
	c.Writer.Flush()

	// 保存AI回复到数据库
	aiMessage := models.ChatMessage{
		SessionID: request.SessionID,
		Role:      "assistant",
		Content:   fullResponse,
	}
	cc.DB.Create(&aiMessage)
}
