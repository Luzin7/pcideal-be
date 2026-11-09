package dto

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGenerateBuildRecommendationsDTOValidate(t *testing.T) {
	tests := []struct {
		name       string
		dto        GenerateBuildRecommendationsDTO
		wantStatus int
	}{
		{
			name:       "valid gaming request",
			dto:        GenerateBuildRecommendationsDTO{UsageType: "GAMING", BudgetCents: 500000, CpuPreference: "AMD", GpuPreference: "NVIDIA"},
			wantStatus: 0,
		},
		{
			name:       "valid work request with no preferences",
			dto:        GenerateBuildRecommendationsDTO{UsageType: "WORK", BudgetCents: 300000},
			wantStatus: 0,
		},
		{
			name:       "invalid usage type",
			dto:        GenerateBuildRecommendationsDTO{UsageType: "OFFICE", BudgetCents: 500000},
			wantStatus: 400,
		},
		{
			name:       "zero budget",
			dto:        GenerateBuildRecommendationsDTO{UsageType: "GAMING", BudgetCents: 0},
			wantStatus: 400,
		},
		{
			name:       "negative budget",
			dto:        GenerateBuildRecommendationsDTO{UsageType: "GAMING", BudgetCents: -100},
			wantStatus: 400,
		},
		{
			name:       "cpu preference too long",
			dto:        GenerateBuildRecommendationsDTO{UsageType: "GAMING", BudgetCents: 500000, CpuPreference: strings.Repeat("a", 51)},
			wantStatus: 400,
		},
		{
			name:       "gpu preference invalid charset",
			dto:        GenerateBuildRecommendationsDTO{UsageType: "GAMING", BudgetCents: 500000, GpuPreference: "NVIDIA<script>"},
			wantStatus: 400,
		},
		{
			name:       "no-preference allowed",
			dto:        GenerateBuildRecommendationsDTO{UsageType: "GAMING", BudgetCents: 500000, CpuPreference: "no-preference"},
			wantStatus: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.dto.Validate()

			if tt.wantStatus == 0 {
				if err != nil {
					t.Fatalf("expected no validation error, got %v", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("expected validation error with status %d, got nil", tt.wantStatus)
			}
			if err.StatusCode != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, err.StatusCode)
			}
		})
	}
}

func TestGenerateBuildRecommendationsDTOBindBudgetCents(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name    string
		body    string
		wantErr bool
	}{
		{name: "missing budget", body: `{}`, wantErr: true},
		{name: "zero budget", body: `{"budget_cents":0}`, wantErr: true},
		{name: "negative budget", body: `{"budget_cents":-1}`, wantErr: true},
		{name: "positive budget", body: `{"budget_cents":500000}`, wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
			ctx.Request = httptest.NewRequest("POST", "/", strings.NewReader(tt.body))
			ctx.Request.Header.Set("Content-Type", "application/json")

			var request GenerateBuildRecommendationsDTO
			err := ctx.ShouldBindJSON(&request)

			if tt.wantErr && err == nil {
				t.Fatal("expected bind error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("expected no bind error, got %v", err)
			}
		})
	}
}
