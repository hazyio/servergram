import type {
  WSServiceConnectServerData,
  WSServiceDisconnectServerData,
} from "#lib/types/websockets/index.ts";
import { sendMessage } from "./services.svelte";

export function GetConnectedServices() {
  sendMessage({
    action: "get_connected_services",
    data: {},
  });
}

export function ConnectServer(data: WSServiceConnectServerData) {
  sendMessage({
    action: "connect_server",
    data: {
      id: data.id,
    },
  });
}

export function DisconnectServer(data: WSServiceDisconnectServerData) {
  sendMessage({
    action: "disconnect_server",
    data: {
      id: data.id,
    },
  });
}
