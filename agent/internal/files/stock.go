package files

import (
	"io/fs"
	"syscall"

	"hazyio.com/servergram/shared/gen/proto"
)

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
