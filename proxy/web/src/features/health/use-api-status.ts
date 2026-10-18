import { useQuery } from "@tanstack/react-query";
import { getHealth } from "./api";

export function useApiStatus() {
  const query = useQuery({
    queryKey: ["health"],
    queryFn: getHealth,
    refetchInterval: 5000,
    retry: false,
    staleTime: 0,
  });

  return {
    isUnreachable: query.isError,
    isUnhealthy: query.data === false,
    isOnline: query.data === true,
  };
}
