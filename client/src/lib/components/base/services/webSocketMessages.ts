import type {
  WSServiceConnectServerData,
  WSServiceDisconnectServerData,
} from "#lib/types/websockets/index.ts";
import { sendMessage, servicesState } from "./services.svelte";

export function GetConnectedServices() {
  servicesState.gettingConnected = true;
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
