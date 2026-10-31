import {
  WSServiceMessageError as ServerWSError,
  WSServiceMessageConnectedServices,
  WSServiceMessagePassPhraseRequired,
  WSServiceMessageTrustHostFingerPrint,
  type WSServiceConnectedList,
  type WSError,
  type WSServiceMessage,
  type WSServicePassPhraseRequiredData,
  type WSServiceTrustHostFingerPrintData,
} from "#lib/types/websockets/index.ts";
import { castType } from "#lib/utils.ts";
import { toast } from "svelte-sonner";
import { servicesState } from "./services.svelte.ts";
export default function handleWebSocketMessage(message: WSServiceMessage) {
  switch (message.action) {
    case WSServiceMessagePassPhraseRequired:
      castType<WSServicePassPhraseRequiredData>(message.data, (value) => {
        let findIndex = servicesState.pendingPassPhraseCollection.findIndex(
          (item) => item.id === value.id,
        );
        if (findIndex === -1) {
          servicesState.pendingPassPhraseCollection.push(value);
        }
      });

      break;
    case WSServiceMessageTrustHostFingerPrint:
      castType<WSServiceTrustHostFingerPrintData>(message.data, (value) => {
        let findIndexHost = servicesState.pendingHostTrustCollection.findIndex(
          (item) => item.id === value.id,
        );
        if (findIndexHost === -1) {
          servicesState.pendingHostTrustCollection.push(value);
        }
      });

      break;
    case ServerWSError:
      // Handle server error message
      toast.error((message.data as WSError).message);
      break;
    case WSServiceMessageConnectedServices:
      // Handle connected services message
      castType<WSServiceConnectedList>(message.data, (value) => {
        servicesState.services = value.items;
      });
      break;
    default:
      console.warn("Unhandled WebSocket message:", message);
      break;
  }
}
