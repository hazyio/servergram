<script lang="ts">
  import { createHighlighterCore, type HighlighterCore } from "shiki/core";
  import { createOnigurumaEngine } from "shiki/engine/oniguruma";
  import { bundledLanguagesInfo } from "shiki/langs";

  let { code, fileName = "" } = $props();
  let html = $state("");

  const theme = "github-dark";
  const languageLoaders = new Map(
    bundledLanguagesInfo.map((language) => [language.id, language.import]),
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

  function getHighlighter() {
    highlighterPromise ??= createHighlighterCore({
      themes: [import("shiki/themes/github-dark.mjs")],
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

  function escapeHtml(value: string) {
    return value
      .replaceAll("&", "&amp;")
      .replaceAll("<", "&lt;")
      .replaceAll(">", "&gt;");
  }

  let renderId = 0;

  $effect(() => {
    const currentCode = code;
    const language = detectLanguage(fileName);
    const currentRenderId = ++renderId;

    void (async () => {
      try {
        const highlighter = await getHighlighter();
        await ensureLanguage(highlighter, language);
        const result = highlighter.codeToHtml(currentCode, {
          lang: language,
          theme,
        });
        if (currentRenderId === renderId) html = result;
      } catch (error) {
        console.error("Failed to highlight file", error);
        if (currentRenderId === renderId) {
          html = `<pre class="shiki ${theme}">${escapeHtml(currentCode)}</pre>`;
        }
      }
    })();
  });
</script>

{@html html}
