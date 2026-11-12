export type ProcessAction = {
  label: string;
  signal: string;
  variant?: "default" | "destructive";
};

export const stopProcessAction: ProcessAction = {
  label: "Stop",
  signal: "SIGTERM",
};

export const forceStopProcessAction: ProcessAction = {
  label: "Force stop",
  signal: "SIGKILL",
  variant: "destructive",
};

export const processSignalActions: ProcessAction[] = [
  { label: "Terminate", signal: "SIGTERM" },
  { label: "Kill", signal: "SIGKILL", variant: "destructive" },
  { label: "Interrupt", signal: "SIGINT" },
  { label: "Hang up", signal: "SIGHUP" },
  { label: "Quit", signal: "SIGQUIT" },
  { label: "Pause", signal: "SIGSTOP" },
  { label: "Continue", signal: "SIGCONT" },
  { label: "User 1", signal: "SIGUSR1" },
  { label: "User 2", signal: "SIGUSR2" },
  { label: "Abort", signal: "SIGABRT", variant: "destructive" },
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
  signal: string;
};

export type ProcessActionResult = {
  error?: ProcessActionError;
};

/**
 * Stub for the process-signal API.
 *
 * Replace this with the real request once the backend exposes a
 * process-signal endpoint, for example:
 *
 *   return await NewRequestClient().POST("/api/server/send-signal/{id}", {
 *     params: { path: { id: serverId } },
 *     body: { pid, signal },
 *   });
 */
export async function sendProcessSignal(
  _request: ProcessActionRequest,
): Promise<ProcessActionResult> {
  await new Promise((resolve) => setTimeout(resolve, 700));
  return {};
}
