import { isAxiosError } from "axios";
import { apiClient } from "@/lib/api-client";
import {
  type Timeserie,
  type Category,
  type Summary,
  type ApiFailure,
  type ApiSuccess,
  type LogEntry,
  type PaginatedResult,
  type RequestsParams,
} from "@/types/api";

async function get<T>(path: string, params?: RequestsParams): Promise<T> {
  try {
    if (params) {
      path = `${path}?limit=${params.pageSize}&offset=${(params.page - 1) * params.pageSize}&search=${params.searchQuery}&method=${params.method}&since=${params.since}&to=${params.to}`;
    }
    const res = await apiClient.get<ApiSuccess<T>>(path);

    return res.data.data;
  } catch (error: unknown) {
    if (isAxiosError<ApiFailure>(error)) {
      console.log(error.response);
      const message = error.response?.data?.error ?? error.message;

      throw new Error(message, { cause: error });
    }

    throw error;
  }
}

export const api = {
  logs: (params: RequestsParams): Promise<PaginatedResult<LogEntry[]>> =>
    get<PaginatedResult<LogEntry[]>>("/requests", params),
  summary: (): Promise<Summary> => get<Summary>("/requests/summary"),
  timeseries: (): Promise<Timeserie[]> =>
    get<Timeserie[]>("/requests/timeseries"),
  categories: (): Promise<Category[]> =>
    get<Category[]>("/requests/categories"),
};
