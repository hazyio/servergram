import type { components } from "./request-schema";

export type ThemeType = "light" | "dark" | "system";
export type EntryKind = "file" | "directory" | "symlink";

export type FileEntry = components["schemas"]["FileEntry"];
// Mirrors Go's io/fs.FileMode bit layout
export const KIND_DIR = 1;
export const KIND_SYMLINK = 2;
export const MODE_DIR = 0x80000000; // 1 << 31
export const MODE_SYMLINK = 1 << 27;
export const MODE_DEVICE = 1 << 26;
export const MODE_NAMED_PIPE = 1 << 25;
export const MODE_SOCKET = 1 << 24;
export const MODE_SETUID = 1 << 23;
export const MODE_SETGID = 1 << 22;
export const MODE_CHAR_DEV = 1 << 21;
export const MODE_STICKY = 1 << 20;