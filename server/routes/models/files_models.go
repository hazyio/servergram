package models

type FileListRequest struct {
	ID            uint   `path:"id" minimum:"1" required:"true" doc:"The ID of the server"`
	Dir           string `query:"dir" required:"true"  doc:"The directory to list, if empty it defaults to home"`
	IncludeHidden bool   `query:"include_hidden" required:"true"  doc:"Include hidden files"`
}
type WatchDirectoryRequest struct {
	ID            uint   `path:"id" minimum:"1" required:"true" doc:"The ID of the server"`
	Dir           string `query:"dir" required:"true" doc:"The directory to watch, if empty it defaults to home"`
	IncludeHidden bool   `query:"include_hidden" required:"true" doc:"Include hidden files"`
}
