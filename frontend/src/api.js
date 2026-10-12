import axios from "axios";

// VITE_API_BASE_URL can be set in frontend/.env (for example http://localhost:8080/api).
export default axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || "http://localhost:8080/api",
  withCredentials: true,
});
