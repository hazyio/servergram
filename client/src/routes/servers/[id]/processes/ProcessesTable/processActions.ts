import { NewRequestClient } from "#lib/api/request.ts";

export type ProcessActionKind = "stop" | "force-stop" | "signal";

/**
 * `sigNum` is the numeric value of the
 * `SignalProcessRequest.SignalProcessRequestSignals` protobuf enum used by the
 * API. It is not the Unix signal number (for example, SIGTERM is 4 here and 15
 * on Linux).
 */
export type ProcessAction = {
  label: string;
  signal: string;
  sigNum: number;
  kind: ProcessActionKind;
  variant?: "default" | "destructive";
};

export const stopProcessAction: ProcessAction = {
  label: "Stop",
  signal: "SIGTERM",
  sigNum: 4,
  kind: "stop",
};

export const forceStopProcessAction: ProcessAction = {
  label: "Force stop",
  signal: "SIGKILL",
  sigNum: 3,
  kind: "force-stop",
  variant: "destructive",
};

export const processSignalActions: ProcessAction[] = [
  { label: "Terminate", signal: "SIGTERM", sigNum: 4, kind: "signal" },
  {
    label: "Kill",
    signal: "SIGKILL",
    sigNum: 3,
    kind: "signal",
    variant: "destructive",
  },
  { label: "Interrupt", signal: "SIGINT", sigNum: 2, kind: "signal" },
  { label: "Hang up", signal: "SIGHUP", sigNum: 1, kind: "signal" },
  { label: "Quit", signal: "SIGQUIT", sigNum: 5, kind: "signal" },
  { label: "Pause", signal: "SIGSTOP", sigNum: 8, kind: "signal" },
  { label: "Continue", signal: "SIGCONT", sigNum: 9, kind: "signal" },
  { label: "User 1", signal: "SIGUSR1", sigNum: 6, kind: "signal" },
  { label: "User 2", signal: "SIGUSR2", sigNum: 7, kind: "signal" },
  {
    label: "Abort",
    signal: "SIGABRT",
    sigNum: 12,
    kind: "signal",
    variant: "destructive",
  },
];

export type ProcessActionError = {
  type: string;
  status?: number;
  title?: string;
  detail?: string;
};

export type ProcessActionRequest = {
  serverId: number;
  pid: number;
  action: ProcessAction;
};

export type ProcessActionResult = {
  error?: ProcessActionError;
};

export async function sendProcessSignal(
  request: ProcessActionRequest,
): Promise<ProcessActionResult> {
  const client = NewRequestClient();
  const params = { path: { id: request.serverId } };

  switch (request.action.kind) {
    case "stop": {
      const { error } = await client.POST("/api/server/{id}/stop-process/", {
        params,
        body: { pid: request.pid },
      });
      return { error };
    }
    case "force-stop": {
      const { error } = await client.POST(
        "/api/server/{id}/force-stop-process/",
        {
          params,
          body: { pid: request.pid },
        },
      );
      return { error };
    }
    case "signal": {
      const { error } = await client.POST("/api/server/{id}/signal-process/", {
        params,
        body: { pid: request.pid, sig_num: request.action.sigNum },
      });
      return { error };
    }
  }
}
