package models

type FileListRequest struct {
	Body struct {
		Id            uint   `json:"id" required:"true" min:"1" doc:"The server ID"`
		Dir           string `json:"dir" required:"true"`
		IncludeHidden bool   `json:"include_hidden"`
	}
}
