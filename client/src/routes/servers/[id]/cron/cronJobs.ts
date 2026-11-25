import { NewRequestClient } from "#lib/api/request.ts";
import type { components } from "#lib/types/request-schema.d.ts";

export type CronJobStatus = "active" | "inactive" | "unmanaged";

export type CronJob = {
  id: number;
  key?: string;
  hash?: string;
  duplicate?: boolean;
  name: string;
  command: string;
  schedule: string;
  status: CronJobStatus;
  last_run: string | null;
  next_run: string | null;
  run_logs_keep_count?: number;
};

export type CronJobInfo = {
  id: number;
  name: string;
  command: string;
  schedule: string;
  managed: boolean;
  status: CronJobStatus;
  last_run: string | null;
  next_run: string | null;
  run_logs_keep_count: number;
};

export type CronDuplicateGroup = {
  hash: string;
  job: CronJob;
  count: number;
};

export type CreateCronJobInput = {
  name: string;
  command: string;
  schedule: string;
  wrap: boolean;
  run_logs_keep_count: number;
};

export type UpdateCronJobInput = {
  name: string;
  command: string;
  schedule: string;
  wrap: boolean;
  run_logs_keep_count: number;
};

export type CronRunStatus =
  "running" | "succeeded" | "failed" | "canceled" | "timed_out";

export type CronRun = {
  id: number;
  job_id: number;
  status: CronRunStatus;
  trigger: "schedule" | "manual";
  started_at: string;
  finished_at: string | null;
  exit_code: number | null;
};

export type CronRunLog = {
  content: string;
  truncated: boolean;
};

type CronRunApi = components["schemas"]["CronRunMeta"];

const cronRunStatuses: Record<string, CronRunStatus> = {
  running: "running",
  success: "succeeded",
  succeeded: "succeeded",
  failed: "failed",
  cancelled: "canceled",
  canceled: "canceled",
  timed_out: "timed_out",
};

function cronRunFromApi(jobId: number, run: CronRunApi): CronRun {
  const status = cronRunStatuses[run.status ?? ""] ?? "failed";

  return {
    id: run.id ?? 0,
    job_id: jobId,
    status,
    trigger: run.trigger === "manual" ? "manual" : "schedule",
    started_at: run.started_at ?? "",
    finished_at: isCronDateSet(run.finished_at) ? run.finished_at : null,
    // The API omits exit_code when it is 0, so only running runs have no value.
    exit_code: status === "running" ? null : (run.exit_code ?? 0),
  };
}

export type CronApiError = {
  type: string;
  status?: number;
  title?: string;
  detail?: string;
};

export type CronApiResponse<T> = {
  data?: T;
  error?: CronApiError;
};

const knownCronStatuses = new Set<CronJobStatus>(["active", "inactive"]);

export function isCronDateSet(
  value: string | null | undefined,
): value is string {
  if (!value) return false;

  const date = new Date(value);

  return !Number.isNaN(date.getTime()) && date.getUTCFullYear() > 1;
}

export function cronJobFromApi(cron: components["schemas"]["Cron"]): CronJob {
  const status: CronJobStatus = !cron.managed
    ? "unmanaged"
    : knownCronStatuses.has(cron.status as CronJobStatus)
      ? (cron.status as CronJobStatus)
      : "inactive";

  return {
    id: cron.id ?? 0,
    hash: cron.hash ?? "",
    name: cron.name ?? "",
    command: cron.command ?? "",
    schedule: cron.schedule ?? "",
    status,
    last_run: isCronDateSet(cron.last_run_at) ? cron.last_run_at : null,
    next_run: isCronDateSet(cron.next_run_at) ? cron.next_run_at : null,
  };
}

const requestDelayMs = 550;

const minutes = (value: number) => value * 60_000;
const seconds = (value: number) => value * 1_000;
const iso = (offsetMs: number) => new Date(Date.now() + offsetMs).toISOString();

function seedJobs(): CronJob[] {
  return [
    {
      id: 1,
      name: "Database backup",
      command: "/usr/local/bin/backup-db.sh",
      schedule: "0 2 * * *",
      status: "active",
      last_run: iso(minutes(-1_420)),
      next_run: iso(minutes(20)),
    },
    {
      id: 2,
      name: "Log rotation",
      command: "/usr/sbin/logrotate /etc/logrotate.conf",
      schedule: "*/30 * * * *",
      status: "active",
      last_run: iso(minutes(-10)),
      next_run: iso(minutes(20)),
    },
    {
      id: 3,
      name: "Cache warmup",
      command: "php /var/www/app/artisan cache:warm",
      schedule: "*/5 * * * *",
      status: "inactive",
      last_run: iso(minutes(-12)),
      next_run: null,
    },
    {
      id: 4,
      name: "Report generator",
      command: "/opt/reports/generate --weekly",
      schedule: "0 8 * * 1",
      status: "inactive",
      last_run: iso(minutes(-2_880)),
      next_run: null,
    },
    {
      id: 5,
      name: "Health check",
      command: "curl -fsS http://localhost/health",
      schedule: "*/1 * * * *",
      status: "active",
      last_run: iso(seconds(-30)),
      next_run: iso(seconds(30)),
    },
    {
      id: 6,
      name: "Temp cleanup",
      command: "/usr/bin/find /tmp -mtime +7 -delete",
      schedule: "0 3 * * 0",
      status: "inactive",
      last_run: iso(minutes(-7_200)),
      next_run: null,
    },
    {
      id: 7,
      name: "",
      command: "/opt/scripts/nightly-sync.sh",
      schedule: "15 1 * * *",
      status: "unmanaged",
      last_run: iso(minutes(-600)),
      next_run: iso(minutes(45)),
    },
    {
      id: 8,
      name: "",
      command: "/usr/bin/rsync -a /var/www /backup/www",
      schedule: "0 */6 * * *",
      status: "unmanaged",
      last_run: null,
      next_run: null,
    },
  ];
}

const jobsByServer = new Map<number, CronJob[]>();

function jobsFor(serverId: number): CronJob[] {
  let jobs = jobsByServer.get(serverId);
  if (!jobs) {
    jobs = seedJobs();
    jobsByServer.set(serverId, jobs);
  }
  return jobs;
}

function clone<T>(value: T): T {
  return JSON.parse(JSON.stringify(value)) as T;
}

/**
 * Temporary request stub. It deliberately returns the same `{ data, error }`
 * shape as `NewRequestClient()` so these calls can be swapped for the real
 * endpoints once the cron API exists.
 */
async function simulateRequest<T>(run: () => T): Promise<CronApiResponse<T>> {
  await new Promise((resolve) => setTimeout(resolve, requestDelayMs));

  try {
    return { data: clone(run()) };
  } catch (error) {
    return {
      error: {
        type: "about:blank",
        status: 500,
        title: "Request failed",
        detail:
          error instanceof Error
            ? error.message
            : "An unexpected error occurred",
      },
    };
  }
}

export function cronJobDisplayName(job: Pick<CronJob, "name">): string {
  return job.name.trim() || "Unnamed job";
}

export async function createCronJob(
  serverId: number,
  input: CreateCronJobInput,
): Promise<CronApiResponse<CreateCronJobInput>> {
  const { error } = await NewRequestClient().POST("/api/cron/{id}/create", {
    params: { path: { id: serverId } },
    body: input,
  });

  if (error) {
    return { error };
  }

  return { data: input };
}

export async function removeDuplicateCronJobs(
  serverId: number,
): Promise<CronApiResponse<{ removed: number }>> {
  const { data, error } = await NewRequestClient().POST(
    "/api/cron/{id}/remove-duplicates",
    { params: { path: { id: serverId } } },
  );

  if (error) {
    return { error };
  }

  return { data: data ?? { removed: 0 } };
}

export async function updateCronJob(
  serverId: number,
  jobId: number,
  input: UpdateCronJobInput,
): Promise<CronApiResponse<Record<string, never>>> {
  const { error } = await NewRequestClient().POST(
    "/api/cron/{id}/update/{job_id}",
    {
      params: {
        path: {
          id: serverId,
          job_id: jobId,
        },
      },
      body: {
        name: input.name,
        command: input.command,
        schedule: input.schedule,
        wrap: input.wrap,
        run_logs_keep_count: input.run_logs_keep_count,
      },
    },
  );

  if (error) {
    return { error };
  }

  return { data: {} };
}

export async function listCronJobs(
  serverId: number,
): Promise<CronApiResponse<{ items: CronJob[] }>> {
  // Replace with:
  // return NewRequestClient().GET("/api/server/{id}/cron-jobs", {
  //   params: { path: { id: serverId } },
  // });
  return simulateRequest(() => ({ items: jobsFor(serverId) }));
}

export async function getCronJob(
  serverId: number,
  jobId: number,
): Promise<CronApiResponse<CronJob>> {
  return simulateRequest(() => {
    const job = jobsFor(serverId).find((entry) => entry.id === jobId);
    if (!job) throw new Error("Cron job no longer exists");
    return job;
  });
}

type CronInfoApi = components["schemas"]["CronInfoResponse"] & {
  name?: string;
  status?: string;
  last_run_at?: string;
  next_run_at?: string;
};

export async function getCronJobInfo(
  serverId: number,
  jobId: number,
): Promise<CronApiResponse<CronJobInfo>> {
  const { data, error } = await NewRequestClient().GET(
    "/api/cron/{id}/info/{job_id}",
    {
      params: {
        path: {
          id: serverId,
          job_id: jobId,
        },
      },
    },
  );

  if (error) {
    return { error };
  }

  if (!data) {
    return {
      error: {
        type: "about:blank",
        status: 404,
        title: "Cron job not found",
        detail: "The cron job no longer exists",
      },
    };
  }

  const info = data as CronInfoApi;
  const managed = info.managed ?? false;

  // TODO: remove this stub fallback once the info endpoint exposes run state.
  // The info response currently only carries configuration, so status and run
  // times come from the same seeded data used by the other cron stubs.
  const fallback = jobsFor(serverId).find((entry) => entry.id === jobId);
  const status: CronJobStatus = !managed
    ? "unmanaged"
    : knownCronStatuses.has(info.status as CronJobStatus)
      ? (info.status as CronJobStatus)
      : fallback?.status === "active"
        ? "active"
        : "inactive";
  const lastRun = isCronDateSet(info.last_run_at)
    ? info.last_run_at
    : (fallback?.last_run ?? null);
  const nextRun = isCronDateSet(info.next_run_at)
    ? info.next_run_at
    : (fallback?.next_run ?? null);

  return {
    data: {
      id: info.id ?? jobId,
      name: info.name ?? "",
      command: info.command ?? "",
      schedule: info.schedule ?? "",
      managed,
      status,
      last_run: lastRun,
      next_run: nextRun,
      run_logs_keep_count: info.run_logs_keep_count ?? 0,
    },
  };
}

export async function listCronRuns(
  serverId: number,
  jobId: number,
): Promise<CronApiResponse<{ items: CronRun[]; more: boolean }>> {
  const { data, error } = await NewRequestClient().GET(
    "/api/cron/{id}/runs/{job_id}",
    {
      params: {
        path: {
          id: serverId,
          job_id: jobId,
        },
        query: {
          limit: 100,
          offset: 0,
        },
      },
    },
  );

  if (error) {
    return { error };
  }

  return {
    data: {
      items: (data?.runs ?? []).map((run) => cronRunFromApi(jobId, run)),
      more: data?.more ?? false,
    },
  };
}
export async function getCronRunLog(
  serverId: number,
  jobId: number,
  runId: number,
): Promise<CronApiResponse<CronRunLog>> {
  // TODO: replace with a real request once the agent exposes run output.
  return simulateRequest(() => {
    const job = jobsFor(serverId).find((entry) => entry.id === jobId);

    return {
      content: [
        `$ ${job?.command ?? "command"}`,
        `[stub] Run #${runId} output will appear here once run log capture is wired up.`,
      ].join("\n"),
      truncated: false,
    };
  });
}

export async function runCronJob(
  serverId: number,
  jobId: number,
): Promise<CronApiResponse<Record<string, never>>> {
  const { data, error } = await NewRequestClient().GET(
    "/api/cron/{id}/run/{job_id}",
    {
      params: {
        path: {
          id: serverId,
          job_id: jobId,
        },
      },
    },
  );

  if (error) {
    return { error };
  }

  return { data: (data ?? {}) as Record<string, never> };
}

export async function wrapCronJobsWithAgent(
  serverId: number,
): Promise<CronApiResponse<{ fixed: number }>> {
  const { data, error } = await NewRequestClient().POST(
    "/api/cron/{id}/fix-unmanaged",
    { params: { path: { id: serverId } } },
  );

  if (error) {
    return { error };
  }

  return { data: data ?? { fixed: 0 } };
}

export async function deleteCronJob(
  serverId: number,
  jobId: number,
): Promise<CronApiResponse<{ id: number }>> {
  const { error } = await NewRequestClient().DELETE(
    "/api/cron/{id}/delete/{job_id}",
    {
      params: {
        path: {
          id: serverId,
          job_id: jobId,
        },
      },
    },
  );

  if (error) {
    return { error };
  }

  return { data: { id: jobId } };
}
