import axios from "axios";
import camelcaseKeys from "camelcase-keys";
import snakecaseKeys from "snakecase-keys";

const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8090/api/v1";

export const apiClient = axios.create({
  baseURL: API_URL,
});

apiClient.interceptors.response.use((response) => {
  if (response.data && typeof response.data === "object") {
    response.data = camelcaseKeys(response.data, { deep: true });
  }

  return response;
});

apiClient.interceptors.request.use((config) => {
  if (config.params) {
    config.params = snakecaseKeys(config.params, { deep: true });
  }

  if (config.data && typeof config.data === "object") {
    config.data = snakecaseKeys(config.data, { deep: true });
  }

  return config;
});
