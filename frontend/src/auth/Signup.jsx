import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import api from "../api";
import "./auth.css";

export default function Signup() {
  const [form, setForm] = useState({ name: "", email: "", password: "", confirmPassword: "" });
  const [error, setError] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const navigate = useNavigate();
  const submit = async (event) => {
    event.preventDefault(); setError("");
    if (form.password !== form.confirmPassword) { setError("Passwords do not match."); return; }
    setSubmitting(true);
    try { await api.post("/auth/signup", { name: form.name, email: form.email, password: form.password }); navigate("/login", { state: { created: true } }); }
    catch (err) { setError(err.response?.data?.error || "Could not create account."); }
    finally { setSubmitting(false); }
  };
  const change = (key) => (event) => setForm({ ...form, [key]: event.target.value });
  return <main className="auth-page"><section className="auth-intro"><span className="brand-mark">A</span><h1>Bring every job into focus.</h1><p>Set up your manager workspace in less than a minute.</p></section><form className="auth-card" onSubmit={submit}><p className="eyebrow">GET STARTED</p><h2>Create manager account</h2>{error && <p className="error">{error}</p>}<label>Name<input required value={form.name} onChange={change("name")} /></label><label>Email<input required type="email" value={form.email} onChange={change("email")} /></label><label>Password<input required minLength="8" type="password" value={form.password} onChange={change("password")} /></label><label>Confirm password<input required minLength="8" type="password" value={form.confirmPassword} onChange={change("confirmPassword")} /></label><button disabled={submitting}>{submitting ? "Creating…" : "Create account"}</button><p className="auth-switch">Already registered? <Link to="/login">Sign in</Link></p></form></main>;
}
