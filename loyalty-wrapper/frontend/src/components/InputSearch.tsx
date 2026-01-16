import React, { useState } from "react";
import { useSearchParams } from "react-router-dom";

interface InputSearchProps {
  onSearchChange: (newSearch: string) => void;
}

export default function InputSearch({ onSearchChange }: InputSearchProps) {
  const [searchParams, setSearchParams] = useSearchParams();
  const initialSearch = searchParams.get("search") || "";
  const [value, setValue] = useState(initialSearch);

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "Enter") {
      onSearchChange(value);
      if (value) {
        searchParams.set("search", value);
      } else {
        searchParams.delete("search");
      }
      setSearchParams(searchParams);
    }
  };

  return (
    <label className="input my-1 w-full sm:w-auto">
      <input
        type="search"
        className="grow"
        placeholder="Buscador"
        value={value}
        onChange={(e) => setValue(e.target.value)}
        onKeyDown={handleKeyDown}
      />
    </label>
  );
}
