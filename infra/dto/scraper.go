package dto

import "github.com/Luzin7/pcideal-be/internal/core/models"

type ProductLinkToUpdate struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

type PartWithID struct {
	ID   string       `json:"id"`
	Part *models.Part `json:"part"`
}
