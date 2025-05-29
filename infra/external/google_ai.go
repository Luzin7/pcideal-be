package external

import (
	"context"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/genai"
)

type GoogleAIClient struct {
	APIKey string
	Client *genai.Client
}

type BasePrompts struct {
	Collection *mongo.Collection
}

type PromptDocument struct {
	Category string `bson:"category"`
	Content  string `bson:"content"`
}

func NewBasePrompts(client *mongo.Database) *BasePrompts {
	collection := client.Collection("base_prompts")

	return &BasePrompts{
		Collection: collection,
	}
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
		Client: client,
	}, nil
}

var allowedCPUPreferences = []string{"amd", "intel"}
var allowedGPUPreferences = []string{"nvidia", "amd"}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if strings.EqualFold(s, item) {
			return true
		}
	}
	return false
}

func sanitizeInput(input string) string {
	input = strings.TrimSpace(input)
	input = strings.ReplaceAll(input, "\"", "")
	input = strings.ReplaceAll(input, "'", "")

	if len(input) > 2000 {
		input = input[:2000]
	}

	return input
}

func validatePreference(input string, typeOfPreference string) (string, error) {
	input = sanitizeInput(input)

	var allowedValues []string

	switch typeOfPreference {
	case "gpu":
		allowedValues = allowedGPUPreferences
	case "cpu":
		allowedValues = allowedCPUPreferences
	default:
		return "", fmt.Errorf("tipo de preferência inválido: %s", typeOfPreference)
	}

	if !contains(allowedValues, input) {
		return "", fmt.Errorf("preferência de %s inválida: %s", typeOfPreference, input)
	}

	return input, nil
}

func (bp *BasePrompts) GetBasePrompt() (string, error) {
	ctx := context.Background()

	var promptDoc PromptDocument
	err := bp.Collection.FindOne(ctx, bson.M{"category": "pc_builder"}).Decode(&promptDoc)
	if err != nil {
		return "", fmt.Errorf("failed to find base prompt: %w", err)
	}

	return promptDoc.Content, nil
}

func (bp *BasePrompts) BuildComputerPrompt(usageType, cpuPreference string, gpuPreference string, budget int64) (string, error) {
	basePrompt, err := bp.GetBasePrompt()
	if err != nil {
		return "", err
	}

	usageType = sanitizeInput(usageType)

	cpuPreference, err = validatePreference(cpuPreference, "cpu")
	if err != nil {
		cpuPreference = "Não tenho preferência"
		fmt.Printf("Aviso: %v. Usando valor padrão.\n", err)
	}

	gpuPreference, err = validatePreference(gpuPreference, "gpu")
	if err != nil {
		gpuPreference = "Não tenho preferência"
		fmt.Printf("Aviso: %v. Usando valor padrão.\n", err)
	}

	fullPrompt := fmt.Sprintf(`%s

Um cliente quer um computador que será utilizado para %s:

%s, %s, o orçamento dele é de %d reais.`, basePrompt, usageType, cpuPreference, gpuPreference, budget)

	return fullPrompt, nil
}

func (c *GoogleAIClient) GenerateBuilds(prompt string) (string, error) {
	ctx := context.Background()

	result, err := c.Client.Models.GenerateContent(
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
