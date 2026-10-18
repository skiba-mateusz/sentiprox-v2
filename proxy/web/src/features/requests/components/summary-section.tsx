import { useSummary } from "../hooks/queries";
import { Activity, ShieldAlert, Ban, Crosshair } from "lucide-react";
import { QueryBoundary } from "@/components/query-boundary";

export function SummarySection() {
  const query = useSummary();

  return (
    <section className="summary">
      <div className="wrapper flow">
        <h2 className="heading">Summary</h2>
        <QueryBoundary query={query}>
          {(summary) => (
            <div className="flex-group">
              <div className="tile flex-1">
                <div className="flex-group h-full items-center">
                  <Activity size={32} className="summary-icon" />
                  <div>
                    <div className="fs-h4">{summary.totalCount}</div>
                    <div className="text-muted">Total requests</div>
                  </div>
                </div>
              </div>
              <div className="tile flex-1">
                <div className="flex-group h-full items-center">
                  <ShieldAlert
                    size={32}
                    className="summary-icon summary-icon--attacks"
                  />
                  <div>
                    <div className="fs-h4">{summary.attacksCount}</div>
                    <div className="text-muted">Attacks detected</div>
                  </div>
                </div>
              </div>
              <div className="tile flex-1">
                <div className="flex-group h-full items-center">
                  <Crosshair size={32} className="summary-icon" />
                  <div>
                    <div className="fs-h4">{summary.topCategory || "-"}</div>
                    <div className="text-muted">Top vector</div>
                  </div>
                </div>
              </div>
              <div className="tile flex-1">
                <div className="flex-group h-full items-center">
                  <Ban
                    size={32}
                    className="summary-icon summary-icon--blocks"
                  />
                  <div>
                    <div className="fs-h4">{summary.totalBlocks}</div>
                    <div className="text-muted">IPs blocked</div>
                  </div>
                </div>
              </div>
            </div>
          )}
        </QueryBoundary>
      </div>
    </section>
  );
}
