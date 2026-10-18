import { useQuery } from "@tanstack/react-query";
import { api } from "../api";

export const requestKeys = {
  all: ["requests"] as const,
  summary: () => [...requestKeys.all, "summary"] as const,
  categories: () => [...requestKeys.all, "categories"] as const,
  timeseries: () => [...requestKeys.all, "timeseries"] as const,
};

export const useSummary = () => {
  return useQuery({
    queryKey: requestKeys.summary(),
    queryFn: api.summary,
    refetchInterval: 10 * 1000,
  });
};

export const useCategories = () => {
  return useQuery({
    queryKey: requestKeys.categories(),
    queryFn: api.categories,
    refetchInterval: 60 * 1000,
  });
};

export const useTimeseries = () => {
  return useQuery({
    queryKey: requestKeys.timeseries(),
    queryFn: api.timeseries,
    refetchInterval: 30 * 1000,
  });
};
