import { useEffect, useRef, useState } from "react";

export default function SearchableSelect({ value, onChange, options, placeholder, required = false, disabled = false, resetKey }) {
  const selected = options.find((option) => String(option.value) === String(value));
  const [query, setQuery] = useState(selected?.label || "");
  const [open, setOpen] = useState(false);
  const root = useRef(null);
  const previousResetKey = useRef(resetKey);
  useEffect(() => {
    if (selected) setQuery(selected.label);
    else if (previousResetKey.current !== resetKey) setQuery("");
    previousResetKey.current = resetKey;
  }, [selected?.label, resetKey]);
  useEffect(() => { const close = (event) => { if (!root.current?.contains(event.target)) setOpen(false); }; document.addEventListener("pointerdown", close, true); return () => document.removeEventListener("pointerdown", close, true); }, []);
  const matches = options.filter((option) => `${option.label} ${option.description || ""}`.toLowerCase().includes(query.toLowerCase()));
  const choose = (option) => { onChange(String(option.value)); setQuery(option.label); setOpen(false); };
  const updateQuery = (nextQuery) => {
    setQuery(nextQuery);
    setOpen(true);
    if (value && nextQuery !== selected?.label) onChange("");
  };
  return <div className="searchable-picker" ref={root}><input disabled={disabled} value={query} onFocus={() => setOpen(true)} onChange={(event) => updateQuery(event.target.value)} placeholder={placeholder} aria-label={placeholder} />{required && <input type="hidden" required value={value} onChange={() => {}} />}{open && !disabled && <div className="searchable-menu">{matches.length === 0 ? <p>No matching options</p> : matches.map((option) => <button type="button" key={option.value} onClick={() => choose(option)}><strong>{option.label}</strong>{option.description && <span>{option.description}</span>}</button>)}</div>}</div>;
}
