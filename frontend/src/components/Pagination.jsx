export default function Pagination({ page, totalItems, pageSize = 10, onPageChange }) {
  const totalPages = Math.max(1, Math.ceil(totalItems / pageSize));
  if (totalItems <= pageSize) return null;
  return <nav className="pagination" aria-label="Table pagination"><span>Showing {(page - 1) * pageSize + 1}-{Math.min(page * pageSize, totalItems)} of {totalItems}</span><div><button type="button" disabled={page === 1} onClick={() => onPageChange(page - 1)}>Previous</button><strong>{page} / {totalPages}</strong><button type="button" disabled={page === totalPages} onClick={() => onPageChange(page + 1)}>Next</button></div></nav>;
}
