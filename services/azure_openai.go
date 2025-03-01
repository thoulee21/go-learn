package services

import (
	"context"
	"errors"
	"fmt"
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

func (s *AIService) GenerateResponse(messages []ChatMessage) (string, error) {
	// 将我们的消息格式转换为 Azure SDK 的消息格式
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
			return "", fmt.Errorf("unsupported role: %s", msg.Role)
		}
	}

	resp, err := s.client.GetChatCompletions(context.TODO(), azopenai.ChatCompletionsOptions{
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
