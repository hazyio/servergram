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
  },
  servers: {
    index: "/servers",
    add: "/servers/add",
    overview: (id: number) => `/servers/${id}/overview`,
    files: (id: number) => `/servers/${id}/files`,
    process: (id: number) => `/servers/${id}/process`,
    docker: (id: number) => `/servers/${id}/docker`,
    cron: (id: number) => `/servers/${id}/cron`,
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
