export interface Option<T> {
  key: string;
  name: string;
  value: T;
}

interface RadioGroupProps<T> {
  value: Option<T> | null;
  defaultValue?: Option<T>;
  onChange: (option: Option<T> | null) => void;
  options: Option<T>[];
}

export function RadioGroup<T>({
  value,
  onChange,
  options,
}: RadioGroupProps<T>) {
  function handleClick(opt: Option<T>) {
    if (value?.key == opt.key) {
      onChange(null);
    } else {
      onChange(opt);
    }
  }

  return (
    <div className="radio-group">
      {options.map((opt) => (
        <button
          className={`radio-group__btn ${opt.key === value?.key ? "radio-group__btn--active" : ""}`}
          key={opt.key}
          onClick={() => handleClick(opt)}
        >
          {opt.name}
        </button>
      ))}
    </div>
  );
}
