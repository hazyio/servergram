import manifest from "material-icon-theme/dist/material-icons.json";

interface Manifest {
  iconDefinitions: Record<string, { iconPath: string }>;
  file: string;
  folder: string;
  folderExpanded: string;
  fileNames: Record<string, string>;
  fileExtensions: Record<string, string>;
  folderNames: Record<string, string>;
  folderNamesExpanded: Record<string, string>;
}
const m = manifest as unknown as Manifest;

const iconFiles = import.meta.glob(
  "/node_modules/material-icon-theme/icons/*.svg",
  { eager: true, query: "?url", import: "default" },
) as Record<string, string>;

// manifest value -> icon definition id -> svg filename -> bundled url
function urlFor(id: string | undefined): string | undefined {
  const def = id ? m.iconDefinitions[id] : undefined;
  if (!def) return undefined;
  const file = def.iconPath.split("/").pop()!;
  return iconFiles[`/node_modules/material-icon-theme/icons/${file}`];
}

export function getFileIcon(
  name: string,
  isDir: boolean,
  open = false,
): string {
  const lower = name.toLowerCase();

  if (isDir) {
    const named = open
      ? (m.folderNamesExpanded[lower] ?? m.folderNames[lower])
      : m.folderNames[lower];
    return urlFor(named) ?? urlFor(open ? m.folderExpanded : m.folder)!;
  }

  const byName = urlFor(m.fileNames[lower]);
  if (byName) return byName;

  // Longest extension first: "app.test.ts" -> "test.ts" -> "ts"
  const parts = lower.split(".");
  for (let i = 1; i < parts.length; i++) {
    const url = urlFor(m.fileExtensions[parts.slice(i).join(".")]);
    if (url) return url;
  }

  return urlFor(m.file)!;
}
