import { useLogs } from "../hooks/queries";
import { DataTable } from "@/components/ui/data-table";
import { SearchInput } from "@/components/ui/search-input";
import {
  type DateRange,
  DateRangePicker,
} from "@/components/ui/date-range-picker";
import { useEffect, useState } from "react";
import { type Option, RadioGroup } from "@/components/ui/radio-group";
import { columns, options } from "../constants";
import { useRequestsParams } from "../hooks/use-requests-params";
import { Button } from "@/components/ui/button";

export function LogsSection() {
  const { params, page, setPage, onSearch, setParams } = useRequestsParams(20);
  const query = useLogs(params);
  const [range, setRange] = useState<DateRange | null>(null);
  const [option, setOption] = useState<Option<string> | null>(null);

  useEffect(() => {
    setParams({
      method: option ? option.value : "",
      since: range ? range.from : "",
      to: range ? range.to : "",
      page: 1,
    });
  }, [range, option, setParams]);

  useEffect(() => {
    document.querySelector("main")?.scrollTo({
      top: 0,
      behavior: "smooth",
    });
  }, [page]);

  return (
    <section className="logs">
      <div className="wrapper flow">
        <h2 className="heading">Logs</h2>
        <div className="flow">
          <div className="flex-group space-between items-center">
            <SearchInput
              onSearch={onSearch}
              placeholder="IP | Path | Query | Code ..."
            />
            <DateRangePicker onChange={setRange} value={range} />
          </div>
          <RadioGroup options={options} value={option} onChange={setOption} />
        </div>
        {query.isError ? (
          <div className="alert">
            <p>{query.error.message}</p>
            <Button size="small" onClick={() => query.refetch()}>
              Try again
            </Button>
          </div>
        ) : (
          <DataTable
            columns={columns}
            rows={query.data?.items ?? []}
            getRowKey={(l) => l.id}
            page={page}
            pageSize={params.pageSize}
            total={query.data?.total ?? 0}
            onPageChange={setPage}
            isLoading={query.isPending}
            isFetching={query.isFetching}
            emptyText="No logs to display"
          />
        )}
      </div>
    </section>
  );
}
