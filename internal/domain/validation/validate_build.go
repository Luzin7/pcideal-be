package validation

import (
	"fmt"

	"slices"

	"github.com/Luzin7/pcideal-be/internal/core/models"
)

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

func ValidateCPUAndMotherboard(cpu *models.Part, mobo *models.Part) bool {

	if cpu.Specs.Socket == "" || mobo.Specs.Socket == "" {
		fmt.Println("Warning: validation is partial")
		return true
	}

	return validateCompatibility("chipset_socket", cpu.Specs.Chipset, mobo.Specs.Socket)
}
