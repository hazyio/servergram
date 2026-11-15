export type CronJobStatus =
  "running" | "scheduled" | "paused" | "failed" | "unmanaged";

export type CronRetention = {
  metadata_runs: number;
  run_logs: number;
  error_log_days: number;
};

export type CronJob = {
  id: number;
  name: string;
  command: string;
  schedule: string;
  status: CronJobStatus;
  last_run: string | null;
  next_run: string | null;
  retention?: CronRetention;
};

export type CreateCronJobInput = {
  name: string;
  command: string;
  schedule: string;
  wrap_with_agent: boolean;
  retention?: CronRetention;
};

export type UpdateCronJobInput = {
  name: string;
  command: string;
  schedule: string;
  retention?: CronRetention;
};

export type CronAction = "run" | "resume" | "stop";

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
      status: "scheduled",
      last_run: iso(minutes(-1_420)),
      next_run: iso(minutes(20)),
    },
    {
      id: 2,
      name: "Log rotation",
      command: "/usr/sbin/logrotate /etc/logrotate.conf",
      schedule: "*/30 * * * *",
      status: "running",
      last_run: iso(minutes(-10)),
      next_run: iso(minutes(20)),
    },
    {
      id: 3,
      name: "Cache warmup",
      command: "php /var/www/app/artisan cache:warm",
      schedule: "*/5 * * * *",
      status: "failed",
      last_run: iso(minutes(-12)),
      next_run: iso(minutes(3)),
    },
    {
      id: 4,
      name: "Report generator",
      command: "/opt/reports/generate --weekly",
      schedule: "0 8 * * 1",
      status: "paused",
      last_run: iso(minutes(-2_880)),
      next_run: null,
    },
    {
      id: 5,
      name: "Health check",
      command: "curl -fsS http://localhost/health",
      schedule: "*/1 * * * *",
      status: "running",
      last_run: iso(seconds(-30)),
      next_run: iso(seconds(30)),
    },
    {
      id: 6,
      name: "Temp cleanup",
      command: "/usr/bin/find /tmp -mtime +7 -delete",
      schedule: "0 3 * * 0",
      status: "failed",
      last_run: iso(minutes(-7_200)),
      next_run: iso(minutes(2_880)),
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
let nextJobId = 1_000;

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

export function getPrimaryCronAction(status: CronJobStatus): {
  kind: CronAction;
  label: string;
} {
  switch (status) {
    case "running":
      return { kind: "stop", label: "Stop" };
    case "paused":
      return { kind: "resume", label: "Resume" };
    default:
      return { kind: "run", label: "Run" };
  }
}

export async function createCronJob(
  serverId: number,
  input: CreateCronJobInput,
): Promise<CronApiResponse<CronJob>> {
  return simulateRequest(() => {
    const name = input.name.trim();
    const command = input.command.trim();
    const schedule = input.schedule.trim();

    if (!name) throw new Error("Name is required");
    if (!command) throw new Error("Command is required");
    if (!schedule) throw new Error("Run schedule is required");

    const job: CronJob = {
      id: nextJobId++,
      name,
      command,
      schedule,
      status: input.wrap_with_agent ? "scheduled" : "unmanaged",
      last_run: null,
      next_run: input.wrap_with_agent ? iso(minutes(5)) : null,
    };

    if (input.wrap_with_agent && input.retention) {
      job.retention = input.retention;
    }

    jobsFor(serverId).push(job);
    return job;
  });
}

export async function updateCronJob(
  serverId: number,
  jobId: number,
  input: UpdateCronJobInput,
): Promise<CronApiResponse<CronJob>> {
  return simulateRequest(() => {
    const job = jobsFor(serverId).find((entry) => entry.id === jobId);
    if (!job) throw new Error("Cron job no longer exists");

    const name = input.name.trim();
    const command = input.command.trim();
    const schedule = input.schedule.trim();

    if (!name) throw new Error("Name is required");
    if (!command) throw new Error("Command is required");
    if (!schedule) throw new Error("Run schedule is required");

    job.name = name;
    job.command = command;
    job.schedule = schedule;

    if (input.retention && job.status !== "unmanaged") {
      job.retention = input.retention;
    }

    return job;
  });
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
    const live = index === 0 && job.status === "running";
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

export async function performCronAction(
  serverId: number,
  jobId: number,
  action: CronAction,
): Promise<CronApiResponse<CronJob>> {
  return simulateRequest(() => {
    const job = jobsFor(serverId).find((entry) => entry.id === jobId);
    if (!job) throw new Error("Cron job no longer exists");
    if (job.status === "unmanaged") {
      throw new Error("Wrap this cron job with the agent before managing it");
    }

    switch (action) {
      case "run":
        job.status = "running";
        job.last_run = iso(0);
        job.next_run = iso(minutes(1));
        break;
      case "resume":
        job.status = "scheduled";
        job.next_run = iso(minutes(5));
        break;
      case "stop":
        job.status = "paused";
        job.next_run = null;
        break;
    }

    return job;
  });
}

export async function duplicateCronJob(
  serverId: number,
  jobId: number,
): Promise<CronApiResponse<CronJob>> {
  return simulateRequest(() => {
    const jobs = jobsFor(serverId);
    const source = jobs.find((entry) => entry.id === jobId);
    if (!source) throw new Error("Cron job no longer exists");

    const copy: CronJob = {
      ...source,
      id: nextJobId++,
      name: source.name.trim() ? `${source.name} copy` : "",
      status: source.status === "running" ? "scheduled" : source.status,
      last_run: null,
      next_run: source.status === "paused" ? null : iso(minutes(5)),
    };
    jobs.push(copy);
    return copy;
  });
}

export async function wrapCronJobsWithAgent(
  serverId: number,
  jobIds: number[],
): Promise<CronApiResponse<{ items: CronJob[] }>> {
  return simulateRequest(() => {
    const jobs = jobsFor(serverId);
    const wrapped = jobs.filter(
      (job) => jobIds.includes(job.id) && job.status === "unmanaged",
    );

    if (wrapped.length === 0) {
      throw new Error("No unmanaged cron jobs were found");
    }

    for (const job of wrapped) {
      job.status = "scheduled";
      job.next_run = iso(minutes(5));
    }

    return { items: wrapped };
  });
}

export async function deleteCronJob(
  serverId: number,
  jobId: number,
): Promise<CronApiResponse<{ id: number }>> {
  return simulateRequest(() => {
    const jobs = jobsFor(serverId);
    const index = jobs.findIndex((entry) => entry.id === jobId);
    if (index === -1) throw new Error("Cron job no longer exists");

    jobs.splice(index, 1);
    return { id: jobId };
  });
}
