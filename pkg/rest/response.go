package rest

import (
	"encoding/json"
	"fmt"
	"math"
)

type Meta struct {
	Next        *bool   `json:"next,omitempty"`
	Code        int     `json:"code,omitempty"`
	Message     string  `json:"message,omitempty"`
	Errors      string  `json:"errors,omitempty"`
	DataCount   *uint64 `json:"-"`
	Page        uint64  `json:"page,omitempty"`
	PerPage     uint64  `json:"per_page,omitempty"`
	SearchQuery string  `json:"search_query,omitempty"`
	SearchType  string  `json:"search_type,omitempty"`
	FilterType  string  `json:"filter_type,omitempty"`
	SortBy      string  `json:"sort_by,omitempty"`
	SortOrder   string  `json:"sort_order,omitempty"`
	StartDate   string  `json:"start_date,omitempty"`
	EndDate     string  `json:"end_date,omitempty"`

	PrevPage    int   `json:"prev_page,omitempty"`
	NextPage    int   `json:"next_page,omitempty"`
	TotalRecord int64 `json:"data_count,omitempty"`
	TotalPage   int   `json:"total_page,omitempty"`
}

type ListResponse struct {
	DataCount   uint64      `json:"data_count"`
	DataPage    uint64      `json:"data_page"`
	DataPerPage uint64      `json:"data_per_page"`
	FilterType  string      `json:"filter_type,omitempty"`
	SearchQuery string      `json:"search_query,omitempty"`
	SearchType  string      `json:"search_type,omitempty"`
	SortBy      string      `json:"sort_by,omitempty"`
	SortOrder   string      `json:"sort_order,omitempty"`
	StartDate   string      `json:"start_date,omitempty"`
	EndDate     string      `json:"end_date,omitempty"`
	DataSlice   interface{} `json:"data_slice"`
}

func (lr *ListResponse) Builder(query *RequestQuery, count uint64, data interface{}) *ListResponse {
	lr.DataPage = query.PageNum
	lr.DataPerPage = query.PerPageNum
	lr.DataCount = count
	lr.DataSlice = data
	lr.SearchType = query.SearchType
	lr.FilterType = query.FilterType
	lr.SearchQuery = query.SearchQuery
	lr.SortBy = query.SortBy
	lr.SortOrder = query.SortOrder
	lr.StartDate = query.StartDate
	lr.EndDate = query.EndDate

	return lr
}

// swagger:response defaultResponse
type Response struct {
	// in: body
	Code    int             `json:"code,omitempty"`
	Meta    *Meta           `json:"meta,omitempty"`
	Errors  []ErrorResponse `json:"errors,omitempty"`
	Data    interface{}     `json:"data,omitempty"`
	Message string          `json:"message,omitempty"`
}

func (resp *Response) MarshalJSON() ([]byte, error) {
	if resp.Meta != nil && resp.Meta != (&Meta{}) {
		meta := new(Meta)
		m := resp.Meta
		meta.Page = m.Page
		meta.PerPage = m.PerPage
		totalRecord := m.TotalRecord
		if m.DataCount != nil {
			totalRecord = int64(*m.DataCount)
		}
		meta.TotalRecord = totalRecord
		if totalRecord > 0 {
			var totalPage int
			if m.PerPage > 1 {
				totalPage = int(math.Ceil(float64(m.TotalRecord) / float64(m.PerPage)))
			}

			if m.Page > 1 {
				meta.PrevPage = int(m.Page - 1)
			} else {
				meta.PrevPage = int(m.Page)
			}

			if int(m.Page) == totalPage {
				meta.NextPage = int(m.Page)
			} else {
				meta.NextPage = int(m.Page + 1)
			}
			meta.TotalPage = totalPage
		}
		return json.Marshal(struct {
			Meta *Meta `json:"meta,omitempty"`
			Response
		}{
			Meta:     meta,
			Response: *resp,
		})
	}
	return json.Marshal(struct {
		Response
	}{
		Response: *resp,
	})
}

type ErrorResponse struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message,omitempty"`
}

type Errors []ErrorResponse

type ConfigStep struct {
	Next bool        `json:"next"`
	Data interface{} `json:"data"`
}

func (resp *Response) ResponseBuilder(code int, msg string, data interface{}, err []ErrorResponse) (rR *Response) {
	configStep, configStepOK := data.(ConfigStep)
	dataSlice, listResponseOK := data.(*ListResponse)
	resp.Code = code
	resp.Message = msg
	resp.Errors = err
	if configStepOK {
		m := Meta{
			Next: &configStep.Next,
		}
		resp.Meta = &m
		resp.Data = configStep.Data
	} else if listResponseOK {
		count := dataSlice.DataCount
		m := Meta{
			TotalRecord: int64(count),
			PerPage:     dataSlice.DataPerPage,
			Page:        dataSlice.DataPage,
			SearchQuery: dataSlice.SearchQuery,
			SearchType:  dataSlice.SearchType,
			FilterType:  dataSlice.FilterType,
			SortOrder:   dataSlice.SortOrder,
			SortBy:      dataSlice.SortBy,
			StartDate:   dataSlice.StartDate,
			EndDate:     dataSlice.EndDate,
		}
		resp.Meta = &m
		resp.Data = dataSlice.DataSlice
	} else {
		resp.Data = data
	}
	return resp
}

func (lr *ListResponse) BuilderDocList(pageNum uint64, perPageNum uint64, count uint64, sortBy string, sortOrder string, data interface{}) *ListResponse {
	lr.DataPage = pageNum
	lr.DataPerPage = perPageNum
	lr.DataCount = count
	lr.DataSlice = data
	lr.SortBy = sortBy
	lr.SortOrder = sortOrder
	return lr
}

func (resp *Response) ErrorGeneratorFromInterface(errors []interface{}) []ErrorResponse {
	var errResponseFields []ErrorResponse
	for _, v := range errors {
		var errR ErrorResponse
		vError := v.(map[string]interface{})
		errR.Field = fmt.Sprint(vError["field"])
		errR.Message = fmt.Sprint(vError["messages"].([]interface{})[0])
		errResponseFields = append(errResponseFields, errR)
	}
	return errResponseFields
}

func (resp *Response) ErrorGenerator(field string, message string) []ErrorResponse {
	var es []ErrorResponse
	en := ErrorResponse{Field: field, Message: message}
	es = append(es, en)
	return es
}
