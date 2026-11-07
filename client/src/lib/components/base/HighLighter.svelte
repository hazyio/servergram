<script lang="ts">
  import * as Select from "#lib/components/ui/select/index.js";

  import { useLocalStorage } from "#lib/local-storage.ts";
  import { darkModernTheme } from "#lib/themes/dark-modern.ts";
  import { createHighlighterCore, type HighlighterCore } from "shiki/core";
  import { createOnigurumaEngine } from "shiki/engine/oniguruma";
  import { bundledLanguagesInfo } from "shiki/langs";
  import { bundledThemesInfo } from "shiki/themes";

  let {
    code,
    fileName = "",
    class: className = "",
  }: { code: string; fileName?: string; class?: string } = $props();
  const DEFAULT_THEME = "dark-modern";
  // The previous default was persisted automatically on mount, so treat it as
  // unset and let existing sessions pick up the new default once.
  const LEGACY_DEFAULT_THEME = "github-dark";

  let html = $state("");
  let theme = $state(DEFAULT_THEME);

  const themes = [
    { id: "dark-modern", label: "Dark Modern" },
    { id: "github-dark", label: "GitHub Dark" },
    { id: "github-light", label: "GitHub Light" },
    { id: "one-dark-pro", label: "One Dark Pro" },
    { id: "one-light", label: "One Light" },
    { id: "dracula", label: "Dracula" },
    { id: "nord", label: "Nord" },
    { id: "catppuccin-mocha", label: "Catppuccin Mocha" },
    { id: "catppuccin-latte", label: "Catppuccin Latte" },
    { id: "solarized-dark", label: "Solarized Dark" },
    { id: "solarized-light", label: "Solarized Light" },
    { id: "vitesse-dark", label: "Vitesse Dark" },
    { id: "vitesse-light", label: "Vitesse Light" },
    { id: "tokyo-night", label: "Tokyo Night" },
    { id: "rose-pine", label: "Rosé Pine" },
    { id: "rose-pine-dawn", label: "Rosé Pine Dawn" },
  ];

  const storage = useLocalStorage();

  const languageLoaders = new Map(
    bundledLanguagesInfo.map((language) => [language.id, language.import]),
  );
  const themeLoaders = new Map(
    bundledThemesInfo.map((availableTheme) => [
      availableTheme.id,
      availableTheme.import,
    ]),
  );

  // File extensions mapped to Shiki language ids.
  const languageByExtension: Record<string, string> = {
    // Web
    html: "html",
    htm: "html",
    css: "css",
    scss: "scss",
    sass: "sass",
    less: "less",
    js: "javascript",
    mjs: "javascript",
    cjs: "javascript",
    jsx: "jsx",
    ts: "typescript",
    mts: "typescript",
    cts: "typescript",
    tsx: "tsx",
    vue: "vue",
    svelte: "svelte",
    astro: "astro",
    // Data / config
    json: "json",
    jsonc: "jsonc",
    json5: "json5",
    yaml: "yaml",
    yml: "yaml",
    toml: "toml",
    ini: "ini",
    cfg: "ini",
    conf: "ini",
    properties: "ini",
    env: "dotenv",
    xml: "xml",
    svg: "xml",
    csv: "csv",
    tsv: "csv",
    sql: "sql",
    graphql: "graphql",
    gql: "graphql",
    proto: "proto",
    md: "markdown",
    markdown: "markdown",
    mdx: "mdx",
    tex: "tex",
    diff: "diff",
    patch: "diff",
    // Systems / compiled
    go: "go",
    rs: "rust",
    c: "c",
    h: "c",
    cpp: "cpp",
    cc: "cpp",
    cxx: "cpp",
    hpp: "cpp",
    hh: "cpp",
    hxx: "cpp",
    cs: "csharp",
    java: "java",
    kt: "kotlin",
    kts: "kotlin",
    scala: "scala",
    swift: "swift",
    m: "objective-c",
    mm: "objective-cpp",
    zig: "zig",
    nim: "nim",
    dart: "dart",
    lua: "lua",
    luau: "luau",
    asm: "asm",
    s: "asm",
    // Scripting
    py: "python",
    pyw: "python",
    rb: "ruby",
    php: "php",
    pl: "perl",
    pm: "perl",
    r: "r",
    ex: "elixir",
    exs: "elixir",
    erl: "erlang",
    hrl: "erlang",
    hs: "haskell",
    clj: "clojure",
    cljs: "clojure",
    cljc: "clojure",
    sh: "shellscript",
    bash: "shellscript",
    zsh: "shellscript",
    fish: "fish",
    ps1: "powershell",
    psd1: "powershell",
    psm1: "powershell",
    bat: "bat",
    cmd: "bat",
    // Build / infra
    dockerfile: "docker",
    tf: "terraform",
    tfvars: "terraform",
    hcl: "hcl",
    cmake: "cmake",
    mk: "make",
    make: "make",
    nix: "nix",
    prisma: "prisma",
    vim: "viml",
  };

  const languageByFileName: Record<string, string> = {
    dockerfile: "docker",
    containerfile: "docker",
    makefile: "make",
    gnumakefile: "make",
    "cmakelists.txt": "cmake",
    ".env": "dotenv",
    ".bashrc": "shellscript",
    ".bash_profile": "shellscript",
    ".zshrc": "shellscript",
    ".profile": "shellscript",
    "go.mod": "go",
    "go.sum": "go",
  };

  function detectLanguage(fileName: string) {
    const base = fileName.split(/[\\/]/).pop()?.toLowerCase() ?? "";
    if (!base) return "text";

    const exact = languageByFileName[base];
    if (exact) return exact;
    if (base.startsWith("dockerfile") || base.startsWith("containerfile")) {
      return "docker";
    }
    if (base === "makefile" || base.startsWith("makefile.")) return "make";

    const dot = base.lastIndexOf(".");
    const extension = dot >= 0 ? base.slice(dot + 1) : "";
    return languageByExtension[extension] ?? "text";
  }

  let highlighterPromise: Promise<HighlighterCore> | undefined;
  const loadedLanguages = new Set<string>();
  const loadedThemes = new Set([DEFAULT_THEME]);

  function getHighlighter() {
    highlighterPromise ??= createHighlighterCore({
      themes: [darkModernTheme],
      langs: [],
      engine: createOnigurumaEngine(import("shiki/wasm")),
    });
    return highlighterPromise;
  }

  async function ensureLanguage(
    highlighter: HighlighterCore,
    language: string,
  ) {
    if (language === "text" || loadedLanguages.has(language)) return;

    const loader = languageLoaders.get(language);
    if (!loader) return;

    await highlighter.loadLanguage(loader);
    loadedLanguages.add(language);
  }

  async function ensureTheme(highlighter: HighlighterCore, themeName: string) {
    if (loadedThemes.has(themeName)) return;

    const loader = themeLoaders.get(themeName);
    if (!loader) return;

    await highlighter.loadTheme(loader);
    loadedThemes.add(themeName);
  }

  function escapeHtml(value: string) {
    return value
      .replaceAll("&", "&amp;")
      .replaceAll("<", "&lt;")
      .replaceAll(">", "&gt;");
  }

  let themeLoaded = false;

  $effect(() => {
    if (themeLoaded) return;
    themeLoaded = true;

    const savedTheme = storage.get<string>("file-viewer-theme", DEFAULT_THEME);
    const initialTheme =
      savedTheme === LEGACY_DEFAULT_THEME ? DEFAULT_THEME : savedTheme;
    if (themes.some((option) => option.id === initialTheme)) {
      theme = initialTheme;
    }
  });

  $effect(() => {
    storage.save("file-viewer-theme", theme);
  });

  let renderId = 0;

  $effect(() => {
    const currentCode = code;
    const language = detectLanguage(fileName);
    const currentTheme = theme;
    const currentRenderId = ++renderId;

    void (async () => {
      try {
        const highlighter = await getHighlighter();
        await Promise.all([
          ensureLanguage(highlighter, language),
          ensureTheme(highlighter, currentTheme),
        ]);
        const result = highlighter.codeToHtml(currentCode, {
          lang: language,
          theme: currentTheme,
        });
        if (currentRenderId === renderId) html = result;
      } catch (error) {
        console.error("Failed to highlight file", error);
        if (currentRenderId === renderId) {
          html = `<pre class="shiki ${currentTheme}">${escapeHtml(currentCode)}</pre>`;
        }
      }
    })();
  });
</script>

<div
  class={`relative min-w-0 ${className} [&_pre]:m-0 [&_pre]:min-h-[inherit] [&_pre]:w-full [&_pre]:max-w-full [&_pre]:overflow-x-hidden [&_pre]:whitespace-pre-wrap [&_pre]:break-all [&_pre]:p-4 [&_pre]:pt-12 [&_pre]:text-xs [&_pre]:leading-relaxed`}
>
  <div class="absolute top-2 right-2 z-10">
    <Select.Root type="single" bind:value={theme}>
      <Select.Trigger
        size="sm"
        class="w-40 border bg-popover! text-popover-foreground! shadow-md"
        aria-label="Syntax highlighting theme"
      >
        <Select.Value />
      </Select.Trigger>
      <Select.Content>
        {#each themes as option (option.id)}
          <Select.Item value={option.id} label={option.label}>
            {option.label}
          </Select.Item>
        {/each}
      </Select.Content>
    </Select.Root>
  </div>

  {@html html}
</div>
