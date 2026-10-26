import { createBrowserRouter, RouterProvider } from "react-router-dom";
import { DashboardRoot } from "./routes/dashboard/root";
import { DashboardSummary } from "@/features/requests/components/dashboard";
import { DashboardLogs } from "@/features/requests/components/dashboard-logs";
import { DashboardBlocklist } from "@/features/requests/components/dashboard-blocklist";

const router = createBrowserRouter([
  {
    path: "/",
    element: <DashboardRoot />,
    children: [
      {
        path: "",
        element: <DashboardSummary />,
      },
      {
        path: "/logs",
        element: <DashboardLogs />,
      },
      {
        path: "/blocklist",
        element: <DashboardBlocklist />,
      },
    ],
  },
]);

export const AppRouter = () => {
  return <RouterProvider router={router} />;
};
