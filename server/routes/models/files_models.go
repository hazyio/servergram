package models

type FileListRequest struct {
	ID            uint   `path:"id" minimum:"1" doc:"The ID of the server"`
	Dir           string `query:"dir" required:"true" doc:"The directory to list"`
	IncludeHidden bool   `query:"include_hidden" doc:"Include hidden files"`
}
