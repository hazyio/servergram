interface SeverInfo {
  id: number;
  name: string;
}
export const info = $state<SeverInfo>({
  id: 0,
  name: "",
});

export function updateServerInfo(id: number, name: string) {
  info.id = id;
  info.name = name;
}
