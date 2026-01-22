package utils

import (
	"strconv"
	"strings"
)

func ParseIntSlice(value string) ([]int, error) {
	if value == "" {
		return []int{}, nil
	}

	parts := strings.Split(value, ",")
	result := make([]int, 0, len(parts))

	for _, part := range parts {
		id, err := strconv.Atoi(part)
		if err != nil {
			return nil, err
		}
		result = append(result, id)
	}

	return result, nil
}
