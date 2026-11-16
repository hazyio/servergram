package files

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"log"

	"github.com/fsnotify/fsnotify"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"hazyio.com/servergram/agent/internal/common"
	"hazyio.com/servergram/shared/gen/proto"
)

func handleListFiles(req *proto.ListFilesRequest, stream proto.FilesService_ListFilesServer) error {

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
	batch := make([]*proto.FileEntry, 0, 50)
	for _, e := range entries {
		if stream.Context().Err() != nil {
			return stream.Context().Err()
		}
		if stripDirFromPath(newPath, e.Name())[0] == '.' && !req.IncludeHidden {
			continue
		}
		info, _ := e.Info()
		batch = append(batch, toEntry(info))
		if len(batch) == 50 {
			if err := stream.Send(&proto.ListFilesChunkResponse{Entries: batch, ResolvedPath: newPath}); err != nil {
				return err
			}
			batch = make([]*proto.FileEntry, 0, 50)
		}
	}
	if err := stream.Send(&proto.ListFilesChunkResponse{Entries: batch, ResolvedPath: newPath, Done: true}); err != nil {
		return err
	}

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
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to access directory %q: %v", req.Path, err))
	}
	if !info.IsDir() {
		return nil, status.Errorf(codes.InvalidArgument, "path %q is not a directory", req.Path)
	}
	err = os.RemoveAll(newPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, status.Errorf(codes.NotFound, "directory %q not found", req.Path)
		}
		if errors.Is(err, fs.ErrPermission) {
			return nil, status.Errorf(codes.PermissionDenied, "cannot remove directory %q", req.Path)
		}
		log.Printf("failed to remove directory %q: %v", req.Path, err)
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to remove directory %q: %v", req.Path, err))
	}

	return &proto.SuccessResponse{Success: true}, nil
}

func handleWatchDirectory(req *proto.WatchDirectoryRequest, stream proto.FilesService_WatchDirectoryServer) error {
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
func handlePathProperties(req *proto.PathPropertiesRequest) (*proto.PathPropertiesResponse, error) {
	newPath, err := validateAndCleanPath(req.Path)
	if err != nil {
		return nil, err
	}
	_, err = probePathExists(newPath)
	if err != nil {
		return nil, err
	}
	readPath := filepath.Join(newPath, req.Name)
	entry, err := probePathExists(readPath)
	if err != nil {
		return nil, err
	}
	return &proto.PathPropertiesResponse{Entry: toEntry(entry)}, nil
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

func handleReadFile(req *proto.ReadFileRequest, stream proto.FilesService_ReadFileServer) error {

	newPath, err := validateAndCleanPath(req.Path)

	if err != nil {
		return err
	}
	_, err = probePathExists(newPath)
	if err != nil {

		return err
	}
	filePath := filepath.Join(newPath, req.Name)
	file, err := probePathExists(filePath)
	if err != nil {
		return err
	}
	if file.IsDir() {

		return status.Errorf(codes.InvalidArgument, "path '%q' is of type directory, a file was expected ", req.Path)
	}

	if !req.IsDownload && file.Size() > req.MaxBytesSize {
		// Do not stream files larger than the requested max size, disregard if IsDownload is true
		return status.Errorf(codes.OutOfRange, "File cannot be read, too large, set IsDownload=true to download %q", req.Name)
	}

	openFile, err := os.Open(filePath)
	if err != nil {
		return err

	}
	defer openFile.Close()
	if !req.IsDownload {
		// if download is not requested, verify it's readable text
		isText, err := common.IsReadableText(openFile)

		if err != nil {
			return err
		}
		if !isText {
			return status.Errorf(codes.FailedPrecondition, "file %q is not readable text, set IsDownload=true instead to download file", req.Name)
		}
		if _, err := openFile.Seek(0, io.SeekStart); err != nil {
			// reset pointer, IsReadableText moved it.
			return status.Error(codes.Internal, "seek failed")
		}
	}

	buf := make([]byte, 64*1024)
	for {
		select {
		case <-stream.Context().Done():
			return stream.Context().Err()
		default:
			n, err := openFile.Read(buf)
			if err != nil && err != io.EOF {
				return err
			}
			if errors.Is(err, io.EOF) {
				// send done and close stream
				stream.Send(&proto.ReadFileResponse{Done: true})
				return nil
			}
			if n > 0 {
				if sendErr := stream.Send(&proto.ReadFileResponse{Data: buf[:n], Done: false}); sendErr != nil {
					return sendErr
				}
			}
		}

	}

}

func handleWriteFile(stream proto.FilesService_WriteFileServer) error {
	// First message must be metadata
	first, err := stream.Recv()
	if err != nil {
		return err
	}
	meta := first.GetMeta()
	if meta == nil {
		return status.Error(codes.InvalidArgument, "first message must be meta")
	}

	newPath, err := validateAndCleanPath(meta.Path)

	if err != nil {
		return err
	}
	_, err = probePathExists(newPath)
	if err != nil {
		return err
	}
	if err := validateEntryName(meta.Name); err != nil {
		return err
	}
	dest := filepath.Join(newPath, filepath.Base(meta.Name))

	fileExists, _ := probePathExists(dest)
	if fileExists != nil {
		return status.Errorf(codes.AlreadyExists, "file name %q, already exists", dest)
	}

	// Write to a temp file, rename at the end, so failed uploads leave no partial file
	tmp, err := os.CreateTemp(filepath.Dir(dest), ".upload-*")
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	defer os.Remove(tmp.Name())
	w := bufio.NewWriterSize(tmp, 128*1024)
	var written int64
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break // client finished sending
		}
		if err != nil {
			return err // client cancelled or connection died
		}
		n, err := w.Write(msg.GetChunk())
		if err != nil {
			return status.Error(codes.Internal, err.Error())
		}
		written += int64(n)
	}
	if err := w.Flush(); err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	if err := tmp.Close(); err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	if err := os.Rename(tmp.Name(), dest); err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	return stream.SendAndClose(&proto.WriteFileResponse{BytesWritten: written})

}
