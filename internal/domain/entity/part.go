package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PartType string

const (
	TypeCPU  PartType = "CPU"
	TypeGPU  PartType = "GPU"
	TypeMobo PartType = "MOTHERBOARD"
	TypeRAM  PartType = "RAM"
	TypePSU  PartType = "PSU"
	TypeSSD  PartType = "SSD"
)

type Specs struct {
	Socket     string `bson:"socket,omitempty" json:"socket,omitempty"`
	MemoryType string `bson:"memory_type,omitempty" json:"memory_type,omitempty"`
	// FormFactor string `bson:"form_factor,omitempty" json:"form_factor,omitempty"` // ATX, mATX
	Wattage     int `bson:"wattage,omitempty" json:"wattage,omitempty"`             // O que a Fonte ENTREGA
	MinPSUWatts int `bson:"min_psu_watts,omitempty" json:"min_psu_watts,omitempty"` // O que a GPU PEDE
	VramGB      int `bson:"vram_gb,omitempty" json:"vram_gb,omitempty"`             // GPU
	CapacityGB  int `bson:"capacity_gb,omitempty" json:"capacity_gb,omitempty"`     // RAM/SSD
	// GpuTier    string `bson:"gpu_tier,omitempty" json:"gpu_tier,omitempty"`
}

type Part struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	Type          PartType           `bson:"type"`
	Brand         string             `bson:"brand"`
	Model         string             `bson:"model"`
	URL           string             `bson:"url"`
	Store         string             `bson:"store"`
	AffiliatedURL string             `bson:"affiliated_url"`
	PriceCents    int64              `bson:"price_cents"`
	Specs         Specs              `bson:"specs"`
	IsParsed      bool               `bson:"is_parsed"`
	UpdatedAt     time.Time          `bson:"updated_at"`
}

func NewPart(ty PartType, brand string, model string, url string, store string, affiliatedURL string, priceCents int64, specs Specs) *Part {
	return &Part{
		ID:            primitive.NewObjectID(),
		Type:          ty,
		Brand:         brand,
		Model:         model,
		URL:           url,
		Store:         store,
		AffiliatedURL: affiliatedURL,
		PriceCents:    priceCents,
		Specs:         specs,
		IsParsed:      true,
		UpdatedAt:     time.Now(),
	}
}
