<script lang="ts">
  import Header from "#lib/components/base/header/header.svelte";
  import { requestClient } from "#lib/api/request.ts";
  $effect(() => {
    async function fetchData() {
      const { data, error } = await requestClient.GET("/api/service/health");
    }
    fetchData();
  });
</script>

<main>
  <Header />
</main>
