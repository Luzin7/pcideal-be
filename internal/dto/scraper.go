package dto

import "github.com/Luzin7/pcideal-be/internal/domain/entity"

type ProductLinkToUpdate struct {
	ID  string `json:"id"`
	Url string `json:"url"`
}

type PartWithID struct {
	ID   string       `json:"id"`
	Part *entity.Part `json:"part"`
}
