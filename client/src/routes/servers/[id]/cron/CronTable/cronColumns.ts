import type { DataTableFeatures } from "#lib/components/base/data-table-features.js";
import SortableTableHeader from "#lib/components/base/SortableTableHeader.svelte";
import { createColumnHelper, renderComponent } from "@tanstack/svelte-table";
import { cronJobDisplayName, type CronJob } from "../cronJobs.ts";
import CronActions from "./CronActions.svelte";
import CronJobName from "./CronJobName.svelte";
import CronRelativeTime from "./CronRelativeTime.svelte";
import CronSchedule from "./CronSchedule.svelte";
import CronStatusBadge from "./CronStatusBadge.svelte";

const columnHelper = createColumnHelper<DataTableFeatures, CronJob>();

const sortableHeader =
  (text: string) =>
  ({
    column,
  }: {
    column: {
      getToggleSortingHandler: () => undefined | ((event: unknown) => void);
    };
  }) =>
    renderComponent(SortableTableHeader, {
      text,
      onclick: column.getToggleSortingHandler(),
    });

export const getColumns = (
  serverId: number,
  onUpdated: (job: CronJob) => void,
  onDeleted: (id: number) => void,
) =>
  columnHelper.columns([
    columnHelper.accessor("status", {
      header: sortableHeader("Status"),
      filterFn: (row, _columnId, filterValue) => {
        if (!filterValue || filterValue === "all") return true;
        if (filterValue === "unmanaged") {
          return row.original.status === "unmanaged";
        }
        return row.original.status === filterValue;
      },
      cell: (info) =>
        renderComponent(CronStatusBadge, { status: info.getValue() }),
    }),
    columnHelper.accessor("name", {
      id: "job",
      header: sortableHeader("Job"),
      filterFn: (row, _columnId, filterValue) => {
        const query = String(filterValue ?? "")
          .trim()
          .toLowerCase();
        if (!query) return true;
        return `${cronJobDisplayName(row.original)} ${row.original.command} ${row.original.schedule}`
          .toLowerCase()
          .includes(query);
      },
      cell: (info) =>
        renderComponent(CronJobName, {
          job: info.row.original,
          serverId,
          duplicate: info.row.original.duplicate,
        }),
    }),
    columnHelper.accessor("schedule", {
      header: sortableHeader("Schedule"),
      cell: (info) =>
        renderComponent(CronSchedule, { schedule: info.getValue() }),
    }),
    columnHelper.accessor("last_run", {
      id: "lastRun",
      header: sortableHeader("Last run"),
      cell: (info) =>
        renderComponent(CronRelativeTime, {
          datetime: info.getValue(),
          empty: "Never",
        }),
    }),
    columnHelper.accessor("next_run", {
      id: "nextRun",
      header: sortableHeader("Next run"),
      cell: (info) =>
        renderComponent(CronRelativeTime, {
          datetime: info.getValue(),
          empty: "Not scheduled",
        }),
    }),
    columnHelper.display({
      id: "actions",
      header: "",
      cell: (info) =>
        renderComponent(CronActions, {
          serverId,
          job: info.row.original,
          onUpdated,
          onDeleted,
        }),
    }),
  ]);
