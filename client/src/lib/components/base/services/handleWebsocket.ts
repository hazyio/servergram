import {
  type WSError,
  type WSServiceConnectedList,
  type WSServiceHostFingerPrintMismatchData,
  type WSServiceMessage,
  type WSServicePassPhraseRequiredData,
  type WSServiceTrustHostFingerPrintData,
} from "#lib/types/websockets/index.ts";
import { castType } from "#lib/utils.ts";
import { toast } from "svelte-sonner";
import { ConnectedServices, servicesState } from "./services.svelte.ts";
export default function handleWebSocketMessage(message: WSServiceMessage) {
  switch (message.action) {
    case "pass_phrase_required":
      castType<WSServicePassPhraseRequiredData>(message.data, (value) => {
        let findIndex = ConnectedServices.findByIdInData(
          servicesState.pendingPassPhraseCollection,
          value.id,
        );
        if (!findIndex) {
          servicesState.pendingPassPhraseCollection.push(value);
        }
      });

      break;
    case "trust_host_finger_print":
      castType<WSServiceTrustHostFingerPrintData>(message.data, (value) => {
        let findIndex = ConnectedServices.findByIdInData(
          servicesState.pendingHostTrustCollection,
          value.id,
        );
        if (!findIndex) {
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
        servicesState.gettingConnected = false;
      });
      break;
    case "host_finger_print_mismatch":
      castType<WSServiceHostFingerPrintMismatchData>(message.data, (value) => {
        let findIndex = ConnectedServices.findByIdInData(
          servicesState.pendingHostTrustMismatchCollection,
          value.id,
        );
        if (!findIndex) {
          servicesState.pendingHostTrustMismatchCollection.push(value);
        }
      });
      break;
    default:
      console.warn("Unhandled WebSocket message:", message);
      break;
  }
}
