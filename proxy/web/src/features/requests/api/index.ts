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
  type BlocklistEntry,
} from "@/types/api";

function handleError(error: unknown): never {
  if (isAxiosError<ApiFailure>(error)) {
    const message = error.response?.data?.error ?? error.message;
    throw new Error(message, { cause: error });
  }

  throw error;
}

async function get<T>(path: string, params?: RequestsParams): Promise<T> {
  try {
    if (params) {
      path = `${path}?limit=${params.pageSize}&offset=${(params.page - 1) * params.pageSize}&search=${params.searchQuery}&method=${params.method}&since=${params.since}&to=${params.to}`;
    }
    const res = await apiClient.get<ApiSuccess<T>>(path);

    return res.data.data;
  } catch (error: unknown) {
    handleError(error);
  }
}

async function remove(path: string): Promise<void> {
  try {
    await apiClient.delete(path);
  } catch (error: unknown) {
    handleError(error);
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
  blocklist: (): Promise<BlocklistEntry[]> =>
    get<BlocklistEntry[]>("/requests/blocklist"),
  blocklistDelete: (ip: string): Promise<void> => remove(`/requests/blocklist/${ip}`)
};
