import type { RequestsParams } from "@/types/api";
import { useCallback, useState } from "react";

export function useRequestsParams(pageSize = 20) {
  const [params, setParamsState] = useState<RequestsParams>({
    page: 1,
    pageSize,
    searchQuery: "",
    method: "",
    since: "",
    to: "",
  });

  const setParams = useCallback((patch: Partial<RequestsParams>) => {
    setParamsState((prev) => ({ ...prev, ...patch }));
  }, []);

  const setPage = useCallback(
    (page: number) => setParams({ page }),
    [setParams],
  );

  const onSearch = useCallback(
    (searchQuery: string) => setParams({ searchQuery, page: 1 }),
    [setParams],
  );

  return { params, page: params.page, setPage, onSearch, setParams };
}
