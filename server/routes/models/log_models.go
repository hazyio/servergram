package models

type ListLogRequest struct {
	Body struct {
		StartIndex uint `json:"start" required:"true" minLength:"1" doc:"The start index for collecting logs"`
		Collect    uint `json:"collect" required:"true" min:"1" doc:"The number of logs to collect"`
	}
}
type LogListItem struct {
	Id        uint   `json:"id"`
	Action    string `json:"action"`
	Message   string `json:"message"`
	CreatedAt string `json:"created_at"`
}

type LogListResponse struct {
	Body struct {
		Items []LogListItem `json:"items"`
	}
}
