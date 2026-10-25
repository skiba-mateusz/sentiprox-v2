import { Loader2 } from "lucide-react";

interface LoaderProps {
  stacked?: boolean;
}

export function Loader({ stacked }: LoaderProps) {
  return (
    <div>
      <Loader2
        className={`loader-icon ${stacked ? "loader-icon--stacked" : ""}`}
      />
    </div>
  );
}
