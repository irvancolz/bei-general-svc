package model

import "time"

type FilterQueryParameter struct {
	QueryList map[string][]string
	EndDate   time.Time
	Now       time.Time
	OrderBy   string
	Order     string
	Limit     int
	Offset    int
}


