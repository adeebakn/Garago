export default function SearchInput({ value, onChange, placeholder = "Search records" }) {
  return <input className="module-search" value={value} onChange={(event) => onChange(event.target.value)} placeholder={placeholder} aria-label={placeholder} />;
}
