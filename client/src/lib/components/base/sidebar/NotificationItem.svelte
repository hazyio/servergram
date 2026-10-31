<script lang="ts" module>
  export interface NotificationItemProps {
    id: number;
    title: string;
    message: string;
  }
</script>

<script lang="ts">
  let { id, title, message }: NotificationItemProps = $props();
</script>

<article id={`notification-${id}`} class="mt-2 rounded-md border bg-muted p-3">
  <h4 class="text-sm font-medium">{title}</h4>
  <p class="mt-1 wrap-break-word text-sm text-muted-foreground">{message}</p>
</article>
