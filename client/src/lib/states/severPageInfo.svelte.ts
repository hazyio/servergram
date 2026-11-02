interface SeverInfo {
  id: number;
  name: string;
}
export const severPageInfo = $state<SeverInfo>({
  id: 0,
  name: "Loading server",
});

export function updatePageServerInfo(id: number, name: string) {
  severPageInfo.id = id;
  severPageInfo.name = name;
}
