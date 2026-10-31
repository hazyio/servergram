// services.svelte.ts
import {
  type ConnectServerWSMessage,
  type ServiceWSMessage,
  ConnectToSSH
} from "#lib/types/websockets/index.ts";
import { getWebSocketUrl } from "#lib/utils.ts";
import { toast } from "svelte-sonner";
import type { ServiceItemProps } from "./ServiceItem.svelte";
let ws: WebSocket | undefined;

export const servicesState = $state({
  services: [] as ServiceItemProps[],
  connecting: true,
  error: null as string | null,
});

export function startWebSocket<T>(
  onOpen: () => void,
  onMessage: (message: T) => void,
) {
  if (ws?.readyState === WebSocket.OPEN) {
    console.log("WebSocket already open");
    return;
  }
  ws = new WebSocket(getWebSocketUrl("/ws/services"));
  ws.onopen = () => {
    onOpen();
    servicesState.connecting = false;
    servicesState.error = null;
  };
  ws.onerror = (e) => {
    console.log(e);
    servicesState.error = "Could not connect, webSocket error";
  };
  ws.onmessage = (e) => {
    onMessage(JSON.parse(e.data));
  };
}

export function stopWebSocket() {
  ws?.close(1000, "bye");
  ws = undefined;
}
export function checkWs(): boolean {
  if (!ws || ws.readyState !== WebSocket.OPEN) {
    servicesState.error = "WebSocket is not connected";
    toast(
      `WebSocket is not connected, click red floating circle to reconnect.`,
      {
        position: "top-right",
      },
    );
    return false;
  }
  return true;
}

export function startServerConnection(serviceId: number) {
  if (!checkWs()) return;
  let data: ServiceWSMessage = {
    action: ConnectToSSH,
    data: {
      id: serviceId,
    } as ConnectServerWSMessage,
  };
  ws?.send(JSON.stringify(data));
}
