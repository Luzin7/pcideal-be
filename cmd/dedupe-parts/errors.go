package main

import (
	"fmt"
	"strings"
)

type MissingConfirmationError struct {
	Provided string
}

func (e *MissingConfirmationError) Error() string {
	return fmt.Sprintf("--apply requires --confirm=%s (got %q)", deleteConfirmation, e.Provided)
}

type MissingConfigError struct {
	Variable string
}

func (e *MissingConfigError) Error() string {
	return fmt.Sprintf("%s is required", e.Variable)
}

type DeleteGroupError struct {
	Store string
	URL   string
	Err   error
}

func (e *DeleteGroupError) Error() string {
	return fmt.Sprintf("delete duplicates for store=%s url=%s: %v", e.Store, e.URL, e.Err)
}

func (e *DeleteGroupError) Unwrap() error {
	return e.Err
}

type VerificationError struct {
	Reasons []string
}

func (e *VerificationError) Error() string {
	return "verification failed: " + strings.Join(e.Reasons, "; ")
}

type MissingKeyError struct {
	Count int64
}

func (e *MissingKeyError) Error() string {
	return fmt.Sprintf(
		"cannot apply: %d parts have missing or empty store/url; fix them before deduplicating",
		e.Count,
	)
}
