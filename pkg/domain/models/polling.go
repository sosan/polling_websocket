package models

const (
	UpdateCommand = "update"
	CreateCommnad = "create"
)

type ResponsePollingActionID struct {
	Data                   []RequestGoogleAction `json:"data,omitempty"`
	Rows                   *int64                `json:"rows,omitempty"`
	RowsBeforeLimitAtLeast *int64                `json:"rows_before_limit_at_least,omitempty"`
	Statistics             *Statistics           `json:"statistics,omitempty"`
	Meta                   []Meta                `json:"meta,omitempty"`
}
