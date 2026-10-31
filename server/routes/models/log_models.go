package models




type ListLogRequest struct {
	StartIndex uint      `query:"start-index" required:"true" min:"1" doc:"The start index for collecting logs"`
	Limit      uint      `query:"limit" min:"3" required:"true" max:"200" doc:"The maximum number of logs to return"`
}
type LogListItem struct {
	Id        uint   `json:"id"`
	Action    string `json:"action"`
	Message   string `json:"message"`
	CreatedAt string `json:"created_at"`
}

type LogListResponse struct {
	Body struct {
		Items         []LogListItem `json:"items"`
		NextIndex     uint          `json:"next_index"`
		PreviousIndex uint          `json:"previous_index"`
	}
}
