package models

type FileListRequest struct {
	ID            uint   `path:"id" minimum:"1" required:"true" doc:"The ID of the server"`
	Dir           string `query:"dir" doc:"The directory to watch, if empty it defaults to home"`
	IncludeHidden bool   `query:"include_hidden" required:"true"  doc:"Include hidden files"`
}
type WatchDirectoryRequest struct {
	ID            uint   `path:"id" minimum:"1" required:"true" doc:"The ID of the server"`
	Dir           string `query:"dir" required:"true" doc:"The directory to watch, if empty it defaults to home"`
	IncludeHidden bool   `query:"include_hidden" required:"true" doc:"Include hidden files"`
}

type CreateDirectoryRequest struct {
	Body struct {
		Path string `json:"path" required:"true" doc:"The path of the directory to create"`
		Name string `json:"name" required:"true" doc:"The name of the new directory"`
	} `json:"body"`
	ID uint `path:"id" minimum:"1" required:"true" doc:"The ID of the server"`
}
type RemoveDirectoryRequest struct {
	Body struct {
		Path string `json:"path" required:"true" doc:"The path of the directory to remove"`
	} `json:"body"`
	ID uint `path:"id" minimum:"1" required:"true" doc:"The ID of the server"`
}
type CreateFileRequest struct {
	Body struct {
		Path string `json:"path" required:"true" doc:"The path of the file to create"`
		Name string `json:"name" required:"true" doc:"The name of the new file"`
	} `json:"body"`
	ID uint `path:"id" minimum:"1" required:"true" doc:"The ID of the server"`
}

type DeleteFileRequest struct {
	Body struct {
		Path string `json:"path" required:"true" doc:"The path of the file to delete"`
		Name string `json:"name" required:"true" doc:"The name of the file to delete"`
	} `json:"body"`
	ID uint `path:"id" minimum:"1" required:"true" doc:"The ID of the server"`
}

type RenamePathRequest struct {
	Body struct {
		Path    string `json:"path" required:"true" doc:"The current path of the file/directory"`
		OldName string `json:"old_name" required:"true" doc:"The current name of the file/directory"`
		NewName string `json:"new_name" required:"true" doc:"The new name for the file/directory"`
	} `json:"body"`
	ID uint `path:"id" minimum:"1" required:"true" doc:"The ID of the server"`
}
