import { SummarySection } from "./summary-section";
import { TimeseriesSection } from "./timeseries-section";
import { CategoriesSection } from "./categories-section";

export function DashboardSummary() {
  return (
    <>
      <SummarySection />
      <TimeseriesSection />
      <CategoriesSection />
    </>
  );
}
