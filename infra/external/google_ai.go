package external

import (
	"context"
	"fmt"

	"google.golang.org/genai"
)

type GoogleAIClient struct {
	APIKey string
	client *genai.Client
}

func NewGoogleAIClient(APIKey string) (*GoogleAIClient, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  APIKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	return &GoogleAIClient{
		APIKey: APIKey,
		client: client,
	}, nil
}

func (c *GoogleAIClient) GenerateBuilds(prompt string) (string, error) {
	ctx := context.Background()

	result, err := c.client.Models.GenerateContent(
		ctx,
		"gemini-2.0-flash",
		genai.Text(prompt),
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate builds: %w", err)
	}

	return result.Text(), nil
}

func (c *GoogleAIClient) Close() error {
	return nil
}
