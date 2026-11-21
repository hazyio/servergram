import {
  MODE_CHAR_DEV,
  MODE_DEVICE,
  MODE_DIR,
  MODE_NAMED_PIPE,
  MODE_SETGID,
  MODE_SETUID,
  MODE_SOCKET,
  MODE_STICKY,
  MODE_SYMLINK,
  type FileEntry,
} from "./types";

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

export const STREAM_CONNECT_TIMEOUT_MS = 20_000;

/**
 * Stops an EventSource from retrying forever when it never connects.
 * Browsers reconnect EventSources automatically, so closing it on timeout is
 * the only way to give up on the connection.
 */
export function addEventSourceTimeout(
  source: EventSource,
  onTimeout: () => void,
  timeoutMs = STREAM_CONNECT_TIMEOUT_MS,
): void {
  let timer: ReturnType<typeof setTimeout> | undefined = setTimeout(() => {
    timer = undefined;
    if (source.readyState === EventSource.OPEN) return;

    source.close();
    onTimeout();
  }, timeoutMs);

  source.addEventListener(
    "open",
    () => {
      if (timer === undefined) return;
      clearTimeout(timer);
      timer = undefined;
    },
    { once: true },
  );

  const close = source.close.bind(source);
  source.close = () => {
    if (timer !== undefined) {
      clearTimeout(timer);
      timer = undefined;
    }
    close();
  };
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

export function formatMode(entry: FileEntry): string {
  const mode = (entry.mode ?? 0) >>> 0; // force unsigned 32-bit
  const has = (flag: number) => (mode & flag) !== 0;
  console.log(mode);
  let type = "-";
  if (has(MODE_DIR)) type = "d";
  else if (has(MODE_SYMLINK)) type = "l";
  else if (has(MODE_NAMED_PIPE)) type = "p";
  else if (has(MODE_SOCKET)) type = "s";
  else if (has(MODE_DEVICE)) type = has(MODE_CHAR_DEV) ? "c" : "b";

  const triplet = (shift: number, special: boolean, sChar: string) => {
    const r = mode & (1 << (shift + 2)) ? "r" : "-";
    const w = mode & (1 << (shift + 1)) ? "w" : "-";
    const x = (mode & (1 << shift)) !== 0;
    const xc = special ? (x ? sChar : sChar.toUpperCase()) : x ? "x" : "-";
    return r + w + xc;
  };

  return (
    type +
    triplet(6, has(MODE_SETUID), "s") +
    triplet(3, has(MODE_SETGID), "s") +
    triplet(0, has(MODE_STICKY), "t")
  );
}
export function joinPath(dir: string, name: string): string {
  const base = dir.endsWith("/") ? dir.slice(0, -1) : dir;
  return `${base}/${name}`; // dir "/" -> "" + "/name" = "/name"
}
