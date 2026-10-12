import { createContext, useContext, useEffect, useState } from "react";
import api from "../api";

const AuthContext = createContext(null);

export function AuthProvider({ children }) {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const responseInterceptor = api.interceptors.response.use(
      (response) => response,
      (error) => {
        if (error.response?.status === 401) setUser(null);
        return Promise.reject(error);
      },
    );
    api.get("/auth/me").then((response) => setUser(response.data.user)).catch(() => setUser(null)).finally(() => setLoading(false));
    return () => api.interceptors.response.eject(responseInterceptor);
  }, []);

  return <AuthContext.Provider value={{ user, setUser, loading }}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  return useContext(AuthContext);
}
