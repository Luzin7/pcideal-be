package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Specs representa as especificações técnicas de uma peça
type specs struct {
	// CPU e Motherboard
	Socket           string  `bson:"socket,omitempty" json:"socket,omitempty"`           // AM4, LGA1700, etc
	TDP              int64   `bson:"tdp,omitempty" json:"tdp,omitempty"`                 // em Watts
	ClockSpeed       float64 `bson:"clock_speed,omitempty" json:"clock_speed,omitempty"` // em GHz
	BoostClock       float64 `bson:"boost_clock,omitempty" json:"boost_clock,omitempty"` // em GHz
	MemoryMaxSpeed   int64   `bson:"memory_max_speed,omitempty" json:"memory_max_speed,omitempty"`
	Chipset          string  `bson:"chipset,omitempty" json:"chipset,omitempty"`
	CpuCompatibility string  `bson:"cpu_compatibility,omitempty" json:"cpu_compatibility,omitempty"`

	// RAM e Motherboard
	MemoryType  string `bson:"memory_type,omitempty" json:"memory_type,omitempty"` // DDR4, DDR5
	MemorySlots int64  `bson:"memory_slots,omitempty" json:"memory_slots,omitempty"`

	// RAM e SSD
	Capacity int64 `bson:"capacity,omitempty" json:"capacity,omitempty"` // em GB
	Speed    int64 `bson:"speed,omitempty" json:"speed,omitempty"`       // MHz para RAM

	// GPU e SSD
	Interface string `bson:"interface,omitempty" json:"interface,omitempty"` // PCIe 3.0, PCIe 4.0

	// GPU
	PowerSupply int64  `bson:"power_supply,omitempty" json:"power_supply,omitempty"` // Fonte recomendada em W
	VideoMemory string `bson:"video_memory,omitempty" json:"video_memory,omitempty"`

	// Motherboard e Case
	FormFactor string `bson:"form_factor,omitempty" json:"form_factor,omitempty"` // ATX, mATX, ITX
}

// type benchmark struct {
// 	CinebenchR23 int `bson:"cinebench_r23,omitempty" json:"cinebench_r23,omitempty"`
// 	TimeSpy      int `bson:"3dmark_timespy,omitempty" json:"3dmark_timespy,omitempty"`
// }

type Part struct {
	ID    primitive.ObjectID `bson:"_id" json:"id"`
	Type  string             `bson:"type" json:"type"`   // CPU, GPU, MOTHERBOARD, RAM, SSD, PSU, CASE
	Brand string             `bson:"brand" json:"brand"` // AMD, Intel, NVIDIA, etc
	Model string             `bson:"model" json:"model"`
	Specs specs              `bson:"specs" json:"specs"`
	// Benchmark   benchmark `bson:"benchmark" json:"benchmark"`
	PriceCents    int64     `bson:"price_cents" json:"price_cents"`
	URL           string    `bson:"url" json:"url"`
	AffiliatedURL string    `bson:"affiliate_url" json:"affiliate_url"`
	Store         string    `bson:"store" json:"store"` // Kabum, Amazon, etc
	UpdatedAt     time.Time `bson:"updated_at" json:"updated_at"`
}

func NewPart(partType, brand, model string, specs specs, priceCents int64, url string, affiliateUrl string, store string) *Part {
	return &Part{
		ID:            primitive.NewObjectID(),
		Type:          partType,
		Brand:         brand,
		Model:         model,
		Specs:         specs,
		PriceCents:    priceCents,
		URL:           url,
		AffiliatedURL: affiliateUrl,
		Store:         store,
		UpdatedAt:     time.Now(),
	}
}
