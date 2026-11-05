package files

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"log"

	"github.com/fsnotify/fsnotify"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"hazyio.com/servergram/shared/gen/proto"
)

func stripDirFromPath(dir, path string) string {
	return strings.Replace(strings.Replace(path+"/", dir, "", 1), dir, "", 1)
}

func handleListFiles(req *proto.ListFilesRequest, stream proto.Files_ListFilesServer) error {
	entries, err := os.ReadDir(req.Dir)
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
	if fs.ValidPath(req.Path) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid directory path %q", req.Path)
	}
	newDirPath := filepath.Join(req.Path, string(req.Name))
	err := os.Mkdir(newDirPath, 0755)
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
	if fs.ValidPath(req.Path) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid directory path %q", req.Path)
	}
	info, err := os.Stat(req.Path)
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
	if fs.ValidPath(req.Path) {
		return status.Errorf(codes.InvalidArgument, "invalid directory path %q", req.Path)
	}
	info, err := os.Stat(req.Path)
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

	for {
		select {
		case event, ok := <-watcher.Events:
			if ok {
				name := stripDirFromPath(event.Name, req.Path) // remove the directory from the path to save bytes
				if event.Has(fsnotify.Create) {
					log.Println("created:", event.Name)
					stream.Send(&proto.WatchDirectoryResponse{Kind: proto.WatchDirectoryKind_CREATED, Path: name})
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
