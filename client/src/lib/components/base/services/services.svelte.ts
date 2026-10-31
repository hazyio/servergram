import {
  type WSServiceConnectedListItem,
  type WSServiceConnectServerData,
  type WSServiceMessage,
  type WSServicePassPhraseRequiredData,
  type WSServiceTrustHostFingerPrintData,
  type WSServiceHostFingerPrintMismatchData,
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
  pendingHostTrustMismatchCollection:
    [] as WSServiceHostFingerPrintMismatchData[],
  tries: 0,
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
    servicesState.tries = 0;
  };
  ws.onerror = (e) => {
    console.log(e);
    servicesState.error = "Could not connect, webSocket error";
    servicesState.connecting = false;
    servicesState.tries += 1;
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
    toastError(`WebSocket is not connected, click red button to reconnect.`);
    return false;
  }
  return true;
}
export function toastError(message: string) {
  toast.error(message, {
    duration: 15000,
  });
}

export class ConnectedServices {
  static findServer(serviceId: number) {
    return servicesState.services.find(
      (service) => service.id === serviceId && service.type === "server",
    );
  }
  static findByIdInData<T extends { id: number }>(collection: T[], id: number) {
    return collection.find((item) => item.id === id);
  }

  static startServerConnection(serviceId: number, serviceName: string) {
    if (!checkWs()) {
      return;
    }
    let data: WSServiceMessage = {
      action: "connect_server",
      data: {
        id: serviceId,
        name: serviceName,
      } as WSServiceConnectServerData,
    };
    let findService = ConnectedServices.findServer(serviceId);
    if (findService) {
      toastError(`Service with id ${serviceId} already connected or pending`);
      return;
    }
    ws?.send(JSON.stringify(data));
  }
}
