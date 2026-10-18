import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  Cell,
  Tooltip,
  ResponsiveContainer,
} from "recharts";
import type { Category } from "@/types/api";
import { catColor, FONT_MONO, tooltipStyle } from "@/utils/colors";

const yAxisTick = {
  fill: "#C7CEDC",
  fontSize: 12,
  fontFamily: FONT_MONO,
};

const ROW_HEIGHT = 34;

interface CategoryChartProps {
  data: Category[];
}

export default function CategoryChart({ data }: CategoryChartProps) {
  return (
    <ResponsiveContainer
      width="100%"
      height={Math.max(120, data.length * ROW_HEIGHT + 20)}
    >
      <BarChart
        data={data}
        layout="vertical"
        margin={{ top: 4, right: 20, left: 6, bottom: 4 }}
      >
        <XAxis type="number" hide allowDecimals={false} />
        <YAxis
          type="category"
          dataKey="name"
          width={94}
          tick={yAxisTick}
          tickLine={false}
          axisLine={false}
        />
        <Tooltip
          cursor={{ fill: "rgba(255,255,255,0.03)" }}
          contentStyle={tooltipStyle}
          labelStyle={{ color: "#E6EAF2" }}
          itemStyle={{ color: "#8A94A8" }}
          formatter={(value, _name, item) => [
            <span style={{ color: catColor(item.payload.name) }}>{value}</span>,
            "Attacks",
          ]}
        />
        <Bar dataKey="count" name="Attacks" radius={[0, 4, 4, 0]} barSize={16}>
          {data.map((d) => (
            <Cell key={d.name} fill={catColor(d.name)} />
          ))}
        </Bar>
      </BarChart>
    </ResponsiveContainer>
  );
}
