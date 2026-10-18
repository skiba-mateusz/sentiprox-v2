import { useMemo } from "react";
import {
  AreaChart,
  Area,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from "recharts";
import type { Timeserie } from "@/types/api";
import { FONT_MONO, tooltipStyle } from "@/utils/colors";

const axisTick = {
  fill: "#8A94A8",
  fontSize: 11,
  fontFamily: FONT_MONO,
};

interface TrafficChartProps {
  data: Timeserie[];
}

export default function TrafficChart({ data }: TrafficChartProps) {
  const chartData = useMemo(
    () =>
      data.map((p) => ({
        time: p.time,
        benign: Math.max(p.total - p.attacks, 0),
        attack: p.attacks,
      })),
    [data],
  );

  return (
    <ResponsiveContainer width="100%" height={240}>
      <AreaChart
        data={chartData}
        margin={{ top: 8, right: 24, left: 0, bottom: 0 }}
      >
        <defs>
          <linearGradient id="gBenign" x1="0" y1="0" x2="0" y2="1">
            <stop offset="0%" stopColor="#3FB68B" stopOpacity={0.45} />
            <stop offset="100%" stopColor="#3FB68B" stopOpacity={0.04} />
          </linearGradient>
          <linearGradient id="gAttack" x1="0" y1="0" x2="0" y2="1">
            <stop offset="0%" stopColor="#E8634A" stopOpacity={0.7} />
            <stop offset="100%" stopColor="#E8634A" stopOpacity={0.05} />
          </linearGradient>
        </defs>
        <CartesianGrid
          stroke="#243049"
          strokeDasharray="2 5"
          vertical={false}
        />
        <XAxis
          dataKey="time"
          tick={axisTick}
          interval="preserveStartEnd"
          minTickGap={40}
          tickLine={false}
          axisLine={{ stroke: "#243049" }}
        />
        <YAxis
          tick={axisTick}
          tickLine={false}
          axisLine={false}
          width={34}
          allowDecimals={false}
        />
        <Tooltip
          contentStyle={tooltipStyle}
          labelStyle={{ color: "#E6EAF2" }}
        />
        <Legend
          iconType="circle"
          wrapperStyle={{ fontSize: 12, fontFamily: FONT_MONO }}
        />
        <Area
          type="monotone"
          dataKey="benign"
          name="Benign"
          stackId="1"
          stroke="#3FB68B"
          fill="url(#gBenign)"
          strokeWidth={1.5}
        />
        <Area
          type="monotone"
          dataKey="attack"
          name="Attacks"
          stackId="1"
          stroke="#E8634A"
          fill="url(#gAttack)"
          strokeWidth={1.5}
        />
      </AreaChart>
    </ResponsiveContainer>
  );
}
