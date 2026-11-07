package files

import (
	"errors"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"hazyio.com/servergram/shared/gen/proto"
)

func validateAndCleanPath(path string) (string, error) {
	if !filepath.IsAbs(path) {
		return "", status.Errorf(codes.InvalidArgument, "directory must be absolute: %q", path)
	}
	clean := filepath.Clean(path)
	info, err := os.Stat(clean)
	if err != nil || !info.IsDir() {
		return "", status.Errorf(codes.InvalidArgument, "invalid directory %q", path)
	}
	return clean, nil
}

func toKind(info fs.FileInfo) proto.FileEntryKind {
	if info.IsDir() {
		return proto.FileEntryKind_DIR
	}
	if info.Mode()&fs.ModeSymlink != 0 {
		return proto.FileEntryKind_SYMLINK
	}
	if info.Mode().IsRegular() {
		return proto.FileEntryKind_FILE
	}
	return proto.FileEntryKind_OTHER
}
func ownerIDs(info fs.FileInfo) (uid, gid uint32) {
	if st, ok := info.Sys().(*syscall.Stat_t); ok {
		return st.Uid, st.Gid
	}
	return 0, 0
}

func toEntry(info fs.FileInfo) *proto.FileEntry {
	name := info.Name()
	size := info.Size()
	mode := uint32(info.Mode())
	modifiedUnixNs := info.ModTime().UnixNano()
	uid, gid := ownerIDs(info)

	return &proto.FileEntry{
		Name:           name,
		Size:           size,
		Kind:           toKind(info),
		Mode:           mode,
		Uid:            uid,
		Gid:            gid,
		ModifiedUnixNs: modifiedUnixNs,
		LinkBroken:     false,
		IsHidden:       name[0] == '.',
	}
}

func stripDirFromPath(dir, path string) string {
	// replace dir
	s := strings.Replace(path, dir, "", 1)
	return strings.TrimPrefix(s, "/")
}

// validateEntryName rejects names that could escape the target directory or
// produce an unusable path.
func validateEntryName(name string) error {
	if name == "" {
		return status.Error(codes.InvalidArgument, "file name cannot be empty")
	}
	if name == "." || name == ".." {
		return status.Errorf(codes.InvalidArgument, "invalid file name %q", name)
	}
	if len(name) > 255 {
		return status.Errorf(codes.InvalidArgument, "file name %q is too long", name)
	}
	if strings.ContainsAny(name, `/\`) || strings.ContainsRune(name, '\x00') || strings.ContainsRune(name, ':') || strings.ContainsRune(name, '*') || strings.ContainsRune(name, '?') || strings.ContainsRune(name, '"') || strings.ContainsRune(name, '<') || strings.ContainsRune(name, '>') || strings.ContainsRune(name, '|') || strings.Contains(name, "..") {
		return status.Errorf(codes.InvalidArgument, "invalid file name %q", name)
	}

	return nil
}

func probePathExists(path string) (fs.FileInfo, error) {
	info, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, status.Errorf(codes.NotFound, "path %q not found", path)
		}
		if errors.Is(err, fs.ErrPermission) {
			return nil, status.Errorf(codes.PermissionDenied, "cannot access path %q", path)
		}
		log.Printf("failed to access path %q: %v", path, err)
		return nil, status.Error(codes.Internal, "unknown error")
	}
	return info, nil
}
