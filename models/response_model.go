package models

// ? Outputs structures

// ResponseDTO structure
type ResponseDTO struct {
	Data    map[string]interface{} `json:"data"`
	Success bool                   `json:"success"`
	Status  int                    `json:"status"`
	Message string                 `json:"message"`
}