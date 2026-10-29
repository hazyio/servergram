<script lang="ts">
  import * as Popover from "#lib/components/ui/popover/index.js";
  import { Button } from "../../ui/button";
  let open = $state(false);
</script>

<div
  class=" fixed bottom-2 right-2 rounded-full p-3 h-15 w-15 flex justify-center items-center"
>
  <Popover.Root {open} onOpenChange={(value) => {}}>
    <Popover.Trigger
      ><Button
        title="connected services"
        variant="outline"
        class="h-15 w-15 rounded-full"
      >
        3
      </Button></Popover.Trigger
    >
    <Popover.Content align="end"
      >Place content for the popover here.</Popover.Content
    >
  </Popover.Root>
</div>
