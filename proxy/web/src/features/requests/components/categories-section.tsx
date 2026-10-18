import { useCategories } from "../hooks/queries";
import { QueryBoundary } from "@/components/query-boundary";
import CategoryChart from "./category-chart";

export function CategoriesSection() {
  const query = useCategories();

  return (
    <section className="categories">
      <div className="wrapper">
        <h2 className="heading">Attack Vectors</h2>
        <QueryBoundary
          query={query}
          isEmpty={(d) => d.length === 0}
          emptyText="No attacks classified yet"
        >
          {(categories) => <CategoryChart data={categories} />}
        </QueryBoundary>
      </div>
    </section>
  );
}
