<script lang="ts">
  import { createHighlighterCore } from 'shiki/core';
  import { createOnigurumaEngine } from 'shiki/engine/oniguruma';

  let { code,  fileName = '' } = $props();
  let html = $state('');

  $effect(() => {
    (async () => {
      const hl = await createHighlighterCore({
        themes: [import('shiki/themes/github-dark.mjs')],
        langs: [import('shiki/langs/rust.mjs')],
        engine: createOnigurumaEngine(import('shiki/wasm'))
      });
      html = hl.codeToHtml(code, { lang, theme: 'github-dark' });
    })();
  });
</script>

{@html html}