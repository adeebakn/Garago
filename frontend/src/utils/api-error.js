export function apiErrorMessage(error, fallback) {
  return error?.response?.data?.error || fallback;
}
