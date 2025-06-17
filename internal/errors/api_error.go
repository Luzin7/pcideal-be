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
		502, // Bad Gateway
	)
}

func ErrInternalServerError() *ErrService {
	return New(
		"Internal server error",
		500,
	)
}
