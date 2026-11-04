package models

type FileListRequest struct {
	ID            uint   `path:"id" minimum:"1" doc:"The ID of the server"`
	Dir           string `query:"dir" doc:"The directory to list, if empty it defaults to home"`
	IncludeHidden bool   `query:"include_hidden" doc:"Include hidden files"`
}
