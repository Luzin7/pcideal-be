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

	lowerProductName := strings.ToLower(productName)
	words := strings.Fields(lowerProductName)

	andFilter := pm.buildFilter(productType, productBrand, words, true, true)
	parts, err := pm.executeQuery(andFilter)
	if err != nil {
		return nil, err
	}

	if len(parts) > 0 {
		log.Printf("found %s using first type search", productName)
		return parts, nil
	}

	orFilter := pm.buildFilter(productType, productBrand, words, false, true)
	parts, err = pm.executeQuery(orFilter)
	if err != nil {
		log.Printf("found %s using second type search", productName)
		return nil, err
	}

	if len(parts) > 0 {
		return parts, nil
	}

	noBrandFilter := pm.buildFilter(productType, "", words, false, false)
	parts, err = pm.executeQuery(noBrandFilter)
	if err != nil {
		return nil, err
	}

	if len(parts) > 0 {
		log.Printf("found %s using third type seach", productName)
		return parts, nil
	}

	return nil, fmt.Errorf("no results found for type=%s, brand=%s, name=%s", productType, productBrand, productName)
}

func (pm *PartMatchingService) buildFilter(productType, productBrand string, words []string, useAnd, useBrand bool) bson.M {
	var wordsFilter []bson.M
	op := "$or"
	if useAnd {
		op = "$and"
	}

	for _, w := range words {
		if len(w) <= 2 {
			continue
		}
		escapedWord := regexp.QuoteMeta(w)
		regexPattern := fmt.Sprintf("(?i).*%s.*", escapedWord)
		wordsFilter = append(wordsFilter, bson.M{
			"model": bson.M{"$regex": regexPattern},
		})
	}

	filter := bson.M{
		"type": op,
	}

	filter["type"] = strings.ToUpper(productType)
	filter[op] = wordsFilter

	if useBrand && productBrand != "" {
		filter["brand"] = strings.ToLower(productBrand)
	}

	return filter
}

func (pm *PartMatchingService) executeQuery(filter bson.M) ([]models.Part, error) {
	cursor, err := pm.collection.Find(context.Background(), filter, options.Find().SetLimit(200))
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil, err
	}
	defer cursor.Close(context.Background())

	var parts []models.Part
	if err = cursor.All(context.Background(), &parts); err != nil {
		log.Printf("Error decoding results: %v", err)
		return nil, err
	}

	return parts, nil
}

func (pm *PartMatchingService) FindBestMatch(targetName string, parts []models.Part) *models.Part {
	if len(parts) == 0 {
		log.Printf("parts é vazio")
		return nil
	}

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

	log.Printf("melhor distancia para %s foi de %d com a string %s", targetName, bestDistance, bestMatch.Model)

	return bestMatch
}
