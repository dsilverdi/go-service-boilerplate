package rest

import "time"

type RequestQuery struct {
	Limit         uint64    `json:"limit" schema:"limit"`
	Offset        uint64    `json:"offset" schema:"offset"`
	Page          string    `json:"page" schema:"page"`
	PerPage       string    `json:"per_page" schema:"per_page"`
	SearchQuery   string    `json:"search_query" schema:"search_query"`
	SearchType    string    `json:"search_type" schema:"search_type"`
	FilterType    string    `json:"filter_type" schema:"filter_type"`
	SearchColumn  string    `json:"search_column" schema:"search_column"`
	SortBy        string    `json:"sort_by" schema:"sort_by"`
	SortColumn    string    `json:"sort_column" schema:"sort_column"`
	SortOrder     string    `json:"sort_order" schema:"sort_order"`
	ExpiresString string    `json:"expires_string" schema:"expires"`
	Expires       int64     `json:"expires" schema:"expires"`
	Signature     string    `json:"signature" schema:"signature"`
	StartDate     string    `json:"start_date" schema:"start_date"`
	EndDate       string    `json:"end_date" schema:"end_date"`
	StartDateDate time.Time `json:"start_date_date" schema:"-"`
	EndDateDate   time.Time `json:"end_date_date" schema:"-"`
	PageNum       uint64    `json:"-" schema:"-"`
	PerPageNum    uint64    `json:"-" schema:"-"`
	TotalPage     int64     `json:"total_page" schema:"-"`
}
