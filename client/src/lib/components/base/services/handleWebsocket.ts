import {
  type WSError,
  type WSServiceConnectedList,
  type WSServiceMessage,
  type WSServicePassPhraseRequiredData,
  type WSServiceTrustHostFingerPrintData,
} from "#lib/types/websockets/index.ts";
import { castType } from "#lib/utils.ts";
import { toast } from "svelte-sonner";
import { servicesState } from "./services.svelte.ts";
export default function handleWebSocketMessage(message: WSServiceMessage) {
  switch (message.action) {
    case "pass_phrase_required":
      castType<WSServicePassPhraseRequiredData>(message.data, (value) => {
        let findIndex = servicesState.pendingPassPhraseCollection.findIndex(
          (item) => item.id === value.id,
        );
        if (findIndex === -1) {
          servicesState.pendingPassPhraseCollection.push(value);
        }
      });

      break;
    case "trust_host_finger_print":
      castType<WSServiceTrustHostFingerPrintData>(message.data, (value) => {
        let findIndexHost = servicesState.pendingHostTrustCollection.findIndex(
          (item) => item.id === value.id,
        );
        if (findIndexHost === -1) {
          servicesState.pendingHostTrustCollection.push(value);
        }
      });

      break;
    case "error":
      // Handle server error message
      toast.error((message.data as WSError).message);
      break;
    case "connected_services":
      // Handle connected services message
      castType<WSServiceConnectedList>(message.data, (value) => {
        servicesState.services = value.items;
      });
      break;
    case "host_finger_print_mismatch":
      // Handle host finger print mismatch message
      toast.error("Host finger print mismatch");
      break;
    default:
      console.warn("Unhandled WebSocket message:", message);
      break;
  }
}
