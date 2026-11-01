<script lang="ts">
  import appPaths from "#lib/appPaths.ts";
  import Button from "#lib/components/ui/button/button.svelte";
  import ArrowLeft from "@lucide/svelte/icons/arrow-left";
  import RotateCcw from "@lucide/svelte/icons/rotate-ccw";
  import TriangleAlert from "@lucide/svelte/icons/triangle-alert";
  import { page } from "$app/state";

  const isNotFound = $derived(page.status === 404);
</script>

<svelte:head>
  <title>{page.status} | ServerGram</title>
  <meta
    name="description"
    content="The requested ServerGram page could not be loaded."
  />
</svelte:head>

<main class="grid min-h-[76vh] place-items-center p-6 sm:p-8 lg:p-16">
  <section class="w-full max-w-2xl" aria-labelledby="error-title">
    <p
      class="m-0 flex items-center gap-2.5 font-mono text-[0.7rem] font-semibold tracking-[0.08em] text-muted-foreground"
    >
      <span
        class="size-[0.45rem] rounded-full bg-chart-3 ring-4 ring-chart-3/15"
        aria-hidden="true"
      ></span>
      SERVERGRAM <span class="text-border">/</span> SYSTEM RESPONSE
    </p>

    <div
      class="relative my-6 flex min-h-40 items-center justify-between overflow-hidden border-y border-border bg-[linear-gradient(90deg,color-mix(in_oklch,var(--primary)_5%,transparent)_1px,transparent_1px)] bg-size-[3rem_100%] sm:my-9 sm:min-h-48"
      aria-label="HTTP status {page.status}"
    >
      <div
        class="absolute inset-x-0 top-4 flex justify-between gap-4 font-mono text-[0.6rem] font-semibold tracking-[0.06em] text-muted-foreground sm:text-[0.7rem]"
      >
        <span>HTTP STATUS</span>
        <span>{isNotFound ? "RESOURCE NOT FOUND" : "REQUEST INTERRUPTED"}</span>
      </div>
      <span
        class="font-mono text-[5rem] font-semibold leading-none text-foreground sm:text-[7rem]"
        aria-hidden="true">{page.status}</span
      >
      <span
        class="grid size-13 place-items-center border border-chart-3/35 text-chart-3"
        aria-hidden="true"><TriangleAlert class="size-5" /></span
      >
      <span
        class="absolute right-0 bottom-0 h-0.5 w-20 bg-primary"
        aria-hidden="true"
      ></span>
    </div>

    <h1
      id="error-title"
      class="m-0 font-mono text-2xl font-semibold leading-tight sm:text-[1.75rem]"
    >
      {isNotFound ? "This page isn't here." : "We couldn't load this view."}
    </h1>
    <p class="mt-3 max-w-xl text-[0.95rem] leading-[1.7] text-muted-foreground">
      {isNotFound
        ? "The address may be outdated, or the page may have moved."
        : "Something interrupted this request. Try again, or head back to your dashboard."}
    </p>

    <div class="mt-6 flex flex-wrap gap-2.5">
      <Button href={appPaths.home}>
        <ArrowLeft data-icon="inline-start" />
        Back to dashboard
      </Button>
      <Button variant="outline" onclick={() => location.reload()}>
        <RotateCcw data-icon="inline-start" />
        Try again
      </Button>
    </div>

    <p
      class="mt-10 flex flex-wrap gap-x-3 gap-y-1.5 border-t border-border pt-3 font-mono text-[0.7rem] font-semibold tracking-[0.04em] text-muted-foreground"
    >
      <span>REQUESTED PATH</span>
      <code
        class="min-w-0 wrap-anywhere font-normal tracking-normal text-foreground"
      >
        {page.url.pathname}
      </code>
    </p>
  </section>
</main>
