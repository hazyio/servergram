package models

import (
	"github.com/danielgtaylor/huma/v2"
	"hazyio.com/servergram/shared/gen/proto"
)

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

type PathPropertiesRequest struct {
	Path string `query:"path" required:"true" doc:"The current path of the file/directory"`
	Name string `query:"name" required:"true" doc:"The name of the file/directory"`
	ID   uint   `path:"id" minimum:"1" required:"true" doc:"The ID of the server"`
}

type PathPropertiesResponse struct {
	Body struct {
		Enity *proto.FileEntry
	}
}

type ReadFileRequest struct {
	Path         string `query:"path" required:"true" doc:"The path of the file to read"`
	Name         string `query:"name" required:"true" doc:"The name of the file to read"`
	MaxBytesSize int64  `query:"max_bytes_size" required:"true" doc:"The maximum number of bytes to read from the file"`
	ID           uint   `path:"id" minimum:"1" required:"true" doc:"The ID of the server"`
}
type DownloadFileRequest struct {
	Path string `query:"path" required:"true" doc:"The path of the file to read"`
	Name string `query:"name" required:"true" doc:"The name of the file to read"`
	ID   uint   `path:"id" minimum:"1" required:"true" doc:"The ID of the server"`
}

type UploadFileRequest struct {
	ID      uint `path:"id" minimum:"1" required:"true" doc:"The ID of the server"`
	RawBody huma.MultipartFormFiles[struct {
		File huma.FormFile `form:"file" required:"true"`
	}]
}
