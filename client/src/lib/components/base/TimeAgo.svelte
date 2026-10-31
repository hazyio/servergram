<script lang="ts">
  let { datetime }: { datetime: string | Date } = $props();

  const relativeTimeFormatter = new Intl.RelativeTimeFormat(undefined, {
    numeric: "auto",
  });
  const units = [
    { name: "year", seconds: 31_536_000 },
    { name: "month", seconds: 2_592_000 },
    { name: "day", seconds: 86_400 },
    { name: "hour", seconds: 3_600 },
    { name: "minute", seconds: 60 },
    { name: "second", seconds: 1 },
  ] as const;

  let now = $state(Date.now());
  let date = $derived(datetime instanceof Date ? datetime : new Date(datetime));
  let isValid = $derived(!Number.isNaN(date.getTime()));
  let relativeTime = $derived.by(() => {
    const difference = (date.getTime() - now) / 1000;
    const unit = units.find(({ seconds }) => Math.abs(difference) >= seconds);
    const selectedUnit = unit ?? units[units.length - 1];

    return relativeTimeFormatter.format(
      Math.round(difference / selectedUnit.seconds),
      selectedUnit.name,
    );
  });

  $effect(() => {
    const interval = setInterval(() => {
      now = Date.now();
    }, 60_000);

    return () => clearInterval(interval);
  });
</script>

{#if isValid}
  <time datetime={date.toISOString()} title={date.toLocaleString()}>
    {relativeTime}
  </time>
{:else}
  <span>{String(datetime)}</span>
{/if}
