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
};

export default appPaths;
