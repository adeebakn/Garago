import { useEffect, useState } from "react";
import { Link, useLocation, useNavigate } from "react-router-dom";
import api from "../api";
import { useAuth } from "./AuthContext";
import "./auth.css";

export default function Login() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const { setUser } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();
  const [created] = useState(() => Boolean(location.state?.created));
  const [reset] = useState(() => Boolean(location.state?.reset));
  useEffect(() => {
    if (location.state?.created || location.state?.reset) navigate("/login", { replace: true, state: null });
  }, [location.state, navigate]);
  const submit = async (event) => {
    event.preventDefault(); setError(""); setSubmitting(true);
    try {
      const response = await api.post("/auth/login", { email, password });
      setUser(response.data.user); navigate("/");
    } catch (err) { setError(err.response?.data?.error || "Could not sign in."); }
    finally { setSubmitting(false); }
  };
  return <main className="auth-page"><section className="auth-intro"><p className="eyebrow">Garago</p><h1>Run your workshop with clarity.</h1><p>Customers, services, invoices, and payments—one focused workspace.</p></section><form className="auth-card" onSubmit={submit}><p className="eyebrow">WELCOME </p><h2>Sign in to Garago</h2>{created && <p className="success">Account created successfully. Please sign in.</p>}{reset && <p className="success">Password reset successfully. Please sign in with your new password.</p>}{error && <p className="error">{error}</p>}<label>Email<input required type="email" value={email} onChange={(e) => setEmail(e.target.value)} /></label><label>Password<input required type="password" value={password} onChange={(e) => setPassword(e.target.value)} /></label><p className="auth-switch auth-forgot"><Link to="/forgot-password">Forgot password?</Link></p><button disabled={submitting}>{submitting ? "Signing in…" : "Sign in"}</button><p className="auth-switch">New user? <Link to="/signup">Create an account</Link></p></form></main>;
}
