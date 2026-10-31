// services.svelte.ts
import {
  type ConnectServerWSMessage,
  type ServiceWSMessage,
  ConnectToSSH,
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

export function removeServiceFunction(serviceId: number) {
  servicesState.services = servicesState.services.filter(
    (service) => service.id !== serviceId,
  );
}

export function checkWs(): boolean {
  if (!ws || ws.readyState !== WebSocket.OPEN) {
    servicesState.error = "WebSocket is not connected";
    toast.error(
      `WebSocket is not connected, click red floating circle to reconnect.`,
    );
    return false;
  }
  return true;
}

export function startServerConnection(serviceId: number, name: string) {
  if (!checkWs()) return;
  let data: ServiceWSMessage = {
    action: ConnectToSSH,
    data: {
      id: serviceId,
    } as ConnectServerWSMessage,
  };
  let findService = servicesState.services.find((s) => s.id === serviceId);
  if (findService) {
    toast.error(`Service with id ${serviceId} already connected or pending`);
    return;
  }
  servicesState.services.push({ id: serviceId, status: "pending", name: name });
  ws?.send(JSON.stringify(data));
}
