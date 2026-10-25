import { useId, useState } from "react";
import { RadioGroup, type Option } from "./radio-group";

export interface DateRange {
  from: string;
  to: string;
}

interface DateRangePickerProps {
  value: DateRange | null;
  onChange: (range: DateRange) => void;
}

const toInputDate = (d: Date) => d.toLocaleDateString("sv-SE");

const daysAgo = (n: number) => {
  const d = new Date();
  d.setDate(d.getDate() - n);
  return toInputDate(d);
};

export function DateRangePicker({ value, onChange }: DateRangePickerProps) {
  const id = useId();
  const today = toInputDate(new Date());
  const [option, setOption] = useState<Option<DateRange> | null>(null);

  const presets: Option<DateRange>[] = [
    {
      key: "today",
      name: "Today",
      value: { from: today, to: today },
    },
    {
      key: "7-days",
      name: "7 days",
      value: { from: daysAgo(6), to: today },
    },
    {
      key: "30-days",
      name: "30 days",
      value: { from: daysAgo(29), to: today },
    },
  ];

  function handleOptionChange(opt: Option<DateRange> | null) {
    setOption(opt);

    if (opt === null) {
      onChange({ from: "", to: "" });
    } else {
      onChange(opt.value);
    }
  }

  const currentFrom = value?.from || "";
  const currentTo = value?.to || "";

  return (
    <div className="flex-group items-center">
      <label htmlFor={`${id}-from`}>From</label>
      <input
        id={`${id}-from`}
        type="date"
        value={currentFrom}
        max={currentTo || today}
        onChange={(e) => {
          setOption(null);
          onChange({ from: e.target.value, to: currentTo });
        }}
      />
      <label htmlFor={`${id}-to`}>To</label>
      <input
        id={`${id}-to`}
        type="date"
        value={currentTo}
        min={currentFrom || undefined}
        max={today}
        onChange={(e) => {
          setOption(null);
          onChange({ from: currentFrom, to: e.target.value });
        }}
      />
      <RadioGroup
        value={option}
        onChange={handleOptionChange}
        options={presets}
      />
    </div>
  );
}
