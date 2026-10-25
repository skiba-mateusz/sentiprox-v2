import type { UseQueryResult } from "@tanstack/react-query";
import type { ReactNode } from "react";
import { Loader } from "./ui/loader";
import { Button } from "./ui/button";

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
  if (query.isPending) return <Loader />;

  if (query.isError) {
    return (
      <div className="alert">
        <p>{query.error.message}</p>
        <Button size="small" onClick={() => query.refetch()}>
          Try again
        </Button>
      </div>
    );
  }

  if (isEmpty?.(query.data)) return <div className="message">{emptyText}</div>;

  return <>{children(query.data)}</>;
}
