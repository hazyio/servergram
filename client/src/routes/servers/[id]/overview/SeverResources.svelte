<script lang="ts">
  import CardTitle from "#lib/components/ui/card/card-title.svelte";
  import Card from "#lib/components/ui/card/card.svelte";
  import { CardContent, CardHeader } from "#lib/components/ui/card/index.ts";
  import { Progress } from "#lib/components/ui/progress/index.js";
</script>

<Card>
  <CardHeader><CardTitle>Resources</CardTitle></CardHeader>
  <CardContent class="space-y-6">
    <div class="space-y-2">
      <div class="flex">
        <span class="grow">Cpu</span>
        <span class="text-[12px] text-muted-foreground">50% / 4 cores</span>
      </div>
      <Progress value={50} />
    </div>
    <div class="space-y-2">
      <div class="flex">
        <span class="grow">Memory</span>
        <div class="text-[12px] text-muted-foreground">
          <span>70%</span> <span>(23GB/32GB)</span>
        </div>
      </div>
      <Progress value={70} />
    </div>
    <div class="space-y-2">
      <div class="flex">
        <span class="grow">Disk</span>
        <span class="text-[12px] text-muted-foreground">80% (400GB/500GB)</span>
      </div>
      <Progress value={80} />
    </div>
    <div class="space-y-2">
      <div class="flex">
        <span class="grow">Swap</span>
        <span class="text-[12px] text-muted-foreground">60% (6GB/10GB)</span>
      </div>
      <Progress value={60} />
    </div>
  </CardContent>
</Card>
