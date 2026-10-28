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
    
  },
  databases: {
    index: "/servers",
    add: "/databases/add",
  },
};

export default appPaths;
