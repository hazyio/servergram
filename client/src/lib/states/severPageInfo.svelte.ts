interface SeverInfo {
  id: number;
  name: string;
  lastVisitedDir: string;
}
export const severPageInfo = $state<SeverInfo>({
  id: 0,
  name: "Loading server",
  lastVisitedDir: ""
});

export function updatePageServerInfo(id: number, name: string) {
  severPageInfo.id = id;
  severPageInfo.name = name;
}
