<script lang="ts">
  import * as DropdownMenu from "#lib/components/ui/dropdown-menu/index.js";
  import Database from "@lucide/svelte/icons/database";
  import Server from "@lucide/svelte/icons/server";
  import { Button } from "../../ui/button";
  import appPaths from "#lib/appPaths.ts";
</script>

<DropdownMenu.Root>
  <DropdownMenu.Trigger>
    <Button>Add Service</Button>
  </DropdownMenu.Trigger>
  <DropdownMenu.Content>
    <DropdownMenu.Group>
      <DropdownMenu.Item>
        <a href={appPaths.servers.add} class="gap-2 flex items-center">
          <Server />
          Server
        </a>
      </DropdownMenu.Item>
      <DropdownMenu.Item>
        <a href={appPaths.databases.add} class="gap-2 flex items-center">
          <Database />
          Database
        </a>
      </DropdownMenu.Item>
    </DropdownMenu.Group>
  </DropdownMenu.Content>
</DropdownMenu.Root>
