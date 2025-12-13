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

type UpdateProductsResponse struct {
	Store   string        `json:"store"`
	Total   int           `json:"total"`
	Success int           `json:"success"`
	Failed  int           `json:"failed"`
	Results []*PartWithID `json:"results"`
}
