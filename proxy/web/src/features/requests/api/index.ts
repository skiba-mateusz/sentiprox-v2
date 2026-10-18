import { isAxiosError } from "axios";
import { apiClient } from "@/lib/api-client";
import {
  type Timeserie,
  type Category,
  type Summary,
  type ApiFailure,
  type ApiSuccess,
} from "@/types/api";

async function get<T>(path: string): Promise<T> {
  try {
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
  summary: (): Promise<Summary> => get<Summary>("/requests/summary"),
  timeseries: (): Promise<Timeserie[]> =>
    get<Timeserie[]>("/requests/timeseries"),
  categories: (): Promise<Category[]> =>
    get<Category[]>("/requests/categories"),
};
