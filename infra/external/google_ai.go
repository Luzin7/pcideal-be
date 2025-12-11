package external

import (
	"context"
	"fmt"

	"github.com/Luzin7/pcideal-be/infra/http/presenters"
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

func (c *GoogleAIClient) getBasePrompt(ctx context.Context, category string) (string, error) {
	var promptDoc PromptDocument
	err := c.BasePrompts.Collection.FindOne(ctx, map[string]string{"category": category}).Decode(&promptDoc)
	if err != nil {
		return "", fmt.Errorf("failed to get base prompt: %w", err)
	}
	return promptDoc.Content, nil
}

func (c *GoogleAIClient) GeneratePcBuildAnalysis(ctx context.Context, build *presenters.RecommendationBuild) (string, error) {
	// Validar build
	if build == nil {
		return "", fmt.Errorf("build is nil")
	}
	if build.Parts.CPU == nil || build.Parts.GPU == nil || build.Parts.RAM == nil ||
		build.Parts.PrimaryStorage == nil || build.Parts.PSU == nil {
		return "", fmt.Errorf("one or more build parts are nil")
	}

	basePrompt, err := c.getBasePrompt(ctx, "build_analysis")
	if err != nil {
		return "", fmt.Errorf("failed to get base prompt: %w", err)
	}

	budget := float64(build.Budget) / 100
	buildValue := float64(build.BuildValue) / 100
	remainingBudget := budget - buildValue

	prompt := fmt.Sprintf(`%s

## Build Configuration

**Build Type:** %s
**Budget:** R$ %.2f | **Build Cost:** R$ %.2f | **Remaining:** R$ %.2f

**Components:**
- CPU: %s %s (Performance: %d/10)
- GPU: %s %s (%dGB VRAM, Performance: %d/10)
- RAM: %dGB @ %dMHz
- Storage: %dGB SSD
- PSU: %dW 80+ %s

Provide a brief, user-friendly summary in Portuguese (PT-BR) explaining what this PC can do.`,
		basePrompt,
		build.BuildType,
		budget,
		buildValue,
		remainingBudget,
		// CPU
		build.Parts.CPU.Brand,
		build.Parts.CPU.Model,
		build.Parts.CPU.Specs.PerformanceScore,
		// GPU
		build.Parts.GPU.Brand,
		build.Parts.GPU.Model,
		build.Parts.GPU.Specs.VramGB,
		build.Parts.GPU.Specs.PerformanceScore,
		// RAM
		build.Parts.RAM.Specs.CapacityGB,
		build.Parts.RAM.Specs.MemorySpeedMHz,
		// Storage
		build.Parts.PrimaryStorage.Specs.CapacityGB,
		// PSU
		build.Parts.PSU.Specs.Wattage,
		c.getEfficiencyRatingName(build.Parts.PSU.Specs.EfficiencyRating),
	)

	config := &genai.GenerateContentConfig{
		Temperature:     genai.Ptr(float32(0.7)),
		MaxOutputTokens: 500,
	}

	resp, err := c.Client.Models.GenerateContent(ctx, "gemini-2.0-flash", []*genai.Content{
		{
			Role: "user",
			Parts: []*genai.Part{
				{Text: prompt},
			},
		},
	}, config)
	if err != nil {
		return "", fmt.Errorf("google AI API error: %w", err)
	}

	if len(resp.Candidates) == 0 {
		return "", fmt.Errorf("no candidates in response")
	}
	if resp.Candidates[0].Content == nil || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no content parts in response")
	}

	analysis := resp.Candidates[0].Content.Parts[0].Text
	return analysis, nil
}

func (c *GoogleAIClient) getEfficiencyRatingName(rating int8) string {
	switch rating {
	case 1:
		return "Bronze"
	case 2:
		return "Silver"
	case 3:
		return "Gold"
	case 4:
		return "Platinum"
	case 5:
		return "Titanium"
	default:
		return "Standard"
	}
}

func (c *GoogleAIClient) Close() error {
	return nil
}
