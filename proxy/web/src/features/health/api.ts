import { apiClient } from "@/lib/api-client";
import type { ApiSuccess } from "@/types/api";

export async function getHealth(): Promise<boolean> {
  const res = await apiClient.get<ApiSuccess<{ status: string }>>("/health", {
    timeout: 3000,
  });

  return res.data.data.status === "ok";
}
