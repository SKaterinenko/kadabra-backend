package core

import (
	"errors"
	"fmt"
)

func BuildSQLError(err error) error {
	return fmt.Errorf("build sql: %w", err)
}

func QueryError(err error) error {
	return fmt.Errorf("query error: %w", err)
}

func ScanError(err error) error {
	return fmt.Errorf("scan error: %w", err)
}

func RowsError(err error) error {
	return fmt.Errorf("rows error: %w", err)
}

var (
	NoRecordDelete = errors.New("no record delete")
	NoGetRecord    = errors.New("record no found")
	NoUpdate       = errors.New("no fields to update")
)
