import { NewRequestClient } from "#lib/api/request.ts";
import type { components } from "#lib/types/request-schema.d.ts";

export type CronJobStatus = "active" | "inactive" | "unmanaged";

export type CronRetention = {
  metadata_runs: number;
  run_logs: number;
  error_log_days: number;
};

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
  retention?: CronRetention;
};

export type CronJobInfo = {
  id: number;
  name: string;
  command: string;
  schedule: string;
  managed: boolean;
  retention: CronRetention;
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
  meta_keep_count: number;
  run_logs_keep_count: number;
  error_logs_keep_days: number;
};

export type UpdateCronJobInput = {
  name: string;
  command: string;
  schedule: string;
  wrap: boolean;
  retention: CronRetention;
};

export type CronRunStatus = "running" | "succeeded" | "failed" | "canceled";

export type CronRun = {
  id: number;
  job_id: number;
  status: CronRunStatus;
  trigger: "schedule" | "manual";
  started_at: string;
  finished_at: string | null;
  exit_code: number | null;
};

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
        meta_keep_count: input.retention.metadata_runs,
        run_logs_keep_count: input.retention.run_logs,
        error_logs_keep_days: input.retention.error_log_days,
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

  return {
    data: {
      id: info.id ?? jobId,
      name: info.name ?? "",
      command: info.command ?? "",
      schedule: info.schedule ?? "",
      managed: info.managed ?? false,
      retention: {
        metadata_runs: info.meta_keep_count ?? 0,
        run_logs: info.run_logs_keep_count ?? 0,
        error_log_days: info.error_logs_keep_days ?? 0,
      },
    },
  };
}

const runsByJob = new Map<number, CronRun[]>();

function seedRuns(job: CronJob): CronRun[] {
  const statuses: CronRunStatus[] = [
    "succeeded",
    "succeeded",
    "failed",
    "succeeded",
    "canceled",
    "succeeded",
    "succeeded",
    "failed",
  ];

  return Array.from({ length: 12 }, (_, index) => {
    const live = index === 0 && job.status === "active";
    const status: CronRunStatus = live
      ? "running"
      : statuses[(index + job.id) % statuses.length];
    const startedAt = iso(minutes(-(index * 45 + 5)));
    const durationMs = seconds(45 + ((job.id * 17 + index * 29) % 180));

    return {
      id: job.id * 1_000 + (12 - index),
      job_id: job.id,
      status,
      trigger: (index + job.id) % 5 === 0 ? "manual" : "schedule",
      started_at: startedAt,
      finished_at: live
        ? null
        : new Date(Date.parse(startedAt) + durationMs).toISOString(),
      exit_code: live
        ? null
        : status === "succeeded"
          ? 0
          : status === "canceled"
            ? 130
            : 1,
    };
  });
}

function runsFor(job: CronJob): CronRun[] {
  let runs = runsByJob.get(job.id);
  if (!runs) {
    runs = seedRuns(job);
    runsByJob.set(job.id, runs);
  }
  return runs;
}

export async function listCronRuns(
  serverId: number,
  jobId: number,
): Promise<CronApiResponse<{ items: CronRun[] }>> {
  return simulateRequest(() => {
    const job = jobsFor(serverId).find((entry) => entry.id === jobId);
    if (!job) throw new Error("Cron job no longer exists");

    if (job.status === "unmanaged") return { items: [] };

    return { items: runsFor(job) };
  });
}

function buildRunLog(job: CronJob, run: CronRun): string {
  const started = Date.parse(run.started_at);
  const line = (offsetMs: number, level: string, message: string) =>
    `${new Date(started + offsetMs).toISOString()} ${level.padEnd(5)} ${message}`;

  const lines = [
    line(0, "INFO", `Starting "${cronJobDisplayName(job)}" (run #${run.id})`),
    line(
      20,
      "INFO",
      `Trigger: ${run.trigger === "manual" ? "manual" : "schedule"}`,
    ),
    line(80, "INFO", `Command: ${job.command}`),
    line(160, "INFO", "Loading configuration from /etc/servergram/agent.yaml"),
    line(420, "INFO", "Connected to local agent socket"),
  ];

  if (run.status === "running") {
    lines.push(
      line(900, "INFO", `Executing ${job.command}`),
      line(1_200, "INFO", "Streaming output to the agent"),
      line(1_500, "INFO", "Run is still in progress..."),
    );
    return lines.join("\n");
  }

  if (run.status === "succeeded") {
    lines.push(
      line(900, "INFO", `Executing ${job.command}`),
      line(2_400, "INFO", "Task finished successfully"),
      line(2_600, "INFO", "Cleaning up temporary files"),
      line(2_700, "INFO", `Exited with code ${run.exit_code}`),
    );
  } else if (run.status === "failed") {
    lines.push(
      line(900, "INFO", `Executing ${job.command}`),
      line(1_900, "ERROR", "Command failed"),
      line(1_950, "ERROR", "stderr: connection to upstream service timed out"),
      line(1_980, "ERROR", "Retrying once in 5s..."),
      line(6_900, "ERROR", "Retry failed, aborting"),
      line(7_000, "ERROR", `Exited with code ${run.exit_code}`),
    );
  } else {
    lines.push(
      line(900, "INFO", `Executing ${job.command}`),
      line(1_400, "WARN", "Run was canceled before completion"),
      line(1_450, "INFO", `Exited with code ${run.exit_code}`),
    );
  }

  return lines.join("\n");
}

export async function getCronRunLog(
  serverId: number,
  jobId: number,
  runId: number,
): Promise<CronApiResponse<{ content: string }>> {
  return simulateRequest(() => {
    const job = jobsFor(serverId).find((entry) => entry.id === jobId);
    if (!job) throw new Error("Cron job no longer exists");
    if (job.status === "unmanaged") {
      throw new Error("This cron job is not managed by the agent");
    }

    const run = runsFor(job).find((entry) => entry.id === runId);
    if (!run) throw new Error("Run no longer exists");

    return { content: buildRunLog(job, run) };
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

export async function duplicateCronJob(
  serverId: number,
  jobId: number,
): Promise<CronApiResponse<Record<string, never>>> {
  const { data, error } = await NewRequestClient().PUT(
    "/api/cron/{id}/duplicate/{job_id}",
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
