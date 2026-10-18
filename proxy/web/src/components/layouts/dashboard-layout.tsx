import { Header } from "./header";
import { Sidebar } from "./sidebar";

export function DashboardLayout({ children }: React.PropsWithChildren) {
  return (
    <>
      <Header />
      <Sidebar />
      <main>{children}</main>
    </>
  );
}
