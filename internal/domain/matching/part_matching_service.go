package matching

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/Luzin7/pcideal-be/internal/core/models"
	"github.com/Luzin7/pcideal-be/internal/core/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PartMatchingService struct {
	collection *mongo.Collection
}

func NewPartMatchingService(client *mongo.Database) *PartMatchingService {
	collection := client.Collection("parts")

	return &PartMatchingService{
		collection,
	}
}

func (pm *PartMatchingService) FindParts(productName, productType, productBrand string) ([]models.Part, error) {
	if productName == "" {
		return nil, fmt.Errorf("productName is empty")
	}
	if productType == "" {
		return nil, fmt.Errorf("productType is empty")
	}
	if productBrand == "" {
		return nil, fmt.Errorf("productBrand is empty")
	}

	lowerProductName := strings.ToLower(productName)
	escapedProductName := regexp.QuoteMeta(lowerProductName)

	words := strings.Fields(escapedProductName)
	var wordsToFilter []bson.M
	for _, w := range words {
		escapedWord := regexp.QuoteMeta(w)
		regexPattern := fmt.Sprintf("(?i)\\b.*%s.*\\b", escapedWord)
		wordsToFilter = append(wordsToFilter, bson.M{
			"model": bson.M{
				"$regex": regexPattern,
			},
		})
	}

	filter := bson.M{
		"type":  strings.ToUpper(productType),
		"brand": strings.ToUpper(productBrand),
		"$and":  wordsToFilter,
	}

	cursor, err := pm.collection.Find(context.Background(), filter, options.Find().SetLimit(20))
	if err != nil {
		log.Printf("Error executing AND regex query: %v", err)
		return nil, err
	}
	defer cursor.Close(context.Background())

	var parts []models.Part
	if err = cursor.All(context.Background(), &parts); err != nil {
		log.Printf("Error decoding AND regex results: %v", err)
		return nil, err
	}

	if len(parts) == 0 {
		return nil, fmt.Errorf("no results found for type=%s, brand=%s, name=%s", productType, productBrand, productName)
	}

	return parts, nil
}

func (pm *PartMatchingService) FindBestMatch(targetName string, parts []models.Part) *models.Part {
	if len(parts) == 0 {
		log.Printf("parts é vazio")
		return nil
	}

	log.Print(targetName)
	log.Print(parts)

	normalizedTarget := util.NormalizeString(targetName)
	bestMatch := &parts[0]
	bestDistance := util.Levenshtein(normalizedTarget, util.NormalizeString(parts[0].Model))

	for i := 1; i < len(parts); i++ {
		normalizedCandidate := util.NormalizeString(parts[i].Model)
		distance := util.Levenshtein(normalizedTarget, normalizedCandidate)

		if distance < bestDistance {
			bestDistance = distance
			bestMatch = &parts[i]
		}
	}

	// threshold := len(normalizedTarget) / 2
	// if bestDistance > threshold {
	// 	return nil
	// }

	log.Printf("melhor distancia %d", bestDistance)

	return bestMatch
}
