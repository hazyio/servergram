<script lang="ts">
  import type { ThemeType } from "#lib/types/index.ts";
  import * as Popover from "#lib/components/ui/popover/index.js";
  import { Button } from "#lib/components/ui/button/index.js";

  import Check from "@lucide/svelte/icons/check";
  import Sun from "@lucide/svelte/icons/sun";
  import Moon from "@lucide/svelte/icons/moon";
  import SunMoon from "@lucide/svelte/icons/sun-moon";

  import { useLocalStorage } from "#lib/local-storage.ts";
  let currentTheme: ThemeType = $state("system");

  function applyTheme(theme: ThemeType) {
    const resolvedTheme =
      theme === "system"
        ? window.matchMedia("(prefers-color-scheme: dark)").matches
          ? "dark"
          : "light"
        : theme;
    const htmlElement = document.documentElement;
    htmlElement.classList.remove("dark", "light", "system");
    htmlElement.classList.add(resolvedTheme);
  }
  function setTheme(newTheme: ThemeType) {
    applyTheme(newTheme);
    useLocalStorage().save("theme", newTheme);
    currentTheme = newTheme;
  }
  $effect(() => {
    const getSavedTheme = useLocalStorage().get("theme", "system");
    setTheme(getSavedTheme);
  });
</script>

<div>
  <Popover.Root>
    <Popover.Trigger>
      <Button variant="outline">
        {#if currentTheme === "light"}
          <Sun class={"h-7 w-7"} />
        {:else if currentTheme === "dark"}
          <Moon class={"h-7 w-7"} />
        {:else}
          <SunMoon class={"h-7 w-7"} />
        {/if}
      </Button>
    </Popover.Trigger>
    <Popover.Content class={"max-w-40 p-1"}>
      <div class={"flex flex-col gap-2"}>
        <button
          class={"flex cursor-pointer items-center gap-2 rounded-md px-2 py-2 hover:bg-muted/60"}
          onclick={() => setTheme("light")}
        >
          <Sun class={"h-5 w-5"} />
          <div class={"grow"}>Light</div>
          <Check
            class={`h-5 w-5 ${currentTheme === "light" ? "visible" : "invisible"}`}
          />
        </button>
        <button
          class={"flex cursor-pointer items-center gap-2 rounded-md px-2 py-2 hover:bg-muted/60"}
          onclick={() => setTheme("dark")}
        >
          <Moon class={"h-5 w-5"} />
          <div class={"grow"}>Dark</div>
          <Check
            class={`h-5 w-5 ${currentTheme === "dark" ? "visible" : "invisible"}`}
          />
        </button>
        <button
          class={"flex cursor-pointer items-center gap-2 rounded-md px-2 py-2 hover:bg-muted/60"}
          onclick={() => setTheme("system")}
        >
          <SunMoon class={"h-5 w-5"} />
          <div class={"grow"}>System</div>
          <Check
            class={`h-5 w-5 ${currentTheme === "system" ? "visible" : "invisible"}`}
          />
        </button>
      </div>
    </Popover.Content>
  </Popover.Root>
</div>
