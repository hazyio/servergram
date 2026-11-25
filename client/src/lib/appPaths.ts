const appPaths = {
  home: "/",
  sshKeys: {
    index: "/ssh-keys",
    resources: (id: number) => `/ssh-keys/${id}/resources`,
  },
  groups: {
    index: "/groups",
    resources: (id: number) => `/groups/${id}/resources`,
  },
  account: {
    login: "/account/login",
    sessions: "/sessions",
  },
  servers: {
    index: "/",
    add: "//add",
    overview: (id: number) => `//${id}/overview`,
    files: (id: number) => `//${id}/files`,
    processes: (id: number) => `//${id}/processes`,
    docker: (id: number) => `//${id}/docker`,
    cron: (id: number) => `//${id}/cron`,
    cronAdd: (id: number) => `//${id}/cron/add`,
    cronJob: (id: number, jobId: number) => `//${id}/cron/${jobId}`,
  },
  databases: {
    index: "/databases",
    add: "/databases/add",
  },
  cronJobs: {
    index: "/cron-jobs",
    add: "/cron-jobs/add",
  },
  logs: {
    index: "/logs",
  },
};

export default appPaths;
