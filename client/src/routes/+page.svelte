<script lang="ts">
  import Header from "#lib/components/base/header/header.svelte";
  import { NewRequestClient } from "#lib/api/request.ts";
  $effect(() => {
    async function fetchData() {
      const { data, error } = await NewRequestClient().GET("/api/files/health");
    }
    fetchData();
  });
</script>

<main>
  <Header />
</main>
