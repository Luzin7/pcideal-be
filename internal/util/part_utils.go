package util

import (
	"time"

	"github.com/Luzin7/pcideal-be/internal/domain/entity"
)

func PartNeedToUpdate(part *entity.Part) bool {
	if part != nil && time.Since(part.UpdatedAt) >= 2*time.Hour {
		return true
	}
	return false
}
