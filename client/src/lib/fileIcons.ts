const icons = import.meta.glob("./assets/icons/*.svg", {
  eager: true,
  query: "?url",
  import: "default",
}) as Record<string, string>;

const iconUrl = (name: string) => icons[`./assets/icons/${name}.svg`];

const byExtension: Record<string, string> = {
  js: "javascript",
  mjs: "javascript",
  ts: "typescript",
  html: "html",
  css: "css",
  go: "go",
  svelte: "svelte",
  json: "json",
  md: "markdown",
  log: "log",
  txt: "document",
  yaml: "yaml",
  yml: "yaml",
  toml: "toml",
  xml: "xml",
  py: "python",
  rs: "rust",
  php: "php",
  rb: "ruby",
  java: "java",
  c: "c",
  cpp: "cpp",
  cs: "csharp",
  vue: "vue",
  sh: "console",
  ps1: "powershell",
  pdf: "pdf",
  png: "image",
  jpg: "image",
  jpeg: "image",
  gif: "image",
  svg: "image",
  webp: "image",
  mp4: "video",
  mkv: "video",
  mp3: "audio",
  wav: "audio",
  zip: "zip",
  gz: "zip",
  tar: "zip",
  "7z": "zip",
  pem: "certificate",
  crt: "certificate",
  key: "key",
  lock: "lock",
  sql: "database",
  woff: "font",
  woff2: "font",
  ttf: "font",
};

const byFileName: Record<string, string> = {
  ".env": "tune",
  dockerfile: "docker",
  license: "license",
  "readme.md": "readme",
  "go.mod": "go-mod",
};

const byFolderName: Record<string, string> = {
  src: "folder-src",
  dist: "folder-dist",
  logs: "folder-log",
  log: "folder-log",
  assets: "folder-resource",
  images: "folder-images",
  public: "folder-public",
  config: "folder-config",
  docs: "folder-docs",
  test: "folder-test",
  tests: "folder-test",
  tmp: "folder-temp",
  temp: "folder-temp",
  node_modules: "folder-node",
  ".git": "folder-git",
};

export function getFileIcon(name: string, isDir: boolean): string {
  const lower = name.toLowerCase();
  if (isDir) {
    return iconUrl(byFolderName[lower] ?? "folder-base") ?? iconUrl("folder-base");
  }
  const ext = lower.includes(".") ? lower.split(".").pop()! : "";
  return (
    iconUrl(byFileName[lower] ?? byExtension[ext] ?? "document") ??
    iconUrl("document")
  );
}
