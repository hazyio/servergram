<script lang="ts">
  import TimeAgo from "#lib/components/base/TimeAgo.svelte";
  import { Badge } from "#lib/components/ui/badge/index.ts";
  import * as Card from "#lib/components/ui/card/index.js";
  import type { components } from "#lib/types/request-schema.js";
  import Monitor from "@lucide/svelte/icons/monitor";
  import Smartphone from "@lucide/svelte/icons/smartphone";
  import Tablet from "@lucide/svelte/icons/tablet";
  import SessionAction from "./SessionAction.svelte";

  let {
    session,
    current = false,
    onRemoved,
  }: {
    session: components["schemas"]["Session"];
    current?: boolean;
    onRemoved?: () => void;
  } = $props();

  let deviceIcon = $derived.by(() => {
    const os = (session.os ?? "").toLowerCase();
    if (os.includes("ipad") || os.includes("tablet")) return Tablet;
    if (
      os.includes("android") ||
      os.includes("ios") ||
      os.includes("iphone") ||
      os.includes("phone")
    )
      return Smartphone;
    return Monitor;
  });
</script>

<Card.Root class="h-full">
  <Card.Header>
    <Card.Title class="flex items-center gap-2">
      {@const Icon = deviceIcon}
      <Icon class="size-4 shrink-0 text-muted-foreground" />
      <span class="truncate">{session.os || "Unknown OS"}</span>
      {#if current}
        <Badge variant="secondary" class="ms-auto">This device</Badge>
      {/if}
    </Card.Title>
    <Card.Description
      class="line-clamp-2 wrap-break-word"
      title={session.user_agent}
    >
      {session.user_agent || "Unknown user agent"}
    </Card.Description>
  </Card.Header>
  <Card.Content class="mt-auto grid gap-1 text-sm text-muted-foreground">
    <p>
      Signed in
      {#if session.created_at}
        <TimeAgo datetime={session.created_at} />
      {:else}
        unknown
      {/if}
    </p>
    <p>
      Expires
      {#if session.expires_at}
        <TimeAgo datetime={session.expires_at} />
      {:else}
        unknown
      {/if}
    </p>
  </Card.Content>
  <Card.Footer>
    <SessionAction id={session.id} os={session.os} {current} {onRemoved} />
  </Card.Footer>
</Card.Root>
