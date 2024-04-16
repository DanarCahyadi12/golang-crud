package models

type Response[T any] struct {
	Message  string    `json:"message,omitempty"`
	Metadata *Metadata `json:"metadata,omitempty"`
	Data     T         `json:"data,omitempty"`
}
type ErrorResponse struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Status  string `json:"status,omitempty"`
}

type Metadata struct {
	PageNumber     int    `json:"page_number,omitempty"`
	PageSize       int64  `json:"page_size,omitempty"`
	TotalItemCount int64  `json:"total_item_count,omitempty"`
	Next           string `json:"next"`
	Prev           string `json:"prev"`
}

func (e ErrorResponse) Error() string {
	return e.Message
}
