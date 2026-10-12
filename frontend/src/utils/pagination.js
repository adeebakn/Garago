export function pageRows(items, page, pageSize = 10) {
  return items.slice((page - 1) * pageSize, page * pageSize);
}
