package request

import (
	"net/url"
	"strconv"
)

func GetPagination(query url.Values) (page int, limit int) {
	limit, _ = strconv.Atoi(query.Get("limit"))
	if limit <= 0 {
		limit = 30
	}
	if limit > 50 {
		limit = 50
	}

	page, _ = strconv.Atoi(query.Get("page"))
	if page <= 0 {
		page = 1
	}

	return page, limit
}

func GetCoordinates(query url.Values) (lat *float64, lng *float64) {
	latStr := query.Get("lat")
	lngStr := query.Get("lng")

	if latStr != "" && lngStr != "" {
		l1, err1 := strconv.ParseFloat(latStr, 64)
		l2, err2 := strconv.ParseFloat(lngStr, 64)
		if err1 == nil && err2 == nil {
			return &l1, &l2
		}
	}
	return nil, nil
}
