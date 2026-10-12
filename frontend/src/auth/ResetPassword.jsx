import { useState } from "react";
import { Link, useNavigate, useSearchParams } from "react-router-dom";
import api from "../api";
import { apiErrorMessage } from "../utils/api-error";
import "./auth.css";

export default function ResetPassword() {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [error, setError] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const token = searchParams.get("token") || "";

  async function submit(event) {
    event.preventDefault();
    setError("");
    if (!token) return setError("This password reset link is invalid or incomplete.");
    if (password !== confirmPassword) return setError("Passwords do not match.");
    setSubmitting(true);
    try {
      await api.post("/auth/reset-password", { token, password });
      navigate("/login", { replace: true, state: { reset: true } });
    } catch (requestError) {
      setError(apiErrorMessage(requestError, "Unable to reset password. Please try again."));
    } finally {
      setSubmitting(false);
    }
  }

  return <main className="auth-page"><section className="auth-intro"><span className="brand-mark">Garago</span><h1>Choose a new password.</h1><p>Your new password will sign out every existing session.</p></section><form className="auth-card" onSubmit={submit}><p className="eyebrow">PASSWORD RESET</p><h2>Set a new password</h2>{error && <p className="error">{error}</p>}<label>New password<input required minLength="8" type="password" autoComplete="new-password" value={password} onChange={(event) => setPassword(event.target.value)} /></label><label>Confirm new password<input required minLength="8" type="password" autoComplete="new-password" value={confirmPassword} onChange={(event) => setConfirmPassword(event.target.value)} /></label><button disabled={submitting}>{submitting ? "Saving…" : "Reset password"}</button><p className="auth-switch"><Link to="/login">Back to sign in</Link></p></form></main>;
}
