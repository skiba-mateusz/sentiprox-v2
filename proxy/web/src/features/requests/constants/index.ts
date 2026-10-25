import type { Column } from "@/components/ui/data-table";
import type { Option } from "@/components/ui/radio-group";
import type { LogEntry } from "@/types/api";

export const columns: Column<LogEntry>[] = [
  {
    key: "time",
    header: "Time",
    cell: (l) => new Date(l.timestamp).toLocaleString(),
  },
  { key: "ip", header: "IP", cell: (l) => l.ip ?? "" },
  { key: "method", header: "Method", cell: (l) => l.method ?? "" },
  { key: "path", header: "Path", cell: (l) => l.path ?? "" },
  { key: "query", header: "Query", cell: (l) => l.query ?? "" },
  {
    key: "content-type",
    header: "Content Type",
    cell: (l) => l.contentType ?? "-",
  },
  { key: "category", header: "Category", cell: (l) => l.category ?? "-" },
  { key: "action", header: "Action", cell: (l) => l.action ?? "-" },
  { key: "score", header: "Score", cell: (l) => l.score.toFixed(2) ?? "-" },
  { key: "source", header: "Source", cell: (l) => l.source ?? "-" },
  {
    key: "status-code",
    header: "Status Code",
    cell: (l) => l.statusCode ?? "-",
  },
];

export const options: Option<string>[] = [
  {
    key: "get",
    name: "GET",
    value: "GET",
  },
  {
    key: "post",
    name: "POST",
    value: "POST",
  },
  {
    key: "put",
    name: "PUT",
    value: "PUT",
  },
  {
    key: "patch",
    name: "PATCH",
    value: "PATCH",
  },
  {
    key: "delete",
    name: "DELETE",
    value: "DELETE",
  },
];
