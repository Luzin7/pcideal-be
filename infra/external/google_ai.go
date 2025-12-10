package external

import (
	"context"
	"fmt"
	"strings"

	"github.com/Luzin7/pcideal-be/internal/domain/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/genai"
)

type GoogleAIClient struct {
	APIKey      string
	Client      *genai.Client
	BasePrompts *BasePrompts
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

func NewGoogleAIClient(APIKey string, db *mongo.Database) (*GoogleAIClient, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  APIKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	basePrompts := NewBasePrompts(db)

	return &GoogleAIClient{
		APIKey:      APIKey,
		Client:      client,
		BasePrompts: basePrompts,
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

	if len(input) > 200 {
		input = input[:200]
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

	userPreference := fmt.Sprintf("customer has %s preference for their %s", input, typeOfPreference)
	return userPreference, nil
}

func (bp *BasePrompts) getBasePrompt() (string, error) {
	ctx := context.Background()

	var promptDoc PromptDocument
	err := bp.Collection.FindOne(ctx, bson.M{"category": "pc_builder"}).Decode(&promptDoc)
	if err != nil {
		return "", fmt.Errorf("failed to find base prompt: %w", err)
	}

	return promptDoc.Content, nil
}

func (c *GoogleAIClient) BuildComputerPrompt(usageType string, cpuPreference string, gpuPreference string, budget int64) (string, error) {
	basePrompt, err := c.BasePrompts.getBasePrompt()
	if err != nil {
		return "", err
	}

	usageType = sanitizeInput(usageType)

	cpuPreference, err = validatePreference(cpuPreference, "cpu")
	if err != nil {
		cpuPreference = "the client have no preferences of CPU brand"
		fmt.Printf("Aviso: %v. Usando valor padrão.\n", err)
	}

	gpuPreference, err = validatePreference(gpuPreference, "gpu")
	if err != nil {
		gpuPreference = "the client have no preferences of GPU brand"
		fmt.Printf("Aviso: %v. Usando valor padrão.\n", err)
	}

	fullPrompt := fmt.Sprintf(`%s

The client wants a PC that will be used for %s:

%s, %s, his budget is %d BRL.`, basePrompt, usageType, cpuPreference, gpuPreference, budget)

	return fullPrompt, nil
}

// TODO: Descomentar quando entity.AIBuildResponse for criada
// func CleanAndParseGeminiResponse(raw string) (*entity.AIBuildResponse, error) {
// 	re := regexp.MustCompile("(?s)```json\\n(.*?)\\n```")
// 	match := re.FindStringSubmatch(raw)

// 	var cleaned string
// 	if len(match) > 1 {
// 		cleaned = match[1]
// 	} else {
// 		cleaned = raw
// 	}

// 	cleaned = strings.ReplaceAll(cleaned, "\\n", "")
// 	cleaned = strings.ReplaceAll(cleaned, "\\\"", "\"")
// 	cleaned = strings.ReplaceAll(cleaned, "\\\\", "\\")

// 	var result *entity.AIBuildResponse
// 	err := json.Unmarshal([]byte(cleaned), &result)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to unmarshal cleaned response: %w", err)
// 	}

// 	return result, nil
// }

// func (c *GoogleAIClient) GenerateBuilds(prompt string) (*entity.AIBuildResponse, error) {
// 	ctx := context.Background()

// 	rawResponse, err := c.Client.Models.GenerateContent(
// 		ctx,
// 		"gemini-2.0-flash",
// 		genai.Text(prompt),
// 		nil,
// 	)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to generate builds: %w", err)
// 	}

// 	cleanedJSON, err := CleanAndParseGeminiResponse(rawResponse.Text())
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to clean raw response: %s", err)
// 	}

// 	return cleanedJSON, nil
// }

func (c *GoogleAIClient) GeneratePcBuildAnalysis(ctx context.Context, part *entity.Part) (string, error) {
	// TODO: Implementar análise de build via IA
	return "", fmt.Errorf("not implemented yet")
}

func (c *GoogleAIClient) Close() error {
	return nil
}
