interface SeverInfo {
  id: number;
  name: string;
  lastVisitedDir: string;
}
export const severPageInfo = $state<SeverInfo>({
  id: 0,
  name: "Loading server",
  lastVisitedDir: "",
});

interface CronJobInfo {
  id: number;
  name: string;
}
export const cronJobPageInfo = $state<CronJobInfo>({
  id: 0,
  name: "",
});

export function updatePageServerInfo(id: number, name: string) {
  severPageInfo.id = id;
  severPageInfo.name = name;
}

export function updateCronJobPageInfo(id: number, name: string) {
  cronJobPageInfo.id = id;
  cronJobPageInfo.name = name;
}
