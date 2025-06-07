package validation

import (
	"context"
	"fmt"
	"log"
	"math"
	"strings"
	"unicode"

	"slices"

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

func (v *ValidateBuild) findParts(productName string, productType string) ([]models.Part, error) {
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
		"$or":  wordsToFilter,
	}

	log.Printf("productName: %s", productName)
	log.Printf("filter: %+v", filter)

	cursor, err := v.collection.Find(context.Background(), filter)
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

	log.Printf("Resultados encontrados: %v", parts)

	return parts, nil
}

func (v *ValidateBuild) findBestMatch(targetName string, parts []models.Part) *models.Part {
	if len(parts) == 0 {
		log.Printf("parts é vazio")
		return nil
	}

	log.Print(targetName)
	log.Print(parts)

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

	// threshold := len(normalizedTarget) / 2
	// if bestDistance > threshold {
	// 	return nil
	// }

	log.Printf("melhor distancia %d", bestDistance)

	return bestMatch
}

type compatibilityMap struct {
	name    string
	mapping map[string][]string
}

var compatibilityMaps = []compatibilityMap{
	{
		name: "chipset_socket",
		mapping: map[string][]string{
			// Intel - 13th & 12th Gen (Raptor Lake/Alder Lake)
			"Z790": {"LGA 1700"},
			"Z690": {"LGA 1700"},
			"B760": {"LGA 1700"},
			"B660": {"LGA 1700"},
			"H770": {"LGA 1700"},
			"H670": {"LGA 1700"},
			"H610": {"LGA 1700"},

			// Intel - 11th Gen (Rocket Lake)
			"Z590": {"LGA 1200"},
			"B560": {"LGA 1200"},
			"H570": {"LGA 1200"},
			"H510": {"LGA 1200"},

			// AMD - AM5 (Ryzen 7000 Series)
			"X670E": {"AM5"},
			"X670":  {"AM5"},
			"B650E": {"AM5"},
			"B650":  {"AM5"},

			// AMD - AM4 (Ryzen 5000/3000 Series)
			"X570": {"AM4"},
			"B550": {"AM4"},
			"A520": {"AM4"},
			"X470": {"AM4"},
			"B450": {"AM4"},
		},
	},
}

// validateCompatibility checks if a given value is compatible with a specific key in a compatibility map.
// It takes three parameters:
//   - validationType: string that specifies which compatibility map to use
//   - key: string representing the key to check in the compatibility map
//   - value: string to validate against the compatible values
//
// The function returns true if:
//   - The value is found in the list of compatible values for the given key
//
// Returns false if the value is not in the list of compatible values for the given key.
func validateCompatibility(validationType string, key string, value string) bool {
	for _, cMap := range compatibilityMaps {
		if cMap.name == validationType {
			compatibleValues, exists := cMap.mapping[key]
			if !exists {
				fmt.Printf("Warning: Key %s not found to %s", key, validationType)
				return false
			}
			return slices.Contains(compatibleValues, value)
		}
	}
	return false
}

func (v *ValidateBuild) ValidateCPUAndMotherboard(cpu string, mobo string) bool {
	cpuParts, err := v.findParts(cpu, "cpu")
	if err != nil {
		log.Printf("erro ao cpu: %s", err)
		return false
	}

	moboParts, err := v.findParts(mobo, "motherboard")
	if err != nil {
		log.Printf("erro ao mobo: %s", err)
		return false
	}

	bestCPUMatch := v.findBestMatch(cpu, cpuParts)
	if bestCPUMatch == nil {
		log.Printf("erro ao cpu mathc: %s", err)
		return false
	}

	bestMoboMatch := v.findBestMatch(mobo, moboParts)
	if bestMoboMatch == nil {
		log.Printf("erro ao cpu mathc: %s", err)
		return false
	}

	if bestCPUMatch.Specs.Socket == "" || bestMoboMatch.Specs.Socket == "" {
		log.Printf("a busca foi de %s e foi encontrado %s, assim como %s encontrou %s", cpu, bestCPUMatch.Model, mobo, bestMoboMatch.Model)
		fmt.Println("Warning: validation is partial")
		return true
	}

	log.Printf("a busca foi de %s e foi encontrado %s, assim como %s encontrou %s", cpu, bestCPUMatch.Model, mobo, bestMoboMatch.Model)

	return validateCompatibility("chipset_socket", bestMoboMatch.Specs.Chipset, bestCPUMatch.Specs.Socket)
}
