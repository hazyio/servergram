import type { EntryKind } from "./types";

export { cn } from "cn";

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export type WithoutChild<T> = T extends { child?: any } ? Omit<T, "child"> : T;
// eslint-disable-next-line @typescript-eslint/no-explicit-any
export type WithoutChildren<T> = T extends { children?: any }
  ? Omit<T, "children">
  : T;
export type WithoutChildrenOrChild<T> = WithoutChildren<WithoutChild<T>>;
export type WithElementRef<T, U extends HTMLElement = HTMLElement> = T & {
  ref?: U | null;
};

export function getWebSocketUrl(path: string = "/ws"): string {
  const url = new URL(path, window.location.href);
  url.protocol = url.protocol === "https:" ? "wss:" : "ws:";
  return url.toString();
}

export function castType<T>(value: unknown, callback: (value: T) => void) {
  callback(value as T);
}
export function formatBytes(bytes: number) {
  if (bytes <= 0) return "0 B";
  const units = ["B", "KB", "MB", "GB", "TB"];
  const i = Math.min(
    Math.floor(Math.log(bytes) / Math.log(1024)),
    units.length - 1,
  );
  return `${new Intl.NumberFormat(undefined, { maximumFractionDigits: 1 }).format(bytes / 1024 ** i)} ${units[i]}`;
}


export function formatMode(mode: number, kind: EntryKind) {
  const type = kind === "directory" ? "d" : kind === "symlink" ? "l" : "-";
  const bits = "rwxrwxrwx";
  let out = type;
  for (let i = 0; i < 9; i++) {
    out += mode & (1 << (8 - i)) ? bits[i] : "-";
  }
  return out;
}