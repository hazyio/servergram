import {
  type WSServiceConnectedListItem,
  type WSServiceHostFingerPrintMismatchData,
  type WSServiceMessage,
  type WSServicePassPhraseRequiredData,
  type WSServiceTrustHostFingerPrintData,
} from "#lib/types/websockets/index.ts";
import { getWebSocketUrl } from "#lib/utils.ts";
import { toast } from "svelte-sonner";
import handleWebSocketMessage from "./handleWebsocket";
import { ConnectServer, DisconnectServer } from "./webSocketMessages";

const WS_CONNECT_TIMEOUT_MS = 10_000;

let ws: WebSocket | undefined;
let connectTimer: ReturnType<typeof setTimeout> | undefined;

function clearConnectTimer() {
  if (connectTimer === undefined) return;
  clearTimeout(connectTimer);
  connectTimer = undefined;
}

export const servicesState = $state({
  services: [] as WSServiceConnectedListItem[],
  connecting: true,
  error: null as string | null,
  pendingHostTrustCollection: [] as WSServiceTrustHostFingerPrintData[],
  pendingPassPhraseCollection: [] as WSServicePassPhraseRequiredData[],
  pendingHostTrustMismatchCollection:
    [] as WSServiceHostFingerPrintMismatchData[],
  tries: 0,
  gettingConnected: true,
});
export function sendMessage(message: WSServiceMessage) {
  if (!checkWs()) {
    return;
  }
  ws?.send(JSON.stringify(message));
}
export function startWebSocket(onOpen: () => void) {
  if (ws?.readyState === WebSocket.OPEN) {
    console.log("WebSocket already open");
    return;
  }

  clearConnectTimer();
  ws?.close(1000, "reconnecting");
  servicesState.connecting = true;

  const socket = new WebSocket(getWebSocketUrl("/ws/services"));
  ws = socket;

  connectTimer = setTimeout(() => {
    connectTimer = undefined;
    if (socket.readyState === WebSocket.OPEN) return;

    socket.close(4000, "Connection timed out");
    if (ws !== socket) return;

    ws = undefined;
    servicesState.error = "Could not connect, connection timed out";
    servicesState.connecting = false;
    servicesState.tries += 1;
  }, WS_CONNECT_TIMEOUT_MS);

  socket.onopen = () => {
    if (ws !== socket) return;
    clearConnectTimer();
    onOpen();
    servicesState.connecting = false;
    servicesState.error = null;
    servicesState.tries = 0;
  };
  socket.onerror = (e) => {
    if (ws !== socket) return;
    clearConnectTimer();
    console.log(e);
    servicesState.error = "Could not connect, webSocket error";
    servicesState.connecting = false;
    servicesState.tries += 1;
  };
  socket.onclose = () => {
    if (ws !== socket) return;
    clearConnectTimer();
    ws = undefined;
    servicesState.connecting = false;
    servicesState.error ??= "WebSocket connection closed";
  };
  socket.onmessage = (e) => {
    handleWebSocketMessage(JSON.parse(e.data));
  };
}
export function stopWebSocket() {
  clearConnectTimer();
  const socket = ws;
  ws = undefined;
  socket?.close(1000, "bye");
}
export type ServiceState = "connected";
export function checkWs(showError: boolean = false): boolean {
  if (!ws || ws.readyState !== WebSocket.OPEN) {
    servicesState.error = "WebSocket is not connected";
    if (showError) {
      toastError(`WebSocket is not connected, click red button to reconnect.`);
    }
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
  /**
   * Find a server by its service ID.
   */
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

    let findService = ConnectedServices.findServer(serviceId);
    if (findService) {
      toastError(`Service with id ${serviceId} already connected or pending`);
      return;
    }
    ConnectServer({ id: serviceId, name: serviceName });
  }
  static disconnectServer(serviceId: number) {
    if (!checkWs()) {
      return;
    }
    let findService = ConnectedServices.findServer(serviceId);
    if (!findService) {
      toastError(`Server with id ${serviceId} is not connected`);
      return;
    }
    DisconnectServer({ id: serviceId });
  }
  static disconnectDatabase(serviceId: number) {
    if (!checkWs()) {
      return;
    }
    alert("Disconnect database not implemented yet");
  }
}
