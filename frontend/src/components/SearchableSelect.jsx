import { useEffect, useRef, useState } from "react";

export default function SearchableSelect(props) {
  return <SearchableSelectInput key={props.resetKey} {...props} />;
}

function SearchableSelectInput({ value, onChange, options, placeholder, required = false, disabled = false }) {
  const selected = options.find((option) => String(option.value) === String(value));
  const [draft, setDraft] = useState(null);
  const query = draft ?? selected?.label ?? "";
  const [open, setOpen] = useState(false);
  const root = useRef(null);
  useEffect(() => { const close = (event) => { if (!root.current?.contains(event.target)) setOpen(false); }; document.addEventListener("pointerdown", close, true); return () => document.removeEventListener("pointerdown", close, true); }, []);
  const matches = options.filter((option) => `${option.label} ${option.description || ""}`.toLowerCase().includes(query.toLowerCase()));
  const choose = (option) => { onChange(String(option.value)); setDraft(null); setOpen(false); };
  const updateQuery = (nextQuery) => {
    setDraft(nextQuery);
    setOpen(true);
    if (value && nextQuery !== selected?.label) onChange("");
  };
  return <div className="searchable-picker" ref={root}><input disabled={disabled} value={query} onFocus={() => setOpen(true)} onChange={(event) => updateQuery(event.target.value)} placeholder={placeholder} aria-label={placeholder} />{required && <input type="hidden" required value={value} onChange={() => {}} />}{open && !disabled && <div className="searchable-menu">{matches.length === 0 ? <p>No matching options</p> : matches.map((option) => <button type="button" key={option.value} onClick={() => choose(option)}><strong>{option.label}</strong>{option.description && <span>{option.description}</span>}</button>)}</div>}</div>;
}
