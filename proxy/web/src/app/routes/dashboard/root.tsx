import { DashboardLayout } from "@/components/layouts/dashboard-layout";
import { Outlet } from "react-router-dom";

export function DashboardRoot() {
  return (
    <DashboardLayout>
      <Outlet />
    </DashboardLayout>
  );
}
