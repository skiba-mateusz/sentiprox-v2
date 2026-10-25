import type { ReactNode } from "react";
import { Loader } from "./loader";
import { Pagination } from "./pagination";

export interface Column<T> {
  key: string;
  header: string;
  cell: (row: T) => ReactNode;
}

interface DataTableProps<T> {
  columns: Column<T>[];
  rows: T[];
  getRowKey: (row: T) => string | number;
  page: number;
  pageSize: number;
  total: number;
  onPageChange: (page: number) => void;
  isLoading?: boolean;
  isFetching?: boolean;
  emptyText?: string;
}

export function DataTable<T>({
  columns,
  rows,
  getRowKey,
  page,
  pageSize,
  total,
  onPageChange,
  emptyText = "Nothing to display",
  isLoading = false,
  isFetching = false,
}: DataTableProps<T>) {
  if (isLoading) {
    return <Loader />;
  }

  const isChangingPage = isFetching && !isLoading;

  return (
    <div className="flow">
      {total === 0 ? <div className="message">{emptyText}</div> : ""}
      <div className="table-wrapper" style={{ position: "relative" }}>
        {isChangingPage ? <Loader stacked /> : ""}
        <table className="table" style={{ opacity: isChangingPage ? 0.5 : 1 }}>
          <thead className="table__header">
            <tr className="table__row">
              {columns.map((c) => (
                <th key={c.key}>{c.header}</th>
              ))}
            </tr>
          </thead>
          <tbody className="table__body">
            {rows.map((row) => (
              <tr key={getRowKey(row)}>
                {columns.map((c) => (
                  <td key={c.key} style={{}}>
                    {c.cell(row)}
                  </td>
                ))}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      <Pagination
        page={page}
        pageSize={pageSize}
        total={total}
        setPage={onPageChange}
      />
    </div>
  );
}
