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
  schedule: string;
  wrap_with_agent: boolean;
  retention?: CronRetention;
};

export type CronAction = "run" | "resume" | "stop";

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
    const schedule = input.schedule.trim();

    if (!name) throw new Error("Name is required");
    if (!schedule) throw new Error("Run schedule is required");

    const job: CronJob = {
      id: nextJobId++,
      name,
      command: "",
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

export async function listCronJobs(
  serverId: number,
): Promise<CronApiResponse<{ items: CronJob[] }>> {
  // Replace with:
  // return NewRequestClient().GET("/api/server/{id}/cron-jobs", {
  //   params: { path: { id: serverId } },
  // });
  return simulateRequest(() => ({ items: jobsFor(serverId) }));
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

export async function renameCronJob(
  serverId: number,
  jobId: number,
  name: string,
): Promise<CronApiResponse<CronJob>> {
  return simulateRequest(() => {
    const job = jobsFor(serverId).find((entry) => entry.id === jobId);
    if (!job) throw new Error("Cron job no longer exists");

    const trimmed = name.trim();
    if (!trimmed) throw new Error("Name is required");

    job.name = trimmed;
    return job;
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
