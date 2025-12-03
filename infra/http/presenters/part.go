package presenters

import (
	"time"

	"github.com/Luzin7/pcideal-be/internal/domain/entity"
)

type PartPresenter struct {
	ID             string `json:"id"`
	PriceFormatted int64  `json:"price"`
	Specs          any    `json:"specs"`
	LastUpdate     string `json:"last_update"`
}

func ToPartPresenter(p *entity.Part) PartPresenter {
	return PartPresenter{
		ID:             p.ID.Hex(),
		PriceFormatted: p.PriceCents,
		Specs:          p.Specs,
		LastUpdate:     time.Since(p.UpdatedAt).String(),
	}
}
