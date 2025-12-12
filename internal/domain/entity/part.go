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
	Socket           string `bson:"socket,omitempty" json:"socket,omitempty"`
	MemoryType       string `bson:"memory_type,omitempty" json:"memory_type,omitempty"`
	FormFactor       string `bson:"form_factor,omitempty" json:"form_factor,omitempty"`
	MinPSUWatts      int16  `bson:"min_psu_watts,omitempty" json:"min_psu_watts,omitempty"`
	Wattage          int16  `bson:"wattage,omitempty" json:"wattage,omitempty"`
	EfficiencyRating int8   `bson:"efficiency_rating,omitempty" json:"efficiency_rating,omitempty"`
	PerformanceScore int8   `bson:"performance_score,omitempty" json:"performance_score,omitempty"`
	TierScore        int8   `bson:"tier_score,omitempty" json:"tier_score,omitempty"`
	VramGB           int    `bson:"vram_gb,omitempty" json:"vram_gb,omitempty"`
	CapacityGB       int    `bson:"capacity_gb,omitempty" json:"capacity_gb,omitempty"`
	MemorySpeedMHz   int    `bson:"memory_speed_mhz,omitempty" json:"memory_speed_mhz,omitempty"`
	CasLatency       int    `bson:"cas_latency,omitempty" json:"cas_latency,omitempty"`
}

type Part struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Type          PartType           `bson:"type" json:"type"`
	Brand         string             `bson:"brand" json:"brand"`
	Model         string             `bson:"model" json:"model"`
	URL           string             `bson:"url" json:"url"`
	Store         string             `bson:"store" json:"store"`
	AffiliatedURL string             `bson:"affiliate_url" json:"affiliate_url"`
	PriceCents    int64              `bson:"price_cents" json:"price_cents"`
	Specs         Specs              `bson:"specs" json:"specs"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`
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
		UpdatedAt:     time.Now(),
	}
}
