package database

import "fmt"

type UniqueIndexNotReadyError struct {
	CheckErr error
}

func (e *UniqueIndexNotReadyError) Error() string {
	if e.CheckErr != nil {
		return fmt.Sprintf("could not verify unique index %s: %v", StoreURLIndexName, e.CheckErr)
	}

	return fmt.Sprintf(
		"unique index %s is missing; imports are NOT idempotent. Run the duplicate audit: %s",
		StoreURLIndexName, StoreURLDuplicateAuditCommand,
	)
}

func (e *UniqueIndexNotReadyError) Unwrap() error {
	return e.CheckErr
}
