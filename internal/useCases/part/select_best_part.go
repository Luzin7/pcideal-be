package part

import (
	"context"
	"log"

	"github.com/Luzin7/pcideal-be/internal/domain/entity"
	"github.com/Luzin7/pcideal-be/internal/dto"
	"github.com/Luzin7/pcideal-be/internal/util"
)

type partRanking struct {
	score    func(*entity.Part) int
	lowerTie func(*entity.Part) int
}

func performanceScore(part *entity.Part) int {
	return int(part.Specs.PerformanceScore)
}

func chooseBestPart(parts []*entity.Part, ranking partRanking) *entity.Part {
	var best *entity.Part

	for _, part := range parts {
		if best == nil || isBetterCandidate(part, best, ranking) {
			best = part
		}
	}

	return best
}

func isBetterCandidate(candidate *entity.Part, current *entity.Part, ranking partRanking) bool {
	candidateScore := ranking.score(candidate)
	currentScore := ranking.score(current)

	if candidateScore != currentScore {
		return candidateScore > currentScore
	}

	if ranking.lowerTie != nil {
		candidateTie := ranking.lowerTie(candidate)
		currentTie := ranking.lowerTie(current)

		if candidateTie != currentTie {
			return candidateTie < currentTie
		}
	}

	if candidate.PriceCents != current.PriceCents {
		return candidate.PriceCents < current.PriceCents
	}

	candidateCapacity := capacityOf(candidate)
	currentCapacity := capacityOf(current)

	if candidateCapacity != currentCapacity {
		return candidateCapacity > currentCapacity
	}

	return candidate.ID.Hex() < current.ID.Hex()
}

func capacityOf(part *entity.Part) int {
	if part.Specs.CapacityGB > 0 {
		return part.Specs.CapacityGB
	}
	if part.Specs.VramGB > 0 {
		return part.Specs.VramGB
	}
	return int(part.Specs.Wattage)
}

func queueStalePartsUpdate(uc *UpdatePartsUseCase, parts []*entity.Part) {
	var partsToUpdate []dto.ProductLinkToUpdate

	for _, part := range parts {
		if !util.PartNeedToUpdate(part) {
			continue
		}

		partsToUpdate = append(partsToUpdate, dto.ProductLinkToUpdate{
			ID:  part.ID.Hex(),
			Url: part.URL,
		})
	}

	if len(partsToUpdate) == 0 {
		return
	}

	go func() {
		if err := uc.Execute(context.Background(), partsToUpdate, "kabum"); err != nil {
			log.Printf("failed to update stale parts: %v", err)
		}
	}()
}
