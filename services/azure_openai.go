package services

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

type AIService struct {
	client          *azopenai.Client
	deploymentName  string
	maxTokens       int32
	temperature     float32
	topP            float32
	freqPenalty     float32
	presencePenalty float32
	stop            []string
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func NewAIService() (*AIService, error) {
	azureOpenAIEndpoint := os.Getenv("AZURE_OPENAI_ENDPOINT")
	azureOpenAIKey := os.Getenv("AZURE_OPENAI_API_KEY")
	deploymentName := os.Getenv("AZURE_OPENAI_DEPLOYMENT_NAME")

	if azureOpenAIEndpoint == "" || azureOpenAIKey == "" || deploymentName == "" {
		return nil, errors.New(
			"AZURE_OPENAI_ENDPOINT, AZURE_OPENAI_API_KEY, and AZURE_OPENAI_DEPLOYMENT_NAME environment variables must be set",
		)
	}

	cred := azcore.NewKeyCredential(azureOpenAIKey)
	client, err := azopenai.NewClientWithKeyCredential(azureOpenAIEndpoint, cred, nil)
	if err != nil {
		return nil, err
	}

	return &AIService{
		client:          client,
		deploymentName:  deploymentName,
		maxTokens:       800,
		temperature:     0.7,
		topP:            0.95,
		freqPenalty:     0,
		presencePenalty: 0,
		stop:            []string{},
	}, nil
}

func (s *AIService) convertToAzureMessages(messages []ChatMessage) ([]azopenai.ChatRequestMessageClassification, error) {
	azMessages := make([]azopenai.ChatRequestMessageClassification, 0, len(messages))

	for _, msg := range messages {
		switch msg.Role {
		case "system":
			azMessages = append(azMessages, &azopenai.ChatRequestSystemMessage{
				Content: azopenai.NewChatRequestSystemMessageContent(msg.Content),
			})
		case "user":
			azMessages = append(azMessages, &azopenai.ChatRequestUserMessage{
				Content: azopenai.NewChatRequestUserMessageContent(msg.Content),
			})
		case "assistant":
			azMessages = append(azMessages, &azopenai.ChatRequestAssistantMessage{
				Content: azopenai.NewChatRequestAssistantMessageContent(msg.Content),
			})
		default:
			return nil, fmt.Errorf("unsupported role: %s", msg.Role)
		}
	}

	return azMessages, nil
}

func (s *AIService) GenerateResponse(ctx context.Context, messages []ChatMessage) (string, error) {
	// 将我们的消息格式转换为 Azure SDK 的消息格式
	azMessages, err := s.convertToAzureMessages(messages)
	if err != nil {
		return "", err
	}

	resp, err := s.client.GetChatCompletions(ctx, azopenai.ChatCompletionsOptions{
		Messages:         azMessages,
		DeploymentName:   &s.deploymentName,
		MaxTokens:        &s.maxTokens,
		Temperature:      &s.temperature,
		TopP:             &s.topP,
		FrequencyPenalty: &s.freqPenalty,
		PresencePenalty:  &s.presencePenalty,
		Stop:             s.stop,
	}, nil)

	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", errors.New("no response generated")
	}

	return *resp.Choices[0].Message.Content, nil
}

func (s *AIService) GenerateStreamResponse(
	ctx context.Context,
	messages []ChatMessage,
	callback func(chunk string),
) error {
	// 将我们的消息格式转换为 Azure SDK 的消息格式
	azMessages, err := s.convertToAzureMessages(messages)
	if err != nil {
		return err
	}

	// 创建流式请求
	streamResp, err := s.client.GetChatCompletionsStream(
		ctx,
		azopenai.ChatCompletionsStreamOptions{
			Messages:         azMessages,
			DeploymentName:   &s.deploymentName,
			MaxTokens:        &s.maxTokens,
			Temperature:      &s.temperature,
			TopP:             &s.topP,
			FrequencyPenalty: &s.freqPenalty,
			PresencePenalty:  &s.presencePenalty,
			Stop:             s.stop,
		},
		nil,
	)
	if err != nil {
		return err
	}

	// 处理流式响应
	for {
		resp, err := streamResp.ChatCompletionsStream.Read()
		if err != nil {
			// 检查是否是结束的错误
			if errors.Is(err, io.EOF) {
				break // 流已结束，正常退出
			}
			return err
		}

		// 检查并处理响应内容
		if len(resp.Choices) > 0 && resp.Choices[0].Delta != nil && resp.Choices[0].Delta.Content != nil {
			callback(*resp.Choices[0].Delta.Content)
		}
	}

	return nil
}
