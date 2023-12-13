package models

type StandardResponse struct {
	StatusId string `json:"status_id"`
	Message  string `json:"message"`
	Data     any    `json:"data"`
}
