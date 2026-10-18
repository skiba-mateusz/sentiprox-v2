import type { UseQueryResult } from "@tanstack/react-query";
import { Loader2 } from "lucide-react";
import type { ReactNode } from "react";

type Props<T> = {
  query: UseQueryResult<T>;
  children: (data: T) => ReactNode;
  isEmpty?: (data: T) => boolean;
  emptyText?: string;
};

export function QueryBoundary<T>({
  query,
  children,
  isEmpty,
  emptyText = "No data to display",
}: Props<T>) {
  if (query.isPending) return <Loader2 />;

  if (query.isError) {
    return (
      <div role="alert">
        <p>Error: {query.error.message}</p>
        <button onClick={() => query.refetch()}>Try again</button>
      </div>
    );
  }

  if (isEmpty?.(query.data)) return <p>{emptyText}</p>;

  return <>{children(query.data)}</>;
}
