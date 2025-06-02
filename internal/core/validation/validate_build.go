package validation

import (
	"context"
	"fmt"
	"math"
	"regexp"
	"strings"
	"unicode"

	"github.com/Luzin7/pcideal-be/internal/core/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ValidateBuild struct {
	collection *mongo.Collection
}

func NewValidateBuild(client *mongo.Database) *ValidateBuild {
	collection := client.Collection("parts")

	return &ValidateBuild{
		collection,
	}
}

type spec struct {
	model string
}

func normalizeString(s string) string {
	s = strings.ToLower(s)

	var result strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsSpace(r) || unicode.IsDigit(r) {
			result.WriteRune(r)
		}
	}

	return strings.Join(strings.Fields(result.String()), " ")
}

func min(a, b, c int) int {
	return int(math.Min(float64(a), math.Min(float64(b), float64(c))))
}

func levenshtein(a, b string) int {
	if len(a) > len(b) {
		a, b = b, a
	}

	prev := make([]int, len(a)+1)
	curr := make([]int, len(a)+1)

	for i := 0; i <= len(a); i++ {
		prev[i] = i
	}

	for j := 1; j <= len(b); j++ {
		curr[0] = j

		for i := 1; i <= len(a); i++ {
			if a[i-1] == b[j-1] {
				curr[i] = prev[i-1]
			} else {
				curr[i] = min(
					prev[i]+1,
					curr[i-1]+1,
					prev[i-1]+1,
				)
			}
		}

		prev, curr = curr, prev
	}

	return prev[len(a)]
}

func dynamicRegex(productName string) string {
	escaped := regexp.QuoteMeta(strings.TrimSpace(productName))
	return fmt.Sprintf("(?i).*%s.*", escaped)
}

func (v *ValidateBuild) findParts(productName string, productType string) ([]models.Part, error) {
	if len(productName) == 0 {
		return nil, nil
	}

	regexPattern := dynamicRegex(productName)

	filter := bson.M{
		"type": productType,
		"name": bson.M{
			"$regex": regexPattern,
		},
	}

	cursor, err := v.collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var parts []models.Part
	if err = cursor.All(context.Background(), &parts); err != nil {
		return nil, err
	}

	return parts, nil
}

func (v *ValidateBuild) findBestMatch(targetName string, parts []models.Part) *models.Part {
	if len(parts) == 0 {
		return nil
	}

	normalizedTarget := normalizeString(targetName)
	bestMatch := &parts[0]
	bestDistance := levenshtein(normalizedTarget, normalizeString(parts[0].Model))

	for i := 1; i < len(parts); i++ {
		normalizedCandidate := normalizeString(parts[i].Model)
		distance := levenshtein(normalizedTarget, normalizedCandidate)

		if distance < bestDistance {
			bestDistance = distance
			bestMatch = &parts[i]
		}
	}

	threshold := len(normalizedTarget) / 2
	if bestDistance > threshold {
		return nil
	}

	return bestMatch
}

func (v *ValidateBuild) ValidateCPUAndMotherboard(cpu spec, mobo spec) bool {
	cpuParts, err := v.findParts(cpu.model, "cpu")
	if err != nil {
		return false
	}

	moboParts, err := v.findParts(mobo.model, "motherboard")
	if err != nil {
		return false
	}

	bestCPUMatch := v.findBestMatch(cpu.model, cpuParts)
	if bestCPUMatch == nil {
		return false
	}

	bestMoboMatch := v.findBestMatch(mobo.model, moboParts)
	if bestMoboMatch == nil {
		return false
	}

	if bestCPUMatch.Specs.Socket != bestMoboMatch.Specs.Socket {
		return false
	}

	return true
}
