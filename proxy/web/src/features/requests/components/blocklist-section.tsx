import { Button } from "@/components/ui/button";
import { useBlocklist, useBlocklistDelete } from "../hooks/queries";
import { DataTable, type Column } from "@/components/ui/data-table";
import type { BlocklistEntry } from "@/types/api";
import { useState } from "react";

const columns: Column<BlocklistEntry>[] = [
  {
    key: "ip",
    header: "IP",
    cell: (l) => l.ip ?? "",
  },
  {
    key: "ttl",
    header: "TTL (minutes)",
    cell: (l) => formatTtl(l.ttl),
  },
  {
    key: "edit",
    header: "Edit",
    cell: (l) => <UnblockButton ip={l.ip}/>
  }
];

function formatTtl(totalSeconds: number): string {
  const minutes = Math.floor(totalSeconds / 60);
  const seconds = totalSeconds % 60;

  if (minutes === 0) return `${seconds}s`;
  if (seconds === 0) return `${minutes}m`;
  return `${minutes}m ${seconds}s`;
}

export function BlocklistSection() {
  const query = useBlocklist();
  const [page, setPage] = useState(1);

  return (
    <section>
      <div className="wrapper flow">
        <h2 className="heading">Blocklist</h2>
        <DataTable
          columns={columns}
          rows={query.data ?? []}
          getRowKey={(l) => l.ip}
          page={page}
          pageSize={20}
          total={query.data?.length ?? 0}
          onPageChange={setPage}
          isLoading={query.isPending}
          isFetching={query.isFetching}
          emptyText="No blocks to display"
        />
      </div>
    </section>
  );
}

function UnblockButton({ ip }: { ip: string }) {
  const { mutate, isPending } = useBlocklistDelete();

  return (
    <Button
      variant="secondary"
      onClick={() => mutate(ip)}
      disabled={isPending}
    >
      {isPending ? "Unblocking..." : "Unblock"}
    </Button>
  );
}
