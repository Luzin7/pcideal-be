package errors

import (
	"fmt"
)

type ErrService struct {
	StatusCode int
	Message    string
}

func (e *ErrService) Error() string {
	return fmt.Sprintf("%d: %s", e.StatusCode, e.Message)
}

func New(msg string, statusCode int) *ErrService {
	return &ErrService{
		Message:    msg,
		StatusCode: statusCode,
	}
}

func ErrNotFound(ctx string) *ErrService {
	return New(
		fmt.Sprintf("%s not found", ctx),
		404,
	)
}

func ErrAlreadyExists(ctx string) *ErrService {
	return New(
		fmt.Sprintf("%s already exists", ctx),
		409,
	)
}

func ErrScrapingFailed(partId string) *ErrService {
	return New(
		fmt.Sprintf("Failed to scrape product for part %s", partId),
		502,
	)
}

func ErrScrapingFailedBatch(store string, count int) *ErrService {
	return New(
		fmt.Sprintf("Failed to scrape products for store %s (%d parts)", store, count),
		502,
	)
}

func ErrScrapingTimeout(partId string) *ErrService {
	return New(
		fmt.Sprintf("Scraping timeout for part %s", partId),
		504, // Gateway Timeout
	)
}

func ErrScrapingInvalidURL(partId, url string) *ErrService {
	return New(
		fmt.Sprintf("Invalid or inaccessible URL for part %s: %s", partId, url),
		400,
	)
}

func ErrScrapingParseError(partId string) *ErrService {
	return New(
		fmt.Sprintf("Failed to parse scraped data for part %s", partId),
		502, // Bad Gateway
	)
}

func ErrScrapingBlocked(partId string) *ErrService {
	return New(
		fmt.Sprintf("Scraping blocked or rate limited for part %s", partId),
		429, // Too Many Requests
	)
}

func ErrScrapingProductNotFound(partId string) *ErrService {
	return New(
		fmt.Sprintf("Product not found on target site for part %s", partId),
		404,
	)
}

func ErrScrapingNetworkError(partId string) *ErrService {
	return New(
		fmt.Sprintf("Network error while scraping part %s", partId),
		502,
	)
}

func ErrInternalServerError() *ErrService {
	return New(
		"Internal server error",
		500,
	)
}

func ErrMissingIP() *ErrService {
	return New("User IP not provided", 400)
}

func ErrMissingGoal() *ErrService {
	return New("Required parameter 'goal' not provided", 400)
}

func ErrMissingBudget() *ErrService {
	return New("Required parameter 'budget_cents' not provided", 400)
}

func ErrInvalidSince() *ErrService {
	return New("'since' cannot be a future date", 400)
}

func ErrInvalidLimit() *ErrService {
	return New("'limit' must be a positive integer up to the allowed maximum", 400)
}

func ErrInvalidOffset() *ErrService {
	return New("'offset' must be a non-negative integer", 400)
}

func ErrInvalidBudget() *ErrService {
	return New("'budget_cents' must be greater than zero", 400)
}

func ErrInvalidRequestBody() *ErrService {
	return New("Invalid request body", 400)
}

func ErrInvalidUsageType(usageType string) *ErrService {
	return New(
		fmt.Sprintf("Invalid 'usage_type': %s", usageType),
		400,
	)
}

func ErrInvalidPreference(field string) *ErrService {
	return New(
		fmt.Sprintf("Invalid '%s'", field),
		400,
	)
}

func ErrNoCompatibleParts(partType string) *ErrService {
	return New(
		fmt.Sprintf("No compatible %s parts found for the requested budget", partType),
		422,
	)
}

func ErrBuildAttemptLimit() *ErrService {
	return New("Build attempt limit reached for this IP", 429)
}

func ErrBuildAttemptNotFound() *ErrService {
	return New("BuildAttempt not found", 404)
}
