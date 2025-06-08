package matching

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/Luzin7/pcideal-be/internal/core/models"
	"github.com/Luzin7/pcideal-be/internal/core/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

func (pm *PartMatchingService) FindParts(productName string, productType string, productBrand string) ([]models.Part, error) {
	if len(productName) == 0 {
		return nil, nil
	}

	wordsInProductName := strings.Fields(strings.ToLower(productName))
	var wordsToFilter []bson.M
	for _, w := range wordsInProductName {
		wordsToFilter = append(wordsToFilter, bson.M{
			"model": bson.M{
				"$regex":   w,
				"$options": "i",
			},
		})
	}

	filter := bson.M{
		"type": strings.ToUpper(productType),
		"$and": wordsToFilter,
	}

	log.Printf("productName: %s", productName)
	log.Printf("filter: %+v", filter)

	cursor, err := pm.collection.Find(context.Background(), filter)
	if err != nil {
		log.Printf("Erro ao executar a query: %v", err)
		return nil, err
	}
	defer cursor.Close(context.Background())

	var parts []models.Part
	if err = cursor.All(context.Background(), &parts); err != nil {
		log.Printf("Erro ao decodificar os resultados: %v", err)
		return nil, err
	}

	if len(parts) <= 0 {
		log.Printf("Nenhum resultado encontrado para %s %s", productType, productName)
		return nil, fmt.Errorf("nao foi encontrado %s %s", productType, productName)
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
