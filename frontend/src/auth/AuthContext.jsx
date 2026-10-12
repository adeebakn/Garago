import { useEffect, useMemo, useState } from "react";
import api from "../api";
import { AuthContext } from "./auth-context";

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

	const value = useMemo(() => ({ user, setUser, loading }), [user, loading]);
	return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}
