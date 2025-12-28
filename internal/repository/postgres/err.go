package repository

import (
	"errors"
	"fmt"
)

func buildSQLError(err error) error {
	return fmt.Errorf("build sql: %w", err)
}

func queryError(err error) error {
	return fmt.Errorf("query error: %w", err)
}

func scanError(err error) error {
	return fmt.Errorf("scan error: %w", err)
}

func rowsError(err error) error {
	return fmt.Errorf("rows error: %w", err)
}

var (
	NoRecordDelete = errors.New("no record delete")
	NoGetRecord    = errors.New("record no found")
	NoUpdate       = errors.New("no fields to update")
)
