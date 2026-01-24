package web

import (
	"errors"
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

func GetCoordinates(query url.Values) (lat *float64, lng *float64, err error) {
	latStr := query.Get("lat")
	lngStr := query.Get("lng")

	if latStr == "" && lngStr == "" {
		return nil, nil, nil
	}

	if latStr == "" || lngStr == "" {
		return nil, nil, errors.New("both 'lat' and 'lng' parameters are required for location filter")
	}
	l1, err1 := strconv.ParseFloat(latStr, 64)
	l2, err2 := strconv.ParseFloat(lngStr, 64)

	if err1 != nil || err2 != nil {
		return nil, nil, errors.New("invalid format for 'lat' or 'lng'")
	}

	return &l1, &l2, nil
}

func GetPrecision(query url.Values) (precision *int) {
	if query.Get("precision") != "" {
		p, err := strconv.Atoi(query.Get("precision"))
		if err == nil {
			return &p
		}
	}
	return nil
}
