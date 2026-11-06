package files

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"log"

	"github.com/fsnotify/fsnotify"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"hazyio.com/servergram/shared/gen/proto"
)

func handleListFiles(req *proto.ListFilesRequest, stream proto.Files_ListFilesServer) error {

	newPath, err := validateAndCleanPath(req.Dir)
	if err != nil {
		return err
	}
	entries, err := os.ReadDir(newPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return status.Errorf(codes.NotFound, "directory %q not found", req.Dir)
		}
		if errors.Is(err, fs.ErrPermission) {
			return status.Errorf(codes.PermissionDenied, "cannot read %q", req.Dir)
		}

		return status.Error(codes.Internal, "unknown error")
	}
	batch := make([]*proto.FileEntry, 0, 500)
	for _, e := range entries {
		if stripDirFromPath(req.Dir, e.Name())[0] == '.' && !req.IncludeHidden {
			continue
		}
		info, _ := e.Info()
		batch = append(batch, toEntry(info))
		if len(batch) == 500 {
			stream.Send(&proto.ListFilesChunkResponse{Entries: batch, ResolvedPath: req.Dir})
			batch = batch[:0]
		}
	}
	stream.Send(&proto.ListFilesChunkResponse{Entries: batch, ResolvedPath: req.Dir, Done: true})

	return nil
}

func handleCreateDirectory(req *proto.CreateDirectoryRequest) (*proto.SuccessResponse, error) {
	newPath, err := validateAndCleanPath(req.Path)
	if err != nil {
		return nil, err
	}
	if err := validateEntryName(req.Name); err != nil {
		return nil, err
	}
	newDirPath := filepath.Join(newPath, string(req.Name))
	err = os.Mkdir(newDirPath, 0755)
	if err != nil {
		if errors.Is(err, fs.ErrExist) {
			return nil, status.Errorf(codes.AlreadyExists, "directory %q already exists", req.Path)
		}
		if errors.Is(err, fs.ErrPermission) {
			return nil, status.Errorf(codes.PermissionDenied, "cannot create directory %q", req.Path)
		}
		log.Printf("failed to create directory %q: %v", req.Path, err)
		return nil, status.Error(codes.Internal, "unknown error")
	}

	return &proto.SuccessResponse{Success: true}, nil
}

func handleRemoveDirectory(req *proto.RemoveDirectoryRequest) (*proto.SuccessResponse, error) {
	newPath, err := validateAndCleanPath(req.Path)
	if err != nil {
		return nil, err
	}
	info, err := os.Stat(newPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, status.Errorf(codes.NotFound, "directory %q not found", req.Path)
		}
		if errors.Is(err, fs.ErrPermission) {
			return nil, status.Errorf(codes.PermissionDenied, "cannot access directory %q", req.Path)
		}
		log.Printf("failed to access directory %q: %v", req.Path, err)
		return nil, status.Error(codes.Internal, "unknown error")
	}
	if !info.IsDir() {
		return nil, status.Errorf(codes.InvalidArgument, "path %q is not a directory", req.Path)
	}
	err = os.Remove(req.Path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, status.Errorf(codes.NotFound, "directory %q not found", req.Path)
		}
		if errors.Is(err, fs.ErrPermission) {
			return nil, status.Errorf(codes.PermissionDenied, "cannot remove directory %q", req.Path)
		}
		log.Printf("failed to remove directory %q: %v", req.Path, err)
		return nil, status.Error(codes.Internal, "unknown error")
	}

	return &proto.SuccessResponse{Success: true}, nil
}

func handleWatchDirectory(req *proto.WatchDirectoryRequest, stream proto.Files_WatchDirectoryServer) error {
	newPath, err := validateAndCleanPath(req.Path)
	if err != nil {
		return err
	}
	info, err := os.Stat(newPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return status.Errorf(codes.NotFound, "directory %q not found", req.Path)
		}
		if errors.Is(err, fs.ErrPermission) {
			return status.Errorf(codes.PermissionDenied, "cannot access directory %q", req.Path)
		}
		log.Printf("failed to access directory %q: %v", req.Path, err)
		return status.Error(codes.Internal, "unknown error")
	}
	if !info.IsDir() {
		return status.Errorf(codes.InvalidArgument, "path %q is not a directory", req.Path)
	}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return status.Error(codes.Internal, "failed to create watcher")
	}
	defer watcher.Close()
	watcher.Add(req.Path)
	log.Printf("watching directory %q", req.Path)

	go func() {
		// Add all existing files and dir in the directory to the watcher
		entries, _ := os.ReadDir(req.Path)
		for _, e := range entries {
			if stripDirFromPath(req.Path, e.Name())[0] == '.' && !req.IncludeHidden {
				continue
			}
			log.Printf("adding %q to watcher", e.Name())
			watcher.Add(e.Name())
		}

	}()
	getPathEntryInfo := func(path string) (*proto.FileEntry, error) {
		info, err := os.Stat(path)
		if err != nil {
			return nil, err
		}
		return toEntry(info), nil
	}
	for {
		select {
		case event, ok := <-watcher.Events:
			if ok {
				name := stripDirFromPath(req.Path, event.Name) // remove the directory from the path to save bytes
				if event.Has(fsnotify.Create) {
					log.Println("created:", event.Name)
					entry, err := getPathEntryInfo(event.Name)
					if err != nil {
						log.Printf("failed to get entry for %q: %v", event.Name, err)
						continue
					}
					stream.Send(&proto.WatchDirectoryResponse{Kind: proto.WatchDirectoryKind_CREATED, Path: name, Entry: entry})
					watcher.Add(event.Name)
				} else if event.Has(fsnotify.Remove) {
					log.Println("deleted:", event.Name)
					stream.Send(&proto.WatchDirectoryResponse{Kind: proto.WatchDirectoryKind_DELETED, Path: name})
					watcher.Remove(event.Name)

				} else if event.Has(fsnotify.Rename) {
					// send a deleted event for the old name
					log.Println("renamed:", event.Name)
					stream.Send(&proto.WatchDirectoryResponse{Kind: proto.WatchDirectoryKind_DELETED, Path: name})
					watcher.Remove(event.Name)
				}

			}

		case err, ok := <-watcher.Errors:
			if !ok {
				log.Println("watcher closed")
				return nil
			}
			if err != nil {
				log.Println("watcher error:", err)
				return status.Error(codes.Internal, "watcher error")
			}
		}
	}

}

func handleCreateFile(req *proto.CreateFileRequest) (*proto.SuccessResponse, error) {
	newPath, err := validateAndCleanPath(req.Path)
	if err != nil {
		return nil, err
	}

	if err := validateEntryName(req.Name); err != nil {
		return nil, err
	}
	info, err := probePathExists(newPath)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, status.Errorf(codes.InvalidArgument, "path %q is not a directory", req.Path)
	}

	filePath := filepath.Join(newPath, req.Name)
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		if errors.Is(err, fs.ErrExist) {
			return nil, status.Errorf(codes.AlreadyExists, "file %q already exists", req.Name)
		}
		if errors.Is(err, fs.ErrPermission) {
			return nil, status.Errorf(codes.PermissionDenied, "cannot create file %q", req.Name)
		}
		if errors.Is(err, fs.ErrNotExist) {
			return nil, status.Errorf(codes.NotFound, "directory %q not found", req.Path)
		}
		log.Printf("failed to create file %q: %v", filePath, err)
		return nil, status.Error(codes.Internal, "unknown error")
	}
	if err := file.Close(); err != nil {
		log.Printf("failed to close file %q: %v", filePath, err)
		return nil, status.Error(codes.Internal, "unknown error")
	}

	return &proto.SuccessResponse{Success: true}, nil
}

func handleDeleteFile(req *proto.RemoveFileRequest) (*proto.SuccessResponse, error) {
	newPath, err := validateAndCleanPath(req.Path)
	if err != nil {
		return nil, err
	}

	if err := validateEntryName(req.Name); err != nil {
		return nil, err
	}

	info, err := os.Stat(newPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, status.Errorf(codes.NotFound, "directory %q not found", req.Path)
		}
		if errors.Is(err, fs.ErrPermission) {
			return nil, status.Errorf(codes.PermissionDenied, "cannot access directory %q", req.Path)
		}
		log.Printf("failed to access directory %q: %v", req.Path, err)
		return nil, status.Error(codes.Internal, "unknown error")
	}

	if !info.IsDir() {
		return nil, status.Errorf(codes.InvalidArgument, "path %q is not a directory", req.Path)
	}

	filePath := filepath.Join(newPath, req.Name)
	err = os.Remove(filePath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, status.Errorf(codes.NotFound, "file %q not found", filePath)
		}
		if errors.Is(err, fs.ErrPermission) {
			return nil, status.Errorf(codes.PermissionDenied, "cannot delete file %q", req.Name)
		}
		log.Printf("failed to delete file %q: %v", filePath, err)
		return nil, status.Error(codes.Internal, "unknown error")
	}

	return &proto.SuccessResponse{Success: true}, nil
}

func handleRenamePath(req *proto.RenamePathRequest) (*proto.SuccessResponse, error) {
	newPath, err := validateAndCleanPath(req.Path)
	if err != nil {
		return nil, err
	}
	_, err = probePathExists(newPath)
	if err != nil {
		return nil, err
	}
	if req.NewName == "" {
		return nil, status.Errorf(codes.InvalidArgument, "new name cannot be empty")
	}
	if req.OldName == "" {
		return nil, status.Errorf(codes.InvalidArgument, "old name cannot be empty")
	}
	if req.NewName == req.OldName {
		return nil, status.Errorf(codes.InvalidArgument, "new name cannot be the same as old name")
	}
	if err := validateEntryName(req.OldName); err != nil {
		return nil, err
	}
	if err := validateEntryName(req.NewName); err != nil {
		return nil, err
	}
	oldFilePath := filepath.Join(newPath, req.OldName)
	newFilePath := filepath.Join(newPath, req.NewName)
	_, err = probePathExists(oldFilePath)
	if err != nil {
		return nil, err
	}
	// if new path an value, it means it must exists
	probeNewPath, _ := probePathExists(newFilePath)

	if probeNewPath != nil {
		return nil, status.Errorf(codes.AlreadyExists, "rename path %q already exists", req.NewName)
	}
	err = os.Rename(oldFilePath, newFilePath)
	if err != nil {
		return nil, err
	}

	return &proto.SuccessResponse{Success: true}, nil

}
