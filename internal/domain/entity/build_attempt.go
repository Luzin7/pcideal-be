package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BuildAttempt struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	IP        string             `bson:"ip" json:"ip"`
	Goal      string             `bson:"goal" json:"goal"`
	Budget    int64              `bson:"budget" json:"budget"`
	CPUPref   string             `bson:"cpu_pref,omitempty" json:"cpu_pref,omitempty"`
	GPUPref   string             `bson:"gpu_pref,omitempty" json:"gpu_pref,omitempty"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

func NewBuildAttempt(ip string, goal string, budget int64, cpuPref string, gpuPref string) *BuildAttempt {
	return &BuildAttempt{
		ID:        primitive.NewObjectID(),
		IP:        ip,
		Goal:      goal,
		Budget:    budget,
		CPUPref:   cpuPref,
		GPUPref:   gpuPref,
		CreatedAt: time.Now(),
	}
}
