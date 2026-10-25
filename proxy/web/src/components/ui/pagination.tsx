import { ChevronLeft, ChevronRight } from "lucide-react";
import { Button } from "./button";

interface PaginationProps {
  page: number;
  pageSize: number;
  total: number;
  setPage: (page: number) => void;
}

export function Pagination({
  page,
  pageSize,
  total,
  setPage,
}: PaginationProps) {
  const from = total == 0 ? 0 : pageSize * (page - 1) + 1;
  const to = Math.min(from + pageSize - 1, total);
  const pageCount = Math.ceil(total / pageSize);

  return (
    <div className="flex-group space-between">
      <div>
        {from}-{to} from {total}
      </div>
      <div className="flex-group items-center">
        <Button
          variant="circle"
          onClick={() => setPage(page - 1)}
          disabled={page <= 1}
        >
          <ChevronLeft />
        </Button>
        <span>
          {total == 0 ? 0 : page} / {pageCount}
        </span>
        <Button
          variant="circle"
          onClick={() => setPage(page + 1)}
          disabled={page >= pageCount}
        >
          <ChevronRight />
        </Button>
      </div>
    </div>
  );
}
