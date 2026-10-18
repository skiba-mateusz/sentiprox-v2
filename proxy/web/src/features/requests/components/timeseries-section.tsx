import { useTimeseries } from "../hooks/queries";
import { QueryBoundary } from "@/components/query-boundary";
import TrafficChart from "./traffic-chart";

export function TimeseriesSection() {
  const query = useTimeseries();

  return (
    <section className="timeseries">
      <div className="wrapper flow">
        <h2 className="heading">Timeline</h2>
        <QueryBoundary
          query={query}
          isEmpty={(d) => d.length === 0}
          emptyText="No traffic in last hour"
        >
          {(timeseries) => <TrafficChart data={timeseries} />}
        </QueryBoundary>
      </div>
    </section>
  );
}
