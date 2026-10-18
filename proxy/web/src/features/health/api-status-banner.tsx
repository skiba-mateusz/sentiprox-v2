import { useApiStatus } from "./use-api-status";

export function ApiStatusBanner() {
  const { isUnhealthy, isUnreachable } = useApiStatus();

  if (isUnhealthy || isUnreachable) {
    return (
      <div
        className="tile status-dot status-dot--offline"
        style={{ padding: "0.5rem 0.75rem" }}
      >
        Offline
      </div>
    );
  }

  return (
    <div
      className="tile status-dot status-dot--online"
      style={{ padding: "0.5rem 0.75rem" }}
    >
      Online
    </div>
  );
}
