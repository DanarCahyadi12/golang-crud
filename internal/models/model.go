package models

type Response[T any] struct {
	Message string `json:"message,omitempty"`
	Data    T      `json:"data,omitempty"`
}
type ErrorResponse struct {
	Code    int
	Message string
	Status  string
}

func (e ErrorResponse) Error() string {
	return e.Message
}
