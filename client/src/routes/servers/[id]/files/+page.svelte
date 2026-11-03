<script lang="ts">
  import Button from "#lib/components/ui/button/button.svelte";
  import Input from "#lib/components/ui/input/input.svelte";
  import { ChevronLeft, ChevronRight, Settings } from "@lucide/svelte";
  import type { PageProps } from "./$types";
  import SettingOption from "./SettingOption.svelte";

  let { data }: PageProps = $props();
  let settings = $state({ showHidden: false });
</script>

<div class="flex gap-2">
  <div class="flex gap-2">
    <Button variant="outline"><ChevronLeft /></Button>
    <Button variant="outline"><ChevronRight /></Button>
  </div>
  <Input value="/var/www/app" class="grow-2" />
  <Input value="" placeholder="Search files here..." />
  <SettingOption bind:showHidden={settings.showHidden} />
</div>
