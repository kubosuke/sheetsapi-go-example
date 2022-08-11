package handler

type SuccessResult struct {
	Status bool `json:"status"`
	Data   any  `json:"data,omitempty"`
}
