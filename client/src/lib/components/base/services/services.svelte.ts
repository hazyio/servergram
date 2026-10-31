import {
  WSServiceMessageConnectServer,
  WSServiceServerServiceType,
  type WSServiceConnectedListItem,
  type WSServiceConnectServerData,
  type WSServiceMessage,
  type WSServicePassPhraseRequiredData,
  type WSServiceTrustHostFingerPrintData
} from "#lib/types/websockets/index.ts";
import { getWebSocketUrl } from "#lib/utils.ts";
import { toast } from "svelte-sonner";
import handleWebSocketMessage from "./handleWebsocket";
let ws: WebSocket | undefined;

export const servicesState = $state({
  services: [] as WSServiceConnectedListItem[],
  connecting: true,
  error: null as string | null,
  pendingHostTrustCollection: [] as WSServiceTrustHostFingerPrintData[],
  pendingPassPhraseCollection: [] as WSServicePassPhraseRequiredData[],
});

export function startWebSocket(onOpen: () => void) {
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
    handleWebSocketMessage(JSON.parse(e.data));
  };
}
export function stopWebSocket() {
  ws?.close(1000, "bye");
  ws = undefined;
}
export type ServiceState = "connected";
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

export class ConnectedServices {
  static findServer(serviceId: number) {
    return servicesState.services.find(
      (service) =>
        service.id === serviceId && service.type === WSServiceServerServiceType,
    );
  }

  static startServerConnection(serviceId: number, name: string) {
    if (!checkWs()) return;
    let data: WSServiceMessage = {
      action: WSServiceMessageConnectServer,
      data: {
        id: serviceId,
      } as WSServiceConnectServerData,
    };
    let findService = ConnectedServices.findServer(serviceId);
    if (findService) {
      toast.error(`Service with id ${serviceId} already connected or pending`);
      return;
    }
    ws?.send(JSON.stringify(data));
  }
}
