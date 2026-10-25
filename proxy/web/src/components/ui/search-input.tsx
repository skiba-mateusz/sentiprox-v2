import { useEffect, useState } from "react";

interface SearchInputProps {
  onSearch: (value: string) => void;
  placeholder?: string;
  delay?: number;
}

export function SearchInput({
  onSearch,
  placeholder = "Search...",
  delay = 400,
}: SearchInputProps) {
  const [value, setValue] = useState("");

  useEffect(() => {
    const id = setTimeout(() => onSearch(value.trim()), delay);

    return () => clearTimeout(id);
  }, [value, delay, onSearch]);

  return (
    <input
      className="search-input"
      type="search"
      value={value}
      onChange={(e) => setValue(e.target.value)}
      placeholder={placeholder}
    />
  );
}
