package utils

import (
	"fmt"
	"net/url"
	"strconv"
)

func GetLimitOffset(q url.Values) (limit, offset int, err error) {
	if limitStr := q.Get("limit"); limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid limit: %w", err)
		}
	}

	if offsetStr := q.Get("offset"); offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid offset: %w", err)
		}
	}

	return limit, offset, nil
}
